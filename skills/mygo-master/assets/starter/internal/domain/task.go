package domain

import (
	"errors"
	"strings"
)

var ErrEmptyTitle = errors.New("a task needs a title")

type Task struct {
	ID    int
	Title string
	Done  bool
}

func NewTask(id int, title string) (Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Task{}, ErrEmptyTitle
	}
	return Task{ID: id, Title: title}, nil
}

func Remaining(tasks []Task) int {
	count := 0
	for _, task := range tasks {
		if !task.Done {
			count++
		}
	}
	return count
}
