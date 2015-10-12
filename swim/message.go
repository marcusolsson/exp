package main

import (
	"bytes"
	"encoding/json"
	"io"
)

type messageType uint8

const (
	joinType messageType = 1 << iota
	joinResponseType
	queryType
	queryResponseType
)

type messageJoin struct {
	Name    string
	Address string
}

type messageJoinResponse struct {
	Members List
}

type messageQuery struct {
	Name    string
	Updates []Update
	Data    []byte
}

type messageQueryResponse struct {
	Updates []Update
	Ack     bool
}

func encodeMessage(t messageType, m interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(t))

	if err := json.NewEncoder(buf).Encode(m); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func decodeMessage(r io.Reader, out interface{}) error {
	if err := json.NewDecoder(r).Decode(&out); err != nil {
		return err
	}

	return nil
}
