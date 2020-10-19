package chat

import (
	"log"

	"golang.org/x/net/context"
	"time"
	"os"
	"encoding/csv"
	"strconv"
	"sync"

	"fmt"

	"encoding/json"
	"github.com/streadway/amqp"
	
)

//Estructura Json
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


// Servidos con las variables que maneja
type Server struct {
	Asd int
	qret []Paquete
	qprio []Paquete
	qnormal []Paquete
	SegOrd map[string]string
	mux sync.Mutex
	Origen map [string]string
	Destino map[string]string
	

}

func checkError(message string, err error) {
    if err != nil {
        log.Fatal(message, err)
    }
}


//Funcion de referencia
func (s *Server) SayHello(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Receive message body from client: %s", in.Body)

	return &Message{Body: "Hello From the Server!"}, nil
}


//Servicio que recibe una orden y entrega un codigo de seguimiento
func (s *Server) HacerPedido(ctx context.Context, in *Orden) (*Codigo, error) {

	
	//fmt.Println(s.SegOrd["asd"])

	currentTime := time.Now()
	asd := currentTime.Format("2006-01-02 15:04:05")

	//Manejo de archivo

	file, err := os.OpenFile("registroPedidos.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
    checkError("Cannot create file", err)
	defer file.Close()
	
	writer := csv.NewWriter(file)
	defer writer.Flush()


	
	
	j:= strconv.Itoa(int(in.Valor))
	
	tipo := in.Tipo
	var tipostr string
	var i string
	


	switch tipo {
    case 0:
		tipostr = "normal"
		s.Asd = s.Asd + 1
		i= strconv.Itoa(s.Asd)
    case 1:
		tipostr = "prioritario"
		s.Asd = s.Asd + 1
		i= strconv.Itoa(s.Asd)
    case 2:
		tipostr = "retail"
		i = "0"
		
		
    }
	var mensaje = []string{asd,in.Id,tipostr,in.Producto,j,in.Origen,in.Destino,i}
	
	

	//err := writer.Write({asd,in.Id,in.Tipo,in.Producto,in.Valor,in.Origen,in.Destino,0})
	err2 := writer.Write(mensaje)
    checkError("Cannot write to file", err2)

	//Creacion del paquete

	paquete := Paquete{Id: in.Id, Estado: "En bodega",Idseg: i,Intentos: 0,Valor: int32(in.Valor),Tipo: tipostr}
	log.Printf("Se recibio el pedido %s",in.Id)


	//Añade los paquetes a la cola que corresponde
	s.mux.Lock()
	switch tipostr {
	case "normal":
		s.qnormal = append(s.qnormal,paquete)
		s.SegOrd[i] = "En bodega"
	
	case "prioritario":
		s.qprio = append(s.qprio,paquete)
		s.SegOrd[i] = "En bodega"

	case "retail":
		s.qret = append(s.qret,paquete)
		
	}

	s.Origen[in.Id] = in.Origen
	s.Destino[in.Id] = in.Destino
	
	s.mux.Unlock()

	//fmt.Println(len(s.qnormal))
	//fmt.Println(len(s.qprio))
	//fmt.Println(len(s.qret))

	
	return &Codigo{Idcompra: i}, nil
}

// Servicio que le asigna la carga correspondiente al camion cuando "llega"
func (s *Server)LlegoCamion(ctx context.Context, in *Camion) (*Carga, error) {

	
	log.Printf("Llego el camion: %s", in.Id)

	tipocam := in.Tipo

	c:= &Carga{}

	switch tipocam {
	case "normal":

		carga := 0
		wait := false
		flag := false

		for carga != 2{
			//Parte con 1 paquete
			if flag == true{
				break
			}

			//Espera un segundo paquete
			if wait == true {

				time.Sleep(time.Duration(in.Espera) * time.Second)
				flag = true
			}

			//Hay paquetes en la cola prioritaria
			s.mux.Lock()
			if len(s.qprio) > 0{
				//Asigma 2 paquetes y parte el camion
				if len(s.qprio) >=2 && carga == 0{

					log.Printf("Checkpoint 1")
					

					packet1 := s.qprio[0]
					s.qprio = s.qprio[1:]
					packet2 := s.qprio[0]
					s.qprio = s.qprio[1:]

					s.SegOrd[packet1.Idseg] = "En camino"
					s.SegOrd[packet2.Idseg] = "En camino"

					

					c.Paq1 = &packet1
					c.Paq2 = &packet2
					c.Flag = 2
					s.mux.Unlock()
					break
				//Asigna solo un paquete y espera
				} else {
					if carga == 1{
						log.Printf("Checkpoint 2")
						

						
						packet2 := s.qprio[0]
						s.qprio = s.qprio[1:]
						s.SegOrd[packet2.Idseg] = "En camino"

						

						c.Paq2 = &packet2
						c.Flag = 2
					} else {

						
						log.Printf("Checkpoint 3")
						
						packet1 := s.qprio[0]
						log.Printf("Checkpoint 3.1")
						s.qprio = s.qprio[1:]
						s.SegOrd[packet1.Id] = "En camino"

						

						c.Paq1 = &packet1
						c.Flag = 1

					}
				
					carga =  carga + 1
					wait = true
					

				}

				
				
				
			//Hay paquetes en la cola normal
			} else if len(s.qnormal) > 0{
				if len(s.qnormal) >=2 && carga == 0{

					log.Printf("Checkpoint 4")

					

					packet1 := s.qnormal[0]
					s.qnormal = s.qnormal[1:]
					packet2 := s.qnormal[0]
					s.qnormal = s.qnormal[1:]

					s.SegOrd[packet1.Idseg] = "En camino"
					s.SegOrd[packet2.Idseg] = "En camino"

					

					c.Paq1 = &packet1
					c.Paq2 = &packet2
					c.Flag = 2
					s.mux.Unlock()
					break
				//Asigna solo un paquete y espera
				} else {
					if carga == 1{


						log.Printf("Checkpoint 5")
						

						
						packet2 := s.qnormal[0]
						s.qnormal = s.qnormal[1:]
						

						

						c.Paq2 = &packet2
						c.Flag = 2

					} else {

						log.Printf("Checkpoint 6")
						
						packet1 := s.qnormal[0]
						s.qnormal = s.qnormal[1:]
						

						

						c.Paq1 = &packet1
						c.Flag = 1

					}

					carga =  carga + 1
					wait = true
				}


			}
			s.mux.Unlock()

		}
	case "retail":

		carga := 0
		wait := false
		flag := false

		for carga != 2{

			//Parte con 1 paquete
			if flag == true{
				break
			}

			//Espera un segundo paquete
			if wait == true {

				
				time.Sleep(time.Duration(in.Espera) * time.Second)
				flag = true
			}


			//Hay paquetes en la cola retail
			s.mux.Lock()

			if len(s.qret) > 0{
				
				if len(s.qret) >=2 && carga == 0{
				

					
					packet1 := s.qret[0]
					s.qret = s.qret[1:]
					packet2 := s.qret[0]
					s.qret = s.qret[1:]

					
					c.Paq1 = &packet1
					c.Paq2 = &packet2
					c.Flag = 2

					break
					s.mux.Unlock()


				} else {

					if carga == 1{
					

						
						packet2 := s.qret[0]
						s.qret = s.qret[1:]
						

					

						c.Paq2 = &packet2
						c.Flag = 2
				

					} else {

						
						packet1 := s.qret[0]
						s.qret = s.qret[1:]
						

						

						c.Paq1 = &packet1
						c.Flag = 1

					}

					carga =  carga + 1
					wait = true

				}
			//Hay paquetes en la cola prioritaria
			} else if len(s.qprio) > 0{
				//Asigma 2 paquetes y parte el camion
				if len(s.qprio) >=2 && carga == 0{


					

					packet1 := s.qprio[0]
					s.qprio = s.qprio[1:]
					packet2 := s.qprio[0]
					s.qprio = s.qprio[1:]

					s.SegOrd[packet1.Idseg] = "En camino"
					s.SegOrd[packet2.Idseg] = "En camino"

					

					c.Paq1 = &packet1
					c.Paq2 = &packet2
					c.Flag = 2
					s.mux.Unlock()

					break

				} else {

					if carga ==1 {
						

						
						packet2 := s.qprio[0]
						s.qprio = s.qprio[1:]
						s.SegOrd[packet2.Idseg] = "En camino"

						

						c.Paq2 = &packet2
						c.Flag = 2
						

					} else {

						
						packet1 := s.qprio[0]
						s.qprio = s.qprio[1:]
						s.SegOrd[packet1.Id] = "En camino"

						

						c.Paq1 = &packet1
						c.Flag = 1

					}

					carga =  carga + 1
					wait = true

				}
			}
			s.mux.Unlock()
		}
		
		

	}

	if c.Flag == 2{
		c.Origen1 = s.Origen[c.Paq1.Id]
		c.Destino1 = s.Destino[c.Paq1.Id]
	
		c.Origen2 = s.Origen[c.Paq2.Id]
		c.Destino2 = s.Destino[c.Paq2.Id]
	} else {
		c.Origen1 = s.Origen[c.Paq1.Id]
		c.Destino1 = s.Destino[c.Paq1.Id]

	}


	//d:= &Carga{Paq1: &Paquete{ Id: "asd"},
	//Paq2: &Paquete{Id : "qwerty"},
	//Flag: 0}

	log.Printf("El camion %s partio", in.Id)
	return c, nil

}


//Servicio que actualiza el estado de las entregas del camion y lo notifica a financiero
func (s *Server)EntregaCamion(ctx context.Context, in *Entrega) (*Respuesta, error){


	//log.Printf("Recibi %d paquetes", in.Num)
	//log.Printf("Id: %s , Tipo: %s, Estado: %s, Valor %d, Intentos %d",in.Inf1.Id,in.Inf1.Tipo,in.Inf1.Estado,in.Inf1.Valor,in.Inf1.Intentos)

	//user := &User{Name: "Frank", Numero: 2}
	//b, _ := json.Marshal(user)
	
    

	//Conecion a fiananciero por rabbitMQ
	//conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	conn, err := amqp.Dial("amqp://admin:admin@dist01:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()




	q, err := ch.QueueDeclare(
		"hello-queue", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)

	failOnError(err, "Failed to declare a queue")

	//Publica dos mensajes a la cola si entrego 2 paquetes
	if in.Num == 2{
		
		load1 := &Datos{Id: in.Inf1.Id,
		Tipo: in.Inf1.Tipo,
		Estado:in.Inf1.Estado,
		Valor:int(in.Inf1.Valor),
		Intentos:int(in.Inf1.Intentos)}

		

		ld1, _ := json.Marshal(load1)

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        ld1,
			})

		failOnError(err, "Failed to publish a message")

		load2 := &Datos{Id: in.Inf2.Id,
		Tipo: in.Inf2.Tipo,
		Estado:in.Inf2.Estado,
		Valor:int(in.Inf2.Valor),
		Intentos:int(in.Inf2.Intentos)}

		ld2, _ := json.Marshal(load2)
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        ld2,
			})

		failOnError(err, "Failed to publish a message")

		//Actualiza el estado
		s.mux.Lock()

		s.SegOrd[in.Inf1.Idseg] = in.Inf1.Estado
		s.SegOrd[in.Inf2.Idseg] = in.Inf2.Estado

		s.mux.Unlock()

	//Publica 1 mensaje a la cola si entrego 1 paquete
	} else{

		load1 := &Datos{Id: in.Inf1.Id,
		Tipo: in.Inf1.Tipo,
		Estado:in.Inf1.Estado,
		Valor:int(in.Inf1.Valor),
		Intentos:int(in.Inf1.Intentos)}

		ld1, _ := json.Marshal(load1)

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        ld1,
			})

		failOnError(err, "Failed to publish a message")

		//Actualiza el estado
		s.mux.Lock()
		s.SegOrd[in.Inf1.Idseg] = in.Inf1.Estado

		s.mux.Unlock()


	}

	
	//log.Printf(" [x] Sent %s", string(b))
	//failOnError(err, "Failed to publish a message")
	
	return &Respuesta{Ack: "Datos recibidos"}, nil

}

func (s *Server) Estado(ctx context.Context, in *Codigo) (*EstOrden, error) {
	//log.Printf("Receive message body from client: %s", in.Idcompra)

	return &EstOrden{Estado: s.SegOrd[in.Idcompra]}, nil
}