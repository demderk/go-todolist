package router

import (
	domain "go-todolist/internal/domain/task"
	"go-todolist/internal/interfaces/http/task"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(taskRepo *domain.TaskService) chi.Router {
	r := chi.NewRouter()

	th := task.NewTaskHandler(*taskRepo)

	r.Handle("/*", http.FileServer(http.Dir("web")))
	r.Get("/api/nextdate", th.NextDate)

	r.Post("/api/task", th.AddTask)
	r.Get("/api/task", th.GetTask)
	r.Put("/api/task", th.UpdateTask)
	r.Delete("/api/task", th.DeleteTask)
	r.Post("/api/task/done", th.CompleteTask)
	r.Get("/api/tasks", th.GetAllTasks)

	return r
}
