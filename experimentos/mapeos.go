package main

import "fmt"

type person struct {
    name string
	age  int
	mapeo map[string]string
	lista []int
}

func newPerson() *person {

	var p person
	p.name = "asd"
	p.age = 42
	p.mapeo = make(map[string]string)
	return &p
}

func main() {


	s := newPerson()
	s.mapeo["qwer"] = "ty"
	fmt.Println(s.mapeo["qwer"])
	s.mapeo["qwer"] = "uiop"
	fmt.Println(s.mapeo["qwer"])
	
	s.lista = append(s.lista,2)
	//s.lista = append(s.lista,6)


	//pop
	_, s.lista = s.lista[0], s.lista[1:]

	fmt.Println(len(s.lista))



    
}