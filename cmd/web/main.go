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
	log.Println("[main] starting web server...")
	http.ListenAndServe(portNumber, nil)
	log.Println("[main] server is listening on port", portNumber[1:])
}
