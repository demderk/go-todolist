package task

import (
	"time"
)

type Task struct {
	id      int
	date    time.Time
	title   string
	comment string
	repeat  string
}
