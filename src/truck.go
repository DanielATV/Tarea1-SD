package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/tutorialedge/go-grpc-beginners-tutorial/chat"

	"sync"
	
)




func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	go camion(&wg,"1","retail")
	go camion(&wg,"2","retail")
	go camion(&wg,"3","normal")
	wg.Wait()

}	

	

func camion(wg *sync.WaitGroup, id string,tipo string ){
	defer wg.Done()
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)

	response, err := c.LlegoCamion(context.Background(), &chat.Camion{Id: id, Tipo: tipo})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s", response.Paq1.Id)


}