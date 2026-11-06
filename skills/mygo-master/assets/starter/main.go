package main

import (
	"log"
	"path/filepath"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"starter/internal/adapter"
	"starter/internal/usecase"
	"starter/internal/view"
)

func main() {
	dir, err := mygo.App.Path(mygo.PathUserData)
	if err != nil {
		log.Fatal(err)
	}
	tasks, loadErr := usecase.NewTasks(adapter.NewFileStore(filepath.Join(dir, "tasks.json")))
	mygo.App.WhenReady(func() {
		screen := view.LoadFailure(loadErr)
		if loadErr == nil {
			screen = view.New(tasks, mygo.Theme.IsDark(), setDark).View
		}
		mygo.NewWindow(mygo.WindowOptions{
			Title:     "Tasks",
			Width:     560,
			Height:    640,
			MinWidth:  420,
			MinHeight: 420,
			StateKey:  "main",
			Content:   ui.View(screen),
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}

func setDark(dark bool) {
	if dark {
		mygo.Theme.SetSource(mygo.ThemeDark)
		return
	}
	mygo.Theme.SetSource(mygo.ThemeLight)
}
