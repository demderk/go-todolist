package sqlite

import (
	"database/sql"
	"go-todolist/internal/domain/task"
	"os"

	_ "modernc.org/sqlite"
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

type TaskRepo struct{}

func NewTaskRepo() (task.Repository, error) {
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
	defer db.Close()

	if install {
		if _, err := db.Exec(schema); err != nil {
			return nil, err
		}
	}

	return &TaskRepo{}, nil
}
