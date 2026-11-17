package main

import (
	"log"

	"github.com/egoist/mygo"
)

func main() {
	a := newApp(mygo.Theme.SetSource)
	mygo.Theme.SetSource(mygo.ThemeLight)
	mygo.App.WhenReady(func() {
		mygo.NewWindow(dashboardWindow("Dashboard", a.view))
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
