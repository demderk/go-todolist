package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go-todolist/internal/domain/task"
	infrastructure "go-todolist/internal/infrastructure"

	_ "modernc.org/sqlite"
)

type TaskRepoDB struct {
	db *sql.DB
}

func NewTaskRepo(db *sql.DB) (*TaskRepoDB, error) {
	return &TaskRepoDB{db: db}, nil
}

func (r *TaskRepoDB) AddTask(task *task.Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`

	res, err := r.db.Exec(
		query,
		sql.Named("date", task.Date.Format(infrastructure.TimeFormat)),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)

	if err == nil {
		id, err = res.LastInsertId()
	}

	return id, err
}

func (r *TaskRepoDB) GetAllTasks(limit int) ([]*task.Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler LIMIT :limit`

	rows, err := r.db.Query(query, sql.Named("limit", limit))
	if err != nil {
		return nil, fmt.Errorf("databse error: %w", err)
	}
	defer rows.Close()

	result := make([]*task.Task, 0, 50)

	for rows.Next() {

		row := &task.Task{}

		// Тут task.Date это time.Time, может есть более изящный способ это смэтчить?
		// Делать отдельный DTO оверхед, не?...
		var date string

		if err := rows.Scan(&row.Id, &date, &row.Title, &row.Comment, &row.Repeat); err != nil {
			return nil, fmt.Errorf("data read failed: %w", err)
		}

		parsedDate, err := time.Parse(infrastructure.TimeFormat, date)

		if err != nil {
			return nil, fmt.Errorf("incorrect date in database: %w", err)
		}

		row.Date = parsedDate

		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("databse iteration error: %w", err)
	}

	return result, nil
}

func (r *TaskRepoDB) GetTask(id int) (*task.Task, error) {
	query := `SELECT date, title, comment, repeat FROM scheduler WHERE id = :id`

	var date string

	task := &task.Task{Id: id}

	err := r.db.QueryRow(query, sql.Named("id", id)).Scan(&date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, infrastructure.ErrItemNotFound
		}
	}

	parsedDate, err := time.Parse(infrastructure.TimeFormat, date)

	if err != nil {
		return nil, fmt.Errorf("incorrect date in database: %w", err)
	}

	task.Date = parsedDate

	return task, nil
}

func (r *TaskRepoDB) UpdateTask(task *task.Task) error {
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
		sql.Named("date", task.Date.Format(infrastructure.TimeFormat)),
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
		return infrastructure.ErrIDNotFound
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
			return infrastructure.ErrIDNotFound
		}
	} else {
		return err
	}

	return nil
}
