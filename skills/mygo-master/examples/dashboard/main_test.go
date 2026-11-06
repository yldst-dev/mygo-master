package main

import (
	"reflect"
	"testing"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func testApp() *app {
	return newApp(func(mygo.ThemeSource) {})
}

func TestDashboardShowsReferenceFigures(t *testing.T) {
	a := testApp()
	tt := ui.NewTester(a.view, 1440, 900)
	for _, want := range []string{
		"$128,430", "+12.4%", "24,819", "+8.7%", "3,972", "+5.6%", "3.86%", "−0.9%",
		"$128.4K", "42%", "28%", "18%", "12%",
		"#ORD-10051", "Ethan Miller", "$249.00", "Jun 23, 2025", "#ORD-10047", "Jun 21, 2025",
	} {
		if !tt.HasText(want) {
			t.Errorf("missing %q in %q", want, tt.Texts())
		}
	}
	if tt.HasText("#ORD-10046") {
		t.Error("the dashboard lists more than 5 recent orders")
	}
}

func TestMonthDrivesFigures(t *testing.T) {
	a := testApp()
	tt := ui.NewTester(a.view, 1440, 900)
	if err := tt.Click("Mar"); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"$76,300", "$76.3K", "−5.3%", "40%"} {
		if !tt.HasText(want) {
			t.Errorf("missing %q in %q", want, tt.Texts())
		}
	}
	if err := tt.Click("Jan"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("+10.6%") {
		t.Errorf("January does not compare with last December: %q", tt.Texts())
	}
}

func TestSearchFiltersOrders(t *testing.T) {
	a := testApp()
	tt := ui.NewTester(a.view, 1440, 900)
	if err := tt.Click("Search"); err != nil {
		t.Fatal(err)
	}
	tt.Type("noah")
	if !tt.HasText("#ORD-10049") || tt.HasText("#ORD-10051") {
		t.Errorf("search for noah shows %q", tt.Texts())
	}
	if err := tt.Click("Search"); err != nil {
		t.Fatal(err)
	}
	tt.Type("zzz")
	if !tt.HasText("Nothing matches your search.") {
		t.Errorf("no empty state in %q", tt.Texts())
	}
}

func TestNavigation(t *testing.T) {
	a := testApp()
	tt := ui.NewTester(a.view, 1440, 900)
	if err := tt.Click("View All"); err != nil {
		t.Fatal(err)
	}
	if a.page != pageOrders || !tt.HasText("#ORD-10040") || !tt.HasText("12 of 12") {
		t.Errorf("page %d, texts %q", a.page, tt.Texts())
	}
	if err := tt.Click("Customers"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Isabella Clark") || !tt.HasText("$461.00") {
		t.Errorf("customers page shows %q", tt.Texts())
	}
	if err := tt.Click("Analytics"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Monthly Breakdown") || !tt.HasText("$114,260") {
		t.Errorf("analytics page shows %q", tt.Texts())
	}
	if err := tt.Click("Settings"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Appearance") {
		t.Errorf("settings page shows %q", tt.Texts())
	}
}

func TestDarkModeToggle(t *testing.T) {
	var applied mygo.ThemeSource
	a := newApp(func(source mygo.ThemeSource) { applied = source })
	tt := ui.NewTester(a.view, 1440, 900)
	if err := tt.Click("Dark mode"); err != nil {
		t.Fatal(err)
	}
	if !a.dark || applied != mygo.ThemeDark {
		t.Errorf("after the toggle: dark %v, appearance %q", a.dark, applied)
	}
}

func TestCustomerSpendExcludesCancelledOrders(t *testing.T) {
	for _, s := range summarizeCustomers(orders) {
		if s.name == "Olivia Carter" && (s.orders != 2 || s.spent != 129) {
			t.Errorf("Olivia Carter: %d orders, %.2f spent", s.orders, s.spent)
		}
	}
}

func pillTop(tt *ui.Tester) int {
	img := tt.Image()
	for y := 0; y < img.Bounds().Dy(); y++ {
		c := img.RGBAAt(20, y)
		if c.R == lightPalette.navActive.R && c.G == lightPalette.navActive.G && c.B == lightPalette.navActive.B {
			return y
		}
	}
	return -1
}

func TestNavPillSlidesThenRests(t *testing.T) {
	a := testApp()
	tt := ui.NewTester(a.view, 1440, 900)
	if err := tt.Click("Customers"); err != nil {
		t.Fatal(err)
	}
	start := pillTop(tt)
	tt.Advance(60 * time.Millisecond)
	moving := pillTop(tt)
	tt.Advance(time.Second)
	rest := pillTop(tt)
	if start >= moving || moving >= rest {
		t.Errorf("the pill goes %d, %d, %d", start, moving, rest)
	}
	if tt.FramePending() {
		t.Error("the view keeps asking for frames at rest")
	}
}

func TestDropdownKeys(t *testing.T) {
	a := testApp()
	tt := ui.NewTester(a.view, 1440, 900)
	first, second := "#ORD-10051  Completed", "#ORD-10050  Processing"
	if err := tt.Click("Notifications"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText(first) || !a.alertsSeen {
		t.Errorf("the menu does not show: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEscape)
	tt.Advance(time.Second)
	if a.alertsOpen || tt.HasText(first) || !tt.Focused("Notifications") {
		t.Errorf("after Escape: open %v, texts %q", a.alertsOpen, tt.Texts())
	}
	tt.Key(0, ui.KeyDown)
	if !a.alertsOpen || !tt.Focused(first) {
		t.Errorf("Down does not open the menu on its first item: open %v", a.alertsOpen)
	}
	tt.Key(0, ui.KeyDown)
	if !tt.Focused(second) {
		t.Error("Down does not move to the second item")
	}
	tt.Key(0, ui.KeyTab)
	tt.Advance(time.Second)
	if a.alertsOpen {
		t.Error("Tab does not close the menu")
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	tt.Advance(time.Second)
	if a.alertsOpen || a.page != pageOrders || a.query != "#ORD-10050" {
		t.Errorf("after Enter: open %v, page %d, query %q", a.alertsOpen, a.page, a.query)
	}
}

func highlighted(tt *ui.Tester, label string) bool {
	tt.Advance(time.Second)
	r, ok := tt.Find(label)
	if !ok {
		return false
	}
	c := tt.Image().RGBAAt(int(r.X)+3, int(r.Y+r.H/2))
	return c.R == lightPalette.accent.R && c.G == lightPalette.accent.G && c.B == lightPalette.accent.B
}

func TestDropdownHighlightsOneRow(t *testing.T) {
	a := testApp()
	tt := ui.NewTester(a.view, 1440, 900)
	first, second := "#ORD-10051  Completed", "#ORD-10050  Processing"
	if err := tt.Click("Notifications"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if highlighted(tt, first) || highlighted(tt, second) {
		t.Error("a menu the pointer opened starts with a highlighted row")
	}
	r, _ := tt.Find(second)
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	if highlighted(tt, first) || !highlighted(tt, second) {
		t.Error("the pointer does not highlight the row it is on, alone")
	}
	tt.Key(0, ui.KeyUp)
	if !highlighted(tt, first) || highlighted(tt, second) {
		t.Error("the keys do not take the highlight from the row under the pointer")
	}
}

func TestDropdownClosesOnPressOutside(t *testing.T) {
	a := testApp()
	tt := ui.NewTester(a.view, 1440, 900)
	if err := tt.Click("Account"); err != nil {
		t.Fatal(err)
	}
	if !a.accountOpen {
		t.Fatal("the account menu does not open")
	}
	if err := tt.Click("Notifications"); err != nil {
		t.Fatal(err)
	}
	if a.accountOpen || !a.alertsOpen {
		t.Errorf("after pressing the bell: account %v, alerts %v", a.accountOpen, a.alertsOpen)
	}
	if err := tt.Click("Notifications"); err != nil {
		t.Fatal(err)
	}
	if a.alertsOpen {
		t.Error("pressing the bell again does not close its menu")
	}
}

func same(a, b ui.Color) bool {
	return a.R == b.R && a.G == b.G && a.B == b.B
}

func pixel(tt *ui.Tester, x, y float32) ui.Color {
	c := tt.Image().RGBAAt(int(x), int(y))
	return ui.RGB(c.R, c.G, c.B)
}

func TestShadeUpToStaysInTheData(t *testing.T) {
	xs, ys := []float32{0, 1, 2}, []float32{0, 10, 20}
	for _, c := range []struct {
		at     float32
		points int
		lastX  float32
		lastY  float32
	}{{-0.001, 1, 0, 0}, {0, 1, 0, 0}, {1, 2, 1, 10}, {1.5, 3, 1.5, 15}, {2, 3, 2, 20}, {2.001, 3, 2, 20}} {
		sx, sy := shadeUpTo(xs, ys, c.at)
		if len(sx) != c.points || sx[len(sx)-1] != c.lastX || sy[len(sy)-1] != c.lastY {
			t.Errorf("at %v: %v %v", c.at, sx, sy)
		}
	}
	if xs[2] != 2 || ys[2] != 20 {
		t.Errorf("the shade wrote over the data: %v %v", xs, ys)
	}
}

func TestPaletteFadeCoversEveryColor(t *testing.T) {
	var count func(reflect.Type) int
	count = func(typ reflect.Type) int {
		switch {
		case typ == reflect.TypeOf(ui.Color{}):
			return 1
		case typ.Kind() == reflect.Array:
			return typ.Len() * count(typ.Elem())
		case typ.Kind() == reflect.Struct:
			total := 0
			for i := range typ.NumField() {
				total += count(typ.Field(i).Type)
			}
			return total
		}
		return 0
	}
	if want, got := count(reflect.TypeOf(palette{})), len(lightPalette.channels()); want != got {
		t.Errorf("the palette has %d colors, the fade moves %d", want, got)
	}
	if fadedPalette(0) != &lightPalette || fadedPalette(1) != &darkPalette {
		t.Error("the ends of the fade are not the palettes themselves")
	}
}

func TestThemeFadeMovesTogether(t *testing.T) {
	a := testApp()
	tt := ui.NewTester(a.view, 1440, 900)
	field, _ := tt.Find("Search")
	if err := tt.Click("Dark mode"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(60 * time.Millisecond)
	if a.shown == &lightPalette || a.shown == &darkPalette {
		t.Fatal("the theme does not fade")
	}
	if got := pixel(tt, field.X+field.W-20, field.Y+field.H/2); !same(got, a.shown.field) {
		t.Errorf("the search field shows %v while the palette is at %v", got, a.shown.field)
	}
	tt.Advance(time.Second)
	if a.shown != &darkPalette || tt.FramePending() {
		t.Errorf("after the fade: dark %v, pending %v", a.shown == &darkPalette, tt.FramePending())
	}
}

func TestMonthPillSlidesUnderTheLabels(t *testing.T) {
	a := testApp()
	tt := ui.NewTester(a.view, 1440, 900)
	r, _ := tt.Find("Mar")
	y := r.Y - 4
	blueAt := func(label string) bool {
		l, _ := tt.Find(label)
		return same(pixel(tt, l.X+l.W/2, y), lightPalette.accent)
	}
	if !blueAt("Jun") || blueAt("Mar") {
		t.Fatal("the pill does not start on June")
	}
	if err := tt.Click("Mar"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(60 * time.Millisecond)
	if blueAt("Jun") || blueAt("Mar") {
		t.Error("the pill jumps instead of sliding")
	}
	tt.Advance(time.Second)
	if blueAt("Jun") || !blueAt("Mar") || tt.FramePending() {
		t.Errorf("after the slide: pending %v", tt.FramePending())
	}
}

func TestButtonsShowHoverAndRest(t *testing.T) {
	for _, label := range []string{"Notifications", "Account", "Account menu"} {
		a := testApp()
		tt := ui.NewTester(a.view, 1440, 900)
		tt.Advance(time.Second)
		r, ok := tt.Find(label)
		if !ok {
			t.Fatalf("no %q", label)
		}
		x, y := r.X+1, r.Y+r.H/2
		before := pixel(tt, x, y)
		tt.Move(r.X+r.W/2, r.Y+r.H/2)
		tt.Advance(time.Second)
		if got := pixel(tt, x, y); !same(got, lightPalette.hover) {
			t.Errorf("%s under the pointer shows %v", label, got)
		}
		if tt.FramePending() {
			t.Errorf("%s keeps asking for frames under the pointer", label)
		}
		tt.Move(700, 500)
		tt.Advance(time.Second)
		if got := pixel(tt, x, y); !same(got, before) {
			t.Errorf("%s does not go back: %v, was %v", label, got, before)
		}
	}
}

func TestViewAllClearsTheSearch(t *testing.T) {
	a := testApp()
	a.query = "noah"
	tt := ui.NewTester(a.view, 1440, 900)
	if err := tt.Click("View All"); err != nil {
		t.Fatal(err)
	}
	if a.page != pageOrders || a.query != "" || !tt.HasText("12 of 12") {
		t.Errorf("page %d, query %q, texts %q", a.page, a.query, tt.Texts())
	}
}

func TestClearButtonEmptiesTheSearch(t *testing.T) {
	a := testApp()
	a.query = "noah"
	tt := ui.NewTester(a.view, 1440, 900)
	if err := tt.Click("Clear search"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if a.query != "" || !tt.Focused("Search") {
		t.Errorf("after the button: query %q, focused %v", a.query, tt.Focused("Search"))
	}
	if _, shown := tt.Find("Clear search"); shown {
		t.Error("the button stays after the search is empty")
	}
	tt.Type("zzz")
	tt.Key(0, ui.KeyEscape)
	if a.query != "" || !tt.Focused("Search") {
		t.Errorf("after Escape: query %q, focused %v", a.query, tt.Focused("Search"))
	}
}

func TestPageSwitchStartsAtTheTop(t *testing.T) {
	a := testApp()
	tt := ui.NewTester(a.view, 1024, 680)
	first, _ := tt.Find("Total Revenue")
	tt.Wheel(700, 400, 0, 400)
	tt.Advance(time.Second)
	if scrolled, _ := tt.Find("Total Revenue"); scrolled.Y >= first.Y {
		t.Fatalf("the page does not scroll: %v", scrolled.Y)
	}
	if err := tt.Click("Analytics"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if r, _ := tt.Find("Revenue (Monthly)"); r.Y > first.Y+40 || r.Y < first.Y-40 {
		t.Errorf("the next page opens at %v, the top is near %v", r.Y, first.Y)
	}
}

func TestLessMotionJumpsToTheEnd(t *testing.T) {
	a := testApp()
	tt := ui.NewTester(a.view, 1440, 900)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	tt.Advance(time.Second)
	if err := tt.Click("Mar"); err != nil {
		t.Fatal(err)
	}
	r, _ := tt.Find("Mar")
	if !same(pixel(tt, r.X+r.W/2, r.Y-4), lightPalette.accent) {
		t.Error("the month pill slides while the desktop asks for less motion")
	}
	if err := tt.Click("Dark mode"); err != nil {
		t.Fatal(err)
	}
	if a.shown != &darkPalette {
		t.Error("the theme fades while the desktop asks for less motion")
	}
}

func hoverAt(t *testing.T, tt *ui.Tester, label string, fx, fy float32) {
	t.Helper()
	r, ok := tt.Find(label)
	if !ok {
		t.Fatalf("no %q", label)
	}
	tt.Move(r.X+r.W*fx, r.Y+r.H*fy)
}

func TestLineChartsAnswerThePointerAtOnce(t *testing.T) {
	for _, line := range []struct {
		label       string
		from, to    float32
		first, next string
	}{
		{"Monthly revenue", 0.4, 0.59, "Mar: $76.3K", "Apr: $97.7K"},
		{"Total Revenue trend", 0.01, 0.99, "Dec: $61,200", "Jun: $128,430"},
		{"Active Users trend", 0.2, 0.5, "Jan: 18,650", "Mar: 19,700"},
		{"Orders trend", 0.4, 0.8, "Feb: 3,350", "May: 3,761"},
		{"Conversion Rate trend", 0.5, 0.99, "Mar: 3.59%", "Jun: 3.86%"},
	} {
		a := testApp()
		tt := ui.NewTester(a.view, 1440, 900)
		tt.Advance(time.Second)
		hoverAt(t, tt, line.label, line.from, 0.1)
		if !tt.HasText(line.first) {
			t.Errorf("%s shows no %q as the pointer comes", line.label, line.first)
		}
		hoverAt(t, tt, line.label, line.to, 0.9)
		tt.Advance(time.Second)
		if !tt.HasText(line.next) || tt.HasText(line.first) {
			t.Errorf("%s does not follow the pointer to %q", line.label, line.next)
		}
		if tt.FramePending() {
			t.Errorf("%s keeps asking for frames under the pointer", line.label)
		}
		tt.Move(700, 40)
		tt.Advance(time.Second)
		if tt.HasText(line.next) {
			t.Errorf("%s keeps %q once the pointer leaves", line.label, line.next)
		}
		if tt.FramePending() {
			t.Errorf("%s keeps asking for frames after the pointer leaves", line.label)
		}
	}
}

func TestDonutAnswersThePointer(t *testing.T) {
	a := testApp()
	tt := ui.NewTester(a.view, 1440, 900)
	tt.Advance(time.Second)
	r, ok := tt.Find("Sales by category")
	if !ok {
		t.Fatal("no donut")
	}
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	rim := func() ui.Color { return pixel(tt, cx+donutRadius+2, cy) }
	if !same(rim(), lightPalette.card) {
		t.Errorf("the donut reaches past its radius at rest: %v", rim())
	}
	tt.Move(cx+70, cy)
	if !tt.HasText("Electronics: $53.9K") {
		t.Error("the slice under the pointer shows no amount")
	}
	tt.Advance(time.Second)
	if !same(rim(), lightPalette.slices[0]) {
		t.Errorf("the slice under the pointer does not grow: %v", rim())
	}
	tt.Move(cx, cy)
	tt.Advance(time.Second)
	if tt.HasText("Electronics: $53.9K") || !same(rim(), lightPalette.card) {
		t.Error("the slice stays lit over the hole")
	}
	hoverAt(t, tt, "Apparel", 0.5, 0.5)
	if !tt.HasText("Apparel: $36.0K") {
		t.Error("the legend does not light its slice")
	}
}

func TestChartTipLetsThePointerThrough(t *testing.T) {
	a := testApp()
	tt := ui.NewTester(a.view, 1440, 900)
	tt.Advance(time.Second)
	r, ok := tt.Find("Sales by category")
	if !ok {
		t.Fatal("no donut")
	}
	tt.Move(r.X+r.W/2+60, r.Y+r.H/2-40)
	for range 20 {
		tt.Advance(50 * time.Millisecond)
		if !tt.HasText("Electronics: $53.9K") {
			t.Fatal("the tip goes as it comes under the pointer")
		}
	}
	if tt.FramePending() {
		t.Error("the tip keeps asking for frames under the pointer")
	}
}

func TestKeysWalkTheChartsAndLeaveNoTip(t *testing.T) {
	a := testApp()
	tt := ui.NewTester(a.view, 1440, 900)
	tt.Advance(time.Second)
	for range 11 {
		tt.Key(0, ui.KeyTab)
	}
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("Dec: $61,200") {
		t.Errorf("the first trend does not take the keys: %q", tt.Texts())
	}
	if tt.HasText("Dec: 17,800") {
		t.Error("the arrow reaches a trend without the focus")
	}
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("Jan: $67,700") {
		t.Errorf("the arrow does not move along the trend: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Advance(time.Second)
	if tt.HasText("Jan: $67,700") {
		t.Error("the tip stays as the focus leaves the trend")
	}
	for range 40 {
		tt.Key(0, ui.KeyTab)
		if tt.HasText("Electronics: $53.9K") {
			return
		}
	}
	t.Error("the keys do not reach the legend")
}

func TestChartTipFadesInAndSlides(t *testing.T) {
	a := testApp()
	tt := ui.NewTester(a.view, 1440, 900)
	tt.Advance(time.Second)
	hoverAt(t, tt, "Monthly revenue", 0.4, 0.5)
	tt.Advance(40 * time.Millisecond)
	r, ok := tt.Find("Mar: $76.3K")
	if !ok {
		t.Fatal("no tip")
	}
	if same(pixel(tt, r.X-4, r.Y+r.H/2), lightPalette.ink) {
		t.Error("the tip shows in full at once")
	}
	tt.Advance(time.Second)
	r, _ = tt.Find("Mar: $76.3K")
	if !same(pixel(tt, r.X-4, r.Y+r.H/2), lightPalette.ink) {
		t.Errorf("the tip settles at %v", pixel(tt, r.X-4, r.Y+r.H/2))
	}
	from := r.X + r.W/2
	hoverAt(t, tt, "Monthly revenue", 0.59, 0.5)
	tt.Advance(40 * time.Millisecond)
	r, ok = tt.Find("Apr: $97.7K")
	if !ok {
		t.Fatal("the tip does not change")
	}
	on := r.X + r.W/2
	tt.Advance(time.Second)
	r, _ = tt.Find("Apr: $97.7K")
	to := r.X + r.W/2
	if from >= on || on >= to {
		t.Errorf("the tip does not slide: %v, %v, %v", from, on, to)
	}
}

func TestChartTipStaysInsideItsChart(t *testing.T) {
	for _, tip := range []struct {
		label  string
		fx, fy float32
		text   string
	}{
		{"Monthly revenue", 0.99, 0.5, "Jun: $128.4K"},
		{"Monthly revenue", 0.01, 0.5, "Jan: $67.7K"},
		{"Total Revenue trend", 0.99, 0.5, "Jun: $128,430"},
		{"Conversion Rate trend", 0.01, 0.5, "Dec: 3.42%"},
		{"Sales by category", 0.1, 0.45, "Home & Living: $23.1K"},
	} {
		a := testApp()
		tt := ui.NewTester(a.view, 1440, 900)
		tt.Advance(time.Second)
		hoverAt(t, tt, tip.label, tip.fx, tip.fy)
		tt.Advance(time.Second)
		box, _ := tt.Find(tip.label)
		r, ok := tt.Find(tip.text)
		if !ok {
			t.Errorf("no %q in %q", tip.text, tt.Texts())
			continue
		}
		left, right := r.X-tipPadding, r.X+r.W+tipPadding
		if right > box.X+box.W+1 {
			t.Errorf("%q ends at %v, past %s at %v", tip.text, right, tip.label, box.X+box.W)
		}
		if right-left <= box.W && left < box.X-1 {
			t.Errorf("%q starts at %v, before %s at %v", tip.text, left, tip.label, box.X)
		}
	}
}
