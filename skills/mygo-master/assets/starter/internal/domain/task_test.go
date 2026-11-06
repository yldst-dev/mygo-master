package domain

import (
	"errors"
	"testing"
)

func TestNewTaskTrimsItsTitleAndRefusesABlankOne(t *testing.T) {
	task, err := NewTask(3, "  Write the report ")
	if err != nil || task.ID != 3 || task.Title != "Write the report" || task.Done {
		t.Errorf("task %+v, err %v", task, err)
	}
	if _, err := NewTask(4, "   "); !errors.Is(err, ErrEmptyTitle) {
		t.Errorf("a blank title gives %v", err)
	}
}

func TestRemainingCountsWhatIsNotDone(t *testing.T) {
	if got := Remaining([]Task{{ID: 1}, {ID: 2, Done: true}, {ID: 3}}); got != 2 {
		t.Errorf("remaining %d", got)
	}
}
