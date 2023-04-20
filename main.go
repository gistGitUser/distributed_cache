package main

import (
	"cacher/internal"
	"fmt"
	"log"
	"net"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
}

type Server struct {
	ServerOpts
	cache internal.Cacher
}

func NewServer(opts ServerOpts, c internal.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}

	log.Printf("server starting on port [%s]\n", s.ListenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %s\n", err)
			continue
		}
		go s.handleConnection(conn)
	}

	return nil
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("read error: %s\n", err)
			break
		}
		msg := buf[:n]
		fmt.Println(string(msg))
	}

}

func main() {
	opts := ServerOpts{
		ListenAddr: ":5000",
		IsLeader:   true,
	}

	server := NewServer(opts, internal.NewCache())
	log.Fatal(server.Start())
}
