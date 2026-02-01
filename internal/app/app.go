package app

import (
	router "go-todolist/internal/interfaces/http"
	"net/http"
)

func BuildRouter() http.Handler {
	return router.NewRouter()
}
