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

const (
	FIRST_U_NAME = "Matheus Souza"
	FIRST_U_TAG  = "5B320FE6"
	FIRST_E_NAME = "Arduino Uno"
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

func DeleteUserById(id string) error {
	if id == "1" {
		return errors.New("Cannot delete the first user!")
	}
	stmt, _ := db.Prepare(`DELETE FROM user WHERE id=?`)
	_, err := stmt.Exec(id)
	if err != nil {
		log.Println("[db-handler] Error while deleting a user:", err)
		return errors.New("Error deleting user!")
	} else {
		return nil
	}
}

func DeleteUserByTag(tag string) error {
	if tag == FIRST_U_TAG {
		return errors.New("Cannot delete the first user!")
	}
	stmt, _ := db.Prepare(`DELETE FROM user WHERE RFIDTag=?`)
	_, err := stmt.Exec(tag)
	if err != nil {
		log.Println("[db-handler] Error while deleting a user:", err)
		return errors.New("Error deleting user!")
	} else {
		return nil
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

func UpdateUser(id string, j []byte) error {
	var u dbtypes.User
	json.Unmarshal(j, &u)
	stmt, _ := db.Prepare(`UPDATE User SET Name=?, isAdmin=? WHERE id=?`)
	_, err := stmt.Exec(u.Name, u.IsAdmin, id)
	if err != nil {
		log.Println("[db-handler] Error while updating a user:", err)
		return err
	} else {
		return nil
	}
}

func CreateEquipment(j []byte) (int64, error) {
	var e dbtypes.Equipment
	json.Unmarshal(j, &e)
	stmt, _ := db.Prepare(
		`INSERT INTO equipment (
			Name, 
			isInUse, 
			isBlocked, 
			user_ID
		) VALUES (?, ?, ?, ?)`,
	)
	res, err := stmt.Exec(e.Name, false, false, 0)
	id, err := res.LastInsertId()
	if err != nil {
		log.Println("[db-handler] Error while creating a new equipment:", err)
		return 0, errors.New("error creating the equipment!")
	} else {
		return id, nil
	}
}

func FindEquipmentById(id string) ([]byte, error) {
	var e_slice []*dbtypes.Equipment
	q := `SELECT * FROM equipment WHERE id=?`
	rows, _ := db.Query(q, id)
	for rows.Next() {
		e := &dbtypes.Equipment{}
		rows.Scan(&e.Id, &e.Name, &e.IsInUse, &e.IsBlocked, &e.UserId)
		e_slice = append(e_slice, e)
	}
	defer rows.Close()
	if len(e_slice) != 1 {
		log.Println("[db-handler] Error while searching for equipment Id:", id)
		return []byte{}, errors.New("Length of the equipment slice is different than 1")
	} else {
		bytestream, _ := json.Marshal(e_slice[0])
		return bytestream, nil
	}
}

func FindEquipmentByString(s string) ([]byte, error) {
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

func UpdateEquipment(id string, j []byte) error {
	var e dbtypes.Equipment
	json.Unmarshal(j, &e)
	stmt, _ := db.Prepare(
		`UPDATE equipment SET
			Name=?,
			isInUse=?,
			isBlocked=?,
			user_ID=?
		WHERE id=?`)
	_, err := stmt.Exec(e.Name, e.IsInUse, e.IsBlocked, e.UserId, id)
	if err != nil {
		log.Println("[db-handler] Error while updating a user:", err)
		return err
	} else {
		return nil
	}
}

func DeleteEquipmentById(id string) error {
	stmt, _ := db.Prepare(`DELETE FROM equipment WHERE id=?`)
	_, err := stmt.Exec(id)
	if err != nil {
		log.Println("[db-handler] Error while deleting a equipment:", err)
		return errors.New("Error deleting equipment!")
	} else {
		return nil
	}
}

func RentEquipment(id string, j []byte) error {
	var u dbtypes.User
	var blocked bool
	json.Unmarshal(j, &u)

	query := `SELECT isBlocked FROM equipment WHERE id=?`
	rows, err := db.Query(query, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Scan(&blocked)
	if blocked {
		return errors.New("Cant rent a blocked equipment!")
	}

	insert := `INSERT INTO log (user_ID, equipment_ID, UsageDate) VALUES (?, ?, ?)`
	db.Exec(insert, u.Id, id, time.Now().UTC())

	update := `UPDATE equipment SET isInUse=?, user_ID=? WHERE id=?`
	stmt, _ := db.Prepare(update)
	_, err = stmt.Exec(true, u.Id, id)
	if err != nil {
		return errors.New("Cant update the equipment!")
	}

	return nil
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

	stmt, _ = db.Prepare(
		`CREATE TABLE IF NOT EXISTS log (
			ID INTEGER PRIMARY KEY,
			user_ID INTEGER,
			equipment_ID INTEGER,
			UsageDate TIMESTAMP,
			FOREIGN KEY (user_ID) REFERENCES user(ID),
			FOREIGN KEY (equipment_ID) REFERENCES equipment(ID)
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
	stmt.Exec(FIRST_U_NAME, true, time.Now().UTC(), FIRST_U_TAG)

	// create the first equipment
	stmt, _ = db.Prepare(
		`INSERT INTO equipment (
			Name, 
			isInUse, 
			isBlocked, 
			user_ID
		) VALUES (?, ?, ?, ?)`,
	)
	stmt.Exec(FIRST_E_NAME, false, false, 0)
}
