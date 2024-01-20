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

	mux.Route("/user", func(r chi.Router) {
		r.Post("/", createUser)

		r.Route("/id/{token}", func(r chi.Router) {
			r.Get("/", getUserById)
			r.Put("/", updateUser)
			r.Delete("/", deleteUserById)
		})

		r.Route("/tag/{tag}", func(r chi.Router) {
			r.Get("/", getUserByTag)
			r.Delete("/", deleteUserByTag)
		})
	})

	mux.Route("/equip", func(r chi.Router) {
		r.Post("/", createEquipment)

		r.Route("/id/{token}", func(r chi.Router) {
			r.Get("/", getEquipmentById)
			r.Put("/", updateEquipment)
			r.Delete("/", deleteEquipmentById)
		})

		r.Route("/name/{string}", func(r chi.Router) {
			r.Get("/", getEquipmentByString)
		})
	})

	mux.Route("/request", func(r chi.Router) {
		r.Route("/rent/{id}", func(r chi.Router) {
			r.Post("/", rentEquipment)
		})

		r.Route("/return/{id}", func(r chi.Router) {
			r.Post("/", returnEquipment)
		})
	})

	return mux
}

func getUserById(w http.ResponseWriter, r *http.Request) {
	u, err := dbhandler.FindUserById(chi.URLParam(r, "token"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Write(u)
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	err := dbhandler.UpdateUser(chi.URLParam(r, "token"), body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func deleteUserById(w http.ResponseWriter, r *http.Request) {
	err := dbhandler.DeleteUserById(chi.URLParam(r, "token"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
}

func getUserByTag(w http.ResponseWriter, r *http.Request) {
	u, err := dbhandler.FindUserByTag(chi.URLParam(r, "tag"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Write(u)
	}
}

func deleteUserByTag(w http.ResponseWriter, r *http.Request) {
	err := dbhandler.DeleteUserByTag(chi.URLParam(r, "tag"))
	if err != nil {
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

func createEquipment(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	id, err := dbhandler.CreateEquipment(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.Write([]byte(strconv.FormatInt(id, 10)))
	}
}

func getEquipmentById(w http.ResponseWriter, r *http.Request) {
	u, err := dbhandler.FindEquipmentById(chi.URLParam(r, "token"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Write(u)
	}
}

func getEquipmentByString(w http.ResponseWriter, r *http.Request) {
	equipments, err := dbhandler.FindEquipmentByString(chi.URLParam(r, "string"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Write(equipments)
	}
}

func updateEquipment(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	err := dbhandler.UpdateEquipment(chi.URLParam(r, "token"), body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func deleteEquipmentById(w http.ResponseWriter, r *http.Request) {
	err := dbhandler.DeleteEquipmentById(chi.URLParam(r, "token"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
}

func rentEquipment(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	err := dbhandler.RentEquipment(chi.URLParam(r, "id"), body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func returnEquipment(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	err := dbhandler.ReturnEquipment(chi.URLParam(r, "id"), body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}
