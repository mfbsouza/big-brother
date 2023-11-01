package main

import (
	"log"
	"net/http"
)

const portNumber = ":8080"

func main() {
	log.Println("[main] starting web server...")
	http.ListenAndServe(portNumber, newRouter())
	log.Println("[main] server is listening on port", portNumber[1:])
}
