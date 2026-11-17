package main

import (
	"testing"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func TestDashboardWindowHidesTheTitleBarAndKeepsItsSize(t *testing.T) {
	options := dashboardWindow("Acme", func(*ui.Context) {})
	if options.Title != "Acme" || options.StateKey == "" || options.Content == nil {
		t.Errorf("title %q, state key %q, content %v", options.Title, options.StateKey, options.Content)
	}
	if options.TitleBarStyle != mygo.TitleBarHidden {
		t.Errorf("the title bar is %q, not hidden: the sidebar and the top bar draw it", options.TitleBarStyle)
	}
	if options.MinWidth <= 0 || options.MinWidth >= options.Width || options.MinHeight <= 0 || options.MinHeight >= options.Height {
		t.Errorf("minimum %vx%v, default %vx%v", options.MinWidth, options.MinHeight, options.Width, options.Height)
	}
}
