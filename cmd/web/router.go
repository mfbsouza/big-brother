package main

import (
	"net/http"

	"github.com/mfbsouza/big-brother-web/pkg/handlers"
	"github.com/go-chi/chi/v5"
)

func newRouter() http.Handler {
	mux := chi.NewRouter()
	mux.Get("/", handlers.Home)
	mux.Get("/about", handlers.About)
	return mux
}
