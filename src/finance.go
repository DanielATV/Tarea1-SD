package main

import (
	"fmt"
	"log"
	"encoding/json"
	"github.com/streadway/amqp"
	"os"
	"encoding/csv"
	"strconv"
)

type Datos struct {

	Id string
	Tipo string
	Estado string
	Intentos int
	Valor int
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}


func main() {

	
	
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello-queue", // name
		false,   // durable
		false,   // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	res := Datos{}
    
	fmt.Println(res)
	
	//var slc1 []string
	var dummy string
	var dummy2 string

	go func() {
		for d := range msgs {
			json.Unmarshal(d.Body, &res)

			file, err := os.OpenFile("RegistroFinanciero.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			checkError("Cannot create file", err)
			

			writer := csv.NewWriter(file)
			

			dummy = strconv.Itoa(int(res.Intentos))
			dummy2 = strconv.Itoa(int(res.Valor))

			//falta ganancia, perdida por item

			slc1 := []string{res.Id,res.Tipo,res.Estado,dummy,dummy2}
			err2 := writer.Write(slc1)
			checkError("Cannot write to file", err2)
			
			writer.Flush()
			file.Close()

			fmt.Println(res.Id)
			
			fmt.Println(res.Tipo)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func checkError(message string, err error) {
    if err != nil {
        log.Fatal(message, err)
    }
}