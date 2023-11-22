package user

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mfbsouza/big-brother/internal/db-types"
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
	var u dbtypes.User
	// check the user is present in the cache
	_, exists := uc[t]
	if !exists {
		log.Printf("[user] checking token %s in the database\n", t)
		// TODO: variable base URL
		requestURL := fmt.Sprintf("http://localhost:3030/user/%s", t)
		res, err := http.Get(requestURL)
		if err != nil {
			log.Println("[user] failed requesting data to the database", err)
			return false
		} else if res.StatusCode != http.StatusOK {
			log.Println("[user] user not found in the database")
			return false
		} else {
			body, _ := io.ReadAll(res.Body)
			json.Unmarshal(body, &u)
			uc[t] = &User{
				name: u.Name,
			}
			log.Printf("[user] username %s with token %s added to the cache\n",
				u.Name, t)
			return true
		}
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
