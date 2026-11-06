package usecase

import (
	"fmt"
	"slices"

	"starter/internal/domain"
)

type TaskStore interface {
	Load() ([]domain.Task, error)
	Save(tasks []domain.Task) error
}

type Tasks struct {
	store TaskStore
	items []domain.Task
	next  int
}

func NewTasks(store TaskStore) (*Tasks, error) {
	items, err := store.Load()
	if err != nil {
		return nil, fmt.Errorf("load tasks: %w", err)
	}
	next := 1
	for _, item := range items {
		next = max(next, item.ID+1)
	}
	return &Tasks{store: store, items: items, next: next}, nil
}

func (t *Tasks) List() []domain.Task {
	return t.items
}

func (t *Tasks) Remaining() int {
	return domain.Remaining(t.items)
}

func (t *Tasks) Add(title string) error {
	task, err := domain.NewTask(t.next, title)
	if err != nil {
		return err
	}
	if err := t.commit(append(slices.Clone(t.items), task)); err != nil {
		return err
	}
	t.next++
	return nil
}

func (t *Tasks) Toggle(id int) error {
	updated := slices.Clone(t.items)
	for i := range updated {
		if updated[i].ID == id {
			updated[i].Done = !updated[i].Done
			return t.commit(updated)
		}
	}
	return fmt.Errorf("no task %d", id)
}

func (t *Tasks) Remove(id int) error {
	return t.commit(slices.DeleteFunc(slices.Clone(t.items), func(task domain.Task) bool {
		return task.ID == id
	}))
}

func (t *Tasks) commit(updated []domain.Task) error {
	if err := t.store.Save(updated); err != nil {
		return fmt.Errorf("save tasks: %w", err)
	}
	t.items = updated
	return nil
}
