package task

type TaskRepository interface {
	AddTask(task Task) (int64, error)
	GetAllTasks(limit int) ([]Task, error)
	GetTask(id int) (Task, error)
	UpdateTask(task Task) error
	DeleteTask(id int) error
}
