package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"go-todolist/internal/domain/task"
	"time"

	infrastructure "go-todolist/internal/infrastructure/db"

	_ "modernc.org/sqlite"
)

const (
	timeFormat = "20060102"
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

type TaskRepoDB struct {
	db *sql.DB
}

func NewTaskRepo(db *sql.DB) (*TaskRepoDB, error) {
	return &TaskRepoDB{db: db}, nil
}

func (t *TaskRepoDB) CreateTables() error {
	if _, err := t.db.Exec(schema); err != nil {
		return err
	}
	return nil
}

func (r *TaskRepoDB) AddTask(task task.Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`

	res, err := r.db.Exec(
		query,
		sql.Named("date", task.Date.Format(timeFormat)),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)

	if err == nil {
		id, err = res.LastInsertId()
	}

	return id, err
}

func (r *TaskRepoDB) GetAllTasks(limit int) ([]task.Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler LIMIT :limit`

	rows, err := r.db.Query(query, sql.Named("limit", limit))
	if err != nil {
		return nil, fmt.Errorf("databse error: %w", err)
	}
	defer rows.Close()

	result := make([]task.Task, 0, 50)

	for rows.Next() {
		// Тут может лучше структурой сделать?
		var id int
		var date string
		var title string
		var comment string
		var repeat string

		if err := rows.Scan(&id, &date, &title, &comment, &repeat); err != nil {
			return nil, fmt.Errorf("data read failed: %w", err)
		}

		parsedDate, err := time.Parse(timeFormat, date)

		if err != nil {
			return nil, fmt.Errorf("incorrect date in database: %w", err)
		}

		row := task.Task{
			Id:      id,
			Date:    parsedDate,
			Title:   title,
			Comment: comment,
			Repeat:  repeat,
		}

		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("databse iteration error: %w", err)
	}

	return result, nil
}

func (r *TaskRepoDB) GetTask(id int) (task.Task, error) {
	query := `SELECT date, title, comment, repeat FROM scheduler WHERE id = :id`

	var date string
	var title string
	var comment string
	var repeat string

	err := r.db.QueryRow(query, sql.Named("id", id)).Scan(&date, &title, &comment, &repeat)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return task.Task{}, infrastructure.ErrItemNotFound
		}
	}

	parsedDate, err := time.Parse(timeFormat, date)

	row := task.Task{
		Id:      id,
		Date:    parsedDate,
		Title:   title,
		Comment: comment,
		Repeat:  repeat,
	}

	return row, nil
}

func (r *TaskRepoDB) UpdateTask(task task.Task) error {
	query := `
	UPDATE
		scheduler
	SET 
		date = :date,
		title = :title,
		comment = :comment,
		repeat = :repeat
	WHERE
		id = :id
	`

	res, err := r.db.Exec(
		query,
		sql.Named("id", task.Id),
		sql.Named("date", task.Date.Format(timeFormat)),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)

	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return infrastructure.ErrBadArgument
	}

	return nil
}

func (r *TaskRepoDB) DeleteTask(id int) error {
	query := `DELETE FROM scheduler WHERE id = :id;`

	res, err := r.db.Exec(
		query,
		sql.Named("id", id),
	)

	if err != nil {
		return err
	}

	if count, err := res.RowsAffected(); err != nil {
		if count == 0 {
			return infrastructure.ErrBadArgument
		}
	} else {
		return err
	}

	return nil
}
