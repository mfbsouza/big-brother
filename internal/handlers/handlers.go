package handlers

import (
	"fmt"
	"net/http"

	"github.com/mfbsouza/big-brother-web/internal/render"
)

func Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "home.html")
}

func About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "about.html")
}

func Login(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "login.html")
}

func LoginCheckCredentials(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("i am here!\n")
	r.ParseForm()
	for key, value := range r.Form {
		fmt.Printf("key: %s, value = %s \n", key, value)
	}
	fmt.Fprintln(w, "Done!")
}
