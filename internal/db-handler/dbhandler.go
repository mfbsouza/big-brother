package dbhandler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mfbsouza/big-brother/internal/db-types"
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

func FindUser(id string) ([]byte, error) {
	var u dbtypes.User
	var cnt int = 0
	log.Println("[db-handler] looking database for id:", id)
	q := `SELECT * FROM user WHERE id=?`
	rows, _ := db.Query(q, id)
	for rows.Next() {
		rows.Scan(&u.Id, &u.Name, &u.IsAdmin, &u.RegistrationDate, &u.RFIDTag)
		cnt += 1
	}
	if cnt == 1 {
		t, _ := json.Marshal(u)
		return t, nil
	} else if cnt == 0 {
		return nil, errors.New("no user found")
	} else {
		return nil, errors.New("more than one row")
	}
}

func CreateUser(j []byte) int64 {
	var u dbtypes.User
	json.Unmarshal(j, &u)
	log.Println("[db-handler] creating user for:", u.Name)
	statement, _ := db.Prepare(
		`INSERT INTO User (
			Name, 
			isAdmin, 
			RegistrationDate, 
			RFIDTag
		) VALUES (?, ?, ?, ?)`,
	)
	res, err := statement.Exec(u.Name, u.IsAdmin, time.Now().UTC(), u.RFIDTag)
	id, err := res.LastInsertId()
	if err != nil {
		log.Println("[db-handler] error creating user:", err)
		return 0
	} else {
		log.Println("[db-handler] user created in the database!")
		return id
	}
}

func UpdateUser(id string, j []byte) bool {
	var u dbtypes.User
	json.Unmarshal(j, &u)
	log.Println("[db-handler] updating user for:", u.Name)
	statement, _ := db.Prepare(`UPDATE User SET Name=?, isAdmin=? WHERE id=?`)
	_, err := statement.Exec(u.Name, u.IsAdmin, id)
	if err != nil {
		log.Println("[db-handler] error updating user:", err)
		return false
	} else {
		log.Println("[db-handler] user updated in the database!")
		return true
	}
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
