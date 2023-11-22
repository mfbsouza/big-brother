package dbtypes

import (
	"time"
)

type User struct {
	Id               int       `json:"id"`
	Name             string    `json:"name"`
	IsAdmin          bool      `json:"admin"`
	RegistrationDate time.Time `json:"r_date"`
	RFIDTag          string    `json:"tag"`
}
