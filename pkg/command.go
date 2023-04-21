package pkg

import "time"

type Command string

const (
	CMDSet Command = "SET"
	CMDGet Command = "GET"
)

type Message struct {
	CMD   Command
	Key   []byte
	Value []byte
	TTL   time.Duration
}

//func parseCommand(raw []byte) (*Message, error) {
//
//}
