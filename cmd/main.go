package main

import (
	"github.com/mrhakimov/broadcasts-source/pkg/brb"
	"github.com/mrhakimov/broadcasts-source/pkg/cebrb"
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

	http.HandleFunc("/brb/broadcast", brb.Broadcast)
	http.HandleFunc("/brb/clear", brb.Clear)
	http.HandleFunc("/brb/deliver/send", brb.Send)
	http.HandleFunc("/brb/deliver/echo", brb.Echo)
	http.HandleFunc("/brb/deliver/ready", brb.Ready)

	http.HandleFunc("/cebrb/broadcast", cebrb.Broadcast)
	http.HandleFunc("/cebrb/clear", cebrb.Clear)
	http.HandleFunc("/cebrb/deliver/init", cebrb.Init)
	http.HandleFunc("/cebrb/deliver/witness", cebrb.Witness)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
