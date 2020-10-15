package chat

import (
	"log"

	"golang.org/x/net/context"
	"github.com/tutorialedge/go-grpc-beginners-tutorial/camion"
	"google.golang.org/grpc"
)

type Server struct {
}

func (s *Server) SayHello(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Receive message body from client: %s", in.Body)

	return &Message{Body: "Hello From the Server!"}, nil
}

func (s *Server) HacerPedido(ctx context.Context, in *Orden) (*Codigo, error) {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := camion.NewChatServiceClient(conn)
	response, err := c.SayHelloCamion(context.Background(), &camion.Message{Body: "Hello From logistica!"})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from camion: %s", response.Body)
	log.Printf("El destino del pedido %s es %s",in.Id, in.Destino)


	return &Codigo{Idcompra: "1234"}, nil
}
