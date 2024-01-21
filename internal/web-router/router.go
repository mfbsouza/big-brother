package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mfbsouza/big-brother/internal/db-types"
	"github.com/mfbsouza/big-brother/internal/page-renderer"
	"github.com/mfbsouza/big-brother/internal/user-manager"
)

// function NewRouter creates a http.Handler with all
// the routes supported by this web server
func NewRouter() http.Handler {
	mux := chi.NewRouter()
	mux.Get("/", home)
	mux.Get("/about", about)
	mux.Get("/inuse", inUse)
	mux.Get("/free", free)
	mux.Get("/logout", logout)
	mux.Get("/equip/insert", insertPage)
	mux.Get("/equip/remove", removePage)
	mux.Get("/equip/block", blockPage)
	mux.Get("/equip/unblock", unblockPage)
	mux.Get("/user/remove", removeUserPage)
	mux.Get("/user/admin", adminUserPage)
	mux.Get("/user/update", updateUserPage)
	mux.Get("/log/user", logUserPage)
	mux.Get("/log/equip", logEquipPage)
	mux.Post("/login", signIn)
	mux.Post("/equip/insert", insertData)
	mux.Post("/equip/remove", removeData)
	mux.Post("/equip/block", blockData)
	mux.Post("/equip/unblock", unblockData)
	mux.Post("/user/remove", removeUserData)
	mux.Post("/user/admin", adminUserData)
	mux.Post("/user/update", updateUserData)
	mux.Post("/log/user", logUserData)
	mux.Post("/log/equip", logEquipData)
	mux.Post("/search", searchData)
	return mux
}

func logout(w http.ResponseWriter, r *http.Request) {
	if user.VerifyClearance(r) {
		sId := user.GetSessionId(r)
		user.DeleteSessionId(sId)
		render.RenderTemplate(w, "login.html", nil)
	} else {
		render.RenderTemplate(w, "login.html", nil)
	}
}

func updateUserPage(w http.ResponseWriter, r *http.Request) {
	if user.VerifyClearance(r) {
		render.RenderTemplate(w, "update-user.html", nil)
	} else {
		render.RenderTemplate(w, "login.html", nil)
	}
}

func adminUserPage(w http.ResponseWriter, r *http.Request) {
	if user.VerifyClearance(r) && user.IsAdmin(r) {
		render.RenderTemplate(w, "upgrade-user.html", nil)
	} else if !user.IsAdmin(r) {
		render.RenderTemplate(w, "permission-error.html", nil)
	} else {
		render.RenderTemplate(w, "login.html", nil)
	}
}

func logEquipPage(w http.ResponseWriter, r *http.Request) {
	if user.VerifyClearance(r) {
		render.RenderTemplate(w, "log-equip.html", nil)
	} else {
		render.RenderTemplate(w, "login.html", nil)
	}
}

func logUserPage(w http.ResponseWriter, r *http.Request) {
	if user.VerifyClearance(r) {
		render.RenderTemplate(w, "log-user.html", nil)
	} else {
		render.RenderTemplate(w, "login.html", nil)
	}
}

func removeUserPage(w http.ResponseWriter, r *http.Request) {
	if user.VerifyClearance(r) && user.IsAdmin(r) {
		render.RenderTemplate(w, "del-user.html", nil)
	} else if !user.IsAdmin(r) {
		render.RenderTemplate(w, "permission-error.html", nil)
	} else {
		render.RenderTemplate(w, "login.html", nil)
	}
}

func unblockPage(w http.ResponseWriter, r *http.Request) {
	if user.VerifyClearance(r) && user.IsAdmin(r) {
		render.RenderTemplate(w, "unblock-equip.html", nil)
	} else if !user.IsAdmin(r) {
		render.RenderTemplate(w, "permission-error.html", nil)
	} else {
		render.RenderTemplate(w, "login.html", nil)
	}
}

func blockPage(w http.ResponseWriter, r *http.Request) {
	if user.VerifyClearance(r) && user.IsAdmin(r) {
		render.RenderTemplate(w, "block-equip.html", nil)
	} else if !user.IsAdmin(r) {
		render.RenderTemplate(w, "permission-error.html", nil)
	} else {
		render.RenderTemplate(w, "login.html", nil)
	}
}

func removePage(w http.ResponseWriter, r *http.Request) {
	if user.VerifyClearance(r) && user.IsAdmin(r) {
		render.RenderTemplate(w, "del-equip.html", nil)
	} else if !user.IsAdmin(r) {
		render.RenderTemplate(w, "permission-error.html", nil)
	} else {
		render.RenderTemplate(w, "login.html", nil)
	}
}

func insertPage(w http.ResponseWriter, r *http.Request) {
	if user.VerifyClearance(r) && user.IsAdmin(r) {
		render.RenderTemplate(w, "new-equip.html", nil)
	} else if !user.IsAdmin(r) {
		render.RenderTemplate(w, "permission-error.html", nil)
	} else {
		render.RenderTemplate(w, "login.html", nil)
	}
}

func searchData(w http.ResponseWriter, r *http.Request) {
	var equipments []dbtypes.Equipment
	if user.VerifyClearance(r) {
		err := r.ParseForm()
		if err != nil {
			log.Println("[web-router] failed parsing form:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		s := r.FormValue("search-string")
		if len(s) == 0 {
			log.Println("[web-router] failed reading 'search-string' key from form")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		requestURL := fmt.Sprintf("http://localhost:3030/equip/name/%s", s)
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
			render.SimpleRenderTemplate(w, "equip-list.html", equipments)
		}
	} else {
		render.RenderTemplate(w, "login.html", nil)
	}
}

func unblockData(w http.ResponseWriter, r *http.Request) {
	var e dbtypes.Equipment
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
	res, err := http.Get(requestURL)
	if err != nil {
		log.Println("[web-router] failed doing get call to read equipment")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if res.StatusCode != http.StatusOK {
		log.Println("[web-router] failed to read equipment at the database")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		body, _ := io.ReadAll(res.Body)
		json.Unmarshal(body, &e)
	}

	if !e.IsBlocked {
		io.WriteString(w, "<h3>Equipment already unblocked</h3>")
		return
	}

	e.IsBlocked = false
	bytestream, err := json.Marshal(e)
	if err != nil {
		log.Println("[web-router] failed converting struct to json string")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest("PUT", requestURL, bytes.NewBuffer(bytestream))
	if err != nil {
		log.Println("[web-router] Error creating request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err = client.Do(req)
	if err != nil {
		log.Println("[web-router] Error making the request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Println("[web-router] failed to update equipment from the database")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		io.WriteString(w, "<h3>Success! Equipment was unblocked</h3>")
	}
}

func blockData(w http.ResponseWriter, r *http.Request) {
	var e dbtypes.Equipment
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
	res, err := http.Get(requestURL)
	if err != nil {
		log.Println("[web-router] failed doing get call to read equipment")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if res.StatusCode != http.StatusOK {
		log.Println("[web-router] failed to read equipment at the database")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		body, _ := io.ReadAll(res.Body)
		json.Unmarshal(body, &e)
	}

	if e.IsBlocked {
		io.WriteString(w, "<h3>Equipment already blocked</h3>")
		return
	}

	e.IsBlocked = true
	bytestream, err := json.Marshal(e)
	if err != nil {
		log.Println("[web-router] failed converting struct to json string")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest("PUT", requestURL, bytes.NewBuffer(bytestream))
	if err != nil {
		log.Println("[web-router] Error creating request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err = client.Do(req)
	if err != nil {
		log.Println("[web-router] Error making the request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Println("[web-router] failed to update equipment from the database")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		io.WriteString(w, "<h3>Success! Equipment was blocked</h3>")
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
	e := dbtypes.Equipment{
		Id:        0,
		Name:      name,
		IsInUse:   false,
		IsBlocked: false,
		UserId:    0,
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
		log.Println("[web-router] failed to remove equipment from the database")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		io.WriteString(w, "<h3>Success! Equipment was removed from the database</h3>")
	}
}

func removeUserData(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("[web-router] failed parsing form:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id := r.FormValue("id-user")
	if len(id) == 0 {
		log.Println("[web-router] failed reading 'id-user' key from form")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestURL := fmt.Sprintf("http://localhost:3030/user/id/%s", id)
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
		log.Println("[web-router] failed to remove user from the database")
		// w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "<h3>Error! Cannot remove User from the database</h3>")
		return
	} else {
		io.WriteString(w, "<h3>Success! User was removed from the database</h3>")
	}
}

func updateUserData(w http.ResponseWriter, r *http.Request) {
	var u dbtypes.User
	err := r.ParseForm()
	if err != nil {
		log.Println("[web-router] failed parsing form:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id := r.FormValue("id-user")
	if len(id) == 0 {
		log.Println("[web-router] failed reading 'id-user' key from form")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	name := r.FormValue("name-user")
	if len(name) == 0 {
		log.Println("[web-router] failed reading 'name-user' key from form")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestURL := fmt.Sprintf("http://localhost:3030/user/id/%s", id)
	res, err := http.Get(requestURL)
	if err != nil {
		log.Println("[web-router] failed doing GET call to read user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if res.StatusCode != http.StatusOK {
		log.Println("[web-router] failed to read user at the database")
		// w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "<h3>Error! User does not exists in the database</h3>")
		return
	} else {
		body, _ := io.ReadAll(res.Body)
		json.Unmarshal(body, &u)
	}

	if u.IsAdmin && !user.IsAdmin(r) {
		io.WriteString(w, "<h3>You don't have permissions to modify Administrator data!</h3>")
		return
	}

	u.Name = name
	bytestream, err := json.Marshal(u)
	if err != nil {
		log.Println("[web-router] failed converting struct to json string")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest("PUT", requestURL, bytes.NewBuffer(bytestream))
	if err != nil {
		log.Println("[web-router] Error creating request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err = client.Do(req)
	if err != nil {
		log.Println("[web-router] Error making the request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Println("[web-router] failed to update user from the database")
		// w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "<h3>Error! Failed updating database</h3>")
		return
	} else {
		io.WriteString(w, "<h3>Success! User data updated</h3>")
	}
}

func adminUserData(w http.ResponseWriter, r *http.Request) {
	var u dbtypes.User
	err := r.ParseForm()
	if err != nil {
		log.Println("[web-router] failed parsing form:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id := r.FormValue("id-user")
	if len(id) == 0 {
		log.Println("[web-router] failed reading 'id-user' key from form")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestURL := fmt.Sprintf("http://localhost:3030/user/id/%s", id)
	res, err := http.Get(requestURL)
	if err != nil {
		log.Println("[web-router] failed doing GET call to read user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if res.StatusCode != http.StatusOK {
		log.Println("[web-router] failed to read user at the database")
		// w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "<h3>Error! User does not exists in the database</h3>")
		return
	} else {
		body, _ := io.ReadAll(res.Body)
		json.Unmarshal(body, &u)
	}

	if u.IsAdmin {
		io.WriteString(w, "<h3>User is already a Administrator</h3>")
		return
	}

	u.IsAdmin = true
	bytestream, err := json.Marshal(u)
	if err != nil {
		log.Println("[web-router] failed converting struct to json string")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest("PUT", requestURL, bytes.NewBuffer(bytestream))
	if err != nil {
		log.Println("[web-router] Error creating request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err = client.Do(req)
	if err != nil {
		log.Println("[web-router] Error making the request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Println("[web-router] failed to update user from the database")
		// w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "<h3>Error! Failed updating database</h3>")
		return
	} else {
		io.WriteString(w, "<h3>Success! User is now a Administrator</h3>")
	}
}

func logUserData(w http.ResponseWriter, r *http.Request) {
	var logs []dbtypes.Log
	err := r.ParseForm()
	if err != nil {
		log.Println("[web-router] failed parsing form:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id := r.FormValue("id-user")
	if len(id) == 0 {
		log.Println("[web-router] failed reading 'id-user' key from form")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestURL := fmt.Sprintf("http://localhost:3030/request/log/user/%s", id)
	res, err := http.Get(requestURL)
	if err != nil {
		log.Println("[web-router] failed requesting data to the database", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		body, _ := io.ReadAll(res.Body)
		err := json.Unmarshal(body, &logs)
		if err != nil {
			log.Println("[web-router] failed parsing JSON", err)
			// w.WriteHeader(http.StatusBadRequest)
			// return
		}
		render.SimpleRenderTemplate(w, "log-user-list.html", logs)
	}
}

func logEquipData(w http.ResponseWriter, r *http.Request) {
	var logs []dbtypes.Log
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

	requestURL := fmt.Sprintf("http://localhost:3030/request/log/equip/%s", id)
	res, err := http.Get(requestURL)
	if err != nil {
		log.Println("[web-router] failed requesting data to the database", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		body, _ := io.ReadAll(res.Body)
		err := json.Unmarshal(body, &logs)
		if err != nil {
			log.Println("[web-router] failed parsing JSON", err)
			// w.WriteHeader(http.StatusBadRequest)
			// return
		}
		render.SimpleRenderTemplate(w, "log-equip-list.html", logs)
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
		// w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "<h3>User already logged in</h3>")
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
