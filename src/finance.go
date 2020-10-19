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

	var perdidasTot float32
	var gananciasTot float32
	var balance float32

	perdidasTot = float32(0)
	gananciasTot = float32(0)
	balance = float32(0)

	var gananciaPaq int
	var perdidaPaq int

	var dumb float32
	
	flag := 0
	

	//conn, err := amqp.Dial("amqp://admin:admin@dist01:5672/")
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

	var parche string

	go func() {
		for d := range msgs {
			
			json.Unmarshal(d.Body, &res)

			

			file, err := os.OpenFile("RegistroFinanciero.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			checkError("Cannot create file", err)
			

			writer := csv.NewWriter(file)
			

			
			//Determina las gancias y perdidas del paquete
			switch res.Tipo{
			case "retail":
				gananciaPaq = res.Valor


				
			case "prioritario":
				flag = 1
				if res.Estado == "Recibido"{

					dumb = float32(res.Valor)*1.3
					
				} else{
					dumb = float32(res.Valor)*0.3
				}



			case "normal":
				if res.Estado == "Recibido"{
					gananciaPaq = res.Valor

				} else{

					gananciaPaq = 0

				}
				
			}

			perdidaPaq = (res.Intentos-1)*10

			if flag == 1{

				dummy = strconv.FormatFloat(float64(dumb), 'f',2, 32)
				dummy2 = strconv.Itoa(perdidaPaq)

				gananciasTot = gananciasTot + dumb
				flag = 0
				

			} else{

				dummy = strconv.Itoa(gananciaPaq)
				dummy2 = strconv.Itoa(perdidaPaq)

				gananciasTot = gananciasTot + float32(gananciaPaq)
				

			}

			perdidasTot = perdidasTot + float32(perdidaPaq)

			balance = gananciasTot - perdidasTot

			

			//falta ganancia, perdida por item

			parche= strconv.Itoa(res.Intentos)

			slc1 := []string{res.Id,res.Tipo,res.Estado,parche,dummy,dummy2}
			err2 := writer.Write(slc1)
			checkError("Cannot write to file", err2)
			
			writer.Flush()
			file.Close()

			fmt.Println("Las ganancias totales son de: ", gananciasTot)
			fmt.Println("Las perdidas totales son de: ", perdidasTot)
			fmt.Println("El balance final es: ", balance)
			
			
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