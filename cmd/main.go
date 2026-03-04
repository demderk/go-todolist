package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"go-todolist/internal/app"
	"go-todolist/internal/infrastructure/db/sqlite"
)

func main() {
	port := setupPort()

	db, err := app.BuildBD()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	taskRepo, err := sqlite.NewTaskRepo(db)
	if err != nil {
		panic(err)
	}

	router := app.BuildRouter(taskRepo)

	http.ListenAndServe(fmt.Sprintf(":%v", port), router)
}

func setupPort() int {
	startupPort := os.Getenv("TODO_PORT")
	port := 7540
	if startupPort != "" {
		_, err := strconv.Atoi(os.Getenv("TODO_PORT"))
		if err != nil {
			fmt.Println("invalid port, app port was set to DEFAULT (7540)")
			port = 7540
		}
	}
	return port
}
