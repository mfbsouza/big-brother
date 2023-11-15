package main

import (
	"log"
	"net/http"

	"github.com/mfbsouza/big-brother/internal/web-router"
)

const portNumber = ":8080"

func main() {
	log.Println("[main] starting web server...")
	http.ListenAndServe(portNumber, router.NewRouter())
	log.Println("[main] server is listening on port", portNumber[1:])
}
