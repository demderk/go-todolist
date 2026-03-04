package app

import (
	"database/sql"
	"net/http"
	"os"

	"go-todolist/internal/domain/task"
	router "go-todolist/internal/interfaces/http"
)

func BuildRouter(repo task.TaskRepository) http.Handler {
	taskService := task.NewTaskService(repo)

	return router.NewRouter(taskService)
}

func BuildBD() (*sql.DB, bool, error) {
	dbFile := "./data/scheduler.db"
	_, err := os.Stat(dbFile)

	install := false
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, false, err
	}
	return db, install, nil
}
