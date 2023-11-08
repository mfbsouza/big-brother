package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mfbsouza/big-brother-web/internal/handlers"
)

func newRouter() http.Handler {
	mux := chi.NewRouter()
	mux.Get("/", handlers.Home)
	mux.Get("/about", handlers.About)
	mux.Get("/login", handlers.Login)
	mux.Post("/login", handlers.LoginCheckCredentials)
	return mux
}
