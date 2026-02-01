package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Handle("/*", http.FileServer(http.Dir("web")))

	return r
}
