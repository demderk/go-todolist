package task

import (
	"errors"
	"time"
)

var (
	invalidRepeat = errors.New("invalid repeat instruction")
	repeatIsEmpty = errors.New("repeat is empty")
)

type Task struct {
	Id      int
	Date    time.Time
	Title   string
	Comment string
	Repeat  string
}
