package main

import (
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func dashboardWindow(title string, view func(*ui.Context)) mygo.WindowOptions {
	return mygo.WindowOptions{
		Title:         title,
		Width:         1440,
		Height:        900,
		MinWidth:      1024,
		MinHeight:     680,
		StateKey:      "main",
		TitleBarStyle: mygo.TitleBarHidden,
		Content:       ui.View(view),
	}
}
