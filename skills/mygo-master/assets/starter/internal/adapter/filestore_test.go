package adapter

import (
	"os"
	"path/filepath"
	"testing"

	"starter/internal/domain"
)

func TestFileStoreKeepsWhatItSaves(t *testing.T) {
	store := NewFileStore(filepath.Join(t.TempDir(), "tasks.json"))
	if tasks, err := store.Load(); err != nil || tasks != nil {
		t.Fatalf("a missing file gives %v, %v", tasks, err)
	}
	want := []domain.Task{{ID: 1, Title: "한글 제목", Done: true}}
	if err := store.Save(want); err != nil {
		t.Fatal(err)
	}
	got, err := store.Load()
	if err != nil || len(got) != 1 || got[0] != want[0] {
		t.Errorf("got %+v, %v", got, err)
	}
}

func TestFileStoreLeavesOnlyItsFile(t *testing.T) {
	directory := t.TempDir()
	store := NewFileStore(filepath.Join(directory, "tasks.json"))
	for range 3 {
		if err := store.Save([]domain.Task{{ID: 1, Title: "Kept"}}); err != nil {
			t.Fatal(err)
		}
	}
	names, err := os.ReadDir(directory)
	if err != nil || len(names) != 1 || names[0].Name() != "tasks.json" {
		t.Errorf("the folder holds %v, %v", names, err)
	}
}

func TestFileStoreCleansUpWhenItCannotReplaceTheFile(t *testing.T) {
	directory := t.TempDir()
	target := filepath.Join(directory, "tasks.json")
	if err := os.Mkdir(target, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := NewFileStore(target).Save(nil); err == nil {
		t.Fatal("saving over a directory gives no error")
	}
	names, err := os.ReadDir(directory)
	if err != nil || len(names) != 1 {
		t.Errorf("the folder holds %v, %v", names, err)
	}
}

func TestFileStoreReportsABrokenFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.json")
	if err := os.WriteFile(path, []byte("{"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := NewFileStore(path).Load(); err == nil {
		t.Error("a broken file gives no error")
	}
}
