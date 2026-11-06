package view

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"starter/internal/domain"
	"starter/internal/usecase"
)

type memoryStore struct {
	saved []domain.Task
	fail  error
}

func (s *memoryStore) Load() ([]domain.Task, error) { return s.saved, nil }

func (s *memoryStore) Save(tasks []domain.Task) error {
	if s.fail != nil {
		return s.fail
	}
	s.saved = tasks
	return nil
}

func start(t *testing.T, store *memoryStore, setDark func(bool)) (*App, *ui.Tester) {
	t.Helper()
	tasks, err := usecase.NewTasks(store)
	if err != nil {
		t.Fatal(err)
	}
	a := New(tasks, false, setDark)
	tt := ui.NewTester(a.View, 560, 640)
	tt.Advance(time.Second)
	return a, tt
}

func twoTasks() *memoryStore {
	return &memoryStore{saved: []domain.Task{{ID: 1, Title: "First"}, {ID: 2, Title: "Second"}}}
}

func TestAddingTypesAndSubmitsAndClearsTheDraft(t *testing.T) {
	store := &memoryStore{}
	_, tt := start(t, store, func(bool) {})
	if err := tt.Click("New task"); err != nil {
		t.Fatal(err)
	}
	tt.Type("Write")
	tt.Key(0, ui.KeyEnter)
	if !tt.HasText("Write") || !tt.HasText("1 left") || len(store.saved) != 1 {
		t.Errorf("after Enter: texts %q, saved %+v", tt.Texts(), store.saved)
	}
	tt.Advance(time.Second)
	if tt.FramePending() {
		t.Error("keeps asking for frames at rest")
	}
}

func TestABlankTitleShowsTheProblemUntilTheNextGoodOne(t *testing.T) {
	_, tt := start(t, &memoryStore{}, func(bool) {})
	if err := tt.Click("Add"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText(domain.ErrEmptyTitle.Error()) {
		t.Fatalf("no problem shown: %q", tt.Texts())
	}
	if err := tt.Click("New task"); err != nil {
		t.Fatal(err)
	}
	tt.Type("Good")
	if err := tt.Click("Add"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if tt.HasText(domain.ErrEmptyTitle.Error()) || !tt.HasText("Good") {
		t.Errorf("after a good title: %q", tt.Texts())
	}
}

func TestASaveThatFailsIsShownAndKeepsTheDraft(t *testing.T) {
	store := &memoryStore{}
	a, tt := start(t, store, func(bool) {})
	store.fail = errors.New("disk full")
	if err := tt.Click("New task"); err != nil {
		t.Fatal(err)
	}
	tt.Type("Kept")
	if err := tt.Click("Add"); err != nil {
		t.Fatal(err)
	}
	if a.draft != "Kept" || a.problem == "" {
		t.Errorf("draft %q, problem %q", a.draft, a.problem)
	}
}

func TestRemovingCollapsesTheRowAndTheNeighborSlidesUp(t *testing.T) {
	_, tt := start(t, twoTasks(), func(bool) {})
	first, _ := tt.Find("First")
	second, ok := tt.Find("Second")
	if !ok {
		t.Fatal("no second row")
	}
	from := second.Y
	if err := tt.Click("Remove First"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(60 * time.Millisecond)
	img, background := tt.Image(), ui.LightTheme().Background
	drawn := false
	for x := int(first.X); x < int(first.X+first.W); x++ {
		p := img.RGBAAt(x, int(first.Y+first.H/2))
		drawn = drawn || p.R != background.R || p.G != background.G || p.B != background.B
	}
	if !drawn {
		t.Error("the row vanishes at once instead of collapsing")
	}
	second, _ = tt.Find("Second")
	moving := second.Y
	tt.Advance(time.Second)
	second, _ = tt.Find("Second")
	if from <= moving || moving <= second.Y {
		t.Errorf("the neighbor goes %v, %v, %v", from, moving, second.Y)
	}
	if tt.HasText("First") || tt.FramePending() {
		t.Errorf("after the removal: texts %q, pending %v", tt.Texts(), tt.FramePending())
	}
}

func TestADoubleClickRemovesOneRow(t *testing.T) {
	store := &memoryStore{saved: []domain.Task{{ID: 1, Title: "First"}, {ID: 2, Title: "Second"}, {ID: 3, Title: "Third"}}}
	_, tt := start(t, store, func(bool) {})
	r, ok := tt.Find("Remove First")
	if !ok {
		t.Fatal("no button")
	}
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.ClickAt(x, y)
	tt.Advance(120 * time.Millisecond)
	tt.ClickAt(x, y)
	tt.Advance(time.Second)
	if len(store.saved) != 2 || !tt.HasText("Second") {
		t.Fatalf("the second click takes the next row: %+v", store.saved)
	}
	tt.ClickAt(x, y)
	tt.Advance(time.Second)
	if len(store.saved) != 1 {
		t.Errorf("a later click does nothing: %+v", store.saved)
	}
}

func TestRemovingByKeyboardGivesTheFocusBack(t *testing.T) {
	_, tt := start(t, twoTasks(), func(bool) {})
	for range 12 {
		tt.Key(0, ui.KeyTab)
		if tt.Focused("Remove First") {
			break
		}
	}
	if !tt.Focused("Remove First") {
		t.Fatal("the keys do not reach the button")
	}
	tt.Key(0, ui.KeyEnter)
	tt.Advance(time.Second)
	if tt.HasText("First") || !tt.Focused("New task") {
		t.Errorf("after the key: texts %q, focus on the input %v", tt.Texts(), tt.Focused("New task"))
	}
}

func TestALongTitleKeepsTheButtonInTheWindow(t *testing.T) {
	title := strings.Repeat("unbroken", 60)
	_, tt := start(t, &memoryStore{saved: []domain.Task{{ID: 1, Title: title}}}, func(bool) {})
	r, ok := tt.Find("Remove " + title)
	if !ok {
		t.Fatal("no button")
	}
	if r.X+r.W > 560 {
		t.Errorf("the button ends at %v, past the window", r.X+r.W)
	}
}

func TestRemovingAtTheEndOfAScrolledListDoesNotJump(t *testing.T) {
	store := &memoryStore{}
	for i := 1; i <= 40; i++ {
		store.saved = append(store.saved, domain.Task{ID: i, Title: fmt.Sprintf("Task %d", i)})
	}
	_, tt := start(t, store, func(bool) {})
	for range 12 {
		tt.Wheel(280, 400, 0, 600)
		tt.Advance(time.Second)
	}
	above, ok := tt.Find("Task 38")
	if !ok {
		t.Fatalf("not at the end: %q", tt.Texts())
	}
	if err := tt.Click("Remove Task 40"); err != nil {
		t.Fatal(err)
	}
	moved, _ := tt.Find("Task 38")
	if jump := moved.Y - above.Y; jump > 6 || jump < -6 {
		t.Errorf("the rows above jump by %v", jump)
	}
	tt.Advance(time.Second)
	if tt.FramePending() {
		t.Error("keeps asking for frames at rest")
	}
}

func TestLoadFailureShowsTheReason(t *testing.T) {
	tt := ui.NewTester(LoadFailure(errors.New("broken file")), 560, 640)
	if !tt.HasText("broken file") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestTogglingUpdatesTheCount(t *testing.T) {
	store := twoTasks()
	_, tt := start(t, store, func(bool) {})
	if err := tt.Click("First"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if !tt.HasText("1 left") || !store.saved[0].Done {
		t.Errorf("texts %q, saved %+v", tt.Texts(), store.saved)
	}
}

func TestRowHoverFadesAndGoesBack(t *testing.T) {
	_, tt := start(t, twoTasks(), func(bool) {})
	r, ok := tt.Find("First")
	if !ok {
		t.Fatal("no row")
	}
	x, y := 25, int(r.Y+r.H/2)
	rest := tt.Image().RGBAAt(x, y)
	tt.Move(float32(x), float32(y))
	tt.Advance(40 * time.Millisecond)
	moving := tt.Image().RGBAAt(x, y)
	tt.Advance(time.Second)
	hovered := tt.Image().RGBAAt(x, y)
	if moving == rest || moving == hovered || tt.FramePending() {
		t.Errorf("colors %v, %v, %v, pending %v", rest, moving, hovered, tt.FramePending())
	}
	tt.Move(300, 620)
	tt.Advance(time.Second)
	if got := tt.Image().RGBAAt(x, y); got != rest {
		t.Errorf("does not go back: %v, was %v", got, rest)
	}
}

func TestDarkSwitchFadesTheWindowAndTellsTheHost(t *testing.T) {
	var told []bool
	_, tt := start(t, twoTasks(), func(dark bool) { told = append(told, dark) })
	light := lightTheme().Background
	dark := darkTheme().Background
	pixel := func() ui.Color {
		p := tt.Image().RGBAAt(4, 4)
		return ui.RGB(p.R, p.G, p.B)
	}
	if err := tt.Click("Dark mode"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(60 * time.Millisecond)
	if got := pixel(); got == light.SRGB() || got == dark.SRGB() {
		t.Errorf("the window does not fade: %v", got)
	}
	tt.Advance(time.Second)
	if got := pixel(); got != dark.SRGB() || tt.FramePending() {
		t.Errorf("after the fade: %v, pending %v", got, tt.FramePending())
	}
	if len(told) != 1 || !told[0] {
		t.Errorf("told %v", told)
	}
}

func TestThemeFadeCoversEveryColorAndEndsExactly(t *testing.T) {
	want := 0
	kind := reflect.TypeOf(ui.Theme{})
	for i := range kind.NumField() {
		if kind.Field(i).Type == reflect.TypeOf(ui.Color{}) && kind.Field(i).IsExported() {
			want++
		}
	}
	covered := themeColors(ui.LightTheme())
	if len(covered) > want {
		t.Errorf("the fade moves %d colors of %d", len(covered), want)
	}
	if got, base := fadedTheme(0), lightTheme(); got.Background != base.Background {
		t.Errorf("the start of the fade is %v", got.Background)
	}
	if got, base := fadedTheme(1), darkTheme(); got.Background != base.Background || !got.Dark {
		t.Errorf("the end of the fade is %v", got.Background)
	}
}

func TestNothingMovesWithLessMotion(t *testing.T) {
	_, tt := start(t, twoTasks(), func(bool) {})
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	tt.Advance(time.Second)
	if err := tt.Click("Remove First"); err != nil {
		t.Fatal(err)
	}
	if tt.HasText("First") {
		t.Errorf("the row stays while the desktop asks for less motion: %q", tt.Texts())
	}
}

func TestEveryTargetIsLargeEnough(t *testing.T) {
	_, tt := start(t, twoTasks(), func(bool) {})
	if small := smallTargets(tt, 28, "Dark mode", "New task", "Add", "First", "Remove First"); len(small) > 0 {
		t.Errorf("targets under 28 DIP: %v", small)
	}
}

func smallTargets(tt *ui.Tester, least float32, labels ...string) []string {
	var small []string
	for _, label := range labels {
		r, ok := tt.Find(label)
		if !ok || r.W < least || r.H < least {
			small = append(small, fmt.Sprintf("%s %.1fx%.1f", label, r.W, r.H))
		}
	}
	return small
}
