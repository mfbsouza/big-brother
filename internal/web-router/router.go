package router

import (
	"log"
	"net/http"
	"io"
	"encoding/json"
	"bytes"
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/mfbsouza/big-brother/internal/page-renderer"
	"github.com/mfbsouza/big-brother/internal/user-manager"
	"github.com/mfbsouza/big-brother/internal/db-types"
)

// function NewRouter creates a http.Handler with all
// the routes supported by this web server
func NewRouter() http.Handler {
	mux := chi.NewRouter()
	mux.Get("/", home)
	mux.Get("/about", about)
	mux.Get("/inuse", inUse)
	mux.Get("/free", free)
	mux.Get("/equip/insert", insertPage)
	mux.Get("/equip/remove", removePage)
	mux.Post("/login", signIn)
	mux.Post("/equip/insert", insertData)
	mux.Post("/equip/remove", removeData)
	return mux
}

func removePage(w http.ResponseWriter, r *http.Request) {
	if user.VerifyClearance(r) && user.IsAdmin(r) {
		render.RenderTemplate(w, "del-equip.html", nil)
	} else {
		render.RenderTemplate(w, "login.html", nil)
	}
}

func insertPage(w http.ResponseWriter, r *http.Request) {
	if user.VerifyClearance(r) && user.IsAdmin(r) {
		render.RenderTemplate(w, "new-equip.html", nil)
	} else {
		render.RenderTemplate(w, "login.html", nil)
	}
}

func insertData(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("[web-router] failed parsing form:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	name := r.FormValue("equip-name")
	if len(name) == 0 {
		log.Println("[web-router] failed reading 'equip-name' key from form")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	e := dbtypes.Equipment {
		Id: 0,
		Name: name,
		IsInUse: false,
		IsBlocked: false,
		UserId: 0,
	}
	bytestream, err := json.Marshal(e)
	if err != nil {
		log.Println("[web-router] failed converting struct to json string")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	requestURL := "http://localhost:3030/equip"
	res, err := http.Post(requestURL, "application/json", bytes.NewBuffer(bytestream))
	if err != nil {
		log.Println("[web-router] failed doing post call to create new equipment")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if res.StatusCode != http.StatusOK {
		log.Println("[web-router] failed to create new equipment at the database")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		io.WriteString(w, "<h3>Success! Equipment was added to the database</h3>")
	}
}

func removeData(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("[web-router] failed parsing form:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id := r.FormValue("id-equip")
	if len(id) == 0 {
		log.Println("[web-router] failed reading 'id-equip' key from form")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestURL := fmt.Sprintf("http://localhost:3030/equip/id/%s", id)
	req, err := http.NewRequest("DELETE", requestURL, nil)
	if err != nil {
		log.Println("[web-router] Error creating request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println("[web-router] Error making the request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Println("[web-router] failed to create new equipment at the database")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		io.WriteString(w, "<h3>Success! Equipment was removed from the database</h3>")
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	var equipments []dbtypes.Equipment
	if user.VerifyClearance(r) {
		requestURL := "http://localhost:3030/equip/free"
		res, err := http.Get(requestURL)
		if err != nil {
			log.Println("[web-router] failed requesting data to the database", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		// } else if res.StatusCode != http.StatusOK {
		// 	log.Println("[web-router] no free equipment!")
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	return
		} else {
			body, _ := io.ReadAll(res.Body)
			err := json.Unmarshal(body, &equipments)
			if err != nil {
				log.Println("[web-router] failed parsing JSON", err)
				// w.WriteHeader(http.StatusBadRequest)
				// return
			}
			render.RenderTemplate(w, "home.html", equipments)
		}
	} else {
		render.RenderTemplate(w, "login.html", nil)
	}
}

func inUse(w http.ResponseWriter, r *http.Request) {
	var equipments []dbtypes.Equipment
	if user.VerifyClearance(r) {
		requestURL := "http://localhost:3030/equip/inuse"
		res, err := http.Get(requestURL)
		if err != nil {
			log.Println("[web-router] failed requesting data to the database", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		// } else if res.StatusCode != http.StatusOK {
		// 	log.Println("[web-router] no free equipment!")
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	return
		} else {
			body, _ := io.ReadAll(res.Body)
			err := json.Unmarshal(body, &equipments)
			if err != nil {
				log.Println("[web-router] failed parsing JSON", err)
				// w.WriteHeader(http.StatusBadRequest)
				// return
			}
			render.RenderTemplate(w, "equipment-list.html", equipments)
		}
	} else {
		render.RenderTemplate(w, "login.html", nil)
	}
}

func free(w http.ResponseWriter, r *http.Request) {
	var equipments []dbtypes.Equipment
	if user.VerifyClearance(r) {
		requestURL := "http://localhost:3030/equip/free"
		res, err := http.Get(requestURL)
		if err != nil {
			log.Println("[web-router] failed requesting data to the database", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		// } else if res.StatusCode != http.StatusOK {
		// 	log.Println("[web-router] no free equipment!")
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	return
		} else {
			body, _ := io.ReadAll(res.Body)
			err := json.Unmarshal(body, &equipments)
			if err != nil {
				log.Println("[web-router] failed parsing JSON", err)
				// io.WriteString(w, "<h3>No free equipment!</h3>")
				// return
			}
			render.RenderTemplate(w, "equipment-list.html", equipments)
		}
	} else {
		render.RenderTemplate(w, "login.html", nil)
	}
}

func about(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "about.html", nil)
}

func signIn(w http.ResponseWriter, r *http.Request) {
	// check if the user is already logged in
	if user.VerifyClearance(r) {
		log.Println("[web-router] user already logged in")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get the user token from the form
	err := r.ParseForm()
	if err != nil {
		log.Println("[web-router] failed parsing form:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	token := r.FormValue("token")
	if len(token) == 0 {
		log.Println("[web-router] failed reading 'token' key from form")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate user token
	exists := user.ValidateUserToken(token)
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
	}

	// create session token
	s := user.CreateSessionId(token)

	// set the cookie
	http.SetCookie(w, s.Cookie)
}
