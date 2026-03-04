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

	db, new, err := app.BuildBD()
	if err != nil {
		panic(err)
	}

	taskRepo, err := sqlite.NewTaskRepo(db)
	if err != nil {
		panic(err)
	}

	if new {
		taskRepo.CreateTables()
	}

	router := app.BuildRouter(taskRepo)

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
