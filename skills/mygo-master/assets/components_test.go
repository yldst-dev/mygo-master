package main

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func kitScreen(width, height int, build func(c *ui.Context, p *palette)) *ui.Tester {
	tt := ui.NewTester(func(c *ui.Context) {
		c.SetTheme(restThemes[0])
		ui.Column(c).Fill().Padding(24).Gap(16).Children(func() { build(c, &lightPalette) })
	}, width, height)
	tt.Advance(time.Second)
	return tt
}

func kitTester(build func(c *ui.Context, p *palette)) *ui.Tester {
	return kitScreen(640, 480, build)
}

func kitPixel(tt *ui.Tester, x, y float32) ui.Color {
	c := tt.Image().RGBAAt(int(x), int(y))
	return ui.RGB(c.R, c.G, c.B)
}

func kitSame(a, b ui.Color) bool {
	return a.R == b.R && a.G == b.G && a.B == b.B
}

func kitHover(t *testing.T, tt *ui.Tester, label string, fx, fy float32) {
	t.Helper()
	r, ok := tt.Find(label)
	if !ok {
		t.Fatalf("no %q", label)
	}
	tt.Move(r.X+r.W*fx, r.Y+r.H*fy)
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
			chevron(c, p, chevronIcon, open)
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
		searchField(c, p, &query, searchIcon, "Search here")
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
		themeSwitch(c, p, &dark, "Dark mode", sunIcon, moonIcon, func() { told++ })
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

func TestMetricCardShowsItsFigureAndDeltaAndFollowsThePointer(t *testing.T) {
	value, delta := "$128,430", 12.4
	trend := []float64{1, 3, 2, 5, 4}
	tt := kitScreen(400, 300, func(c *ui.Context, p *palette) {
		metricCard(c, p, metricView{
			label: "Revenue", icon: sunIcon, value: value, against: "vs. last month", delta: delta,
			up: arrowUpIcon, down: arrowDownIcon, trend: trend,
			describe: func(i int) string { return fmt.Sprintf("m%d", i) },
		}).Width(280)
	})
	for _, want := range []string{"Revenue", "$128,430", "+12.4%", "vs. last month"} {
		if !tt.HasText(want) {
			t.Errorf("missing %q in %q", want, tt.Texts())
		}
	}
	kitHover(t, tt, "Revenue trend", 0.99, 0.5)
	if !tt.HasText("m4") {
		t.Errorf("the trend shows no tip under the pointer: %q", tt.Texts())
	}
	tt.Move(380, 280)
	value, delta = "$127,000", -0.9
	tt.Frame()
	tt.Advance(time.Second)
	if !tt.HasText("$127,000") || !tt.HasText("−0.9%") || tt.HasText("+12.4%") || tt.FramePending() {
		t.Errorf("the card does not follow its figures: %q, pending %v", tt.Texts(), tt.FramePending())
	}
}

func TestMoreButtonAndTextLinkAnswerClicksAndStayLarge(t *testing.T) {
	more, opened := 0, ""
	tt := kitTester(func(c *ui.Context, p *palette) {
		if moreButton(c, p, "View All", arrowIcon).Clicked() {
			more++
		}
		table(c, p, "Orders", "Nothing here.", []column{{"ID", 1}}, 1, func(row, col int) {
			if textLink(c, p, "#A-1").Clicked() {
				opened = "#A-1"
			}
		})
	})
	if small := smallTargets(tt, 28, "View All", "#A-1"); len(small) > 0 {
		t.Errorf("targets under 28 DIP: %v", small)
	}
	if err := tt.Click("View All"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("#A-1"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if more != 1 || opened != "#A-1" || tt.FramePending() {
		t.Errorf("more %d, opened %q, pending %v", more, opened, tt.FramePending())
	}
}

func TestChevronTurnsOverWhenTheMenuIsOpen(t *testing.T) {
	open := false
	tt := kitTester(func(c *ui.Context, p *palette) {
		chevron(c, p, chevronIcon, open)
	})
	shape := ui.Rect{X: 24, Y: 24, W: 18, H: 18}
	closed := kitCrop(tt, shape)
	open = true
	tt.Frame()
	tt.Advance(time.Second)
	if bytes.Equal(closed, kitCrop(tt, shape)) || tt.FramePending() {
		t.Errorf("the chevron does not turn: pending %v", tt.FramePending())
	}
}
