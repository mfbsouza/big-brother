package dbhandler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mfbsouza/big-brother/internal/db-types"
)

var db *sql.DB

func LoadDatabase(path string) {
	var err error
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		os.Create(path)
		log.Println("[db-handler] Populating database for the first time")
		defer populateDatabase()
	}

	db, err = sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalln("[db-handler] Error opening database", err)
	}
}

func CloseDatebase() {
	db.Close()
}

func FindUserById(id string) ([]byte, error) {
	var u_slice []*dbtypes.User
	q := `SELECT * FROM user WHERE id=?`
	rows, _ := db.Query(q, id)
	for rows.Next() {
		u := &dbtypes.User{}
		rows.Scan(&u.Id, &u.Name, &u.IsAdmin, &u.RegistrationDate, &u.RFIDTag)
		u_slice = append(u_slice, u)
	}
	defer rows.Close()
	if len(u_slice) != 1 {
		log.Println("[db-handler] Error while searching for user Id:", id)
		return []byte{}, errors.New("Length of the user slice is different than 1")
	} else {
		bytestream, _ := json.Marshal(u_slice[0])
		return bytestream, nil
	}
}

func FindUserByTag(tag string) ([]byte, error) {
	var u_slice []*dbtypes.User
	q := `SELECT * FROM user WHERE RFIDTag=?`
	rows, _ := db.Query(q, tag)
	for rows.Next() {
		u := &dbtypes.User{}
		rows.Scan(&u.Id, &u.Name, &u.IsAdmin, &u.RegistrationDate, &u.RFIDTag)
		u_slice = append(u_slice, u)
	}
	defer rows.Close()
	if len(u_slice) != 1 {
		log.Println("[db-handler] Error while searching for user Tag:", tag)
		return []byte{}, errors.New("Length of the user slice is different than 1")
	} else {
		bytestream, _ := json.Marshal(u_slice[0])
		return bytestream, nil
	}
}

func CreateUser(j []byte) int64 {
	var u dbtypes.User
	json.Unmarshal(j, &u)
	stmt, _ := db.Prepare(
		`INSERT INTO User (
			Name, 
			isAdmin, 
			RegistrationDate, 
			RFIDTag
		) VALUES (?, ?, ?, ?)`,
	)
	res, err := stmt.Exec(u.Name, u.IsAdmin, time.Now().UTC(), u.RFIDTag)
	id, err := res.LastInsertId()
	if err != nil {
		log.Println("[db-handler] Error while creating a new user:", err)
		return 0
	} else {
		return id
	}
}

func UpdateUser(id string, j []byte) bool {
	var u dbtypes.User
	json.Unmarshal(j, &u)
	stmt, _ := db.Prepare(`UPDATE User SET Name=?, isAdmin=? WHERE id=?`)
	_, err := stmt.Exec(u.Name, u.IsAdmin, id)
	if err != nil {
		log.Println("[db-handler] Error while updating a user:", err)
		return false
	} else {
		return true
	}
}

func FindEquipment(s string) ([]byte, error) {
	var e_slice []*dbtypes.Equipment
	sub_string := strings.ToLower(s)
	stmt, _ := db.Prepare(
		`SELECT * FROM equipment WHERE LOWER(Name) LIKE ?`,
	)
	rows, err := stmt.Query("%" + sub_string + "%")
	if err != nil {
		return []byte{}, err
	}
	defer rows.Close()
	for rows.Next() {
		e := &dbtypes.Equipment{}
		rows.Scan(&e.Id, &e.Name, &e.IsInUse, &e.IsBlocked, &e.UserId)
		e_slice = append(e_slice, e)
	}
	bytestream, err := json.Marshal(e_slice)
	if err != nil {
		log.Println("[db-handler] Error converting to JSON byte stream:", err)
		return []byte{}, err
	} else if len(e_slice) == 0 {
		return []byte{}, errors.New("no equipment found!")
	} else {
		return bytestream, nil
	}
}

func populateDatabase() {
	// create tables
	stmt, _ := db.Prepare(
		`CREATE TABLE IF NOT EXISTS user (
			ID INTEGER PRIMARY KEY, 
			Name TEXT, 
			isAdmin BOOL, 
			RegistrationDate TIMESTAMP, 
			RFIDTag TEXT)`,
	)
	stmt.Exec()

	stmt, _ = db.Prepare(
		`CREATE TABLE IF NOT EXISTS equipment (
			ID INTEGER PRIMARY KEY, 
			Name TEXT, 
			isInUse BOOL, 
			isBlocked BOOL, 
			user_ID INTEGER,
			FOREIGN KEY (user_ID) REFERENCES user(ID)
		)`,
	)
	stmt.Exec()

	// create the first user
	stmt, _ = db.Prepare(
		`INSERT INTO user (
			Name, 
			isAdmin, 
			RegistrationDate, 
			RFIDTag
		) VALUES (?, ?, ?, ?)`,
	)
	stmt.Exec("Matheus Souza", true, time.Now().UTC(), "5B320FE6")

	// create the first equipment
	stmt, _ = db.Prepare(
		`INSERT INTO equipment (
			Name, 
			isInUse, 
			isBlocked, 
			user_ID
		) VALUES (?, ?, ?, ?)`,
	)
	stmt.Exec("Arduino Uno", false, false, 0)
}
