package main

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/chart"
)

var kitArrow = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75"><path d="m6 9 6 6 6-6"/></svg>`))

func kitTester(build func(c *ui.Context, p *palette)) *ui.Tester {
	tt := ui.NewTester(func(c *ui.Context) {
		c.SetTheme(restThemes[0])
		ui.Column(c).Fill().Padding(24).Gap(16).Children(func() { build(c, &lightPalette) })
	}, 640, 480)
	tt.Advance(time.Second)
	return tt
}

func TestCardTitleAndStatusPillShowTheirText(t *testing.T) {
	tt := kitTester(func(c *ui.Context, p *palette) {
		card(c, p).Padding(16).Children(func() {
			cardTitle(c, "Orders")
			statusPill(c, p.good, "Completed")
		})
	})
	if !tt.HasText("Orders") || !tt.HasText("Completed") {
		t.Errorf("the card shows %q", tt.Texts())
	}
}

func TestDropdownOpensFromTheKeysChoosesAndGivesTheFocusBack(t *testing.T) {
	open, chosen := false, ""
	entries := []string{"Name", "Date", "Reset"}
	tt := kitTester(func(c *ui.Context, p *palette) {
		button := menuButton(c, p, "Sort", &open).Height(40).Padding(0, 12).Gap(8).Radius(10)
		button.Children(func() {
			ui.Text(c, "Sort")
			chevron(c, p, kitArrow, open)
		})
		list := make([]menuEntry, len(entries))
		for i, name := range entries {
			list[i] = menuEntry{label: name, separated: i == len(entries)-1, choose: func() { chosen = name }}
		}
		dropdown(c, p, button, &open, false, list)
	})
	if err := tt.Click("Sort"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if !open || !tt.HasText("Name") {
		t.Fatalf("the menu does not open: open %v, texts %q", open, tt.Texts())
	}
	if small := smallTargets(tt, 28, entries...); len(small) > 0 {
		t.Errorf("menu rows under 28 DIP: %v", small)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Advance(time.Second)
	if open || tt.HasText("Name") || !tt.Focused("Sort") {
		t.Errorf("after Escape: open %v, focus on Sort %v", open, tt.Focused("Sort"))
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("Date") {
		t.Error("Down does not walk the rows")
	}
	tt.Key(0, ui.KeyEnter)
	tt.Advance(time.Second)
	if open || chosen != "Date" {
		t.Errorf("after Enter: open %v, chosen %q", open, chosen)
	}
}

func TestSegmentedChoosesAndRests(t *testing.T) {
	selected := 0
	names := []string{"Day", "Week", "Month"}
	tt := kitTester(func(c *ui.Context, p *palette) {
		segmented(c, p, "range", "Range", &selected, names, 60)
	})
	if small := smallTargets(tt, 28, names...); len(small) > 0 {
		t.Errorf("segments under 28 DIP: %v", small)
	}
	if err := tt.Click("Month"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if selected != 2 || tt.FramePending() {
		t.Errorf("selected %d, pending %v", selected, tt.FramePending())
	}
}

func TestTableShowsItsRowsAndTheEmptyMessage(t *testing.T) {
	cells := [][]string{{"Ada", "36"}, {"Grace", "45"}}
	rows := len(cells)
	tt := kitTester(func(c *ui.Context, p *palette) {
		table(c, p, "People", "Nothing here.", []column{{"Name", 2}, {"Age", 1}}, rows, func(row, col int) {
			cellText(c, cells[row][col])
		})
	})
	if _, ok := tt.AccessByRole(ui.RoleTable, "People"); !ok {
		t.Error("the table has no accessible name")
	}
	for _, want := range []string{"Name", "Age", "Ada", "45"} {
		if !tt.HasText(want) {
			t.Errorf("missing %q in %q", want, tt.Texts())
		}
	}
	rows = 0
	tt.Frame()
	if !tt.HasText("Nothing here.") || tt.HasText("Ada") {
		t.Errorf("an empty table shows %q", tt.Texts())
	}
}

func TestSearchFieldTypesClearsAndKeepsTheClearButtonLarge(t *testing.T) {
	query := ""
	tt := kitTester(func(c *ui.Context, p *palette) {
		searchField(c, p, &query, kitArrow, "Search here")
	})
	if err := tt.Click("Search"); err != nil {
		t.Fatal(err)
	}
	tt.Type("abc")
	tt.Advance(time.Second)
	if query != "abc" {
		t.Fatalf("the query is %q", query)
	}
	if small := smallTargets(tt, 28, "Clear search"); len(small) > 0 {
		t.Errorf("the clear button is under 28 DIP: %v", small)
	}
	if err := tt.Click("Clear search"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if query != "" {
		t.Errorf("the query stays: %q", query)
	}
}

func TestThemeSwitchTurnsAndTellsTheCaller(t *testing.T) {
	dark, told := false, 0
	tt := kitTester(func(c *ui.Context, p *palette) {
		themeSwitch(c, p, &dark, "Dark mode", kitArrow, kitArrow, func() { told++ })
	})
	if small := smallTargets(tt, 28, "Dark mode"); len(small) > 0 {
		t.Errorf("the switch is under 28 DIP: %v", small)
	}
	if err := tt.Click("Dark mode"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if !dark || told != 1 || tt.FramePending() {
		t.Errorf("dark %v, told %d, pending %v", dark, told, tt.FramePending())
	}
}

func TestChartTipFollowsThePointerAndLeavesNoTip(t *testing.T) {
	labels := []string{"A: 1", "B: 2", "C: 3"}
	xs, ys := []float32{0.1, 0.5, 0.9}, []float32{0.2, 0.5, 0.8}
	tt := kitTester(func(c *ui.Context, p *palette) {
		chart.Plot(c, chart.Axes{}, func(area ui.Element) {
			selected := pickPoint(c, area, xs, func(i int) string { return labels[i] })
			if selected >= 0 {
				chartTip(c, p, area.Bounds().W, xs[selected], ys[selected], labels[selected])
			}
		}).Label("Sales").Height(200)
	})
	r, ok := tt.Find("Sales")
	if !ok {
		t.Fatal("no chart")
	}
	tt.Move(r.X+r.W*0.5, r.Y+r.H*0.5)
	if !tt.HasText("B: 2") {
		t.Errorf("no tip under the pointer: %q", tt.Texts())
	}
	tt.Move(r.X+r.W*0.9, r.Y+r.H*0.5)
	tt.Advance(time.Second)
	if !tt.HasText("C: 3") || tt.HasText("B: 2") || tt.FramePending() {
		t.Errorf("the tip does not follow: %q, pending %v", tt.Texts(), tt.FramePending())
	}
	tt.Move(r.X+r.W+200, r.Y-100)
	tt.Advance(time.Second)
	if tt.HasText("C: 3") || tt.FramePending() {
		t.Errorf("a tip stays after the pointer left: %q", tt.Texts())
	}
}

func TestChartTipLeavesWhenTheFocusLeaves(t *testing.T) {
	labels := []string{"A: 1", "B: 2", "C: 3"}
	xs, ys := []float32{0.1, 0.5, 0.9}, []float32{0.2, 0.5, 0.8}
	tt := kitTester(func(c *ui.Context, p *palette) {
		chart.Plot(c, chart.Axes{}, func(area ui.Element) {
			selected := pickPoint(c, area, xs, func(i int) string { return labels[i] })
			if selected >= 0 {
				chartTip(c, p, area.Bounds().W, xs[selected], ys[selected], labels[selected])
			}
		}).Label("Sales").Height(200)
		ui.ButtonBase(c).Label("Next").Height(30).Padding(0, 12).Children(func() { ui.Text(c, "Next") })
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("A: 1") {
		t.Fatalf("the keys do not walk the chart: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Advance(time.Second)
	if !tt.Focused("Next") || tt.HasText("A: 1") {
		t.Errorf("focus on Next %v, tip left over %v", tt.Focused("Next"), tt.HasText("A: 1"))
	}
}
