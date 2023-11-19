package dbhandler

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func LoadDatabase(path string) {
	var err error
	// check if database path exists
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		os.Create(path)
	}

	db, err = sql.Open("sqlite3", path)
	if err != nil {
		log.Println("[db-handler] error opening database", err)
	}

	log.Println("populating database", path)
	populateDatabase()
}

func CloseDatebase() {
	db.Close()
}

func populateDatabase() {
	// create tables
	statement, _ := db.Prepare(
		`CREATE TABLE IF NOT EXISTS user (
			ID INTEGER PRIMARY KEY, 
			Name TEXT, 
			isAdmin BOOL, 
			RegistrationDate TIMESTAMP, 
			RFIDTag TEXT)`,
	)
	statement.Exec()

	// create the first user
	statement, _ = db.Prepare(
		`INSERT INTO User (
			Name, 
			isAdmin, 
			RegistrationDate, 
			RFIDTag
		) VALUES (?, ?, ?, ?)`,
	)
	statement.Exec("Matheus Souza", true, time.Now().UTC(), "putrealtaghere")
}
