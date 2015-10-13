package main

import (
	"errors"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// Server contains the server context.
type Server struct {
	BindAddr string

	Members *List
	Self    Member

	GossipInterval time.Duration

	listener net.Listener

	Logger *log.Logger
}

// NewServer returns a new instance of Server.
func NewServer(bindAddr string, interval int, logger *log.Logger) *Server {
	return &Server{BindAddr: bindAddr,
		Members:        NewList(2),
		Self:           Member{Address: bindAddr},
		GossipInterval: time.Duration(interval) * time.Millisecond,
		Logger:         logger,
	}
}

// Start ...
func (s *Server) Start() error {
	l, err := net.Listen("tcp", s.BindAddr)
	if err != nil {
		return err
	}
	s.listener = l

	// Add myself to the local membership list.
	s.Members.Add(s.Self)

	// Start gossiping.
	go s.gossip()

	return nil
}

// Join ...
func (s *Server) Join(addr string) error {
	if addr == "" {
		return errors.New("missing address")
	}

	msg := messageJoin{
		// TODO(marcusolsson): Find better member name.
		Name:    s.BindAddr,
		Address: s.BindAddr,
	}

	// Send join message.
	resp, err := sendJoin(addr, msg)
	if err != nil {
		return err
	}

	// Initialize local membership list.
	clone := resp.Members
	s.Members = &clone

	return nil
}

// Ping ...
func (s *Server) Ping(addr string) error {
	msg := messageQuery{
		Name:    "ping",
		Updates: s.Members.Updates,
	}

	// Send ping message query.
	resp, err := sendQuery(addr, msg)
	if err != nil {
		return err
	}

	// Add the events received from the node.
	s.Members.Merge(resp.Updates)

	return nil
}

// PingReq ...
func (s *Server) PingReq(m Member, target Member) error {
	msg := messageQuery{
		Name:    "ping-req",
		Updates: s.Members.Updates,
		Data:    []byte(target.Address),
	}

	// Send ping-req query.
	resp, err := sendQuery(m.Address, msg)
	if err != nil {
		return err
	}

	if !resp.Ack {
		return errors.New("ack not received")
	}

	// Add the events received from the node.
	s.Members.Merge(resp.Updates)

	return nil
}

// Listen listens for incoming connections.
func (s *Server) Listen() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}

		// Read first byte to determine message type.
		buf := make([]byte, 1)
		if _, err := conn.Read(buf); err != nil {
			return err
		}

		switch messageType(buf[0]) {
		case joinType:
			var m messageJoin
			if err := decodeMessage(conn, &m); err != nil {
				return err
			}
			s.handleJoin(conn, m)
		case queryType:
			var m messageQuery
			if err := decodeMessage(conn, &m); err != nil {
				return err
			}
			switch m.Name {
			case "ping":
				s.handlePing(conn, m)
			case "ping-req":
				s.handlePingReq(conn, m)
			}
		case queryResponseType:
			var m messageQueryResponse
			if err := decodeMessage(conn, &m); err != nil {
				return err
			}
		default:
			s.Logger.Println("unrecognized message type")
		}
	}
}

// gossip runs the SWIM protocol.
func (s *Server) gossip() {
	for {
		<-time.After(s.GossipInterval)

		// Increase round number
		s.Members.IncrementRound()

		// Select one random node to ping.
		m, err := s.Members.Random(1, s.Self)
		if err != nil {
			continue
		}

		node := m[0]

		if err := s.Ping(node.Address); err != nil {
			s.Logger.Println("ping: failed to ping", node.Address)

			k := 3
			randmem, err := s.Members.Random(k, s.Self, node)
			if err != nil {
				s.Logger.Println(err)
			}

			var ok bool
			for err := range s.sendPingReq(node, randmem) {
				if err == nil {
					ok = true
					break
				}
			}

			if !ok {
				s.Logger.Println("ping-req: ack was not received, removing node", node.Address)

				s.Members.Remove(node)
			}
		}
	}
}

func (s *Server) sendPingReq(node Member, members []Member) <-chan error {
	var wg sync.WaitGroup
	wg.Add(len(members))

	ch := make(chan error)

	for _, m := range members {
		go func(m Member) {
			ch <- s.PingReq(m, node)
			wg.Done()
		}(m)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

func (s *Server) handleJoin(w io.Writer, req messageJoin) {
	s.Members.Add(Member{
		Name:    req.Name,
		Address: req.Address,
	})

	s.Logger.Printf("join: member %s", req.Address)

	b, err := encodeMessage(joinResponseType, messageJoinResponse{Members: *s.Members})
	if err != nil {
		s.Logger.Println(err)
	}

	w.Write(b)
}

func (s *Server) handlePing(w io.Writer, req messageQuery) {
	s.Members.Merge(req.Updates)

	b, err := encodeMessage(queryResponseType, messageQueryResponse{Updates: s.Members.Updates})
	if err != nil {
		s.Logger.Println(err)
	}

	w.Write(b)
}

func (s *Server) handlePingReq(w io.Writer, req messageQuery) {
	s.Members.Merge(req.Updates)

	ack := true
	if err := s.Ping(string(req.Data)); err != nil {
		ack = false
	}

	b, err := encodeMessage(queryResponseType, messageQueryResponse{Updates: s.Members.Updates, Ack: ack})
	if err != nil {
		s.Logger.Println(err)
		return
	}

	w.Write(b)
}

func sendJoin(addr string, msg messageJoin) (messageJoinResponse, error) {
	var resp messageJoinResponse

	c, err := NewClient(addr)
	if err != nil {
		return resp, err
	}
	defer c.Close()

	r, err := c.sendMessage(joinType, msg)
	if err != nil {
		return resp, err
	}

	buf := make([]byte, 1)
	if _, err := r.Read(buf); err != nil {
		return resp, err
	}

	if messageType(buf[0]) != joinResponseType {
		return resp, err
	}

	if err := decodeMessage(r, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

func sendQuery(addr string, q messageQuery) (messageQueryResponse, error) {
	var response messageQueryResponse

	c, err := NewClient(addr)
	if err != nil {
		return response, err
	}

	// Send query to node.
	body, err := c.sendMessage(queryType, q)
	if err != nil {
		return response, err
	}

	buf := make([]byte, 1)
	if _, err := body.Read(buf); err != nil {
		return response, err
	}

	if messageType(buf[0]) != queryResponseType {
		return response, errors.New("unrecognized message type")
	}

	if err := decodeMessage(body, &response); err != nil {
		return response, err
	}

	return response, nil
}
