# 모션: 부드럽게 만드는 법

이 문서의 값과 레시피는 실제 대시보드 앱에서 나왔고, 요청자가 직접 써 보고 "충분히 부드럽다"고 판단한 값이다. 값을 새로 고르기 전에 이 값부터 쓴다. 코드는 `scripts/check_snippets.py` 로 컴파일과 동작을 함께 확인했다.

목차: [원칙](#원칙) · [도구 고르기](#도구-고르기) · [값 표](#값-표) · [레시피](#레시피) · [쉬는 화면](#쉬는-화면) · [줄인 움직임](#줄인-움직임) · [넣지 않는 곳](#넣지-않는-곳) · [비용](#비용)

## 원칙

부드러움은 효과를 많이 넣어서 생기지 않는다. 끊기는 곳을 없애서 생긴다.

1. **상태가 바뀔 때만 움직인다.** 장식으로 움직이지 않는다. 호버는 표면색만 바뀌고 크기나 위치는 그대로 둔다.
2. **짧게 움직인다.** 지정값은 전부 100~240ms이고 자주 쓰는 동작일수록 짧다. 300ms 넘는 값을 쓰지 않는다. 일회성으로 주의를 끄는 동작(벨 울림 360ms)만 예외다. 스프링의 `Duration` 은 체감 시간이라 꼬리가 더 남는다. 240ms 스프링으로 220 DIP를 옮기면 마지막 픽셀 변화가 336ms, 프레임 요청 종료가 584ms였다. 곡선 120ms는 136ms에 끝났다. 눈에 띄는 움직임은 지정값 근처에서 끝나므로 이 차이 때문에 값을 줄이지 않는다.
3. **튕기지 않는다.** 스프링은 모두 `Bounce` 0이다. 스프링은 공간을 가로지르는 것에만 쓰고, 색과 페이드는 곡선으로 한다.
4. **중간에 바뀌어도 이어진다.** 커서가 다른 곳으로 가거나 다시 눌러도 멈칫하거나 처음부터 시작하지 않는다. 스프링은 속도를 유지한 채 방향을 튼다.
5. **뛰는 것을 찾아서 없앤다.** 요소 하나를 움직이면 옆 요소가 뛴다. `Transition` 은 그 요소만 움직인다. 폭이 변하는 배지 옆의 글자, 알약 아래의 글자색, 숫자가 바뀌는 칸의 이웃까지 함께 처리한다.
6. **기다리게 하지 않는다.** 커서가 닿으면 바로 반응한다. 내장 `Tooltip` 의 0.6초 대기는 그래프처럼 값을 읽는 곳에 맞지 않는다.
7. **쉬면 흔적이 없다.** 애니메이션이 끝난 화면은 애니메이션을 넣기 전과 픽셀까지 같고, 프레임을 더 요청하지 않으며, CPU를 쓰지 않는다.
8. **겹치는 인상을 남기지 않는다.** 말풍선 글자가 바뀔 때 옛 글자와 새 글자가 오래 겹치면 흐리게 보인다. 작은 요소일수록 전환을 짧게 한다.

## 도구 고르기

| 하고 싶은 것 | 쓸 도구 | 이유 |
|---|---|---|
| 호버, 선택, 포커스에 따른 색 변화 | `progress` 스칼라 + `Mix` / `Alpha` | 목표색이 프레임마다 변해도 멈추지 않는다 |
| 누름, 화살표 살짝 밀기, 꺾쇠 회전 | `AnimatePose` | 뷰를 다시 돌리지 않고 칠하기만 해서 가장 싸다 |
| 형제들 위아래로 미끄러지는 알약 | `AnimateTo` + 부모의 `Draw` | `Absolute` 는 자식 위에 칠해져 글자를 가린다 |
| 내용에 따라 폭이 변하는 배지 | 배지에 `Transition{Size}` + `ClipY`, 이웃에 `Transition{Position}` | 이웃까지 함께 움직여야 뛰지 않는다 |
| 나타나고 사라지는 요소 | `Transition{Enter, Exit}` + `c.Key` | 나갈 때 복사본이 재생된다 |
| 숫자와 글자 변화 | `ContentTransition{Numeric}` + 부모 `ClipY` | 바뀐 글자만 구른다 |
| 아이콘 교체 | `ContentTransition{}` | 옛 아이콘이 줄고 새 아이콘이 커진다 |
| 데이터가 바뀌는 그림 | 데이터에 `AnimateTo` 후 `Draw` | 그래프의 면적, 도넛의 비율이 이어서 변한다 |
| 일회성 주의 끌기 | `Keyframes` + `Run` 카운터 | 첫 빌드에 재생되므로 `Run > 0` 일 때만 건다 |
| 페이지 전환 | 키를 준 `Scroll` + `Enter`, 또는 `ui.Router` | 키가 바뀌면 스크롤도 처음으로 돌아간다 |
| 같은 요소가 다른 자리로 이동 | `Transition{Match}` | 문서 기준. 대시보드에서는 쓰지 않았다 |
| 여러 값을 한 타임라인으로 | `Timeline` | 문서 기준. `examples/timeline` 참고 |
| 끌어서 던지기 | `Draggable` | 문서 기준. `docs/ui/motion.md` 의 Gestures |
| 스크롤 위치 이동 | `ScrollState.ScrollTo` | 문서 기준. `docs/ui/scroll.md` |

`ui.List` 의 행은 목록이 직접 배치하므로 행 자체에 `Transition` 을 걸어도 미끄러지거나 드나들지 않는다. 행이 드나드는 움직임이 필요하면 `Column` 에 `Key` 를 준 행을 쓰고, 행이 아주 많으면 움직임을 포기한다.

## 값 표

`assets/motion.go` 를 프로젝트에 복사하면 아래 값을 그대로 쓸 수 있다. 끝나면 쓰지 않는 값은 지운다.

| 이름 | 값 | 쓰는 곳 |
|---|---|---|
| `pressMotion` | 곡선 100ms | 누름 축소, 누르는 동안의 변화 |
| `fadeMotion` | 곡선 120ms | 호버, 선택, 포커스 표면과 글자색 |
| `quickMotion` | 곡선 150ms | 꺾쇠 회전, 화살표 밀기, 밑줄, 작은 점 |
| `themeMotion` | 곡선 200ms | 창 전체 라이트/다크 전환 |
| `markSpring` | 스프링 180ms | 그래프 말풍선과 표시 점 |
| `pillMotion` | 스프링 240ms | 내비 알약, 세그먼트 알약, 토글 손잡이, 데이터 변화, 숫자 굴림 |
| `pageEnter` | 곡선 180ms, Y 6, 불투명 0.35에서 시작 | 페이지 전환 |
| 드롭다운 | 곡선 160ms, Y 4, 축척 0.96 | 메뉴 열고 닫기. 위로 열리면 `Y` 부호와 `TransformOrigin` 을 뒤집는다 |
| `markMotion` | `markSpring`, 축척 0.9 | 말풍선과 그래프 표시의 등장, 퇴장, 이동 |
| `popMotion` | 곡선 120ms, 축척 0.6 | 입력칸 지우기 버튼처럼 작은 요소 |
| 행 전환 | `pillMotion`, 위치와 크기만, 접힘 | 목록 행의 추가와 삭제. 색은 넣지 않는다 |
| `dotExit` | 곡선 150ms, 축척 0.4 | 알림 점이 사라질 때 |
| `tipDigits` | 곡선 120ms | 말풍선의 숫자 |
| 벨 울림 | 360ms, 14 → -10 → 6 → -3 → 0도 | 메뉴를 열 때 한 번 |

누름 축척은 요소가 작을수록 크게 준다. 넓은 행과 내비 0.98, 일반 버튼 0.96, 작은 세그먼트 0.94, 아주 작은 아이콘 버튼 0.86.

그래프 강조값은 이렇다. 선택 점 1.3배, 선 굵기 +0.5, 면적 불투명도 +0.1, 도넛 조각 바깥 +5pt와 안쪽 -3pt, 나머지 조각은 카드 색 쪽으로 80% 흐림. 흐림을 50%로 하면 옆 조각이 호버된 조각처럼 보였다.

## 레시피

모든 레시피는 `assets/motion.go` 의 값과 `progress`, `pressable` 을 쓴다.

### 호버와 누름

표면은 `Alpha(0)` 일 때 아무것도 칠하지 않으므로 쉬는 화면이 그대로 유지된다.

```go
func softButton(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	button := pressable(ui.ButtonBase(c), 0.96).Height(30).Padding(0, 14).Radius(8).Label(label)
	hover := progress(button, "hover", button.Hovered(), fadeMotion)
	button.Background(t.SurfaceHover.Alpha(hover)).Children(func() {
		ui.Text(c, label).FontSize(14)
	})
	return button
}
```

```go test
func TestSoftButtonFadesThenRests(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { softButton(c, "Save") }, 320, 200)
	tt.Advance(time.Second)
	r, ok := tt.Find("Save")
	if !ok {
		t.Fatal("no button")
	}
	edge := func() color.RGBA { return tt.Image().RGBAAt(int(r.X)+2, int(r.Y+r.H/2)) }
	rest := edge()
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Advance(40 * time.Millisecond)
	moving := edge()
	tt.Advance(time.Second)
	hovered := edge()
	if moving == rest || moving == hovered {
		t.Errorf("no fade: %v, %v, %v", rest, moving, hovered)
	}
	if tt.FramePending() {
		t.Error("keeps asking for frames under the pointer")
	}
	tt.Move(300, 180)
	tt.Advance(time.Second)
	if edge() != rest {
		t.Errorf("does not go back: %v, was %v", edge(), rest)
	}
}
```

```go test
func TestSoftButtonShrinksWhilePressed(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Row(c).Padding(20).Children(func() { softButton(c, "Save") })
	}, 320, 200)
	tt.Advance(time.Second)
	r, ok := tt.Find("Save")
	if !ok {
		t.Fatal("no button")
	}
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Move(x, y)
	tt.Advance(time.Second)
	rest := bytes.Clone(tt.Image().Pix)
	tt.Press(x, y)
	tt.Advance(time.Second)
	if bytes.Equal(tt.Image().Pix, rest) {
		t.Error("the button keeps its size while pressed")
	}
	tt.Release(x, y)
	tt.Advance(time.Second)
	if !bytes.Equal(tt.Image().Pix, rest) || tt.FramePending() {
		t.Errorf("the button does not come back: pending %v", tt.FramePending())
	}
}
```

### 미끄러지는 내비 알약

알약을 `Absolute` 박스로 만들면 자식보다 위에 칠해져 글자를 가린다. 부모의 `Draw` 는 자식보다 먼저 칠해진다.

```go
func sideNav(c *ui.Context, selected *int, labels []string) {
	const height, gap = 40, 4
	t := c.Theme()
	nav := ui.Column(c).Gap(gap)
	pillY := nav.AnimateTo("pill", float32(*selected)*(height+gap), pillMotion)
	nav.Draw(func(pt *ui.Painter, r ui.Rect) {
		pt.Fill(ui.Rect{X: r.X, Y: r.Y + pillY, W: r.W, H: height}, t.SurfaceHover, 10)
	}).Children(func() {
		for i, label := range labels {
			item := pressable(ui.ButtonBase(c), 0.98).Height(height).Padding(0, 14).Radius(10).Justify(ui.Start).Label(label)
			ink := t.TextMuted.Mix(t.Accent, progress(item, "selected", i == *selected, fadeMotion))
			item.Children(func() {
				ui.Text(c, label).TextColor(ink)
			})
			if item.Clicked() {
				*selected = i
			}
		}
	})
}
```

```go test
func TestSideNavPillSlidesThenRests(t *testing.T) {
	selected := 0
	tt := ui.NewTester(func(c *ui.Context) { sideNav(c, &selected, []string{"One", "Two", "Three", "Four"}) }, 240, 240)
	tt.Advance(time.Second)
	pillTop := func() int {
		img, want := tt.Image(), ui.LightTheme().SurfaceHover
		for y := 0; y < img.Bounds().Dy(); y++ {
			p := img.RGBAAt(120, y)
			if p.R == want.R && p.G == want.G && p.B == want.B {
				return y
			}
		}
		return -1
	}
	start := pillTop()
	if err := tt.Click("Four"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(60 * time.Millisecond)
	moving := pillTop()
	tt.Advance(time.Second)
	rest := pillTop()
	if start >= moving || moving >= rest {
		t.Errorf("the pill goes %d, %d, %d", start, moving, rest)
	}
	if tt.FramePending() {
		t.Error("keeps asking for frames at rest")
	}
}
```

### 세그먼트 알약과 글자색

선택이 바뀐 순간 글자색만 먼저 흰색이 되면 알약이 도착하기 전에 글자가 안 보인다. 글자색을 알약과 겹친 정도에서 계산한다.

```go
func monthPills(c *ui.Context, selected *int, labels []string) {
	const width, height = 50, 28
	t := c.Theme()
	picker := ui.SegmentedBase(c.Key("months"), selected, len(labels))
	pillX := picker.Track.AnimateTo("pill", float32(*selected)*width, pillMotion)
	picker.Track.Label("Month").Height(height+2).Radius(8).Border(1, t.Border).Draw(func(pt *ui.Painter, r ui.Rect) {
		pt.Fill(ui.Rect{X: r.X + 1 + pillX, Y: r.Y + 1, W: width, H: height}, t.Accent, 7)
	}).Children(func() {
		for i, label := range labels {
			segment := pressable(picker.Segment(i), 0.94).Size(width, height).Radius(7).Center()
			offset := pillX - float32(i)*width
			covered := max(0, 1-max(offset, -offset)/width)
			segment.Children(func() {
				ui.Text(c, label).FontSize(14).TextColor(t.Text.Mix(t.AccentText, covered))
			})
		}
	})
}
```

### 창 전체 색 전환

`AnimateColor` 와 `Transition{Colors}` 는 목표색이 바뀌는 순간 진행도를 0으로 되돌린다. 테마가 전환되는 동안 목표색은 프레임마다 바뀌므로 그 요소는 전환이 끝날 때까지 멈춰 있다가 한 번에 바뀐다. 색을 애니메이션하지 말고 스칼라를 애니메이션한 뒤 현재 팔레트에서 섞는다. 그러면 모든 요소가 같은 박자로 움직인다.

```go
type swatch struct {
	page, surface, ink ui.Color
}

var (
	lightSwatch = swatch{page: ui.Hex("#F6F7FC"), surface: ui.Hex("#FEFEFF"), ink: ui.Hex("#0B1220")}
	darkSwatch  = swatch{page: ui.Hex("#0B0F1A"), surface: ui.Hex("#141B2C"), ink: ui.Hex("#F3F6FC")}
)

func (s *swatch) channels() []*ui.Color {
	return []*ui.Color{&s.page, &s.surface, &s.ink}
}

func fadedSwatch(darkness float32) *swatch {
	switch darkness {
	case 0:
		return &lightSwatch
	case 1:
		return &darkSwatch
	}
	faded := lightSwatch
	from, to := faded.channels(), darkSwatch.channels()
	for i, channel := range from {
		*channel = channel.Mix(*to[i], darkness)
	}
	return &faded
}

func themedRoot(c *ui.Context, dark bool, build func(s *swatch)) {
	s := fadedSwatch(progress(c.Root(), "theme", dark, themeMotion))
	t := *ui.LightTheme()
	if dark {
		t = *ui.DarkTheme()
	}
	t.Background, t.Surface, t.Text = s.page, s.surface, s.ink
	c.SetTheme(&t)
	build(s)
}
```

```go test
func TestSwatchFadeCoversEveryColorAndEndsExactly(t *testing.T) {
	count := reflect.TypeOf(swatch{}).NumField()
	if got := len(lightSwatch.channels()); got != count {
		t.Errorf("the swatch has %d colors, the fade moves %d", count, got)
	}
	if fadedSwatch(0) != &lightSwatch || fadedSwatch(1) != &darkSwatch {
		t.Error("the ends of the fade are not the swatches themselves")
	}
}
```

```go test
func TestThemeFadeMovesEverythingTogether(t *testing.T) {
	dark := false
	var shown *swatch
	tt := ui.NewTester(func(c *ui.Context) {
		themedRoot(c, dark, func(s *swatch) {
			shown = s
			ui.Column(c).Fill().Background(s.surface)
		})
	}, 200, 120)
	tt.Advance(time.Second)
	dark = true
	tt.Frame()
	tt.Advance(60 * time.Millisecond)
	if shown == &lightSwatch || shown == &darkSwatch {
		t.Fatal("the theme does not fade")
	}
	p := tt.Image().RGBAAt(100, 60)
	if p.R != shown.surface.R || p.G != shown.surface.G || p.B != shown.surface.B {
		t.Errorf("the surface shows %v while the swatch is at %v", p, shown.surface)
	}
	tt.Advance(time.Second)
	if shown != &darkSwatch || tt.FramePending() {
		t.Errorf("after the fade: dark %v, pending %v", shown == &darkSwatch, tt.FramePending())
	}
}
```

새 색을 추가하고 `channels` 에 빠뜨리면 그 색만 전환되지 않으므로 위 테스트로 막는다. 끝값이 팔레트 자체여야 쉬는 화면의 픽셀이 정확히 같다. 앱의 외관 설정(`mygo.Theme.SetSource`)은 전환을 시작할 때 함께 바꾼다.

### 숫자가 구르는 칸과 폭이 변하는 배지

숫자 굴림은 위아래로 글자 높이의 0.6만큼 벗어난다. `Text` 에 `Clip` 을 걸어도 자기 글자는 잘리지 않으므로 부모에 `ClipY` 를 건다. `Clip` 대신 `ClipY` 를 쓰는 이유는 옆으로 밀리는 글자까지 자르지 않기 위해서다.

```go
func rollingFigure(c *ui.Context, text string) {
	ui.Box(c).ClipY().Children(func() {
		ui.Text(c, text).FontSize(29).Bold().NoWrap().ContentTransition(rollDigits)
	})
}

func deltaBadge(c *ui.Context, text string, falling bool) {
	t := c.Theme()
	ui.Row(c).Gap(7).Children(func() {
		pill := ui.Row(c).Height(24).PaddingX(9).Radius(12).ClipY().Transition(growFit)
		fall := progress(pill, "falling", falling, fadeMotion)
		ink := t.Success.Mix(t.Danger, fall)
		pill.Background(ink.Alpha(0.14)).Children(func() {
			ui.Text(c, text).FontSize(12.5).FontWeight(600).TextColor(ink).NoWrap().ContentTransition(rollDigits)
		})
		ui.Text(c, "vs. last month").FontSize(12.5).NoWrap().Transition(followFit)
	})
}
```

```go test
func TestDeltaBadgeLeadsItsNeighborWithoutAJump(t *testing.T) {
	text := "+1.2%"
	tt := ui.NewTester(func(c *ui.Context) {
		rollingFigure(c, text)
		deltaBadge(c, text, false)
	}, 420, 160)
	tt.Advance(time.Second)
	follow, ok := tt.Find("vs. last month")
	if !ok {
		t.Fatal("no neighbor")
	}
	from := follow.X
	text = "+1234567.8%"
	tt.Frame()
	tt.Advance(40 * time.Millisecond)
	follow, _ = tt.Find("vs. last month")
	moving := follow.X
	gap := tt.Image().RGBAAt(int(follow.X)-3, int(follow.Y+follow.H/2))
	page := ui.LightTheme().Background
	if gap.R != page.R || gap.G != page.G || gap.B != page.B {
		t.Errorf("the badge runs under its neighbor: %v", gap)
	}
	tt.Advance(time.Second)
	follow, _ = tt.Find("vs. last month")
	if from >= moving || moving >= follow.X {
		t.Errorf("the neighbor goes %v, %v, %v", from, moving, follow.X)
	}
	if tt.FramePending() {
		t.Error("keeps asking for frames at rest")
	}
}
```

### 나타나고 사라지는 드롭다운

내장 `PopoverMenuButton` 은 움직임을 넣을 곳이 없어 `PopoverBase` 로 직접 짠다. 빠뜨리기 쉬운 점이 네 가지 있다.

- `ui.Overlay` 는 닫혀 있을 때도 매 프레임 불러야 한다. 부모가 빌드되어야 `Exit` 가 재생된다.
- `c.Key(open)` 은 `if !*open { return }` 뒤에 둔다. 메뉴끼리 바로 넘어갈 때 패널이 같은 ID를 쓰면 옛 패널이 잘리고, 키를 먼저 잡으면 다음 요소로 새어 나간다.
- `AutoFocus` 는 `ui.Overlay` 가 돌아온 뒤에 부른다. 안에서 부르면 Esc로 닫을 때 포커스를 돌려줄 요소가 잘못 기록된다.
- 포인터가 새 행에 들어올 때만 포커스를 옮겨야 키보드와 호버가 한 행만 강조한다.

```go
type menuItem struct {
	label  string
	choose func()
}

func softMenu(c *ui.Context, label string, open *bool, upward bool, items []menuItem) {
	t := c.Theme()
	button := pressable(ui.ButtonBase(c), 0.96).Height(34).Padding(0, 10).Radius(8).Label(label).Role(ui.RoleMenuButton)
	if button.Clicked() || !*open && button.Shortcut(0, ui.KeyDown) {
		*open = !*open
	}
	button.Background(t.SurfaceHover.Alpha(progress(button, "hover", *open || button.Hovered(), fadeMotion))).Expanded(*open).Children(func() {
		ui.Text(c, label).FontSize(14)
	})
	lift, origin := float32(-4), float32(0)
	if upward {
		lift, origin = 4, 1
	}
	motion := ui.ElementTransition{Colors: true, Duration: 160 * time.Millisecond, Enter: &ui.Motion{Y: lift, Scale: 0.96}, Exit: &ui.Motion{Y: lift, Scale: 0.96}}
	var first ui.Element
	ui.Overlay(c, func() {
		if !*open {
			return
		}
		ui.PopoverBase(c.Key(open), button, open, func(panel ui.Element) {
			panel.Role(ui.RoleMenu).FocusGroup(ui.Vertical).Margin(4, 0, 0, 0).MinWidth(180).Padding(4).Radius(10).
				Background(t.Background).Border(1, t.Border).Shadow(0, 6, 20, 0, ui.RGBA(0, 0, 0, 0.18)).
				TransformOrigin(0, origin).Transition(motion)
			if panel.Shortcut(0, ui.KeyTab) || panel.Shortcut(ui.Shift, ui.KeyTab) {
				*open = false
			}
			rows := make([]ui.Element, len(items))
			lastHovered := ui.Local(panel, "hovered", func() int { return -1 })
			hovered := -1
			for i, item := range items {
				rows[i] = pressable(ui.ButtonBase(c), 0.98).Label(item.label).Role(ui.RoleMenuItem).FocusRing(false).Padding(6, 8).Radius(8).Justify(ui.Start).Shrink(0)
				if rows[i].Hovered() {
					hovered = i
					if *lastHovered != i {
						rows[i].Focus()
					}
				}
			}
			*lastHovered = hovered
			first = rows[0]
			for i, row := range rows {
				lit := progress(row, "highlight", row.FocusVisible() || row.Focused() && row.Hovered(), fadeMotion)
				row.Background(t.Accent.Alpha(lit)).Children(func() {
					ui.Text(c, items[i].label).SingleLine().TextColor(t.Text.Mix(t.AccentText, lit))
				})
				if row.Clicked() {
					*open = false
					items[i].choose()
				}
			}
		})
	})
	first.AutoFocus()
}
```

```go test
func TestSoftMenuOpensClosesAndGivesFocusBack(t *testing.T) {
	open := false
	chosen := ""
	items := []menuItem{
		{"First", func() { chosen = "First" }},
		{"Second", func() { chosen = "Second" }},
	}
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Row(c).Padding(20).Children(func() { softMenu(c, "More", &open, false, items) })
	}, 400, 300)
	if err := tt.Click("More"); err != nil {
		t.Fatal(err)
	}
	if !open || !tt.HasText("First") {
		t.Fatalf("the menu does not open: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEscape)
	tt.Advance(time.Second)
	if open || tt.HasText("First") || !tt.Focused("More") {
		t.Errorf("after Escape: open %v, texts %q, focused %v", open, tt.Texts(), tt.Focused("More"))
	}
	tt.Key(0, ui.KeyDown)
	if !open || !tt.Focused("First") {
		t.Errorf("Down does not open the menu on its first item: open %v", open)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	tt.Advance(time.Second)
	if open || chosen != "Second" {
		t.Errorf("after Enter: open %v, chosen %q", open, chosen)
	}
	if tt.FramePending() {
		t.Error("keeps asking for frames at rest")
	}
}
```

### 드나드는 목록 행

`rowMotion` 은 위치와 크기만 전환한다. 플래그를 비우면 색까지 전환되어 호버 페이드가 멈칫한다. 행은 `c.Key(id)` 로 만들고 삭제는 `OnClick` 으로 미룬다. 스크롤 안의 목록에는 세 가지를 더 챙긴다.

- 끝까지 내린 상태에서 마지막 행을 지우면 내용 높이가 바로 줄어 위 행들이 튄다. 바깥 열의 최소 높이를 스프링으로 줄인다.
- 지우면 다음 행의 삭제 버튼이 포인터 아래로 올라온다. 더블클릭의 두 번째 클릭이 다음 행을 지우므로 500ms 안의 반복을 무시한다.
- 키보드로 지우면 포커스가 함께 사라진다. 지운 버튼에 포커스가 있었으면 입력칸으로 돌려준다. `KeepFocus` 를 걸면 포인터로 누를 때는 포커스가 옮겨 가지 않아 둘을 구분할 수 있다.

끊을 곳 없는 긴 글자가 있는 행은 글자 쪽에 `MinWidth(0)` 을 줘야 버튼이 창 밖으로 밀리지 않는다. 시작 템플릿의 `internal/view/app.go` 에 모두 들어 있다.

```go
func steadyList(c *ui.Context, build func()) {
	ui.Scroll(c).Grow(1).Children(func() {
		outer := ui.Column(c)
		outer.Children(func() {
			inner := ui.Column(c).Children(build)
			if height := inner.Bounds().H; height > 0 {
				outer.MinHeight(outer.AnimateTo("height", height, pillMotion))
			}
		})
	})
}
```

```go test
func TestSteadyListDoesNotJumpAtTheEnd(t *testing.T) {
	rows := 40
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Fill().Children(func() {
			steadyList(c, func() {
				for i := range rows {
					ui.Row(c.Key(i)).Height(40).ClipY().Transition(rowMotion).Children(func() {
						ui.Text(c, fmt.Sprintf("Row %d", i))
					})
				}
			})
		})
	}, 300, 400)
	tt.Advance(time.Second)
	for range 8 {
		tt.Wheel(150, 200, 0, 600)
		tt.Advance(time.Second)
	}
	above, ok := tt.Find("Row 37")
	if !ok {
		t.Fatalf("not at the end: %q", tt.Texts())
	}
	rows = 39
	tt.Frame()
	moved, _ := tt.Find("Row 37")
	if jump := moved.Y - above.Y; jump > 6 || jump < -6 {
		t.Errorf("the rows above jump by %v", jump)
	}
	tt.Advance(time.Second)
	if tt.FramePending() {
		t.Error("keeps asking for frames at rest")
	}
}
```

### 페이지 전환

`Enter` 는 부모가 이전 프레임에 이미 있었을 때만 재생된다. 그래서 키를 준 요소의 자식이 아니라 키를 준 요소 자신에 `Transition` 을 건다. 첫 프레임이 비지 않도록 `Opacity` 를 0이 아닌 0.35에서 시작한다.

```go
func pageHost(c *ui.Context, page int, build func()) {
	ui.Scroll(c.Key(page)).Grow(1).Padding(20).Transition(pageEnter).Children(func() {
		ui.Column(c).Gap(20).Children(build)
	})
}
```

```go test
func TestPageHostStartsEachPageAtTheTop(t *testing.T) {
	page := 0
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Fill().Children(func() {
			pageHost(c, page, func() {
				ui.Text(c, fmt.Sprintf("Page %d", page))
				ui.Box(c).Height(2000)
			})
		})
	}, 400, 300)
	tt.Advance(time.Second)
	top, ok := tt.Find("Page 0")
	if !ok {
		t.Fatal("no page")
	}
	tt.Wheel(200, 150, 0, 400)
	tt.Advance(time.Second)
	if scrolled, _ := tt.Find("Page 0"); scrolled.Y >= top.Y {
		t.Fatalf("the page does not scroll: %v", scrolled.Y)
	}
	page = 1
	tt.Frame()
	tt.Advance(40 * time.Millisecond)
	if !tt.FramePending() {
		t.Error("the next page does not come in")
	}
	tt.Advance(time.Second)
	if next, _ := tt.Find("Page 1"); next.Y != top.Y || tt.FramePending() {
		t.Errorf("the next page rests at %v, the top is %v", next.Y, top.Y)
	}
}
```

### 일회성 벨 울림

`Keyframes` 는 처음 빌드될 때도 재생되므로 `Run` 이 0일 때는 걸지 않는다.

```go
var bellRing = []ui.Keyframe{
	{Value: 0},
	{At: 70 * time.Millisecond, Value: 14},
	{At: 150 * time.Millisecond, Value: -10},
	{At: 230 * time.Millisecond, Value: 6},
	{At: 300 * time.Millisecond, Value: -3},
	{At: 360 * time.Millisecond, Value: 0},
}

func ringingIcon(c *ui.Context, icon *ui.SVG, rings uint64) {
	bell := ui.Icon(c, icon).Size(26, 26).TransformOrigin(0.5, 0.1)
	if rings > 0 {
		bell.Transform(ui.Rotation(bell.Keyframes("ring", bellRing, ui.KeyframeOptions{Run: rings})))
	}
}
```

```go test
func TestRingingIconWaitsThenRingsOnce(t *testing.T) {
	icon := ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M2 2h20v20H2z"/></svg>`))
	rings := uint64(0)
	tt := ui.NewTester(func(c *ui.Context) { ringingIcon(c, icon, rings) }, 120, 120)
	tt.Advance(40 * time.Millisecond)
	if tt.FramePending() {
		t.Error("the icon rings before anything asks for it")
	}
	rest := tt.Image().RGBAAt(5, 5)
	rings = 1
	tt.Frame()
	if !tt.FramePending() {
		t.Error("the icon does not ring")
	}
	tt.Advance(time.Second)
	if tt.FramePending() || tt.Image().RGBAAt(5, 5) != rest {
		t.Errorf("after the ring: pending %v", tt.FramePending())
	}
}
```

### 그래프 호버

세 가지 그래프에 모두 쓰는 공통 부품이다. `chart.Select` 는 포인터가 떠날 때만 선택을 지우므로 포커스가 빠져나갈 때도 지우도록 감싼다. 말풍선에는 항상 `IgnoreInput(true)` 를 건다. 말풍선이 호버 대상의 바깥에 있으면(도넛 레시피) 커서 아래로 왔을 때 호버를 빼앗아 사라졌다 나타나기를 반복한다. `PassThrough` 는 그 요소에만 적용되고 자식에 물려지지 않아 이 문제를 막지 못한다. 말풍선이 그래프 영역의 자식이면 영역의 호버가 유지되어 깜빡이지 않으므로, 이때 `IgnoreInput` 을 없애는 변이는 `MISSED` 가 정상이다. 위치는 실제 글자 폭과 말풍선 높이로 그래프 안에 가둔다. 막대처럼 위까지 솟는 값에서는 위쪽 가두기가 없으면 제목이나 버튼을 가린다. 색은 내장 말풍선처럼 테마의 `Inverse` 를 먼저 쓴다.

```go
const (
	tipPadding  = 8
	tipFontSize = 12.5
	tipHeight   = 26
	tipGap      = 8
)

func pickPoint(c *ui.Context, area ui.Element, xs []float32, describe func(int) string) int {
	selected := ui.Local(area, "selected", func() int { return -1 })
	chart.Select(c, area, selected, xs, describe)
	if !area.Hovered() && !area.Focused() {
		*selected = -1
	}
	return *selected
}

func chartTip(c *ui.Context, bounds ui.Rect, x, y float32, text string) {
	if bounds.W == 0 || bounds.H == 0 {
		return
	}
	t := c.Theme()
	fill, ink := t.Inverse, t.InverseText
	if fill.A == 0 {
		fill, ink = t.Text, t.Background
	}
	textWidth, _ := c.MeasureText(0, ui.Span{Text: text, Size: tipFontSize})
	reach := (tipPadding + textWidth/2) / bounds.W
	ceiling := 1 - (tipHeight+tipGap)/bounds.H
	chart.Anchor(c.Key("tip"), min(max(x, reach), 1-reach), min(y, ceiling), ui.Center, ui.End, func() {
		ui.Row(c).Height(tipHeight).PaddingX(tipPadding).Margin(0, 0, tipGap, 0).Radius(5).Background(fill).Shadow(0, 2, 8, 0, ui.RGBA(0, 0, 0, 0.2)).
			Role(ui.RoleTooltip).ClipY().Transition(tipFit).Children(func() {
			ui.Text(c, text).FontSize(tipFontSize).TextColor(ink).NoWrap().ContentTransition(tipDigits)
		})
	}).IgnoreInput(true).Transition(markMotion)
}

func hoverLine(c *ui.Context, label string, names []string, values []float64) {
	t := c.Theme()
	scale := chart.NewLinear(chart.Extent(values...))
	xs, ys := make([]float32, len(values)), make([]float32, len(values))
	for i, v := range values {
		xs[i] = float32(i) / float32(len(values)-1)
		ys[i] = 0.1 + 0.7*scale.At(v)
	}
	describe := func(i int) string { return fmt.Sprintf("%s: %.0f", names[i], values[i]) }
	chart.Plot(c, chart.Axes{}, func(area ui.Element) {
		area.Label(label).Radius(6)
		selected := pickPoint(c, area, xs, describe)
		lit := progress(area, "lit", selected >= 0, fadeMotion)
		area.Draw(func(pt *ui.Painter, r ui.Rect) {
			pt.FillPath(chart.Area(r, xs, nil, ys, chart.CurveMonotone), t.Accent.Alpha(0.14+0.1*lit))
			pt.StrokePath(chart.Line(r, xs, ys, chart.CurveMonotone), 2+0.5*lit, t.Accent)
		})
		if selected >= 0 {
			chart.RuleX(c.Key("guide"), xs[selected], 1).Background(t.Border).Transition(markMotion)
			chart.Point(c.Key("mark"), xs[selected], ys[selected], 10, 10).Radius(5).Background(t.Accent).Border(2, t.Background).Transition(markMotion)
			chartTip(c, area.Bounds(), xs[selected], ys[selected], describe(selected))
		}
	}).Size(320, 160)
}
```

```go test
func TestHoverLineAnswersAtOnceAndLeavesNothing(t *testing.T) {
	names, values := []string{"Jan", "Feb", "Mar", "Apr"}, []float64{10, 40, 20, 50}
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(40).Children(func() { hoverLine(c, "Sales", names, values) })
	}, 500, 300)
	tt.Advance(time.Second)
	r, ok := tt.Find("Sales")
	if !ok {
		t.Fatal("no chart")
	}
	tt.Move(r.X+r.W*0.34, r.Y+r.H/2)
	if !tt.HasText("Feb: 40") {
		t.Errorf("no tip as the pointer comes: %q", tt.Texts())
	}
	tt.Move(r.X+r.W*0.99, r.Y+r.H/2)
	tt.Advance(time.Second)
	tip, ok := tt.Find("Apr: 50")
	if !ok {
		t.Fatalf("the tip does not follow: %q", tt.Texts())
	}
	if tip.X+tip.W+tipPadding > r.X+r.W+1 {
		t.Errorf("the tip ends at %v, past the chart at %v", tip.X+tip.W+tipPadding, r.X+r.W)
	}
	for range 20 {
		tt.Move(tip.X+tip.W/2, tip.Y+tip.H/2)
		tt.Advance(50 * time.Millisecond)
		if !tt.HasText("Apr: 50") {
			t.Fatal("the tip goes as it comes under the pointer")
		}
	}
	if tt.FramePending() {
		t.Error("keeps asking for frames under the pointer")
	}
	tt.Move(5, 5)
	tt.Advance(time.Second)
	if tt.HasText("Apr: 50") || tt.FramePending() {
		t.Errorf("after leaving: texts %q, pending %v", tt.Texts(), tt.FramePending())
	}
}
```

위까지 솟는 값에서도 말풍선은 그래프 안에 머문다.

```go test
func TestChartTipStaysBelowTheTopOfItsChart(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(40).Children(func() {
			chart.Plot(c, chart.Axes{}, func(area ui.Element) {
				area.Label("Peak chart")
				chartTip(c, area.Bounds(), 0.5, 1, "Peak")
			}).Size(240, 120)
		})
	}, 400, 240)
	tt.Frame()
	tt.Advance(time.Second)
	box, ok := tt.Find("Peak chart")
	if !ok {
		t.Fatal("no chart")
	}
	tip, ok := tt.Find("Peak")
	if !ok {
		t.Fatalf("no tip: %q", tt.Texts())
	}
	if tip.Y < box.Y {
		t.Errorf("the tip starts at %v, above the chart at %v", tip.Y, box.Y)
	}
}
```

말풍선은 커지며 나타나고, 다른 점으로 갈 때는 미끄러진다.

```go test
func TestChartTipFadesInAndSlides(t *testing.T) {
	names, values := []string{"Jan", "Feb", "Mar", "Apr", "May"}, []float64{10, 40, 20, 50, 30}
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(40).Children(func() { hoverLine(c, "Sales", names, values) })
	}, 500, 300)
	tt.Advance(time.Second)
	r, ok := tt.Find("Sales")
	if !ok {
		t.Fatal("no chart")
	}
	ink := ui.LightTheme().Text
	solid := func(tip ui.Rect) bool {
		p := tt.Image().RGBAAt(int(tip.X)-4, int(tip.Y+tip.H/2))
		return p.R == ink.R && p.G == ink.G && p.B == ink.B
	}
	tt.Move(r.X+r.W*0.5, r.Y+r.H/2)
	tt.Advance(40 * time.Millisecond)
	tip, ok := tt.Find("Mar: 20")
	if !ok {
		t.Fatalf("no tip: %q", tt.Texts())
	}
	if solid(tip) {
		t.Error("the tip shows in full at once")
	}
	tt.Advance(time.Second)
	tip, _ = tt.Find("Mar: 20")
	if !solid(tip) {
		t.Error("the tip does not settle")
	}
	from := tip.X + tip.W/2
	tt.Move(r.X+r.W*0.75, r.Y+r.H/2)
	tt.Advance(40 * time.Millisecond)
	tip, ok = tt.Find("Apr: 50")
	if !ok {
		t.Fatalf("the tip does not change: %q", tt.Texts())
	}
	moving := tip.X + tip.W/2
	tt.Advance(time.Second)
	tip, _ = tt.Find("Apr: 50")
	if to := tip.X + tip.W/2; from >= moving || moving >= to {
		t.Errorf("the tip goes %v, %v, %v", from, moving, to)
	}
}
```

키보드로 다니다 벗어날 때와 다른 요소에 가려졌을 때도 말풍선이 남으면 안 된다. 그래프 영역의 호버는 위에 덮인 요소를 고려하므로 가드 하나로 둘 다 막는다.

```go test
func TestHoverLineFollowsKeysAndHidesUnderAnOverlay(t *testing.T) {
	names, values := []string{"Jan", "Feb", "Mar"}, []float64{10, 40, 20}
	covered := false
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(40).Children(func() { hoverLine(c, "Sales", names, values) })
		ui.Button(c, "After")
		if covered {
			ui.Overlay(c, func() {
				ui.Box(c).Absolute().Left(0).Top(0).Size(500, 300).Background(ui.RGB(255, 255, 255))
			})
		}
	}, 500, 320)
	tt.Advance(time.Second)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("Jan: 10") {
		t.Errorf("the keys do not move the choice: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Advance(time.Second)
	if tt.HasText("Jan: 10") {
		t.Error("the tip stays as the focus leaves")
	}
	r, ok := tt.Find("Sales")
	if !ok {
		t.Fatal("no chart")
	}
	covered = true
	tt.Frame()
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Advance(time.Second)
	if tt.HasText("Feb: 40") {
		t.Error("the chart answers the pointer under an overlay")
	}
}
```

### 도넛 호버

`Hovered` 는 위에 덮인 요소를 고려하고 `PointerPosition` 은 기하만 본다. 둘을 함께 써야 메뉴에 가려진 조각이 반응하지 않는다. 판정 반지름에 커지는 값을 더해야 커진 가장자리에서 깜빡이지 않는다. 범례 행과 도넛이 서로 반응하게 하려면 범례 행을 먼저 만들어 호버를 읽은 뒤 강조값을 계산한다. 말풍선은 맨 마지막 자식(`Absolute` 이면서 `PassThrough`)에 넣어야 범례 위로 칠해진다.

```go
const (
	ringRadius = 91
	ringHole   = 54
	ringGrowth = 5
	ringSink   = 3
	ringFade   = 0.8
)

func hoverDonut(c *ui.Context, names []string, shares []float64, colors []ui.Color) {
	t := c.Theme()
	ui.Row(c).Gap(23).Children(func() {
		donut := ui.Box(c).Size(2*ringRadius, 2*ringRadius).Shrink(0).Center().Label("Share by name")
		shown := make([]float64, len(shares))
		for i, share := range shares {
			shown[i] = float64(donut.AnimateTo(i, float32(share), pillMotion))
		}
		slices := chart.Pie(shown, 0.03)
		active := -1
		if x, y, over := donut.PointerPosition(); over && donut.Hovered() {
			if i, ok := chart.SliceAt(slices, x-ringRadius, y-ringRadius, ringHole, ringRadius+ringGrowth); ok {
				active = i
			}
		}
		lit := make([]float32, len(slices))
		var dim float32
		ui.Column(c).Grow(1).Gap(15).Children(func() {
			rows := make([]ui.Element, len(names))
			for i, name := range names {
				rows[i] = ui.Row(c).Height(20).Gap(10).Focusable().Label(name)
				if rows[i].Hovered() || rows[i].FocusVisible() {
					active = i
				}
			}
			dim = progress(donut, "dim", active >= 0, fadeMotion)
			for i, name := range names {
				lit[i] = progress(donut, name, i == active, fadeMotion)
				rows[i].TextColor(t.Text.Mix(t.TextMuted, dim-lit[i])).Children(func() {
					ui.Box(c).Size(12, 12).Radius(6).Background(colors[i].Mix(t.Surface, ringFade*(dim-lit[i])))
					ui.Text(c, name).FontSize(13.5).Grow(1)
				})
			}
		})
		donut.Draw(func(pt *ui.Painter, r ui.Rect) {
			cx, cy := r.X+r.W/2, r.Y+r.H/2
			for i, slice := range slices {
				pt.FillPath(chart.Sector(cx, cy, ringHole-ringSink*lit[i], ringRadius+ringGrowth*lit[i], slice), colors[i].Mix(t.Surface, ringFade*(dim-lit[i])))
			}
		})
		ui.Box(c).Absolute().Left(0).Top(0).Size(2*ringRadius, 2*ringRadius).PassThrough().Children(func() {
			if active >= 0 {
				x, y := chart.Polar(0.5, 0.5, 0.4, slices[active].Mid())
				chartTip(c, ui.Rect{W: 2 * ringRadius, H: 2 * ringRadius}, x, 1-y, fmt.Sprintf("%s: %.0f%%", names[active], shares[active]))
			}
		})
	})
}
```

```go test
func TestHoverDonutLightsOneSliceBothWays(t *testing.T) {
	names, shares := []string{"North", "South"}, []float64{60, 40}
	colors := []ui.Color{ui.Hex("#0A5FFE"), ui.Hex("#70A9FE")}
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(40).Children(func() { hoverDonut(c, names, shares, colors) })
	}, 500, 300)
	tt.Advance(time.Second)
	r, ok := tt.Find("Share by name")
	if !ok {
		t.Fatal("no donut")
	}
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	tt.Move(cx+70, cy)
	if !tt.HasText("North: 60%") {
		t.Errorf("no tip over the slice: %q", tt.Texts())
	}
	tt.Advance(time.Second)
	for range 10 {
		tt.Advance(50 * time.Millisecond)
		if !tt.HasText("North: 60%") {
			t.Fatal("the tip goes as it sits by the pointer")
		}
	}
	tt.Move(cx, cy)
	tt.Advance(time.Second)
	if tt.HasText("North: 60%") {
		t.Error("the slice stays lit over the hole")
	}
	row, ok := tt.Find("South")
	if !ok {
		t.Fatal("no legend row")
	}
	tt.Move(row.X+row.W/2, row.Y+row.H/2)
	if !tt.HasText("South: 40%") {
		t.Errorf("the legend does not light its slice: %q", tt.Texts())
	}
	tt.Move(0, 0)
	tt.Advance(time.Second)
	if tt.FramePending() {
		t.Error("keeps asking for frames at rest")
	}
}

func TestHoverDonutHidesUnderAnOverlay(t *testing.T) {
	names, shares := []string{"North", "South"}, []float64{60, 40}
	colors := []ui.Color{ui.Hex("#0A5FFE"), ui.Hex("#70A9FE")}
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(40).Children(func() { hoverDonut(c, names, shares, colors) })
		ui.Overlay(c, func() {
			ui.Box(c).Absolute().Left(0).Top(0).Size(500, 300).Background(ui.RGB(255, 255, 255))
		})
	}, 500, 300)
	tt.Advance(time.Second)
	tt.Move(40+ringRadius+70, 40+ringRadius)
	tt.Advance(time.Second)
	if tt.HasText("North: 60%") {
		t.Error("the donut answers the pointer under an overlay")
	}
}
```

## 쉬는 화면

- `Mix(o, 0)`, `Mix(o, 1)` 은 끝값과 정확히 같고 `Alpha(0)` 은 아무것도 칠하지 않는다. 그래서 스칼라 방식은 쉬는 화면의 픽셀을 바꾸지 않는다.
- 호버 표면을 넣으려고 패딩이나 여백을 바꾸면 레이아웃이 달라진다. 패딩을 늘린 만큼 여백을 줄여 쉬는 화면을 맞춘 뒤 픽셀로 확인한다.
- 작업 전에 쉬는 화면을 찍어 두고 작업 후 `scripts/pixdiff.go` 로 비교한다. 차이가 나면 모두 설명할 수 있어야 한다.
- `tt.FramePending()` 이 쉬는 상태에서 true면 프레임을 계속 요청하는 것이다. 호버 중에도, 호버가 끝난 뒤에도 false여야 한다.

## 줄인 움직임

데스크톱이 움직임을 줄여 달라고 하면(`c.Preferences().ReduceMotion`) MyGo가 대부분 처리한다.

- 값과 포즈는 목표로 바로 간다. `Keyframes` 는 마지막 값으로 끝나고, `Timeline` 은 쉴 자리에 바로 간다.
- 위치와 크기 전환은 없어지고, 색과 페이드는 그대로 부드럽게 바뀐다. 접히지 않는 `Enter`/`Exit` 는 제자리에서 페이드된다.
- `Loop` 는 계속 돈다.
- `c.Now()` 로 직접 계산한 움직임은 직접 확인해야 한다.

`tt.SetPreferences(ui.Preferences{ReduceMotion: true})` 로 테스트한다.

## 넣지 않는 곳

- 많은 행을 담은 표와 목록의 호버 말고는 행마다 움직임을 넣지 않는다.
- 초당 여러 번 누르는 키보드 동작에 긴 전환을 넣지 않는다.
- 의미 없는 대기, 등장 연출, 반복 재생을 넣지 않는다.
- 같은 정보를 두 가지 움직임으로 동시에 알리지 않는다.

## 비용

- `AnimatePose`, `ContentTransition`, `Draw` 안에서 `p.AnimationFrame()` 으로 만든 움직임은 칠하기만 한다. `AnimateTo` 값을 뷰가 읽으면 움직이는 동안 매 프레임 뷰가 돈다.
- 큰 화면에서 같은 값으로 움직이는 요소가 많으면 포즈를 부모 하나에 건다.
- 실측 기준으로 대시보드는 포인터를 3초 동안 253번 옮기는 동안에도 프레임 평균 3.5ms였다. 뷰 빌드, 레이아웃, 칠하기는 각각 평균 0.5ms 이하였고 나머지는 대부분 화면 표시 시간이었다. 16.7ms를 넘긴 프레임은 첫 프레임을 포함해 4개였다. 새 작업은 `MYGO_FRAME_STATS=all` 로 재서 이 수준과 비교한다.
