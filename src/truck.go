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

	//preparar archivo

	nombre := "camion" + id + ".csv"
	file, err := os.OpenFile(nombre, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

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

		
		
		ent1 :=0
		ent2 :=0

		var primero *chat.Paquete
		var segundo *chat.Paquete

		var cost1 float32
		var cost2 float32


		flag := response.Flag

		if flag == 2 {
			if response.Paq1.Tipo == "retail"{
				cost1 = float32(response.Paq1.Valor)
			} else{
				cost1 = float32(response.Paq1.Valor) + float32(response.Paq1.Valor)*0.3
			}
			if response.Paq2.Tipo == "pyme"{
				cost2 = float32(response.Paq2.Valor) + float32(response.Paq2.Valor)*0.3
			}else{
				cost2 = float32(response.Paq2.Valor)
			}
		}

		var dest1 string
		var dest2 string
		var origen1 string
		var origen2 string

		if flag == 2 {

			if ( cost2 > cost1){
				primero = response.Paq2
				dest1 = response.Destino2
				origen1 = response.Origen2

				segundo = response.Paq1
				dest2 = response.Destino1
				origen2 = response.Origen1
			}else{
				primero  = response.Paq1
				dest1 = response.Destino1
				origen1 = response.Origen1

				segundo  = response.Paq2
				dest2 = response.Destino2
				origen2 = response.Origen2
			}
			//fmt.Println(primero,segundo)
		} else{
			ent2 = 1
			primero = response.Paq1
			dest1 = response.Destino1
			origen1 = response.Origen1
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
	
				}else if (tope1 == 2){ //fallo

					
					
					if (try1 > tope1 || primero.Valor < 10*(try1)){ // caso de no mas intentos
						
						ent1 = 1
						paquetesent = paquetesent + 1
						primero.Estado = "Rechazado"
						primero.Intentos = try1
						
					}
					try1 = try1+1
					
				} else{
					if (try1 > tope1){ // caso de no mas intentos
						
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
				}else if (tope2 ==2){ //fallo

					
					
					if (try2 > tope2 || segundo.Valor < 10*(try2)){ // caso de no mas intentos
						ent2 = 1
						paquetesent = paquetesent + 1
						segundo.Estado = "Rechazado"
						segundo.Intentos = try2
						
					}
					try2 = try2+1
					
				} else{
					if (try2 > tope2 ){ // caso de no mas intentos
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
			var slc1 = []string{primero.Id,primero.Tipo,aux,origen1,dest1,aux2, entime1}
			fmt.Println(slc1)
			//var slc1 = []string{primero.Id,primero.Idseg,primero.Tipo,aux,aux2,primero.Estado, entime1}

			err := writer.Write(slc1)
			checkError("Cannot write to file", err)
			aux = strconv.Itoa(int(segundo.Valor))
			aux2 = strconv.Itoa(int(segundo.Intentos))
			//var slc2 = []string{segundo.Id,segundo.Idseg,segundo.Tipo,aux,aux2,segundo.Estado, entime2 }
			var slc2 = []string{segundo.Id,segundo.Tipo,aux,origen2,dest2,aux2, entime2}
			err2 := writer.Write(slc2)
			checkError("Cannot write to file", err2)

		}else{
			fmt.Println(primero)
			aux = strconv.Itoa(int(primero.Valor))
			aux2 = strconv.Itoa(int(primero.Intentos))
			//var slc1 = []string{primero.Id,primero.Idseg,primero.Tipo,aux,aux2,primero.Estado, entime1}
			var slc1 = []string{primero.Id,primero.Tipo,aux,origen1,dest1,aux2, entime1}
			err := writer.Write(slc1)
			checkError("Cannot write to file", err)
		}

		
		//log.Printf("Id: %s, Origen: %s, Destino: %s", response.Paq1.Id,response.Origen1,response.Destino1)
	}


}

func checkError(message string, err error) {
    if err != nil {
        log.Fatal(message, err)
    }
}