package main

import (
	"log"
	"net/http"

	"github.com/mfbsouza/big-brother-web/pkg/handlers"
)

const portNumber = ":8080"

func main() {
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)
	log.Println("[main] starting web server at port:", portNumber[1:])
	http.ListenAndServe(portNumber, nil)
}
