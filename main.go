package main

import (
	"cacher/internal"
	"cacher/pkg"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
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

func (s *Server) handleCommand(conn net.Conn, rawCMD []byte) {

	msg, err := parseMessage(rawCMD)
	if err != nil {
		fmt.Println("failed to parse command:", err)
		conn.Write([]byte(err.Error()))
		return
	}

	switch msg.CMD {
	case pkg.CMDSet:
		err = s.handleSetCmd(conn, msg)
	case pkg.CMDGet:
		err = s.handleGetCmd(conn, msg)
	}

	if err != nil {
		fmt.Println("failed to handle command:", err)
		conn.Write([]byte(err.Error()))
	}

}

func (s *Server) handleSetCmd(conn net.Conn, msg *pkg.Message) error {

	if err := s.cache.Set(msg.Key, msg.Value, msg.TTL); err != nil {
		return err
	}

	go s.sendToFollowers(context.TODO(), msg)

	return nil
}

func (s *Server) handleGetCmd(conn net.Conn, msg *pkg.Message) error {
	return nil
}

/*
Эта функция нужна для того чтобы
distributed система работало
*/
func (s *Server) sendToFollowers(ctx context.Context, msg *pkg.Message) error {
	return nil
}

func parseMessage(raw []byte) (*pkg.Message, error) {
	rawStr := string(raw)
	parts := strings.Fields(rawStr)

	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid command:%q", rawStr)
	}

	msg := &pkg.Message{
		CMD: pkg.Command(parts[0]),
		Key: []byte(parts[1]),
	}

	if msg.CMD == pkg.CMDSet {
		if len(parts) != 4 {
			return nil, fmt.Errorf("invalid SET command:%q", rawStr)
		}

		if len(parts) != 4 {
			return nil, fmt.Errorf("invalid SET command:%q", rawStr)
		}

		msg.Value = []byte(parts[2])

		ttl, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, errors.New("SET TTL error")
		}
		msg.TTL = time.Duration(ttl)
	}

	return msg, nil

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
		s.handleCommand(conn, msg)
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
