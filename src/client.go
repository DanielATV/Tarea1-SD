package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/tutorialedge/go-grpc-beginners-tutorial/chat"

	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"
	"strconv"
)




func main() {

	var eleccion int
	var esperar int
	//Solicitar tipo de cliente
	fmt.Println("Indique el tipo de cliente")
	fmt.Println("0: Pyme")
	fmt.Println("1: Retail")
	
	fmt.Scanln(&eleccion)

	//Frecuencia de envios
	fmt.Println("Indique el tiempo entre envios(segundos)")	
	fmt.Scanln(&esperar) 

	

	if eleccion == 0 {
		csvReaderRow(esperar,"pymes.csv",0)

	} else {
		csvReaderRow(esperar,"retail.csv",1)

	}



}

func csvReaderRow(timer int,nombre string, flag int) {

	cont:= 0
	var codseg string 
	var idaux string

	

	var conn *grpc.ClientConn
	//conn, err := grpc.Dial("dist04:9000", grpc.WithInsecure())
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)



	// Open the file
	recordFile, err := os.Open(nombre)
	if err != nil {
		fmt.Println("An error encountered ::", err)
		return
	}

	// Setup the reader
	reader := csv.NewReader(recordFile)

	// Read the records
	header, err := reader.Read()
	if err != nil {
		fmt.Println("An error encountered ::", err)
		return
	}
	fmt.Printf("Headers : %v \n", header)

	//var casa string

	for i:= 0 ;; i = i + 1 {
		record, err := reader.Read()
		if err == io.EOF {
			break // reached end of the file
		} else if err != nil {
			fmt.Println("An error encountered ::", err)
			return
		}

		i, err := strconv.Atoi(record[2])

		var tipoProd int32
		if flag ==0 {
			i, _ := strconv.Atoi(record[5])
			tipoProd = int32(i)
		} else{
			tipoProd = 2
		}

		response, err := c.HacerPedido(context.Background(),
		&chat.Orden{Id: record[0],
			Producto: record[1],
			Valor: int32(i),
			Origen: record[3],
			Destino: record[4],
			Tipo: tipoProd})
		if err != nil {
			log.Fatalf("Error when calling SayHello: %s", err)
		}
		log.Printf("Se entrego la orden: %s", record[0])
		//log.Printf("Response from server: %s", response.Body)


		//Pregunta el estado del pedido cada 3 ordenes si es cliente tipo  pyme
		if (flag == 0 ){

			if cont == 0{

				codseg = response.Idcompra
				idaux = record[0]

			}

			cont = cont +1

			if cont == 3{
				cont = 0
				//Pide el estado de un pedido
				respuesta, _ := c. Estado(context.Background(),&chat.Codigo{Idcompra: codseg})

				log.Printf("La orden %s se encuentra %s", idaux, respuesta.Estado)
			}
			


		}

		//fmt.Printf("Row %d : %v \n", i, record[4])

		//Espera para hacer la siguiente orden
		time.Sleep(time.Duration(timer) * time.Second)
		//casa = string(record[4])
		//fmt.Printf("%s \n",casa)
		
	}
	

	// Note: Each time Read() is called, it reads the next line from the file
	// r1, _ := reader.Read() // Reads the first row, useful for headers
	// r2, _ := reader.Read() // Reads the second row
}