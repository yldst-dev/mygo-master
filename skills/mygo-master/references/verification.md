# 검증: 정확하고 꼼꼼하게 끝내는 법

테스트가 통과했다고 끝난 것이 아니다. 층마다 잡는 것이 다르므로 변경의 크기에 맞게 층을 고른다. 실제로 돌리지 않은 층은 돌린 것처럼 말하지 않는다. 꼼꼼함은 많이 돌리는 것이 아니라 빠뜨리지 않는 것이다. 같은 움직임을 환경만 바꿔 수십 장 찍지 말고, 위험한 조합을 골라 본다.

목차: [어느 층을 돌릴지](#어느-층을-돌릴지) · [작업 전 기준선](#작업-전-기준선) · [정적 검사](#정적-검사) · [Tester 테스트](#tester-테스트) · [변이 검사](#변이-검사) · [필름스트립](#필름스트립) · [쉬는 화면 비교](#쉬는-화면-비교) · [기준 이미지 비교](#기준-이미지-비교) · [실제 창](#실제-창) · [프레임 통계](#프레임-통계) · [교차 컴파일](#교차-컴파일) · [릴리스 빌드](#릴리스-빌드) · [정리](#정리) · [상태 점검표](#상태-점검표) · [보고](#보고)

아래에서 `$SKILL` 은 이 스킬 폴더, `$WORK` 는 세션의 임시 폴더다. 둘 다 `export` 해 둔다. 검증 도구와 결과물은 저장소에 넣지 않는다. `go test -overlay` 와 `go build -overlay` 로 임시 파일을 끼워 넣으면 저장소가 그대로 유지된다.

## 어느 층을 돌릴지

| 층 | 잡는 것 | 못 잡는 것 | 언제 |
|---|---|---|---|
| 정적 검사 | 형식, 빌드 값 저장 실수 | 동작 | 항상 |
| 최신 버전 확인 | 낡은 API, 놓친 새 기능 | 동작 | 작업 시작할 때 |
| 의존성 방향 | 안쪽이 바깥을 아는 실수 | 층 안의 설계 | 클린 아키텍처 앱이면 항상 |
| Tester 테스트 | 동작, 키보드, 끝값, 프레임 요청 | 눈으로 봐야 아는 것 | 항상 |
| 변이 검사 | 힘이 없는 테스트 | 테스트가 없는 동작 | 동작을 새로 넣었을 때 |
| 필름스트립 | 겹침, 튐, 잘림, 칠하는 순서, 꼬리 | 실제 창 차이 | 움직임을 넣었을 때 |
| 쉬는 화면 비교 | 의도하지 않은 화면 변화 | 움직임 | 기존 화면을 건드렸을 때 |
| 기준 이미지 비교 | 시안과의 차이 | 시안에 없는 상태 | 시안이 있을 때 |
| 실제 창 | Metal 표시, 입력 경로, 접근성 트리, 창 포커스 | 실제 마우스의 감각 | 창에 보이는 변경이면 |
| 프레임 통계 | 느린 프레임 | 감각 | 움직임을 넣었을 때 |
| 누를 곳 크기 | 28 DIP보다 작은 버튼 | 모양의 세련됨 | 화면을 만들거나 바꿨을 때([design.md](design.md#크기를-테스트로-지킨다)) |
| 네이티브 확인 | 웹뷰, 웹 프로세스, 열린 소켓 | 화면 품질 | 새 앱의 릴리스 빌드([design.md](design.md#완전한-네이티브)) |
| 교차 컴파일 | 다른 플랫폼 빌드 오류 | 다른 플랫폼 실행 | 항상 |
| 릴리스 빌드 | 패키징, 실제 앱 실행, 쉬는 CPU | 나머지 동작 | 요청자가 앱을 볼 때 |

## 작업 전 기준선

기존 화면을 건드리기 전에 먼저 찍어 둔다. 고친 뒤에는 되돌릴 기준이 없다.

```sh
go test -count=1 ./...
OVERLAY=$(python3 $SKILL/scripts/overlay.py internal/view $WORK/rest $SKILL/assets/film_test.go $WORK/rest_test.go)
REST_DIR=$WORK/rest/base go test -overlay $OVERLAY -run '^TestRestShots$' -count=1 ./internal/view
```

첫 번째 인자는 테스트를 넣을 패키지 폴더다. 클린 아키텍처 앱은 `internal/view`, 평평한 앱은 `.` 이다. 패키지 이름은 스크립트가 맞춰 준다. `rest_test.go` 는 프로젝트마다 쓴다. 모든 페이지를 라이트와 다크, 기본 크기와 최소 크기로 찍는다.

```go test
func TestRestShots(t *testing.T) {
	directory := os.Getenv("REST_DIR")
	if directory == "" {
		t.Skip("set REST_DIR to write resting screens")
	}
	for _, shot := range []struct {
		name          string
		width, height int
	}{{"default", 1440, 900}, {"smallest", 1024, 680}} {
		a := newApp(func(mygo.ThemeSource) {})
		saveRest(t, directory, shot.name, a.view, shot.width, shot.height)
	}
}
```

## 정적 검사

```sh
gofmt -l .
go vet ./...
go tool mygo vet .
go tool mygo migrate-ui .
CGO_ENABLED=0 go test -count=1 -timeout 120s ./...
```

`gofmt -l` 과 `migrate-ui` 미리 보기는 출력이 없어야 한다. `migrate-ui` 가 파일을 나열하면 예전 API 형태가 남아 있다는 뜻이다. `-timeout` 은 멈춘 테스트를 스택과 함께 실패로 바꿔 준다. `go test` 가 CPU 0%로 멈추면 거의 언제나 메인 스레드로 넘어가는 호출 때문이다. `references/pitfalls.md` 를 본다.

## Tester 테스트

`ui.NewTester` 는 창 없이 뷰를 돌린다. 시계가 멈춰 있어 `tt.Advance` 만큼만 시간이 간다. `time.Sleep` 을 쓰지 않는다. 그래서 결과가 기계 속도와 무관하다.

- 찾기: `tt.Find`, `tt.HasText`, `tt.Texts`, `tt.Focused`. 글자가 없는 요소는 `Label` 로 찾는다.
- 행동: `tt.Click`, `tt.Move`, `tt.Press`, `tt.Release`, `tt.Key`, `tt.Type`, `tt.Scroll`, `tt.Wheel`.
- 환경: `tt.SetSize`, `tt.SetScale`, `tt.SetDark`, `tt.SetPreferences`, `tt.SetFocused`.
- 픽셀: `tt.Image` 의 `RGBAAt`. 기본 배율 1에서는 DIP와 픽셀 좌표가 같다.
- 입력과 `tt.Frame` 은 뷰가 요청한 프레임을 최대 20개 돌린다. 값을 읽기 전에 `tt.Advance(time.Second)` 로 끝까지 보낸다.
- 포즈 프레임은 칠하기만 하므로 `tt.Find` 의 값이 갱신되지 않는다. 포즈 결과는 픽셀로 비교하거나 `tt.Frame` 뒤에 읽는다.

움직임은 세 지점으로 확인한다. 시작과 중간과 끝이 다르고 순서가 맞으면 뛰지 않고 움직인 것이다. 끝에서는 프레임을 더 요청하면 안 된다.

```go test
func TestMovesThroughTheMiddleAndRests(t *testing.T) {
	selected := 0
	labels := []string{"One", "Two", "Three"}
	tt := ui.NewTester(func(c *ui.Context) { monthPills(c, &selected, labels) }, 320, 120)
	tt.Advance(time.Second)
	r, ok := tt.Find("Three")
	if !ok {
		t.Fatal("no segment")
	}
	accent := ui.LightTheme().Accent
	blue := func(label string) bool {
		box, _ := tt.Find(label)
		p := tt.Image().RGBAAt(int(box.X+box.W/2), int(r.Y-4))
		return p.R == accent.R && p.G == accent.G && p.B == accent.B
	}
	if !blue("One") {
		t.Fatal("the pill does not start on the first label")
	}
	if err := tt.Click("Three"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(60 * time.Millisecond)
	if blue("One") || blue("Three") {
		t.Error("the pill jumps instead of sliding")
	}
	tt.Advance(time.Second)
	if !blue("Three") || tt.FramePending() {
		t.Errorf("after the slide: pending %v", tt.FramePending())
	}
}
```

줄인 움직임은 끝값으로 바로 가는지 확인한다.

```go test
func TestSideNavJumpsWithLessMotion(t *testing.T) {
	selected := 0
	tt := ui.NewTester(func(c *ui.Context) { sideNav(c, &selected, []string{"One", "Two"}) }, 240, 200)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	tt.Advance(time.Second)
	if err := tt.Click("Two"); err != nil {
		t.Fatal(err)
	}
	if tt.FramePending() {
		t.Error("the pill still moves while the desktop asks for less motion")
	}
}
```

꼭 넣는 테스트는 이렇다.

- 쉬는 상태, 호버 중, 호버가 끝난 뒤에 `FramePending` 이 false다.
- 호버가 끝나면 픽셀이 원래대로 돌아온다.
- 포인터로 되는 일이 키보드로도 된다. Esc로 닫으면 포커스가 돌아온다.
- 데이터의 양 끝(첫 항목, 마지막 항목, 가장 긴 글자)이 경계 안에 있다.
- 요소가 포인터 아래로 와도 반응이 끊기지 않는다.

## 변이 검사

테스트가 코드를 정말 지키는지 코드를 일부러 망가뜨려 확인한다. 새 테스트마다 그 테스트가 지키려는 줄을 하나 이상 망가뜨려 본다. 숫자를 채우려고 변이나 테스트를 늘리지 않는다. 동작 하나에 테스트 하나, 변이 한두 개면 충분하다.

```json
[
  {"file": "pages.go", "old": "}).IgnoreInput(true).Transition(markMotion)", "new": "}).Transition(markMotion)", "test": "TestChartTipLetsThePointerThrough"},
  {"file": "pages.go", "old": "min(max(x, reach), 1-reach)", "new": "min(x, 1-reach)", "test": "TestChartTipStaysInsideItsChart"}
]
```

```sh
python3 $SKILL/scripts/mutate.py . $WORK/mutations.json
```

`caught` 면 테스트가 잡은 것이고, `MISSED` 면 그 테스트에 힘이 없는 것이니 테스트를 고친다. `BROKEN` 은 변이가 컴파일되지 않거나 테스트 이름이 틀린 경우라 아무것도 증명하지 못한다. 변이는 저장소를 건드리지 않는다. 테스트가 고친 파일과 다른 패키지에 있으면 `"package": "internal/view"` 를 더한다. `go:embed` 나 `os.ReadFile` 로 소스를 직접 읽는 테스트는 `-overlay` 를 보지 못하므로 임시 복사본을 실제로 고쳐 확인한다.

`ElementTransition` 에서 `Position: true` 만 지우면 셋 다 꺼진 상태가 되어 오히려 전부 애니메이션된다. 이런 변이는 `Colors: true` 로 바꿔야 의미가 있다.

## 필름스트립

움직임은 눈으로 봐야 안다. 40ms 간격의 프레임을 한 장에 모아 직접 본다.

```go test
func TestFilmSoftButton(t *testing.T) {
	directory := os.Getenv("FILM_DIR")
	if directory == "" {
		t.Skip("set FILM_DIR to write filmstrips")
	}
	tt := ui.NewTester(func(c *ui.Context) { softButton(c, "Save") }, 320, 200)
	tt.Advance(time.Second)
	r, ok := tt.Find("Save")
	if !ok {
		t.Fatal("no button")
	}
	strip := newFilm(t, tt, 2, r, 8)
	strip.grab()
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	strip.grab()
	strip.play(40*time.Millisecond, 6)
	strip.save(filepath.Join(directory, "hover.png"))
}
```

```sh
OVERLAY=$(python3 $SKILL/scripts/overlay.py internal/view $WORK/film $SKILL/assets/film_test.go $WORK/scenes_test.go)
FILM_DIR=$WORK/film/out go test -overlay $OVERLAY -run '^TestFilm' -count=1 ./internal/view
```

저장한 PNG를 Read 도구로 연다. 등장, 이동, 퇴장, 중간에 방향이 바뀌는 경우를 각각 찍고, 최소 창 크기와 다크 모드도 찍는다. 볼 것은 이렇다.

- 첫 프레임이 비거나 튀지 않는가
- 옛 글자와 새 글자가 겹쳐 흐리게 보이는 시간이 길지 않은가
- 꼬리(거의 끝났는데 계속 움직이는 시간)가 길지 않은가
- 글자나 그림자가 잘리지 않는가
- 뒤에 있어야 할 것이 앞에 칠해지지 않는가, 경계 밖으로 나가지 않는가
- 강조된 것과 흐려진 것이 헷갈리지 않는가

실제 창에서 찍은 PNG 여러 장은 `go run $SKILL/scripts/montage.go 결과.png x y w h 열수 파일...` 로 모은다.

## 쉬는 화면 비교

```sh
REST_DIR=$WORK/rest/now go test -overlay $OVERLAY -run '^TestRestShots$' -count=1 ./internal/view
go run $SKILL/scripts/pixdiff.go $WORK/rest/base $WORK/rest/now
```

움직임만 넣었다면 모두 `identical` 이어야 한다. 차이가 나오면 좌표 범위로 원인을 찾고, 의도한 변경이면 보고에 적는다.

## 기준 이미지 비교

시안을 받았으면 맞춘다.

1. 시안과 같은 크기와 배율로 `ui.Render(view, width, height, scale)` 를 찍는다.
2. 색은 눈대중으로 고르지 않고 시안의 픽셀에서 뽑아 팔레트에 넣는다. 크기는 잘라서 나란히 놓고 잰다. 20.2 같은 소수 글자 크기도 시안에 맞으면 쓴다.
3. 시안의 글자와 숫자는 테스트로 고정한다.
4. `go run $SKILL/scripts/pixdiff.go 시안.png 결과.png` 의 평균 차이를 낮춘다. 대시보드는 평균 2.19%에서 끝냈다. 숫자만 믿지 말고 같은 부분을 `montage.go` 로 잘라 나란히 본다.
5. 시안에 없는 페이지, 다크 모드, 좁은 레이아웃, 임의로 채운 데이터는 구분해서 보고한다.

## 실제 창

Tester는 실제 창의 표시, 시스템 입력 경로, 창 포커스를 보지 못한다. 자동화 플러그인을 넣은 시험 빌드를 띄워 돌린다. 저장소는 그대로 유지된다.

```sh
python3 $SKILL/scripts/drive.py build . $WORK/app-e2e
python3 $SKILL/scripts/drive.py start $WORK/app-e2e $WORK/run
```

시나리오는 Python으로 쓴다. 소켓은 작업 폴더 기준 상대 경로로 연다. 경로가 길면 `bind: invalid argument` 가 난다.

```python
import os, sys
sys.dont_write_bytecode = True
sys.path.insert(0, os.path.join(os.environ["SKILL"], "scripts"))
from drive import Client

os.chdir(os.environ["RUN_DIR"])
c = Client("e2e.sock")
c.hover("Monthly revenue", 0.4, 0.5)
c.settle()
print([n.get("label") for n in c.role("tooltip")])
c.shot("hover.png")
c.move(900, 40)
c.settle()
print(len(c.role("tooltip")))
c.click("Orders")
c.settle()
print([n.get("label") for n in c.role("heading")])
```

```sh
python3 $SKILL/scripts/drive.py tree $WORK/run/e2e.sock
python3 $SKILL/scripts/drive.py shot $WORK/run/e2e.sock $WORK/shot.png
kill <pid>
```

알아 둘 점은 이렇다.

- 시험 창이 요청자 화면에 뜬다. 요청자가 작업 중이면 짧게 끝내고, 끝나면 바로 닫는다.
- `kill` 은 앱을 정상 종료시키고 소켓도 지운다. 시험 빌드의 창 상태는 `~/Library/Application Support/<실행 파일 이름>` 에 따로 저장되므로 실제 앱의 창 위치를 건드리지 않는다. 끝나면 이 폴더를 지운다.
- `idle` 은 방금 보낸 포인터 이동의 애니메이션이 시작되기 전에 돌아올 수 있다. `settle` 은 그래서 먼저 잠깐 기다린다.
- 실행 직후에는 `idle` 이 `no_window` 를 낼 수 있다. `Client` 가 창이 생길 때까지 다시 시도한다.
- 포커스에 기대는 동작은 창이 키 윈도우여야 한다. `Focused` 와 `FocusVisible` 이 창 포커스를 함께 보기 때문이다. 합성 Enter와 글자 입력은 키 윈도우가 아니어도 들어갈 수 있으므로, 결과가 이상하면 `windows` 응답의 `focused` 를 확인한다.
- 실제 마우스가 시험 창 위를 지나가면 호버가 섞인다. 말풍선이 두 개로 보이면 이것부터 의심하고 마우스를 치운 채 다시 돌린다.
- 합성 이벤트로 확인한 것이므로 실제 마우스의 감각은 확인한 것이 아니다.

## 프레임 통계

`drive.py start` 는 `MYGO_FRAME_STATS=all` 로 띄운다. 포인터를 몇 초 동안 계속 움직인 뒤 로그를 읽는다.

```sh
python3 $SKILL/scripts/frame_stats.py $WORK/run/e2e.log
```

60Hz 예산은 16.7ms, 120Hz 예산은 8.3ms다. 첫 프레임은 원래 느리다. 예산을 넘는 프레임이 있으면 어느 단계에서 걸렸는지 본다. `present` 나 `other` 만 튀면 시스템 쪽 지연일 수 있고, `build` 나 `layout` 이 계속 크면 뷰를 줄인다. 프레임 수를 줄이려면 포즈를 먼저 검토한다.

## 교차 컴파일

```sh
for target in windows/amd64 windows/arm64 linux/amd64 linux/arm64 darwin/amd64; do
  GOOS=${target%/*} GOARCH=${target#*/} go build -o /dev/null ./... || echo "FAILED $target"
done
```

`-o /dev/null` 을 빼면 실행 파일이 프로젝트에 떨어진다. 컴파일만 한 플랫폼은 실행을 확인하지 못했다고 보고한다.

## 릴리스 빌드

요청자가 쓸 앱이 돌고 있으면 정상 종료하고 새 빌드로 바꾼다.

```sh
go tool mygo build -skip-dmg
osascript -e 'tell application "Dashboard" to quit'
open build/darwin-arm64/Dashboard.app
```

확인할 것은 프로세스 존재, 화면 캡처, 창 위치와 크기가 그대로인지, 쉬는 CPU가 0% 근처인지(`top -l 2 -pid <pid>`), `~/Library/Logs/DiagnosticReports` 에 충돌 기록이 없는지다. 릴리스 빌드에는 자동화가 빠져 있어 호버를 합성할 수 없다. 접근성 조회나 화면 캡처로 확인한 범위만 말한다.

## 정리

- 시험 빌드의 데이터 폴더(`~/Library/Application Support/<실행 파일 이름>`)를 지운다.
- `git status` 로 의도하지 않은 파일이 없는지 본다. 특히 `go build` 가 남긴 실행 파일을 찾는다.
- `pgrep -fl` 로 시험 빌드와 `go tool mygo dev` 의 자식 프로세스가 남아 있지 않은지 본다.
- 코드에 임시 출력이나 주석이 남지 않았는지 본다.

## 상태 점검표

새로 만든 요소마다 해당되는 칸을 확인한다. 전부 테스트로 만들 필요는 없고 필름스트립이나 실제 창으로 본 것도 된다.

| 묶음 | 항목 |
|---|---|
| 요소 상태 | 쉬는 상태, 호버, 누름, 키보드 포커스 링, 선택, 열림, 비활성 |
| 시간 | 등장, 퇴장, 이동, 중간에 방향 전환, 빠르게 드나들기, 더블클릭, 끝난 뒤 프레임 요청 없음 |
| 환경 | 라이트, 다크, 기본 크기, 최소 크기, 창 크기 변경, 배율 1과 2, 줄인 움직임, 창 포커스 없음 |
| 겹침 | 포인터가 요소 아래로 들어옴, 지운 뒤 다음 요소가 포인터 아래로 옴, 메뉴에 가려짐, 말풍선이 카드 밖으로 나감, 이웃 글자 가림 |
| 데이터 | 첫 값, 마지막 값, 가장 긴 글자, 끊을 곳 없는 글자, 빈 데이터, 큰 수, 끝까지 스크롤한 상태에서 삭제, 깨진 저장 파일 |
| 입력 | 포인터, 키보드, 포커스를 가진 요소가 사라질 때, 보조 기술 역할과 이름, 누를 곳 28×28 이상 |
| 한글 | 입력칸이 있으면 `tt.Compose` 로 조합 상태를 확인한다. 조합 중 Enter는 Tester로 증명되지 않으므로 Enter로 제출하는 칸은 실제 입력기로 확인하거나, 못 했다고 적는다 |

## 보고

보고는 요청자의 전역 형식을 따르고 내용은 이렇게 채운다.

- 무엇이 달라졌는지 먼저 말한다.
- 요청 밖으로 바뀐 것을 빠짐없이 적는다. 예를 들어 Tab 정거장이 늘었거나 내장 위젯을 대체해 역할이 바뀌었거나 계산으로 만든 값이 있으면 적는다.
- 찾아서 고친 문제를 적는다.
- 숫자로 말한다. 테스트 수, 변이 결과, 동일한 화면 수, 프레임 통계.
- 무엇으로 확인했는지 구분한다. Tester, 시험 빌드의 실제 창과 합성 이벤트, 릴리스 빌드는 서로 다르다.
- 확인하지 못한 것을 숨기지 않는다. 실제 마우스 감각, 컴파일만 한 플랫폼, 시스템 줄인 움직임 설정, 원인을 못 찾은 현상이 대표적이다.
- 요청받지 않으면 커밋하지 않고, 커밋하지 않았다고 적는다.
