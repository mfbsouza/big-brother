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
	defer rows.Close()
	if cnt == 1 {
		t, _ := json.Marshal(u)
		return t, nil
	} else if cnt == 0 {
		return nil, errors.New("no user found")
	} else {
		return nil, errors.New("more than one row")
	}
}

func FindUserByTag(tag string) ([]byte, error) {
	var u dbtypes.User
	var cnt int = 0
	log.Println("[db-handler] looking database for tag:", tag)
	q := `SELECT * FROM user WHERE RFIDTag=?`
	rows, _ := db.Query(q, tag)
	for rows.Next() {
		rows.Scan(&u.Id, &u.Name, &u.IsAdmin, &u.RegistrationDate, &u.RFIDTag)
		cnt += 1
	}
	defer rows.Close()
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
	stmt, _ := db.Prepare(`UPDATE User SET Name=?, isAdmin=? WHERE id=?`)
	_, err := stmt.Exec(u.Name, u.IsAdmin, id)
	if err != nil {
		log.Println("[db-handler] error updating user:", err)
		return false
	} else {
		log.Println("[db-handler] user updated in the database!")
		return true
	}
}

func FindEquipment(s string) ([]byte, error) {
	var e_slice []*dbtypes.Equipment
	sub_string := strings.ToLower(s)
	log.Println("[db-handler] looking database for sub-string:", sub_string)
	stmt, _ := db.Prepare(
		`SELECT * FROM equipment WHERE LOWER(Name) LIKE ?`,
	)
	rows, err := stmt.Query("%" + sub_string + "%")
	if err != nil {
		log.Println("[db-handler] error looking for sub-string:", err)
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
		log.Println("[db-handler] error converting to JSON byte stream:", err)
		return []byte{}, err
	} else {
		log.Println("[db-handler] FindEquipment succeeded!")
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
