package main
  
import (
    "fmt"
    "time"
)
func main() {
  
    currentTime := time.Now()
  
      
    fmt.Println("YYYY-MM-DD hh:mm:ss : ", currentTime.Format("2006-01-02 15:04:05"))
          
    
}