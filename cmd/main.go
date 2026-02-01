package main

import (
	"fmt"
	"go-todolist/internal/app"
	"net/http"
	"os"
	"strconv"
)

func main() {
	port := setupPort()
	router := app.BuildRouter()
	http.ListenAndServe(fmt.Sprintf(":%v", port), router)
}

func setupPort() int {
	startupPort := os.Getenv("TODO_PORT")
	port := 7540
	if startupPort != "" {
		port, err := strconv.Atoi(os.Getenv("TODO_PORT"))
		if err != nil {
			fmt.Println("invalid port, app port was set to DEFAULT (7540)")
			port = 7540
		}
		if port < 0 || port > 65535 {
			fmt.Println("invalid port, app port was set to DEFAULT (7540)")
			port = 7540
		}
	}
	return port
}
