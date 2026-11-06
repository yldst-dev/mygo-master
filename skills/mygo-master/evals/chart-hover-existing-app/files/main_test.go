package main

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestRegionButtonsChooseTheSeries(t *testing.T) {
	a := newApp()
	tt := ui.NewTester(a.view, 640, 480)
	if err := tt.Click("South"); err != nil {
		t.Fatal(err)
	}
	if a.region != 1 {
		t.Errorf("region %d", a.region)
	}
}
