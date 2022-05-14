package main

import (
	"github.com/mrhakimov/broadcasts-benchmarking/pkg/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	file, err := os.OpenFile("/Users/mukkhakimov/Documents/itmo/thesis/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetFlags(log.Lmicroseconds | log.Ldate)
	log.SetOutput(file)

	http.HandleFunc("/broadcast", handlers.Broadcast)
	http.HandleFunc("/deliver/send", handlers.Send)
	http.HandleFunc("/deliver/echo", handlers.Echo)
	http.HandleFunc("/deliver/ready", handlers.Ready)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
