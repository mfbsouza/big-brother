package user

import (
	"log"
	"time"

	"github.com/google/uuid"
)

type User struct {
	name string
}

type Session struct {
	token  string
	expiry time.Time
}

func (s Session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

// user cache
var uc = make(map[string]*User)

// session cache
var sc = make(map[string]*Session)

func ValidateUserToken(t string) bool {
	// check the user is present in the cache
	_, exists := uc[t]
	if !exists {
		// TODO: ask the database
		log.Printf("[user] checking token %s in the database\n", t)
		uc[t] = &User{
			name: "Matheus Souza",
		}
		log.Printf("[user] username %s with token %s added to the cache\n",
			uc[t].name, t)
		return true
	}
	log.Printf("[user] token %s found in user cache\n", t)
	return exists
}

func ValidateSessionId(id string) bool {
	// check if the session is present in the cache
	_, exists := sc[id]
	return exists
}

func CreateSessionId(t string) (string, time.Time) {
	sessionId := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	sc[sessionId] = &Session{
		token:  t,
		expiry: expiresAt,
	}

	log.Printf("[user] created Session entry for token %s\n", t)
	return sessionId, expiresAt
}

func DeleteSessionId(id string) {
	sc[id] = nil
}
