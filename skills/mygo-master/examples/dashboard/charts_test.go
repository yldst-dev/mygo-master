package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/chart"
)

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

var spikes = []float64{4, 9, 6, 12}

func spikeName(i int) string { return fmt.Sprintf("p%d", i) }

func sparkTester() *ui.Tester {
	return kitScreen(400, 200, func(c *ui.Context, p *palette) {
		sparkline(c, p, spikes, spikeName).Label("Trend").Size(160, 60)
		ui.ButtonBase(c).Label("Next").Height(30).Padding(0, 12).Children(func() { ui.Text(c, "Next") })
	})
}

func TestSparklineAnswersThePointerAtOnceAndLeavesNoTip(t *testing.T) {
	tt := sparkTester()
	kitHover(t, tt, "Trend", 0.1, 0.5)
	if !tt.HasText("p0") {
		t.Fatalf("no tip as the pointer comes: %q", tt.Texts())
	}
	kitHover(t, tt, "Trend", 0.9, 0.5)
	tt.Advance(time.Second)
	if !tt.HasText("p3") || tt.HasText("p0") || tt.FramePending() {
		t.Errorf("the tip does not follow: %q, pending %v", tt.Texts(), tt.FramePending())
	}
	tt.Move(380, 190)
	tt.Advance(time.Second)
	if tt.HasText("p3") || tt.FramePending() {
		t.Errorf("a tip stays after the pointer left: %q, pending %v", tt.Texts(), tt.FramePending())
	}
}

func TestSparklineKeysWalkItAndLeaveNoTip(t *testing.T) {
	tt := sparkTester()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("p0") {
		t.Fatalf("the first arrow shows no tip: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("p1") || tt.HasText("p0") {
		t.Errorf("the arrow does not move along the line: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Advance(time.Second)
	if !tt.Focused("Next") || tt.HasText("p1") {
		t.Errorf("focus on Next %v, tip left over %v", tt.Focused("Next"), tt.HasText("p1"))
	}
}

func TestSparklineTipStaysInsideItsLine(t *testing.T) {
	for _, edge := range []struct {
		fx   float32
		name string
	}{{0.01, "p0"}, {0.99, "p3"}} {
		tt := sparkTester()
		kitHover(t, tt, "Trend", edge.fx, 0.5)
		tt.Advance(time.Second)
		box, _ := tt.Find("Trend")
		r, ok := tt.Find(edge.name)
		if !ok {
			t.Fatalf("no %q in %q", edge.name, tt.Texts())
		}
		left, right := r.X-tipPadding, r.X+r.W+tipPadding
		if left < box.X-1 || right > box.X+box.W+1 {
			t.Errorf("%q spans %v to %v, outside %v to %v", edge.name, left, right, box.X, box.X+box.W)
		}
	}
}

var (
	areaNames  = []string{"A", "B", "C", "D", "E", "F"}
	areaValues = []float64{20, 35, 30, 45, 60, 80}
)

func areaTester(shaded *int) *ui.Tester {
	return kitScreen(640, 400, func(c *ui.Context, p *palette) {
		areaChart(c, p, areaChartView{
			label: "Sales", names: areaNames, values: areaValues, ceiling: 100, step: 25,
			axisLabel: func(v float64) string { return fmt.Sprintf("%.0f units", v) },
			shaded:    *shaded,
			describe:  func(i int) string { return fmt.Sprintf("%s: %.0f", areaNames[i], areaValues[i]) },
		}).Height(300)
	})
}

func TestAreaChartDrawsItsAxesAndTipsTheValueUnderThePointer(t *testing.T) {
	shaded := 5
	tt := areaTester(&shaded)
	for _, want := range []string{"0 units", "25 units", "100 units", "A", "F"} {
		if !tt.HasText(want) {
			t.Errorf("missing %q in %q", want, tt.Texts())
		}
	}
	kitHover(t, tt, "Sales", 0.2, 0.5)
	if !tt.HasText("B: 35") {
		t.Errorf("no tip under the pointer: %q", tt.Texts())
	}
	kitHover(t, tt, "Sales", 0.99, 0.5)
	tt.Advance(time.Second)
	if !tt.HasText("F: 80") || tt.HasText("B: 35") || tt.FramePending() {
		t.Errorf("the tip does not follow: %q, pending %v", tt.Texts(), tt.FramePending())
	}
	box, _ := tt.Find("Sales")
	r, _ := tt.Find("F: 80")
	if r.X+r.W+tipPadding > box.X+box.W+1 {
		t.Errorf("the tip ends at %v, past the chart at %v", r.X+r.W+tipPadding, box.X+box.W)
	}
	tt.Move(630, 390)
	tt.Advance(time.Second)
	if tt.HasText("F: 80") || tt.FramePending() {
		t.Errorf("a tip stays after the pointer left: %q", tt.Texts())
	}
}

func TestAreaChartShadesUpToTheChosenPointAndThenRests(t *testing.T) {
	shaded := 0
	tt := areaTester(&shaded)
	box, _ := tt.Find("Sales")
	x, y := box.X+box.W*(areaInset+areaSpan), box.Y+box.H*0.95
	empty := kitPixel(tt, x, y)
	shaded = 5
	tt.Frame()
	tt.Advance(60 * time.Millisecond)
	if !kitSame(kitPixel(tt, x, y), empty) {
		t.Error("the shade jumps to the last point instead of growing")
	}
	tt.Advance(time.Second)
	if kitSame(kitPixel(tt, x, y), empty) || tt.FramePending() {
		t.Errorf("the shade does not reach the last point, pending %v", tt.FramePending())
	}
}

func TestShadeUpToStaysInTheData(t *testing.T) {
	xs, ys := []float32{0, 1, 2}, []float32{0, 10, 20}
	for _, c := range []struct {
		at     float32
		points int
		lastX  float32
		lastY  float32
	}{{-0.001, 1, 0, 0}, {0, 1, 0, 0}, {1, 2, 1, 10}, {1.5, 3, 1.5, 15}, {2, 3, 2, 20}, {2.001, 3, 2, 20}, {5, 3, 2, 20}} {
		sx, sy := shadeUpTo(xs, ys, c.at)
		if len(sx) != c.points || sx[len(sx)-1] != c.lastX || sy[len(sy)-1] != c.lastY {
			t.Errorf("at %v: %v %v", c.at, sx, sy)
		}
	}
	if xs[2] != 2 || ys[2] != 20 {
		t.Errorf("the shade wrote over the data: %v %v", xs, ys)
	}
}

var (
	donutNames  = []string{"Red", "Green", "Blue"}
	donutShares = []float64{50, 30, 20}
)

func donutTester() *ui.Tester {
	return kitScreen(520, 300, func(c *ui.Context, p *palette) {
		donut(c, p, donutView{
			label: "Split", names: donutNames, shares: donutShares, center: "100", caption: "Total",
			describe: func(i int) string { return fmt.Sprintf("%s: %.0f", donutNames[i], donutShares[i]) },
		})
	})
}

func TestDonutShowsItsLegendAndAnswersThePointer(t *testing.T) {
	tt := donutTester()
	for _, want := range []string{"100", "Total", "Red", "50%", "30%", "20%"} {
		if !tt.HasText(want) {
			t.Errorf("missing %q in %q", want, tt.Texts())
		}
	}
	r, ok := tt.Find("Split")
	if !ok {
		t.Fatal("no donut")
	}
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	rim := func() ui.Color { return kitPixel(tt, cx+donutRadius+2, cy) }
	if kitSame(rim(), lightPalette.slices[0]) {
		t.Fatal("the donut reaches past its radius at rest")
	}
	tt.Move(cx+70, cy)
	if !tt.HasText("Red: 50") {
		t.Error("the slice under the pointer shows no value")
	}
	tt.Advance(time.Second)
	if !kitSame(rim(), lightPalette.slices[0]) || tt.FramePending() {
		t.Errorf("the slice under the pointer does not grow: %v, pending %v", rim(), tt.FramePending())
	}
	dot := func(name string) ui.Color {
		row, _ := tt.Find(name)
		if row.W < 100 {
			t.Fatalf("%q is not a legend row: %+v", name, row)
		}
		return kitPixel(tt, row.X+6, row.Y+row.H/2)
	}
	if !kitSame(dot("Red"), lightPalette.slices[0]) || kitSame(dot("Blue"), lightPalette.slices[2]) {
		t.Errorf("the legend does not dim the other slices: red %v, blue %v", dot("Red"), dot("Blue"))
	}
	tt.Move(cx+30, cy)
	tt.Advance(time.Second)
	if tt.HasText("Red: 50") || kitSame(rim(), lightPalette.slices[0]) {
		t.Error("the slice stays lit over the hole")
	}
	if !kitSame(dot("Blue"), lightPalette.slices[2]) {
		t.Errorf("the legend stays dimmed with nothing lit: blue %v", dot("Blue"))
	}
	kitHover(t, tt, "Green", 0.5, 0.5)
	if !tt.HasText("Green: 30") {
		t.Error("the legend does not light its slice")
	}
}

func TestDonutLegendTakesTheKeys(t *testing.T) {
	tt := donutTester()
	for range 6 {
		tt.Key(0, ui.KeyTab)
		if tt.HasText("Red: 50") {
			tt.Key(0, ui.KeyTab)
			if !tt.HasText("Green: 30") || tt.HasText("Red: 50") {
				t.Errorf("the next row does not take the tip: %q", tt.Texts())
			}
			return
		}
	}
	t.Error("the keys do not reach the legend")
}

func TestDonutKeepsTheTipWhileThePointerRestsOnIt(t *testing.T) {
	tt := donutTester()
	r, _ := tt.Find("Split")
	tt.Move(r.X+r.W/2+60, r.Y+r.H/2-40)
	for range 20 {
		tt.Advance(50 * time.Millisecond)
		if !tt.HasText("Red: 50") {
			t.Fatal("the tip goes as it comes under the pointer")
		}
	}
	if tt.FramePending() {
		t.Error("the tip keeps asking for frames under the pointer")
	}
}

func TestDonutLegendKeepsALongNameOnOneLine(t *testing.T) {
	names := []string{"Red", "Green", "A legend name that is far too long to fit beside its slice"}
	tt := kitScreen(520, 300, func(c *ui.Context, p *palette) {
		donut(c, p, donutView{
			label: "Split", names: names, shares: []float64{50, 30, 20}, center: "100", caption: "Total",
			describe: func(i int) string { return names[i] },
		})
	})
	row, ok := tt.Find(names[2])
	if !ok {
		t.Fatalf("no legend row in %q", tt.Texts())
	}
	background := kitPixel(tt, row.X+row.W-2, row.Y-4)
	for y := row.Y; y < row.Y+3; y++ {
		for x := row.X + 26; x < row.X+26+60; x++ {
			if !kitSame(kitPixel(tt, x, y), background) {
				t.Fatalf("a second line of the long name shows at the top of its row, at %v,%v", x, y)
			}
		}
	}
}
