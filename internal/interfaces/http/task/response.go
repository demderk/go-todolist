package task

import (
	"strconv"

	"go-todolist/internal/domain/task"
	"go-todolist/internal/infrastructure"
)

type TaskResponseDTO struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func ToTaskResponseDTO(task *task.Task) TaskResponseDTO {
	return TaskResponseDTO{
		Id:      strconv.Itoa(task.Id),
		Date:    task.Date.Format(infrastructure.TimeFormat),
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}
}
