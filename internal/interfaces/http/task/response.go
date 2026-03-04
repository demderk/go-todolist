package task

import (
	"strconv"

	"go-todolist/internal/domain/task"
)

type TaskResponseDTO struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func NewTaskResponseDTO(task task.Task) TaskResponseDTO {
	return TaskResponseDTO{
		Id:      strconv.Itoa(task.Id),
		Date:    task.Date.Format("20060102"),
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}
}
