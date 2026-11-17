package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

var frameNav = []navEntry{{homeIcon, "Home"}, {chartIcon, "Reports"}, {usersIcon, "People"}}

func frameAccount(open *bool, chosen *string) accountView {
	return accountView{
		label: "Account menu", initials: "AB", name: "Ada B", role: "Owner", chevron: chevronIcon, open: open,
		entries: []menuEntry{
			{label: "Profile", choose: func() { *chosen = "Profile" }},
			{label: "Sign out", separated: true, choose: func() { *chosen = "Sign out" }},
		},
	}
}

func sidebarTester(selected *int, open *bool, chosen *string) *ui.Tester {
	return kitScreen(520, 640, func(c *ui.Context, p *palette) {
		sidebar(c, p, "Acme", frameNav, *selected, func(i int) { *selected = i }, frameAccount(open, chosen)).Height(560)
	})
}

func kitPillTop(tt *ui.Tester) int {
	for y := 0; y < tt.Image().Bounds().Dy(); y++ {
		if kitSame(kitPixel(tt, 40, float32(y)), lightPalette.navActive) {
			return y
		}
	}
	return -1
}

func kitCrop(tt *ui.Tester, r ui.Rect) []byte {
	img := tt.Image()
	var out []byte
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			c := img.RGBAAt(x, y)
			out = append(out, c.R, c.G, c.B)
		}
	}
	return out
}

func TestSidebarShowsItsPartsAndTheSelectedPillSlidesThenRests(t *testing.T) {
	selected, open, chosen := 0, false, ""
	tt := sidebarTester(&selected, &open, &chosen)
	for _, want := range []string{"Acme", "Home", "Reports", "People", "Ada B", "AB", "Owner"} {
		if !tt.HasText(want) {
			t.Errorf("missing %q in %q", want, tt.Texts())
		}
	}
	if small := smallTargets(tt, 28, "Home", "Reports", "People", "Account menu"); len(small) > 0 {
		t.Errorf("targets under 28 DIP: %v", small)
	}
	if err := tt.Click("People"); err != nil {
		t.Fatal(err)
	}
	start := kitPillTop(tt)
	tt.Advance(60 * time.Millisecond)
	moving := kitPillTop(tt)
	tt.Advance(time.Second)
	rest := kitPillTop(tt)
	if selected != 2 || start >= moving || moving >= rest || tt.FramePending() {
		t.Errorf("selected %d, the pill goes %d, %d, %d, pending %v", selected, start, moving, rest, tt.FramePending())
	}
}

func TestSidebarItemsShowHoverAndRest(t *testing.T) {
	selected, open, chosen := 0, false, ""
	tt := sidebarTester(&selected, &open, &chosen)
	r, _ := tt.Find("Reports")
	x, y := r.X+1, r.Y+r.H/2
	before := kitPixel(tt, x, y)
	tt.Move(r.X+r.W/2, y)
	tt.Advance(time.Second)
	if got := kitPixel(tt, x, y); !kitSame(got, lightPalette.hover) || tt.FramePending() {
		t.Errorf("an item under the pointer shows %v, pending %v", got, tt.FramePending())
	}
	tt.Move(700, 500)
	tt.Advance(time.Second)
	if got := kitPixel(tt, x, y); !kitSame(got, before) {
		t.Errorf("the item does not go back: %v, was %v", got, before)
	}
	h, _ := tt.Find("Home")
	tt.Move(h.X+h.W/2, h.Y+h.H/2)
	tt.Advance(time.Second)
	if kitSame(kitPixel(tt, h.X+1, h.Y+h.H/2), lightPalette.hover) {
		t.Error("the selected item shows a hover over its pill")
	}
}

func TestAccountMenuOpensChoosesAndCloses(t *testing.T) {
	selected, open, chosen := 0, false, ""
	tt := sidebarTester(&selected, &open, &chosen)
	button, _ := tt.Find("Account menu")
	tt.Move(button.X+button.W/2, button.Y+button.H/2)
	tt.Advance(time.Second)
	arrow := ui.Rect{X: button.X + button.W - 24, Y: button.Y + 14, W: 18, H: 18}
	closed := kitCrop(tt, arrow)
	if err := tt.Click("Account menu"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if bytes.Equal(closed, kitCrop(tt, arrow)) {
		t.Error("the chevron does not turn as the menu opens")
	}
	if !open || !tt.HasText("Profile") || !tt.HasText("Sign out") {
		t.Fatalf("the menu does not open: open %v, texts %q", open, tt.Texts())
	}
	if small := smallTargets(tt, 28, "Sign out"); len(small) > 0 {
		t.Errorf("menu rows under 28 DIP: %v", small)
	}
	if err := tt.Click("Sign out"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if open || chosen != "Sign out" || tt.HasText("Profile") {
		t.Errorf("after choosing: open %v, chosen %q", open, chosen)
	}
}

func bellTester(open, seen *bool, rings *uint64, chosen *string) *ui.Tester {
	return kitScreen(520, 300, func(c *ui.Context, p *palette) {
		topbar(c, p, func() {
			ui.Spacer(c)
			notificationBell(c, p, bellView{
				label: "Notifications", icon: bellIcon, open: open, seen: seen, rings: rings,
				entries: []menuEntry{{label: "New order", choose: func() { *chosen = "New order" }}},
			})
		})
	})
}

func TestBellRingsOncePerOpenAndItsDotGoes(t *testing.T) {
	open, seen, rings, chosen := false, false, uint64(0), ""
	tt := bellTester(&open, &seen, &rings, &chosen)
	bell, ok := tt.Find("Notifications")
	if !ok {
		t.Fatal("no bell")
	}
	if small := smallTargets(tt, 28, "Notifications"); len(small) > 0 {
		t.Errorf("the bell is under 28 DIP: %v", small)
	}
	dot := func() ui.Color { return kitPixel(tt, bell.X+bell.W-5, bell.Y+7) }
	if !kitSame(dot(), lightPalette.accent) || rings != 0 {
		t.Fatalf("at rest the dot shows %v and the bell rang %d times", dot(), rings)
	}
	if err := tt.Click("Notifications"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(170 * time.Millisecond)
	swung := kitCrop(tt, bell)
	tt.Advance(time.Second)
	if bytes.Equal(swung, kitCrop(tt, bell)) {
		t.Error("the bell does not swing as the menu opens")
	}
	if !open || !seen || rings != 1 || kitSame(dot(), lightPalette.accent) {
		t.Errorf("open %v, seen %v, rings %d, dot %v", open, seen, rings, dot())
	}
	tt.Key(0, ui.KeyEscape)
	tt.Advance(time.Second)
	if err := tt.Click("Notifications"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(170 * time.Millisecond)
	again := kitCrop(tt, bell)
	tt.Advance(time.Second)
	if bytes.Equal(again, kitCrop(tt, bell)) {
		t.Error("the bell does not swing the second time")
	}
	if rings != 2 {
		t.Errorf("the bell rang %d times after two openings", rings)
	}
	if err := tt.Click("New order"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if chosen != "New order" || open || tt.FramePending() {
		t.Errorf("chosen %q, open %v, pending %v", chosen, open, tt.FramePending())
	}
}

func TestAppFrameKeepsItsPlacesScrollsAndEntersEachPage(t *testing.T) {
	pageNo := 0
	tt := ui.NewTester(func(c *ui.Context) {
		c.SetTheme(restThemes[0])
		appFrame(c, &lightPalette, pageNo,
			func() { ui.Column(c).Width(120).Children(func() { ui.Text(c, "RAIL") }) },
			func() { topbar(c, &lightPalette, func() { ui.Text(c, "BAR") }) },
			func() {
				ui.Text(c, fmt.Sprintf("Page %d", pageNo))
				ui.Box(c).Height(1200)
			})
	}, 800, 500)
	tt.Advance(time.Second)
	for _, want := range []string{"RAIL", "BAR", "Page 0"} {
		if !tt.HasText(want) {
			t.Errorf("missing %q in %q", want, tt.Texts())
		}
	}
	bar := func(y float32) ui.Color { return kitPixel(tt, 600, y) }
	if kitSame(bar(68), bar(60)) || kitSame(bar(68), bar(80)) {
		t.Errorf("the top bar does not end at 69: %v, %v, %v", bar(60), bar(68), bar(80))
	}
	first, _ := tt.Find("Page 0")
	tt.Wheel(500, 300, 0, 100)
	tt.Advance(time.Second)
	if scrolled, _ := tt.Find("Page 0"); scrolled.Y >= first.Y {
		t.Fatalf("the page does not scroll: %v", scrolled.Y)
	}
	pageNo = 1
	tt.Frame()
	tt.Advance(40 * time.Millisecond)
	entering, ok := tt.Find("Page 1")
	if !ok {
		t.Fatalf("no second page in %q", tt.Texts())
	}
	tt.Advance(time.Second)
	rest, _ := tt.Find("Page 1")
	if entering.Y <= rest.Y || rest.Y > first.Y+2 || rest.Y < first.Y-2 || tt.FramePending() {
		t.Errorf("the page enters at %v, rests at %v, the top is %v, pending %v", entering.Y, rest.Y, first.Y, tt.FramePending())
	}
}

func TestAScreenBuiltFromPartsOnly(t *testing.T) {
	people := []struct{ name, city string }{{"Ada Lovelace", "London"}, {"Grace Hopper", "New York"}, {"Alan Turing", "Wilmslow"}}
	months := []string{"Jan", "Feb", "Mar"}
	revenue := []float64{40, 55, 48}
	pageNo, month, query, dark := 0, 2, "", false
	accountOpen, alertsOpen, alertsSeen, rings, chosen := false, false, false, uint64(0), ""
	tt := ui.NewTester(func(c *ui.Context) {
		c.SetTheme(restThemes[0])
		p := &lightPalette
		appFrame(c, p, pageNo,
			func() {
				sidebar(c, p, "Acme", frameNav, pageNo, func(i int) { pageNo = i }, frameAccount(&accountOpen, &chosen))
			},
			func() {
				topbar(c, p, func() {
					searchField(c, p, &query, searchIcon, "Search people")
					ui.Spacer(c)
					notificationBell(c, p, bellView{label: "Notifications", icon: bellIcon, open: &alertsOpen, seen: &alertsSeen, rings: &rings})
					themeSwitch(c, p, &dark, "Dark mode", sunIcon, moonIcon, func() {}).Margin(0, 2, 0, 29)
				})
			},
			func() {
				ui.Grid(c).Columns(2).GapX(17).GapY(17).Children(func() {
					for i, label := range []string{"Revenue", "Orders"} {
						metricCard(c, p, metricView{
							label: label, icon: databaseIcon, value: dollars(revenue[month] * float64(i+1)), against: "vs. last month",
							delta: 4.5, up: arrowUpIcon, down: arrowDownIcon, trend: revenue,
							describe: func(i int) string { return months[i] },
						})
					}
				})
				ui.Row(c).Gap(18).Children(func() {
					card(c, p).Height(289).Grow(1).Basis(0).MinWidth(0).Padding(15, 19, 17, 22).Gap(26).Children(func() {
						ui.Row(c).Justify(ui.SpaceBetween).Children(func() {
							cardTitle(c, "Revenue by month")
							segmented(c, p, "month", "Month", &month, months, 52)
						})
						areaChart(c, p, areaChartView{
							label: "Monthly revenue", names: months, values: revenue, ceiling: 60, step: 20, axisLabel: axisDollars,
							shaded: month, describe: func(i int) string { return months[i] + ": " + dollars(revenue[i]) },
						})
					})
					card(c, p).Height(289).Width(300).Padding(16, 19, 0, 20).Children(func() {
						cardTitle(c, "Split")
						donut(c, p, donutView{
							label: "Split by kind", names: []string{"A", "B", "C", "D"}, shares: []float64{40, 30, 20, 10},
							center: "100", caption: "Total", describe: func(i int) string { return fmt.Sprint(i) },
						})
					})
				})
				card(c, p).Padding(15, 16, 20, 22).Gap(13).Children(func() {
					ui.Row(c).Justify(ui.SpaceBetween).Children(func() {
						cardTitle(c, "People")
						if moreButton(c, p, "View All", arrowIcon).Clicked() {
							query = ""
						}
					})
					var shown []int
					for i, person := range people {
						if strings.Contains(strings.ToLower(person.name), strings.ToLower(query)) {
							shown = append(shown, i)
						}
					}
					table(c, p, "People", "Nobody matches.", []column{{"Name", 1}, {"City", 1}}, len(shown), func(row, col int) {
						if col == 0 {
							textLink(c, p, people[shown[row]].name)
							return
						}
						cellText(c, people[shown[row]].city)
					})
				})
			})
	}, 1440, 900)
	tt.Advance(time.Second)
	for _, want := range []string{"Acme", "Ada B", "Revenue", "Orders", "$48", "$96", "Revenue by month", "Mar", "Split", "100", "Ada Lovelace", "Wilmslow"} {
		if !tt.HasText(want) {
			t.Errorf("missing %q in %q", want, tt.Texts())
		}
	}
	if small := smallTargets(tt, 28, "Reports", "Account menu", "Notifications", "Dark mode", "Jan", "Feb", "Mar", "View All", "Ada Lovelace"); len(small) > 0 {
		t.Errorf("targets under 28 DIP: %v", small)
	}
	if tt.FramePending() {
		t.Error("the screen keeps asking for frames at rest")
	}
	if err := tt.Click("Jan"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if month != 0 || !tt.HasText("$40") || !tt.HasText("$80") {
		t.Errorf("month %d, texts %q", month, tt.Texts())
	}
	if err := tt.Click("Search"); err != nil {
		t.Fatal(err)
	}
	tt.Type("grace")
	tt.Advance(time.Second)
	if !tt.HasText("Grace Hopper") || tt.HasText("Ada Lovelace") {
		t.Errorf("the table does not follow the search: %q", tt.Texts())
	}
	if err := tt.Click("View All"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if query != "" || !tt.HasText("Ada Lovelace") {
		t.Errorf("View All keeps the search %q", query)
	}
	if err := tt.Click("Reports"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if pageNo != 1 || tt.FramePending() {
		t.Errorf("page %d, pending %v", pageNo, tt.FramePending())
	}
}

func TestSidebarAndTopbarKeepClearOfTheWindowControls(t *testing.T) {
	selected, open, chosen := 0, false, ""
	bellOpen, seen, rings := false, false, uint64(0)
	tt := ui.NewTester(func(c *ui.Context) {
		c.SetTheme(restThemes[0])
		p := &lightPalette
		appFrame(c, p, 0,
			func() {
				sidebar(c, p, "Acme", frameNav, selected, func(i int) { selected = i }, frameAccount(&open, &chosen))
			},
			func() {
				topbar(c, p, func() {
					ui.Spacer(c)
					notificationBell(c, p, bellView{label: "Notifications", icon: bellIcon, open: &bellOpen, seen: &seen, rings: &rings})
				})
			},
			func() {})
	}, 900, 500)
	tt.Advance(time.Second)
	title, _ := tt.Find("Acme")
	bell, _ := tt.Find("Notifications")
	tt.SetTitleBar(ui.TitleBar{Height: 28, Right: 140})
	tt.Advance(time.Second)
	lowered, _ := tt.Find("Acme")
	moved, _ := tt.Find("Notifications")
	if lowered.Y <= title.Y {
		t.Errorf("the title stays at %v under the window buttons, was %v", lowered.Y, title.Y)
	}
	if bell.X-moved.X != 140 {
		t.Errorf("the bell moves %v away from the window buttons, want 140", bell.X-moved.X)
	}
}
