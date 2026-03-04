package task

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	infrastructure "go-todolist/internal/infrastructure"
)

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (ts *TaskService) AddTask(task Task) (int64, error) {
	if err := setupDate(&task); err != nil {
		return 0, fmt.Errorf("incorrect date: %w", err)
	}
	id, err := ts.repo.AddTask(task)
	if err != nil {
		return 0, fmt.Errorf("task service error: %w", err)
	}
	return id, nil
}

func (ts *TaskService) GetAllTasks(limit int) ([]Task, error) {
	res, err := ts.repo.GetAllTasks(limit)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Date.After(res[j].Date)
	})
	return res, nil
}

func (ts *TaskService) GetTask(id int) (Task, error) {
	res, err := ts.repo.GetTask(id)
	if err != nil {
		if errors.Is(err, infrastructure.ErrItemNotFound) {
			return Task{}, infrastructure.ErrItemNotFound
		} else {
			return Task{}, fmt.Errorf("database error: %w", err)
		}
	}
	return res, nil
}

func (ts *TaskService) UpdateTask(task Task) error {
	if err := ts.repo.UpdateTask(task); err != nil {
		return fmt.Errorf("repository error: %w", err)
	}
	return nil
}

func NextDate(now time.Time, dstart string, repeat string) (time.Time, error) {
	if strings.TrimSpace(repeat) == "" {
		return time.Time{}, repeatIsEmpty
	}

	task := strings.Split(repeat, " ")

	date, err := time.Parse(infrastructure.TimeFormat, dstart)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid start date format: %w", err)
	}

	instruction := task[0]
	switch instruction {
	case "d":
		if len(task) != 2 {
			return time.Time{}, invalidRepeat
		}
		if count, err := strconv.Atoi(task[1]); err == nil {
			if count > 400 {
				return time.Time{}, fmt.Errorf("maximum day was exceeded")
			}
			res := addDateRepeating(date, now, 0, 0, count)
			return res, nil
		} else {
			return time.Time{}, fmt.Errorf("nextDate() failed to parse instruction: %w", err)
		}

	case "y":
		res := addDateRepeating(date, now, 1, 0, 0)
		return res, nil
	default:
		log.Printf("invalid nextdate instruction: %v\n", instruction)
		return time.Time{}, invalidRepeat
	}
}

func (ts *TaskService) DeleteTask(id int) error {
	if err := ts.repo.DeleteTask(id); err != nil {
		return fmt.Errorf("repository error: %w", err)
	}
	return nil
}

func (ts *TaskService) CompleteTask(id int) error {
	task, err := ts.GetTask(id)
	if err != nil {
		return fmt.Errorf("repository error: %w", err)
	}

	if len(task.Repeat) != 0 {
		next, err := NextDate(time.Now(), task.Date.Format(infrastructure.TimeFormat), task.Repeat)
		if err != nil {
			return fmt.Errorf("can't get next date: %w", err)
		}
		task.Date = next
		ts.UpdateTask(task)
	} else {
		if err := ts.DeleteTask(id); err != nil {
			return fmt.Errorf("repository error: %w", err)
		}
	}
	return nil
}

func setupDate(task *Task) error {
	now := startOfDay(time.Now())
	startOfDate := startOfDay(task.Date)
	if startOfDate.Before(now) {
		if len(task.Repeat) == 0 {
			// если правила повторения нет, то берём сегодняшнее число
			task.Date = now
		} else {
			// в противном случае, берём вычисленную ранее следующую дату
			next, err := NextDate(now, task.Date.Format(infrastructure.TimeFormat), task.Repeat)
			if err != nil {
				return errors.New("date setup failed")
			}
			task.Date = next
		}
	}
	return nil
}

func addDateRepeating(from time.Time, now time.Time, years int, months int, days int) time.Time {
	res := from.AddDate(years, months, days)

	for {
		if res.After(now) {
			return res
		}
		res = res.AddDate(years, months, days)
	}
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
