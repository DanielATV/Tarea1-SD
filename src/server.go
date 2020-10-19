package main

import (
	"fmt"
	"log"
	"net"

	"github.com/tutorialedge/go-grpc-beginners-tutorial/chat"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Servidor iniciado!")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	

	s := chat.Server{Asd: 0, SegOrd: make(map[string]string), Origen: make(map[string]string), Destino: make(map[string]string)}
	//s.SegOrd= make(map[string]string)
	//s.SegOrd["asd"] = "holamundo"
	

	grpcServer := grpc.NewServer()

	chat.RegisterChatServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

	
}