package main 
  
import (
	"fmt"
	"sync"
		)
  
func display(wg *sync.WaitGroup,str string) { 
	defer wg.Done()
    for w := 0; w < 6; w++ { 
        fmt.Println(str) 
    } 
} 
  
func main() { 
	var wg sync.WaitGroup
	wg.Add(3)
  
    // Calling Goroutine 
	go display(&wg,"hola")
	go display(&wg,"chao") 
	go display(&wg,"mundo")

	wg.Wait()

	
  
 
} 