package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mfbsouza/big-brother/internal/db-handler"
	"github.com/mfbsouza/big-brother/internal/db-router"
)

const portNumber = ":3030"

func main() {
	var db_path string
	if len(os.Args) != 2 {
		fmt.Println("Error: expected more command line arguments")
		fmt.Printf("Syntax: %s </path/to/database.db>\n", os.Args[0])
		os.Exit(1)
	} else {
		db_path = os.Args[1]
	}
	defer dbhandler.CloseDatebase()
	dbhandler.LoadDatabase(db_path)

	log.Println("[main] Server is listening on port", portNumber[1:])
	http.ListenAndServe(portNumber, router.NewRouter())
}
