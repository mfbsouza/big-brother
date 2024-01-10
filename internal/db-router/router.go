package router

import (
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mfbsouza/big-brother/internal/db-handler"
)

func NewRouter() http.Handler {
	mux := chi.NewRouter()
	// user routes
	mux.Route("/user/id/{token}", func(r chi.Router) {
		r.Get("/", getUserById)
		r.Put("/", updateUser)
		r.Delete("/", deleteUserById)
	})
	mux.Route("/user/tag/{tag}", func(r chi.Router) {
		r.Get("/", getUserByTag)
		r.Delete("/", deleteUserByTag)
	})
	mux.Post("/user", createUser)

	// equipment routes
	mux.Route("/equip/{string}", func(r chi.Router) {
		r.Get("/", getEquipment)
	})
	return mux
}

func getUserById(w http.ResponseWriter, r *http.Request) {
	u, e := dbhandler.FindUserById(chi.URLParam(r, "token"))
	if e == nil {
		w.Write(u)
	} else {
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

func deleteUserById(w http.ResponseWriter, r *http.Request) {
	e := dbhandler.DeleteUserById(chi.URLParam(r, "token"))
	if e != nil {
		w.WriteHeader(http.StatusNotFound)
	}
}

func getUserByTag(w http.ResponseWriter, r *http.Request) {
	u, e := dbhandler.FindUserByTag(chi.URLParam(r, "tag"))
	if e == nil {
		w.Write(u)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func deleteUserByTag(w http.ResponseWriter, r *http.Request) {
	e := dbhandler.DeleteUserByTag(chi.URLParam(r, "tag"))
	if e != nil {
		w.WriteHeader(http.StatusNotFound)
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

func getEquipment(w http.ResponseWriter, r *http.Request) {
	equipments, err := dbhandler.FindEquipment(chi.URLParam(r, "string"))
	if err == nil {
		w.Write(equipments)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
