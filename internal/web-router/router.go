package router

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mfbsouza/big-brother/internal/page-renderer"
	"github.com/mfbsouza/big-brother/internal/user-manager"
)

// function NewRouter creates a http.Handler with all
// the routes supported by this web server
func NewRouter() http.Handler {
	mux := chi.NewRouter()
	mux.Get("/", home)
	mux.Get("/about", about)
	mux.Post("/login", signIn)
	return mux
}

func home(w http.ResponseWriter, r *http.Request) {
	if user.VerifyClearance(r) {
		render.RenderTemplate(w, "home.html")
	} else {
		render.RenderTemplate(w, "login.html")
	}
}

func about(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "about.html")
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
