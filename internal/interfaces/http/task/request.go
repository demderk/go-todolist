package task

import (
	"errors"
	"fmt"
	"go-todolist/internal/domain/task"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	timeFormat = "20060102"
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

func (r *TaskRequestDTO) ToTask() (task.Task, error) {
	date, err := time.Parse(timeFormat, r.Date)
	if err != nil {
		return task.Task{}, fmt.Errorf("task convertion failed: %w", err)
	}
	var id int
	if len(r.Id) == 0 {
		id = 0
	} else {
		num, err := strconv.Atoi(r.Id)
		if err != nil {
			return task.Task{}, fmt.Errorf("failed to parse non-nil id")
		}
		id = num
	}
	return task.Task{
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

	if len(task) == 1 {
		if task[0] != "y" {
			return false
		}
	} else if len(task) == 2 {
		if task[0] != "d" {
			return false
		}
	} else {
		return false
	}

	var regex = regexp.MustCompile(`^[dy]$`)
	if !regex.MatchString(task[0]) {
		return false
	}

	return true
}
