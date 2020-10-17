package chat

import (
	"log"

	"golang.org/x/net/context"
	"time"
	"os"
	"encoding/csv"
	"strconv"
	//"fmt"
	"sync"
	
)

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



func (s *Server) SayHello(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Receive message body from client: %s", in.Body)

	return &Message{Body: "Hello From the Server!"}, nil
}

func (s *Server) HacerPedido(ctx context.Context, in *Orden) (*Codigo, error) {

	
	//fmt.Println(s.SegOrd["asd"])

	currentTime := time.Now()
	asd := currentTime.Format("2006-01-02 15:04:05")

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

	//crear paquete

	paquete := Paquete{Id: in.Id, Estado: "En bodega",Idseg: i,Intentos: 0,Valor: int32(in.Valor),Tipo: tipostr}
	log.Printf("El destino del pedido %s es %s",in.Id, in.Destino)

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

	//falta arreglar el retorno
	return &Codigo{Idcompra: "1234"}, nil
}

func (s *Server)LlegoCamion(ctx context.Context, in *Camion) (*Carga, error) {

	//agregar origen/destino a carga

	log.Printf("Llego el camion: %s", in.Id)

	tipocam := in.Tipo

	c:= &Carga{}

	switch tipocam {
	case "normal":

		carga := 0
		wait := false
		flag := false

		for carga != 2{
			//parte con 1 paquete
			if flag == true{
				break
			}

			//termina de esperar
			if wait == true {

				//sleep
				time.Sleep(10 * time.Second) 
				flag = true
			}

			if len(s.qprio) > 0{
				if len(s.qprio) >=2 && carga == 0{


					s.mux.Lock()

					packet1 := s.qprio[0]
					s.qprio = s.qprio[1:]
					packet2 := s.qprio[0]
					s.qprio = s.qprio[1:]

					s.SegOrd[packet1.Idseg] = "En camino"
					s.SegOrd[packet2.Idseg] = "En camino"

					s.mux.Unlock()

					c.Paq1 = &packet1
					c.Paq2 = &packet2
					c.Flag = 2

					break

				} else {
					if carga == 1{
						s.mux.Lock()

						
						packet2 := s.qprio[0]
						s.qprio = s.qprio[1:]
						s.SegOrd[packet2.Idseg] = "En camino"

						s.mux.Unlock()

						c.Paq2 = &packet2
						c.Flag = 2
					} else {
						s.mux.Lock()
						packet1 := s.qprio[0]
						s.qprio = s.qprio[1:]
						s.SegOrd[packet1.Id] = "En camino"

						s.mux.Unlock()

						c.Paq1 = &packet1
						c.Flag = 1

					}
				
					carga =  carga + 1
					wait = true
					

				}

				
				
				

			} else if len(s.qnormal) > 0{
				if len(s.qnormal) >=2 && carga == 0{

					s.mux.Lock()

					packet1 := s.qnormal[0]
					s.qnormal = s.qnormal[1:]
					packet2 := s.qnormal[0]
					s.qnormal = s.qnormal[1:]

					s.SegOrd[packet1.Idseg] = "En camino"
					s.SegOrd[packet2.Idseg] = "En camino"

					s.mux.Unlock()

					c.Paq1 = &packet1
					c.Paq2 = &packet2
					c.Flag = 2

					break

				} else {
					if carga == 1{

						s.mux.Lock()

						
						packet2 := s.qnormal[0]
						s.qnormal = s.qnormal[1:]
						

						s.mux.Unlock()

						c.Paq2 = &packet2
						c.Flag = 2

					} else {

						s.mux.Lock()
						packet1 := s.qnormal[0]
						s.qnormal = s.qnormal[1:]
						

						s.mux.Unlock()

						c.Paq1 = &packet1
						c.Flag = 1

					}

					carga =  carga + 1
					wait = true
				}


			}

		}
	case "retail":

		carga := 0
		wait := false
		flag := false

		for carga != 2{

			//parte con 1 paquete
			if flag == true{
				break
			}

			//termina de esperar
			if wait == true {

				//sleep
				time.Sleep(8 * time.Second) 
				flag = true
			}


			// hay paquete en la cola retail
			if len(s.qret) > 0{
				//revisa camion vacio e inclute 2 paquetes
				if len(s.qret) >=2 && carga == 0{
				

					s.mux.Lock()

					packet1 := s.qret[0]
					s.qret = s.qret[1:]
					packet2 := s.qret[0]
					s.qret = s.qret[1:]

					s.mux.Unlock()

					c.Paq1 = &packet1
					c.Paq2 = &packet2
					c.Flag = 2

					break


				} else {

					if carga == 1{
						s.mux.Lock()

						
						packet2 := s.qret[0]
						s.qret = s.qret[1:]
						

						s.mux.Unlock()

						c.Paq2 = &packet2
						c.Flag = 2
				

					} else {

						s.mux.Lock()
						packet1 := s.qret[0]
						s.qret = s.qret[1:]
						

						s.mux.Unlock()

						c.Paq1 = &packet1
						c.Flag = 1

					}

					carga =  carga + 1
					wait = true

				}
				
			} else if len(s.qprio) > 0{
				if len(s.qprio) >=2 && carga == 0{


					s.mux.Lock()

					packet1 := s.qprio[0]
					s.qprio = s.qprio[1:]
					packet2 := s.qprio[0]
					s.qprio = s.qprio[1:]

					s.SegOrd[packet1.Idseg] = "En camino"
					s.SegOrd[packet2.Idseg] = "En camino"

					s.mux.Unlock()

					c.Paq1 = &packet1
					c.Paq2 = &packet2
					c.Flag = 2
					break

				} else {

					if carga ==1 {
						s.mux.Lock()

						
						packet2 := s.qprio[0]
						s.qprio = s.qprio[1:]
						s.SegOrd[packet2.Idseg] = "En camino"

						s.mux.Unlock()

						c.Paq2 = &packet2
						c.Flag = 2
						

					} else {

						s.mux.Lock()
						packet1 := s.qprio[0]
						s.qprio = s.qprio[1:]
						s.SegOrd[packet1.Id] = "En camino"

						s.mux.Unlock()

						c.Paq1 = &packet1
						c.Flag = 1

					}

					carga =  carga + 1
					wait = true

				}
			}

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
	return c, nil

}


func (s *Server)EntregaCamion(ctx context.Context, in *Entrega) (*Respuesta, error){


	log.Printf("Recibi %d paquetes", in.Num)
	log.Printf("Id: %s , Tipo: %s, Estado: %s, Valor %d, Intentos %d",in.Inf1.Id,in.Inf1.Tipo,in.Inf1.Estado,in.Inf1.Valor,in.Inf1.Intentos)
	return &Respuesta{Ack: "Datos recibidos"}, nil

}
