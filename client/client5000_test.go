package main

import (
	"log"
	"net"
	"testing"
)

func TestClient(t *testing.T) {

	conn, err := net.Dial("tcp", ":5000")
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Write([]byte("SET Foo Bar 2500"))
	if err != nil {
		log.Fatal(err)
	}
}
