# MyGo 사용법

MyGo(`github.com/egoist/mygo`)는 Go로 데스크톱 앱을 만드는 프레임워크다. cgo 없이 어느 컴퓨터에서나 모든 플랫폼용으로 빌드하고, 결과는 Go 실행 파일 하나다. 이 문서는 v0.7.4 기준이다.

목차: [버전과 문서 찾기](#버전과-문서-찾기) · [창의 두 종류](#창의-두-종류) · [새 프로젝트](#새-프로젝트) · [앱 골격](#앱-골격) · [뷰 규칙](#뷰-규칙) · [입력과 동작](#입력과-동작) · [레이아웃](#레이아웃) · [스타일과 테마](#스타일과-테마) · [위젯과 베이스](#위젯과-베이스) · [오버레이와 포커스](#오버레이와-포커스) · [차트](#차트) · [아이콘](#아이콘) · [창과 타이틀바](#창과-타이틀바) · [다른 고루틴](#다른-고루틴) · [접근성](#접근성) · [웹 창](#웹-창) · [빌드와 배포](#빌드와-배포) · [문서 지도](#문서-지도)

## 버전과 문서 찾기

요청자의 기본 지시다. 항상 최신 릴리스의 MyGo로 짠다. MyGo는 며칠 간격으로 새 버전이 나오는 0.x 프레임워크라 기억 속 API는 금방 낡는다. 실제로 v0.6.3에서 v0.7.4까지 9일이 걸렸다.

```sh
go list -m -json github.com/egoist/mygo@latest
grep egoist/mygo go.mod
gh release list -R egoist/mygo --limit 5
gh release view v0.7.4 -R egoist/mygo
```

첫 줄의 `Version` 이 최신이고 `Dir` 이 그 버전의 소스 경로다. 그 아래 `docs/` 에 전체 문서가, `examples/` 에 실행 가능한 예제가, `ui/doc.go` 에 뷰 모델 설명이 있다. 릴리스 노트는 내 프로젝트 버전보다 새로운 것을 모두 읽는다. API 하나를 확인할 때는 이렇게 한다.

```sh
go doc github.com/egoist/mygo/ui.ElementTransition
go doc -all github.com/egoist/mygo/ui.Tester | grep '^func'
```

- **최신은 태그가 붙은 릴리스다.** GitHub `main` 브랜치와 웹사이트 문서에는 아직 릴리스되지 않은 기능이 있을 수 있다. 소스와 문서는 모듈 캐시의 그 버전에서 읽는다.
- **직접 짜기 전에 찾아본다.** 새 버전에 필요한 기능이 들어왔을 수 있다. 이 스킬에서 내장 위젯에 움직임을 넣을 자리가 없다고 적은 곳도 해당 버전에서 다시 확인한다.
- **이 스킬과 소스가 다르면 소스를 따른다.** 최신 릴리스가 v0.7.4보다 새로우면 `python3 scripts/check_snippets.py` 로 레시피가 아직 도는지 먼저 확인하고, 깨진 곳은 새 API로 고쳐 쓴 뒤 스킬도 고치자고 보고한다.

### 기존 프로젝트를 최신으로 올리기

이 지시는 요청자가 직접 한 것이므로 뒤처진 프로젝트는 기능을 고치기 전에 올린다. 요청자가 버전을 고정하라고 했거나 `go.mod` 에 `replace` 가 있으면 올리지 않는다. 기능 변경과 섞이지 않도록 먼저 끝내고, 바뀐 버전과 고친 코드를 보고한다.

```sh
go get github.com/egoist/mygo@latest
go get -tool github.com/egoist/mygo/cmd/mygo@latest
go mod tidy
go tool mygo migrate-ui .
go tool mygo migrate-ui -write .
go tool mygo install-skills .
go tool mygo vet .
go test -count=1 -timeout 120s ./...
```

`migrate-ui` 는 처음에 미리 보기만 한다. 결과를 읽고 적용한다. 파일이 하나도 나오지 않으면 최신 형태라는 뜻이므로 평소 정적 검사에도 넣는다. `0.x` 에서는 작은 버전 번호가 바뀌어도 호환이 깨질 수 있으므로 `docs/ui/migration.md` 와 릴리스 노트를 읽는다. 올린 뒤 테스트가 깨지면 고치되, 단순 수정을 넘어 작업 범위를 크게 벗어나면 멈추고 보고한다. 쉬는 화면도 비교한다.

프로젝트의 `.agents/skills/mygo-maintenance/SKILL.md` 는 `mygo init` 과 `install-skills` 가 넣어 주는 공식 지침이므로 있으면 함께 따른다.

## 창의 두 종류

- **네이티브 UI**: `WindowOptions.Content` 에 `ui.View(a.view)` 를 준다. HTML도 웹뷰도 없고 MyGo가 GPU로 직접 그린다. macOS는 Metal, Windows는 Direct3D 11, Linux는 OpenGL이다. 시스템 위젯을 쓰는 것이 아니라 직접 그리되 접근성 트리는 제공한다.
- **웹 페이지**: `WindowOptions.URL` 을 주면 시스템 웹뷰(WKWebView, WebKitGTK, WebView2)에 프런트엔드를 띄운다.

이 스킬의 모션과 검증 방법은 네이티브 UI에서 나왔다. 요청이 분명히 웹이 아니면 네이티브 UI로 만든다.

## 새 프로젝트

새 앱은 이 스킬 폴더(`$SKILL`)의 스크립트로 시작한다. 최신 릴리스로 `mygo init -template native` 를 돌리고, 클린 아키텍처 시작 템플릿과 모션 값을 넣고, 전체 검사까지 돌린다.

```sh
python3 $SKILL/scripts/new_app.py my-app --name "My App" --module example.com/my-app
cd my-app
go tool mygo dev
```

`--name` 은 영문으로 쓴다. 한글 이름은 macOS의 서명 단계에서 빌드가 실패한다. 창 제목과 화면 글자는 한글로 써도 된다. 폴더는 비어 있거나 없어야 한다. `.git` 만 있어도 `mygo init` 이 거부하므로 저장소 안에서는 하위 폴더에 만든다. `go.mod` 가 요구하는 Go 버전은 `GOTOOLCHAIN=auto` 면 Go가 알아서 내려받는다. 네이티브 UI는 Bun이 필요 없다.

```text
my-app/
├── main.go                  조립과 창 옵션
├── architecture_test.go     의존성 방향 검사
├── internal/
│   ├── domain/              엔티티와 업무 규칙
│   ├── usecase/             앱이 하는 일과 포트
│   ├── adapter/             포트 구현(파일, DB, HTTP)
│   └── view/                화면, 팔레트, 모션, 아이콘
├── go.mod                   MyGo와 CLI 버전을 고정
├── mygo.json                앱 이름, 식별자, 버전
├── resources/icon.png
└── .agents/skills/mygo-maintenance/
```

층의 규칙은 [architecture.md](architecture.md)에 있다. 만들어진 할 일 예제는 구조를 보여 주는 용도이므로 실제 도메인으로 바꾸고 남기지 않는다.

`go tool mygo dev` 는 Go 파일을 고치면 다시 빌드하고 재시작한다. 컴파일에 실패하면 이전 빌드가 계속 돈다.

## 앱 골격

MyGo의 가장 작은 형태다. 일회성 데모가 아니면 위의 구조로 나누고, 이 코드는 동작 원리를 보는 데 쓴다. 상태는 앱 타입 하나에 두고, 뷰는 그 상태의 함수로 만든다. `mygo.Theme.SetSource` 처럼 메인 스레드로 넘어가는 호출은 `go test` 에서 메인 루프가 없어 영원히 멈춘다. 그래서 함수로 주입한다.

```go
type app struct {
	count         int
	dark          bool
	setAppearance func(mygo.ThemeSource)
}

func newApp(setAppearance func(mygo.ThemeSource)) *app {
	return &app{setAppearance: setAppearance}
}

func (a *app) view(c *ui.Context) {
	ui.Column(c).Fill().Center().Gap(12).Children(func() {
		ui.Text(c, fmt.Sprint(a.count)).FontSize(40).Bold()
		if ui.PrimaryButton(c, "Increment").Clicked() {
			a.count++
		}
		ui.Switch(c, &a.dark).Label("Dark mode").OnChange(func() {
			if a.dark {
				a.setAppearance(mygo.ThemeDark)
			} else {
				a.setAppearance(mygo.ThemeLight)
			}
		})
	})
}

func main() {
	a := newApp(mygo.Theme.SetSource)
	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{
			Title:     "Counter",
			Width:     480,
			Height:    360,
			MinWidth:  320,
			MinHeight: 280,
			StateKey:  "main",
			Content:   ui.View(a.view),
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
```

```go test
func TestCounterClicksAndTheme(t *testing.T) {
	var applied mygo.ThemeSource
	a := newApp(func(source mygo.ThemeSource) { applied = source })
	tt := ui.NewTester(a.view, 480, 360)
	if err := tt.Click("Increment"); err != nil {
		t.Fatal(err)
	}
	if a.count != 1 || !tt.HasText("1") {
		t.Errorf("count %d, texts %q", a.count, tt.Texts())
	}
	if err := tt.Click("Dark mode"); err != nil {
		t.Fatal(err)
	}
	if applied != mygo.ThemeDark {
		t.Errorf("appearance %q", applied)
	}
}
```

## 뷰 규칙

- 뷰는 프레임마다 메인 스레드에서 돈다. 한 프레임에 여러 번 돌 수도 있다. 뷰 안에서 파일을 읽거나 요청을 보내거나 작업을 시작하지 않는다. 그런 일은 처리된 동작이나 고루틴에서 한다.
- `ui.Element` 는 그 빌드 패스에서만 유효하다. 앱 상태나 고루틴에 저장하지 않는다. 오래 남겨야 하면 `ui.Handle` 을 쓴다. `go tool mygo vet .` 이 이 실수를 찾아 준다.
- `*ui.Context` 는 창마다 하나로 고정이다. `Children` 안에서는 현재 부모만 바뀐다. 상태에 넣지 말고 도우미에 넘긴다. 사이 시간에 쓸 콜백은 `c.Services()` 를 잡아 둔다.
- 요소의 상태(포커스, 스크롤, 애니메이션)는 형제 사이의 순서나 `Key` 를 따른다. 조건부로 생기는 형제나 순서가 바뀌는 목록에는 키를 준다.
- 키는 생성자에 넘긴다. `ui.Button(c.Key(id), "Save")`, `ui.Row(c.Key(id))` 처럼 쓴다. `ui.Row(c).Key(id)` 는 v0.7.4의 `migrate-ui` 가 이 형태로 바꾸는 예전 형태다. `c.Key` 는 다음 요소 하나에만 걸리고, 소비되지 않으면 그 다음 요소로 넘어간다.
- 같은 부모 아래 같은 키가 둘이면 상태를 공유한다. 앱에서는 로그만 남고 Tester에서는 패닉이 난다.
- 요소마다 값을 기억하려면 `ui.Local(e, "key", init)` 을 쓴다.
- `Draw` 는 칠하기만 한다. 여러 번 불릴 수 있고 빌드보다 늦게 불린다. 상태를 바꾸지 않는다.
- 불변식을 지키고 싶은 값(모델, `ListState`, `Router`)은 앱 필드에 두고 뷰에서 만들지 않는다.

## 입력과 동작

- `Clicked` 는 그 자리에서 묻는다. 만드는 중인 목록을 바꾸는 동작은 `OnClick`, `OnChange` 나 지역 변수로 미뤄 반복이 끝난 뒤 처리한다.
- 값 위젯은 포인터로 묶는다. `ui.TextInput(c, &a.query)` 는 생성자가 돌아올 때 이미 값이 바뀌어 있다. `Changed` 와 `Submitted` 로 변화를 확인한다.
- 조회 함수: `Hovered`, `Pressed`, `Clicked`, `DoubleClicked`, `RightClicked`, `Dragged`, `PointerPosition`, `Focused`, `FocusVisible`, `FocusWithin`.
- `Hovered` 는 위에 덮인 요소를 고려한다. `PointerPosition` 은 기하만 보고, 호출한 요소는 포인터가 움직일 때마다 프레임을 받는다. `Focused` 는 창에 포커스가 있어야 true다.
- `PassThrough` 는 그 요소만 포인터를 통과시킨다. `IgnoreInput(true)` 는 자식까지 통과시키고 키보드 포커스도 막지만 보조 기술에는 보인다.
- 단축키: `c.Shortcut(ui.Cmd, ui.KeyS)` 는 창 전체, `e.Shortcut(0, ui.KeyEscape)` 는 그 요소 안에 포커스가 있을 때. `ui.Cmd` 는 macOS의 Command, 그 외의 Ctrl이다.
- `FocusGroup(ui.Vertical)` 은 Tab 한 칸에 화살표, Home, End 이동을 준다. `KeepFocus` 는 누를 때 포커스를 옮기지 않는다.

## 레이아웃

CSS flexbox와 거의 같다. 단위는 DIP이다.

- 컨테이너: `ui.Column`, `ui.Box`(세로로 쌓고 폭을 늘림), `ui.Row`(가로로 놓고 세로 중앙), `ui.Grid`(`Columns`, `GapX`, `GapY`), `ui.Scroll`, `ui.Spacer`, `ui.Divider`.
- 크기: `Width`, `Height`, `Size`, `WidthPercent`, `HeightPercent`, `MinWidth`, `MaxWidth`, `MinHeight`, `MaxHeight`, `Fill`, `FillWidth`, `FillHeight`, `Grow`, `Shrink`, `Basis`, `AspectRatio`.
- 간격: `Gap`, `Padding`(값 1, 2, 4개), `PaddingX`, `PaddingY`, `Margin`, `MarginX`, `MarginY`. `Margin` 의 `ui.Auto` 는 남는 공간을 차지한다.
- 정렬: `Justify`(`ui.Start`, `ui.Center`, `ui.End`, `ui.SpaceBetween`), `AlignItems`, `AlignSelf`, `Center`, `Wrap`.
- 위치: `Absolute` + `Top`, `Right`, `Bottom`, `Left`(테두리 안쪽에서 잰다), `Attach(at, self)`, 그리고 다른 요소에 붙이는 `AttachTo(target, at, self)`.
- 잘라내기와 숨기기: `Clip`, `ClipX`, `ClipY`, `Invisible`(자리 유지), `Hidden(true)`(자리도 없앰, 내부 상태 유지).
- 칠하는 순서: 배경, `Draw`, 흐름 안의 자식, `Absolute` 자식, `DrawOver`. z-index는 없다. 위에 올릴 것은 나중에 만들거나 `Absolute` 로 만든다.
- 창 폭에 따라 바꾸려면 `width, _ := c.Size()` 로 분기한다. 최소 창 크기에서 글자가 겹치지 않는지 반드시 확인한다.
- 큰 데이터는 `ui.List`(보이는 행만 빌드)를 쓴다.
- 행 안의 긴 글자가 옆 요소를 밀어내면 글자 쪽에 `MinWidth(0)` 을 준다. flexbox처럼 내용 크기 아래로 줄지 않는 것이 기본이다.
- 레이아웃이 이상하면 `Debug()` 나 인스펙터(개발 빌드에서 F12)로 본다.

## 스타일과 테마

- 스타일: `Background`, `Gradient`, `LinearGradient`, `Border`, `BorderWidth`, `BorderColor`, `Radius`, `Shadow`, `InsetShadow`, `Opacity`, `Blur`, `Cursor`.
- 글자: `FontSize`, `FontWeight`, `Bold`, `TextColor`, `NoWrap`, `SingleLine`, `TextAlign`.
- 색: `ui.Hex`, `ui.RGB`, `ui.RGBA`, `ui.Oklch`, `color.Mix(other, t)`, `color.Alpha(a)`(현재 알파에 곱한다).
- 상태에 따른 모양은 CSS 선택자 대신 코드로 쓴다. 다만 뚝 바뀌지 않도록 `references/motion.md` 의 스칼라 방식을 쓴다.
- 테마: `t := *c.Theme()` 를 복사해서 바꾼 뒤 `c.SetTheme(&t)`. 위젯은 `Background`, `Surface`, `SurfaceHover`, `Border`, `Text`, `TextMuted`, `Accent`, `AccentText`, `Danger`, `Success`, `Focus`, `Radius`, `Spacing` 등을 따른다. 일부 영역만 바꾸려면 `ui.Themed(c, t)`.
- 앱 외관은 `mygo.Theme.SetSource(mygo.ThemeDark)` 로 정한다. `ThemeSystem` 은 시스템을 따른다.
- 글꼴을 따로 지정하지 않으면 시스템 글꼴이다. 새 글꼴은 도입 전에 제안한다.

## 위젯과 베이스

완성 위젯(`ui.Button`, `ui.PrimaryButton`, `ui.TextInput`, `ui.Switch`, `ui.Checkbox`, `ui.Select`, `ui.Tabs`, `ui.Table`, `ui.List`, `ui.SearchField`, `ui.Popover`, `ui.Modal`, `ui.Toast` 등)은 `docs/ui/<이름>.md` 를 읽고 쓴다. 전체 목록은 `docs/ui/README.md` 에 있다.

디자인을 맞추거나 움직임을 넣으려면 베이스를 쓴다. 키보드, 포커스, 접근성은 베이스가 맡고 모양은 전부 직접 정한다.

| 베이스 | 하는 일 |
|---|---|
| `ui.ButtonBase` | 포인터, Enter, Space로 `Clicked` |
| `ui.SwitchBase`, `ui.CheckboxBase`, `ui.ToggleBase` | `*bool` 토글 |
| `ui.SegmentedBase` | `Track` 과 `Segment`, `*int` 선택 |
| `ui.TabsBase` | `List` 와 `Tab`, `*int` 선택 |
| `ui.SelectBase`, `ui.ComboboxBase` | 팝업과 항목 |
| `ui.PopoverBase`, `ui.DialogBase` | 앵커 옆 패널, 배경 위 패널 |
| `ui.TooltipBase` | 내장 지연 규칙을 따르는 말풍선 |
| `ui.TextInputBase`, `ui.TextAreaBase` | 꾸밈 없는 입력칸 |
| `ui.SliderBase`, `ui.RadioBase`, `ui.CollapsibleBase` | 그 외 |

내장 위젯에는 움직임을 넣을 자리가 없는 경우가 있다. `ui.PopoverMenuButton` 의 패널, `ui.SearchField` 의 지우기 버튼이 그렇다. 베이스로 다시 짤 때 잃는 것(검색 필드 역할, 입력 중 첫 글자 찾기 등)을 확인하고 보고한다.

표는 머리글과 행 모양까지 디자인에 맞추려면 `Row` 로 짜고 `ui.RoleTable`, `ui.RoleRow`, `ui.RoleColumnHeader`, `ui.RoleCell` 을 붙인다. 행이 많으면 `ui.Table` 을 쓴다.

## 오버레이와 포커스

`ui.Overlay(c, fn)` 은 모든 것 위에 그린다. `ui.PopoverBase(c, anchor, &open, fn)` 은 바깥 누름과 Esc로 닫고, 닫힐 때 포커스를 돌려준다. 바깥을 누르면 그 누름은 아래 요소까지 전달된다. 전체 동작은 `docs/ui/overlays.md` 에 있고, 움직임이 있는 메뉴는 `references/motion.md` 의 드롭다운 레시피를 쓴다. 네이티브 대화상자는 고루틴에서 부른다.

## 차트

`ui/chart` 는 모양이 없는 부품이다. 좌표는 플롯 기준 0~1이고 y는 아래에서 위로 간다.

- 눈금: `chart.NewLinear`(`Nice`, `Ticks`, `At`, `Span`), `chart.NewBand`(`Band`, `At`, `Ticks`), `chart.NewTime`, `chart.Extent`. 직접 만든 눈금은 `chart.Tick{At, Value, Label}`.
- 플롯: `chart.Plot(c, chart.Axes{Left, Bottom, Label}, func(area ui.Element) {...})`. 높이를 주거나 부모를 채운다.
- 표시: `chart.Rect`, `chart.Point`, `chart.RuleX`, `chart.RuleY`, `chart.Anchor`. 일반 요소이므로 스타일과 `Transition` 을 걸 수 있다.
- 경로: `chart.Line`, `chart.Area`, `chart.CurveLinear`, `chart.CurveMonotone`, `chart.Stack`, `chart.Pie`, `chart.Sector`, `chart.SliceAt`, `chart.Polar`.
- 선택: `chart.Select(c, area, &selected, xs, describe)` 가 포인터와 키보드 선택을 준다. 호버 동작은 `references/motion.md` 의 그래프 호버 레시피를 쓴다.

```go
func barChart(c *ui.Context, labels []string, values []float64) {
	t := c.Theme()
	x := chart.NewBand(labels...).Padding(0.3, 0.15)
	_, hi := chart.Extent(values...)
	y := chart.NewLinear(0, hi).Nice(4)
	ticks := y.Ticks(4)
	chart.Plot(c, chart.Axes{Left: ticks, Bottom: x.Ticks()}, func(area ui.Element) {
		area.Draw(func(p *ui.Painter, r ui.Rect) {
			for _, tick := range ticks {
				p.Line(r.X, chart.Y(r, tick.At), r.X+r.W, chart.Y(r, tick.At), 1, t.Border)
			}
		})
		for i, v := range values {
			chart.Rect(c, x.Band(i), y.Span(0, v)).Radius(4).Background(t.Accent).
				Transition(ui.ElementTransition{Animation: &pillMotion})
		}
	}).Label("Sales by month").Height(200).Gap(6)
}
```

## 아이콘

SVG를 `//go:embed` 로 넣고 `ui.MustParseSVG` 로 읽어 패키지 변수에 둔다. `ui.Icon(c, icon).Size(22, 22).TextColor(color)` 로 쓰면 글자색을 따른다. 외부 아이콘 파일을 복사하면 라이선스 파일도 함께 넣고 보고한다.

```text
//go:embed icons/*.svg
var iconFiles embed.FS
```

## 창과 타이틀바

- `StateKey` 는 창 위치와 크기를 기억한다. 저장소는 앱 이름별로 나뉜다. macOS에서는 `~/Library/Application Support/<앱 이름>/window-state.json` 이다.
- `TitleBarStyle: mygo.TitleBarHidden` 이면 뷰가 타이틀바를 그린다. `c.TitleBar()` 가 창 버튼이 차지하는 공간을 주고, 해당 줄에 `DragWindow()` 를 건다.
- 창을 여러 개 띄워도 `Content` 하나를 공유할 수 있다. 요소 상태는 창마다 따로 있다.
- 메뉴, 트레이, 대화상자, 알림 등은 아래 문서 지도를 본다.

## 다른 고루틴

뷰가 상태를 읽으므로 다른 고루틴에서 상태를 고칠 때는 `Window.Update` 로 메인 스레드에 넘긴다. 뷰 안의 `c.Invalidate()`, `c.After(d)` 는 프레임만 요청한다.

```go
type loader struct {
	win   *mygo.Window
	items []string
	err   error
}

func (l *loader) load(fetch func() ([]string, error)) {
	go func() {
		items, err := fetch()
		l.win.Update(func() { l.items, l.err = items, err })
	}()
}
```

## 접근성

- 글자가 없는 요소에는 `Label` 을 준다. 스크린 리더뿐 아니라 Tester의 `Click`, `Find` 도 이 이름을 쓴다.
- 직접 그린 위젯에는 `Role`, `Expanded`, `Checked`, `Value`, `Range`, `Level`, `Description` 을 준다.
- 포커스를 받는 요소는 `FocusVisible` 일 때 링을 보여 준다. 링은 요소의 `Radius` 를 따르므로 모서리가 둥근 요소에는 반드시 준다.
- 포커스를 옮기지 않는 소식은 `c.Announce(text)` 로 알린다.
- 포인터로 되는 일은 키보드로도 되어야 한다. 호버로만 보이는 값은 포커스에서도 보여 준다.

## 웹 창

이 부분은 문서 기준이며 대시보드에서는 쓰지 않았다. `mygo init my-app`(Bun 필요)이 Vite 프런트엔드를 함께 만든다. Go 쪽에서 `mygo.Bind(service)` 와 `mygo.NewEvent[T](name)` 를 선언하면 `mygo generate` 가 TypeScript 클라이언트를 만든다. `docs/getting-started.md`, `docs/frontend.md`, `docs/bindings.md` 를 읽는다.

## 빌드와 배포

```sh
go tool mygo build -skip-dmg
go tool mygo build -platform darwin/universal,windows/amd64,linux/amd64
```

결과물은 `build/<os>-<arch>/` 에 생긴다. 릴리스 빌드에는 인스펙터와 자동화 플러그인이 빠진다. 서명, 공증, 설치 프로그램은 `docs/distribution.md`, 자동 업데이트는 `docs/updates.md`, CI는 `docs/github-actions.md` 를 본다.

`go build ./...` 를 패키지가 하나인 모듈에서 돌리면 실행 파일이 프로젝트에 떨어진다. 빌드 확인은 `go build -o /dev/null ./...` 로 한다.

## 문서 지도

모듈 경로 아래 `docs/` 기준이다.

| 필요한 것 | 문서 |
|---|---|
| 뷰와 수명 | `ui/views.md`, `ui/doc.go`(`ui/` 소스) |
| 레이아웃, 격자, 스크롤 | `ui/layout.md`, `ui/grid.md`, `ui/scroll.md` |
| 글자, 입력, 한글 등 입력기 | `ui/text.md`, `ui/input.md`, `ui/text-input.md` |
| 스타일, 테마, 데스크톱 설정 | `ui/styling.md` |
| 칠하기, 경로, 캔버스, 셰이더 | `ui/drawing.md`, `ui/shaders.md` |
| 모션, 전환, 변형 | `ui/motion.md`, `ui/transitions.md`, `ui/transforms.md` |
| 페이지 이동 | `ui/navigation.md` |
| 오버레이, 대화상자, 말풍선, 토스트 | `ui/overlays.md`, `ui/dialog.md`, `ui/tooltip.md`, `ui/toast.md` |
| 메뉴 | `ui/menu-button.md`, `ui/context-menu.md`, `menus.md` |
| 목록, 표, 트리 | `ui/list.md`, `ui/table.md`, `ui/tree.md`, `ui/outline.md`, `ui/grid-view.md` |
| 끌어다 놓기 | `ui/drag-and-drop.md` |
| 코드 에디터, 마크다운, 서식 편집기 | `ui/code-editor.md`, `ui/markdown.md`, `ui/rich-text-editor.md` |
| 접근성, 오른쪽에서 왼쪽 | `ui/accessibility.md`, `ui/right-to-left.md` |
| 테스트, 인스펙터, 렌더링, 성능 | `ui/testing.md`, `ui/inspector.md`, `ui/rendering.md`, `ui/performance.md` |
| 창, 앱 수명, 시스템 기능 | `windows.md`, `ui/windows.md`, `app.md`, `native.md` |
| 다국어 | `i18n.md` |
| 설정 파일, CLI | `configuration.md`, `cli.md` |
| 공식 플러그인 | `plugins/*.md`(자동화, 유리 효과, 게임, 터미널, 로딩, SQLite 등) |
| 브라우저에서 실행 | `ui/browsers.md` |

움직임 예제는 `examples/timeline`, `examples/confetti`, `examples/effort-slider`, 전체 위젯 둘러보기는 `examples/gallery` 에 있다. 예제는 모듈 경로에서 `go run ./examples/gallery` 로 볼 수 있다.
