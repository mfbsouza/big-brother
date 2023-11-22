package router

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mfbsouza/big-brother/internal/db-handler"
)

func NewRouter() http.Handler {
	mux := chi.NewRouter()
	mux.Route("/user/{token}", func(r chi.Router) {
		r.Get("/", getUser)
		r.Post("/", createUser)
	})
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

func createUser(w http.ResponseWriter, r *http.Request) {
}
