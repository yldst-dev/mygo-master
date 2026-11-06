package usecase

import (
	"errors"
	"testing"

	"starter/internal/domain"
)

type memoryStore struct {
	saved []domain.Task
	fail  error
}

func (s *memoryStore) Load() ([]domain.Task, error) { return s.saved, nil }

func (s *memoryStore) Save(tasks []domain.Task) error {
	if s.fail != nil {
		return s.fail
	}
	s.saved = tasks
	return nil
}

func TestTasksSaveEveryChange(t *testing.T) {
	store := &memoryStore{saved: []domain.Task{{ID: 7, Title: "Kept"}}}
	tasks, err := NewTasks(store)
	if err != nil {
		t.Fatal(err)
	}
	if err := tasks.Add("New"); err != nil {
		t.Fatal(err)
	}
	if len(store.saved) != 2 || store.saved[1].ID != 8 {
		t.Fatalf("saved %+v", store.saved)
	}
	if err := tasks.Toggle(7); err != nil {
		t.Fatal(err)
	}
	if !store.saved[0].Done || tasks.Remaining() != 1 {
		t.Errorf("saved %+v, remaining %d", store.saved, tasks.Remaining())
	}
	if err := tasks.Remove(7); err != nil {
		t.Fatal(err)
	}
	if len(tasks.List()) != 1 || tasks.List()[0].Title != "New" {
		t.Errorf("after removing: %+v", tasks.List())
	}
}

func TestTasksKeepTheirStateWhenSavingFails(t *testing.T) {
	store := &memoryStore{}
	tasks, err := NewTasks(store)
	if err != nil {
		t.Fatal(err)
	}
	store.fail = errors.New("disk full")
	if err := tasks.Add("Lost"); !errors.Is(err, store.fail) {
		t.Errorf("adding gives %v", err)
	}
	if len(tasks.List()) != 0 {
		t.Errorf("the list changed without saving: %+v", tasks.List())
	}
	if err := tasks.Toggle(99); err == nil {
		t.Error("toggling a task that is not there gives no error")
	}
}
