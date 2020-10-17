package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/tutorialedge/go-grpc-beginners-tutorial/chat"

	"sync"
	"fmt"
	"math/rand"
    "time"
    "os"
    "encoding/csv"
    "strconv"
  
	
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
	for  i := 0; i < 5; i++{
		var conn *grpc.ClientConn
		conn, err := grpc.Dial(":9000", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		defer conn.Close()

		c := chat.NewChatServiceClient(conn)


		//solicita la carga
		response, err := c.LlegoCamion(context.Background(), &chat.Camion{Id: id, Tipo: tipo})
		if err != nil {
			log.Fatalf("Error when calling SayHello: %s", err)
		}

		//preparar archivo

		nombre := "camion" + id + ".csv"
		file, err := os.OpenFile(nombre, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		checkError("Cannot create file", err)
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()
		
		ent1 :=0
		ent2 :=0

		var primero *chat.Paquete
		var segundo *chat.Paquete

		flag := response.Flag

		if flag == 2 {

			if ( response.Paq2.Valor > response.Paq1.Valor){
				primero = response.Paq2
				segundo = response.Paq1
			}else{
				primero  = response.Paq1
				segundo  = response.Paq2
			}
			//fmt.Println(primero,segundo)
		} else{
			ent2 = 1
			primero = response.Paq1
		}

		//Asignar tope de reintentos
		var tope1 int32
		var tope2 int32
		if primero.Tipo == "retail"{
			tope1 =3
		}else{
			tope1 =2
		}

		if flag == 2{
			if segundo.Tipo == "retail"{
				tope2 =3
			}else{
				tope2 =2
			}
		}

		//Entregar y ver los re-intentos
		paquetesent := 0
		try1 := int32(1)
		try2 :=int32(1)
		var porcentaje int
		rand.Seed(time.Now().UnixNano()) // seed

		entime1:= "0"
		entime2:= "0"

		for paquetesent < int(flag){
			//Paquete 1
			if (ent1 == 0){
				porcentaje = rand.Intn(100)
				time.Sleep(1 * time.Second) 
				//fmt.Println("intento en paquete 1 con probabilidad ",porcentaje)
				if (porcentaje <=80){ //paquete entregado
					ent1 = 1
					paquetesent = paquetesent+1
					primero.Estado = "Entregado"
					primero.Intentos = try1
					currentTime := time.Now()
					entime1 = currentTime.Format("2006-01-02 15:04:05")
	
				}else{ //fallo

					
					
					if (try1 > tope1 || primero.Valor < 10*(try1)){ // caso de no mas intentos
						
						ent1 = 1
						paquetesent = paquetesent + 1
						primero.Estado = "Rechazado"
						primero.Intentos = try1
						
					}
					try1 = try1+1
					
				}
	
			}
			//Paquete 2
			if (ent2 == 0){
				porcentaje = rand.Intn(100)
				time.Sleep(1 * time.Second) 
				//fmt.Println("intento en paquete 2 con probabilidad ",porcentaje)
				if (porcentaje <=80){ //paquete entregado
					ent2 = 1
					paquetesent = paquetesent+1
					segundo.Estado = "Entregado"
					segundo.Intentos = try2
					currentTime := time.Now()
					
					entime2 = currentTime.Format("2006-01-02 15:04:05")
				}else{ //fallo

					
					
					if (try2 > tope2 || segundo.Valor < 10*(try2)){ // caso de no mas intentos
						ent2 = 1
						paquetesent = paquetesent + 1
						segundo.Estado = "Rechazado"
						segundo.Intentos = try2
						
					}
					try2 = try2+1
					
				}
	
			}
	
		}

		//escrita del archivo
		var aux string
		var aux2 string
		// id id-seg tipo valor intentos estado
		if (flag == 2){
			//fmt.Println(primero,segundo)
			aux = strconv.Itoa(int(primero.Valor))
			aux2 = strconv.Itoa(int(primero.Intentos))
			var slc1 = []string{primero.Id,primero.Idseg,primero.Tipo,aux,aux2,primero.Estado, entime1}
			err := writer.Write(slc1)
			checkError("Cannot write to file", err)
			aux = strconv.Itoa(int(segundo.Valor))
			aux2 = strconv.Itoa(int(segundo.Intentos))
			var slc2 = []string{segundo.Id,segundo.Idseg,segundo.Tipo,aux,aux2,segundo.Estado, entime2 }
			err2 := writer.Write(slc2)
			checkError("Cannot write to file", err2)

		}else{
			fmt.Println(primero)
			aux = strconv.Itoa(int(primero.Valor))
			aux2 = strconv.Itoa(int(primero.Intentos))
			var slc1 = []string{primero.Id,primero.Idseg,primero.Tipo,aux,aux2,primero.Estado, entime1}
			err := writer.Write(slc1)
			checkError("Cannot write to file", err)
		}


		log.Printf("Response from server: %s", response.Paq1.Id)
	}


}

func checkError(message string, err error) {
    if err != nil {
        log.Fatal(message, err)
    }
}