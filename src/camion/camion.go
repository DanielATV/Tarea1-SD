package camion

import (
	"log"

	"golang.org/x/net/context"
)

type Server struct {
}

func (s *Server) SayHelloCamion(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Receive message body from logistica: %s", in.Body)
	return &Message{Body: "Hello From the camion!"}, nil
}
