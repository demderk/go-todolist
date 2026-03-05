package task

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go-todolist/internal/domain/task"
	"go-todolist/internal/infrastructure"
)

type TaskRequestDTO struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func (r *TaskRequestDTO) Validate() error {
	if r.Title == "" {
		return errors.New("title is required")
	}

	if len(r.Repeat) != 0 && !validateRepeat(r.Repeat) {
		return errors.New("wrong repeat instruction")
	}

	return nil
}

func (r *TaskRequestDTO) ToTask() (*task.Task, error) {
	date, err := time.Parse(infrastructure.TimeFormat, r.Date)
	if err != nil {
		return nil, fmt.Errorf("task convertion failed: %w", err)
	}
	var id int
	if len(r.Id) == 0 {
		id = 0
	} else {
		num, err := strconv.Atoi(r.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to parse non-nil id")
		}
		id = num
	}
	return &task.Task{
		Id:      id,
		Date:    date,
		Title:   r.Title,
		Comment: r.Comment,
		Repeat:  r.Repeat,
	}, nil
}

func validateRepeat(repeat string) bool {
	trim := strings.TrimSpace(repeat)
	task := strings.Split(trim, " ")

	if len(task) != 1 && len(task) != 2 {
		return false
	}
	if len(task) == 1 && task[0] != "y" {
		return false
	}
	if len(task) == 2 && task[0] != "d" {
		return false
	}

	var regex = regexp.MustCompile(`^[dy]$`)
	if !regex.MatchString(task[0]) {
		return false
	}

	return true
}
