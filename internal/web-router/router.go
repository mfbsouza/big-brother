package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mfbsouza/big-brother/internal/page-renderer"
)

// function NewRouter creates a http.Handler with all
// the routes supported by this web server
func NewRouter() http.Handler {
	mux := chi.NewRouter()
	mux.Get("/", home)
	mux.Get("/about", about)
	mux.Get("/login", login)
	mux.Post("/login", loginCheckCredentials)
	return mux
}

func home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "home.html")
}

func about(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "about.html")
}

func login(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "login.html")
}

func loginCheckCredentials(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("i am here!\n")
	r.ParseForm()
	for key, value := range r.Form {
		fmt.Printf("key: %s, value = %s \n", key, value)
	}
	fmt.Fprintln(w, "Done!")
}
