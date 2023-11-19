package main

import (
	"fmt"
	"os"

	dbhandler "github.com/mfbsouza/big-brother/internal/db-handler"
)

func main() {
	var db_path string

	if len(os.Args) != 2 {
		fmt.Println("Error: expected more command line arguments")
		fmt.Printf("Syntax: %s </path/to/database.db>\n", os.Args[0])
		os.Exit(1)
	} else {
		db_path = os.Args[1]
	}

	fmt.Println("db_path is", db_path)
	defer dbhandler.CloseDatebase()
	dbhandler.LoadDatabase(db_path)
}
