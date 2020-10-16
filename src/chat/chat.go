package chat

import (
	"log"

	"golang.org/x/net/context"
	"time"
	"os"
	"encoding/csv"
	"strconv"
	"fmt"
	"sync"
	
)

type Server struct {
	Asd int
	qret []Paquete
	qprio []Paquete
	qnormal []Paquete
	SegOrd map[string]string
	mux sync.Mutex
	

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

	
	fmt.Println(s.SegOrd["asd"])

	currentTime := time.Now()
	asd := currentTime.Format("2006-01-02 15:04:05")

	file, err := os.OpenFile("registroPedidos.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
    checkError("Cannot create file", err)
	defer file.Close()
	
	writer := csv.NewWriter(file)
	defer writer.Flush()


	s.Asd = s.Asd + 1
	i:= strconv.Itoa(s.Asd)
	j:= strconv.Itoa(int(in.Valor))
	
	tipo := in.Tipo
	var tipostr string

	switch tipo {
    case 0:
        tipostr = "normal"
    case 1:
        tipostr = "prioritario"
    case 2:
        tipostr = "retail"
    }
	var mensaje = []string{asd,in.Id,tipostr,in.Producto,j,in.Origen,in.Destino,i}
	

	//err := writer.Write({asd,in.Id,in.Tipo,in.Producto,in.Valor,in.Origen,in.Destino,0})
	err2 := writer.Write(mensaje)
    checkError("Cannot write to file", err2)

	//crear paquete

	paquete := Paquete{Id: in.Id, Estado: "En bodega",Idseg: i,Intentos: 0}
	log.Printf("El destino del pedido %s es %s",in.Id, in.Destino)

	s.mux.Lock()
	switch tipostr {
	case "normal":
		s.qnormal = append(s.qnormal,paquete)
	
	case "prioritario":
		s.qprio = append(s.qprio,paquete)
	case "retail":
		s.qret = append(s.qret,paquete)
		
	}
	s.mux.Unlock()

	fmt.Println(len(s.qnormal))
	fmt.Println(len(s.qprio))
	//falta arreglar el retorno
	return &Codigo{Idcompra: "1234"}, nil
}
