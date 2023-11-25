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

type Equipment struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	IsInUse   bool   `json:"in_use"`
	IsBlocked bool   `json:"blocked"`
	UserId    int    `json:"user_id"`
}
