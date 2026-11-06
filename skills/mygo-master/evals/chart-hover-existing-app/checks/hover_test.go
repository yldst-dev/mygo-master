package main

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

const marchNorth = "76"

func shows(tt *ui.Tester, value string) bool {
	for _, text := range tt.Texts() {
		if strings.Contains(text, value) && !strings.HasPrefix(text, value+" ") {
			return true
		}
	}
	return false
}

func startChart(t *testing.T) (*ui.Tester, ui.Rect) {
	t.Helper()
	a := newApp()
	tt := ui.NewTester(a.view, 640, 480)
	tt.Advance(time.Second)
	r, ok := tt.Find("Sales by month")
	if !ok {
		t.Fatal("no chart labeled Sales by month")
	}
	return tt, r
}

func sweepTo(tt *ui.Tester, r ui.Rect, value string) (float32, float32, bool) {
	y := r.Y + r.H*0.75
	for f := float32(0.02); f < 1; f += 0.02 {
		x := r.X + r.W*f
		tt.Move(x, y)
		if shows(tt, value) {
			return x, y, true
		}
	}
	return 0, 0, false
}

func TestTipShowsAtOnce(t *testing.T) {
	tt, r := startChart(t)
	if _, _, ok := sweepTo(tt, r, marchNorth); !ok {
		t.Errorf("no value shows without waiting: %q", tt.Texts())
	}
}

func TestTipDoesNotFlickerUnderThePointer(t *testing.T) {
	tt, r := startChart(t)
	if _, _, ok := sweepTo(tt, r, marchNorth); !ok {
		t.Fatal("no tip")
	}
	tt.Advance(time.Second)
	var tip ui.Rect
	found := false
	for _, text := range tt.Texts() {
		if strings.Contains(text, marchNorth) {
			tip, found = tt.Find(text)
		}
	}
	if !found {
		t.Fatal("no box for the tip")
	}
	tt.Move(tip.X+tip.W/2, tip.Y+tip.H/2)
	last, changes := shows(tt, marchNorth), 0
	for range 40 {
		tt.Advance(25 * time.Millisecond)
		if now := shows(tt, marchNorth); now != last {
			last = now
			changes++
		}
	}
	if changes > 1 {
		t.Errorf("the tip comes and goes %d times under a still pointer", changes)
	}
}

func TestRestsUnderThePointerAndAfterLeaving(t *testing.T) {
	tt, r := startChart(t)
	if _, _, ok := sweepTo(tt, r, marchNorth); !ok {
		t.Fatal("no tip")
	}
	tt.Advance(time.Second)
	if tt.FramePending() {
		t.Error("keeps asking for frames under the pointer")
	}
	tt.Move(4, 4)
	tt.Advance(time.Second)
	if shows(tt, marchNorth) {
		t.Error("the tip stays after the pointer leaves")
	}
	if tt.FramePending() {
		t.Error("keeps asking for frames after the pointer leaves")
	}
}

func TestTipStaysInsideTheChartAtBothEnds(t *testing.T) {
	for _, value := range []string{"42", "83"} {
		tt, r := startChart(t)
		if _, _, ok := sweepTo(tt, r, value); !ok {
			t.Errorf("no tip for %s", value)
			continue
		}
		tt.Advance(time.Second)
		for _, text := range tt.Texts() {
			if !strings.Contains(text, value) {
				continue
			}
			tip, ok := tt.Find(text)
			if !ok {
				continue
			}
			if tip.X < r.X-1 || tip.X+tip.W > r.X+r.W+1 {
				t.Errorf("%q runs from %v to %v, the chart from %v to %v", text, tip.X, tip.X+tip.W, r.X, r.X+r.W)
			}
		}
	}
}

func TestKeyboardWalksTheChartAndLeavesNoTip(t *testing.T) {
	tt, _ := startChart(t)
	reached := false
	for range 8 {
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeyRight)
		tt.Advance(100 * time.Millisecond)
		if shows(tt, "42") {
			reached = true
			break
		}
	}
	if !reached {
		t.Fatalf("the keys never show the first value: %q", tt.Texts())
	}
	tt.Key(ui.Shift, ui.KeyTab)
	tt.Advance(time.Second)
	if shows(tt, "42") {
		t.Error("the tip stays as the focus leaves the chart")
	}
}

func TestTipComesInGradually(t *testing.T) {
	tt, r := startChart(t)
	tt.Advance(time.Second)
	if _, _, ok := sweepTo(tt, r, marchNorth); !ok {
		t.Fatal("no tip")
	}
	moving := tt.FramePending()
	early := append([]uint8(nil), tt.Image().Pix...)
	tt.Advance(time.Second)
	late := tt.Image().Pix
	if !moving {
		t.Error("nothing moves as the tip comes")
	}
	same := true
	for i := range early {
		if early[i] != late[i] {
			same = false
			break
		}
	}
	if same {
		t.Error("the first frame already shows the end of the motion")
	}
}
