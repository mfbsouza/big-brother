package router

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mfbsouza/big-brother/internal/db-handler"
)

func NewRouter() http.Handler {
	mux := chi.NewRouter()
	mux.Route("/user/{token}", func(r chi.Router) {
		r.Get("/", getUser)
		r.Put("/", updateUser)
	})
	mux.Post("/user", createUser)
	return mux
}

func getUser(w http.ResponseWriter, r *http.Request) {
	u, e := dbhandler.FindUser(chi.URLParam(r, "token"))
	if e == nil {
		w.Write(u)
	} else {
		log.Println("[db-router] error search for user:", e)
		w.WriteHeader(http.StatusNotFound)
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	ok := dbhandler.UpdateUser(chi.URLParam(r, "token"), body)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	id := dbhandler.CreateUser(body)
	if id == 0 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.Write([]byte(strconv.FormatInt(id, 10)))
	}
}
