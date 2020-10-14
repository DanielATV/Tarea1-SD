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

	//Lectura de archivo

	if eleccion == 0 {
		csvReaderRow(esperar,"pyme.csv")

	} else {
		csvReaderRow(esperar,"retail.cvs")

	}



}

func csvReaderRow(timer int,nombre string) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)



	// Open the file
	recordFile, err := os.Open("pymes.csv")
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

		response, err := c.HacerPedido(context.Background(),
		&chat.Orden{Id: record[0],
			Producto: record[1],
			Valor: int32(i),
			Origen: record[3],
			Destino: record[4],
			Tipo: 2})
		if err != nil {
			log.Fatalf("Error when calling SayHello: %s", err)
		}
		log.Printf("Response from server: %s", response.Idcompra)
		//log.Printf("Response from server: %s", response.Body)

		//fmt.Printf("Row %d : %v \n", i, record[4])
		time.Sleep(time.Duration(timer) * time.Second)
		//casa = string(record[4])
		//fmt.Printf("%s \n",casa)
		
	}
	

	// Note: Each time Read() is called, it reads the next line from the file
	// r1, _ := reader.Read() // Reads the first row, useful for headers
	// r2, _ := reader.Read() // Reads the second row
}