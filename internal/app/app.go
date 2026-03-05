package app

import (
	"database/sql"
	"net/http"
	"os"

	"go-todolist/internal/domain/task"
	router "go-todolist/internal/interfaces/http"
)

var (
	schema = `
	CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR,
	comment TEXT,
	repeat VARCHAR(128))
	`
)

func BuildRouter(repo task.TaskRepository) http.Handler {
	taskService := task.NewTaskService(repo)

	return router.NewRouter(taskService)
}

func BuildBD() (*sql.DB, error) {
	dbFile := "./data/scheduler.db"
	_, err := os.Stat(dbFile)

	install := false
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}
	if install {
		if err := createTables(db); err != nil {
			db.Close()
			return nil, err
		}
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	if _, err := db.Exec(schema); err != nil {
		return err
	}
	return nil
}
