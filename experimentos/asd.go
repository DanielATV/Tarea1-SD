package main

import (
    "os"
    "log"
    "encoding/csv"
)

var data = [][]string{{"Line1", "Hello Readers of"}, {"Line2", "golangcode.com"}}

func main() {
    file, err := os.OpenFile("result.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
    checkError("Cannot create file", err)
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    for _, value := range data {
        err := writer.Write(value)
        checkError("Cannot write to file", err)
    }
}

func checkError(message string, err error) {
    if err != nil {
        log.Fatal(message, err)
    }
}