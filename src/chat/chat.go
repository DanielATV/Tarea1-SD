package chat

import (
	"log"

	"golang.org/x/net/context"
)

type Server struct {
}

func (s *Server) SayHello(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Receive message body from client: %s", in.Body)
	return &Message{Body: "Hello From the Server!"}, nil
}

func (s *Server) HacerPedido(ctx context.Context, in *Orden) (*Codigo, error) {
	log.Printf("El destino del pedido %s es %s",in.Id, in.Destino)


	return &Codigo{Idcompra: "1234"}, nil
}