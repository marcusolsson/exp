package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestEncodeMessage(t *testing.T) {
	msg := &messageJoin{
		Name:    "test",
		Address: "addr",
	}

	b, err := encodeMessage(joinType, msg)
	if err != nil {
		t.Fatal(err)
	}

	if messageType(b[0]) != joinType {
		t.Fatal("unexpected message type")
	}

	var m messageJoin
	if err := decodeMessage(bytes.NewReader(b[1:]), &m); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, &m) {
		t.Fatalf("expected %s, got %s", msg, m)
	}
}
