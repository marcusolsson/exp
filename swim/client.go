package main

import (
	"io"
	"net"
)

// Client holds the client connection.
type Client struct {
	conn net.Conn
}

// Close closes the client connection.
func (c *Client) Close() {
	c.conn.Close()
}

// NewClient returns a new instance of Client.
func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn}, nil
}

func (c *Client) sendMessage(t messageType, m interface{}) (io.Reader, error) {
	b, err := encodeMessage(t, m)
	if err != nil {
		return nil, err
	}

	c.conn.Write(b)

	return io.Reader(c.conn), nil
}
