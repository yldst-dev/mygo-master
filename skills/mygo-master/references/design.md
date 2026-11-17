# 디자인: 네이티브로, 세련되고 편하게

이 문서의 값은 [`examples/dashboard`](../examples/dashboard)에서 왔다. 요청자가 "세련되고 부드럽다"고 판단하고 디자인 예시로 정한 대시보드 앱의 전체 소스다. 새 화면은 여기 값에서 시작한다. 시안이 있으면 시안을 잰다. 기존 앱은 그 앱의 팔레트와 크기를 따르되, 새로 넣는 누를 곳은 28 DIP 이상으로 만들고 기존의 작은 버튼은 고치자고 제안한다. 움직임 값은 [motion.md](motion.md)에 있다.

목차: [완전한 네이티브](#완전한-네이티브) · [방향](#방향) · [색](#색) · [글자](#글자) · [간격과 모서리](#간격과-모서리) · [누를 곳의 크기](#누를-곳의-크기) · [딱딱하지 않게](#딱딱하지-않게) · [복잡하지 않게](#복잡하지-않게) · [예시 지도](#예시-지도) · [코드로 가져다 쓰기](#코드로-가져다-쓰기) · [점검표](#점검표)

## 완전한 네이티브

- 앱 화면은 `WindowOptions.Content` 에 `ui.View` 를 주고 MyGo 위젯과 베이스로 그린다. HTML, CSS, JavaScript, 웹뷰(`WindowOptions.URL`, 네이티브 화면 안의 웹뷰)는 쓰지 않는다. 웹 페이지를 띄우는 것 자체가 요청일 때만 웹 창을 쓴다.
- 아이콘은 SVG를 `ui.MustParseSVG` 로 읽는다. MyGo가 Go 코드로 해석해서 그린다. 글자는 시스템 글꼴이고 macOS에서는 CoreText가 그린다.
- 창은 실제 macOS 창이다. 타이틀바를 숨겼다면 `c.TitleBar()` 로 창 버튼 자리를 비우고, 끌 수 있는 줄에 `DragWindow()` 를 건다. 크기와 위치는 `StateKey` 로 기억한다.
- MyGo(v0.9.2에서 확인)는 웹 창 기능 때문에 앱을 켤 때 WebKit을 늘 불러온다(`internal/darwin/objc.go`). 그래서 "WebKit이 없다"고 보고하지 않고 "WebKit 코드가 돌지 않는다"까지만 말한다.

macOS에서는 이렇게 확인한다. 릴리스 빌드를 띄운 뒤 프로세스 번호를 `$PID` 에 넣는다.

```sh
grep -rniE "webview|<html|javascript|URL:" --include='*.go' .
find build -path '*.app/Contents/*' -type f
sample "$PID" 1 -file "$TMPDIR/app-sample.txt" >/dev/null
awk '/^Binary Images:/{exit} {print}' "$TMPDIR/app-sample.txt" | grep -cE "WebKit|JavaScriptCore|WebCore"
ps -axo pid,lstart,comm | grep WebContent
lsof -nP -a -p "$PID" -i
```

- 첫 줄은 아무것도 찾지 않는다.
- 앱 묶음에는 실행 파일, Info.plist, PkgInfo, 서명 파일, 아이콘만 있다.
- 샘플 속 WebKit 프레임은 0이다. 대시보드는 스레드 20개를 1초 동안 샘플링해 0건이었다.
- 앱을 켠 시각 뒤에 생긴 WebContent 프로세스가 없다.
- `lsof` 는 아무것도 내지 않는다. 네트워크 기능을 요청받았다면 그 연결만 있다.

## 방향

대시보드의 방향은 "차분한 업무 도구"다.

- 푸른 기가 조금 도는 밝은 회색 바탕에 강조색은 파랑 하나다.
- 면은 얇은 선과 아주 옅은 그림자로 나누고, 모서리는 중간 정도로 둥글다.
- 여백이 넉넉하고 한 화면에 놓는 것이 적다.
- 상태가 바뀌면 100~240ms 안에 튕기지 않고 움직인다.

세련됨은 꾸밈을 더해서가 아니라 덜어 내고 맞춰서 나온다. 다음은 넣지 않는다. 보라나 남색 그라데이션, 그라데이션 글자, 유리 효과, 빛나는 효과, 카드 안의 카드, 왼쪽 굵은 띠 강조, 진한 그림자와 큰 모서리의 조합, 색칠한 아이콘 타일, 고정폭 글꼴 남발.

## 색

| 역할 | 라이트 | 다크 | 쓰는 곳 |
|---|---|---|---|
| 본문 바탕 | `#F6F7FC`, 아래 절반에서 `#EEF1F8` 로 | `#0B0F1A` | 페이지 |
| 사이드바, 상단 막대 | `#FBFBFE` 에서 `#F7F9FE`, `#F7F8FD` | `#0F1524` | 앱 틀 |
| 카드 | `#FEFEFF` | `#141B2C` | 내용 면 |
| 선 | `#E6E9F0` | `#232C40` | 카드 테두리, 구분선, 표 행 |
| 입력칸 | `#F4F5FA`, 테두리 `#DCE1EB` | `#182034`, `#2A3450` | 검색칸, 보조 버튼 |
| 제목, 값 | `#0B1220` | `#F3F6FC` | 카드 제목, 지표 값 |
| 라벨, 메뉴 | `#434A62` | `#C4CCDD` | 내비 글자, 표 머리 |
| 보조 글자 | `#545B75` | `#8E98AE` | 축, 비교 문구 |
| 강조 | `#0A5FFE` | `#3D7BFF` | 선택 알약, 그래프 선, 포커스 |
| 강조 면(흰 글자 아래) | `#0A5FFE` | `#2F6BF2` | 주요 버튼, 체크, 켜진 스위치, 위젯의 `Theme.Accent` |
| 강조 글자 | `#0A52EE` | `#7CAEFF` | 선택된 내비, 링크 |
| 선택 면 | `#E3EBFC` | `#1B2B4D` | 내비 알약 |
| 호버 면 | `#EEF1F9` | `#18213A` | 버튼, 행 |
| 카드 그림자 | `ui.RGBA(15, 23, 42, 0.05)` | `ui.RGBA(0, 0, 0, 0.25)` | 카드 |

- 강조색은 하나다. 선택 표시, 주요 데이터 선, 포커스 링, 링크에만 쓴다. 면을 강조색으로 채우는 곳은 월 선택 알약처럼 작은 선택 표시뿐이다.
- 글자색은 세 단계(제목, 라벨, 보조)로 위계를 만든다. 강조하려고 글자에 다른 색을 칠하지 않는다.
- 회색은 강조색 쪽으로 아주 조금 기운 차가운 값이다. 순회색을 섞어 쓰지 않는다.
- 그래프는 강조색의 명도 단계와 중립 회색으로 칠한다(`#0A5FFE`, `#70A9FE`, `#C4DAFD`, `#E7E8EA`). 범주마다 다른 색상을 늘어놓지 않는다.
- 상태는 옅은 바탕, 진한 글자, 점 하나로 보여 준다. 완료는 초록, 진행은 파랑, 취소는 빨강이다.
- 다크 모드는 뒤집은 값이 아니라 역할마다 다시 고른 값이다. 강조색은 밝게, 그림자는 진하게 바꾼다.
- 글자 대비는 4.5:1 이상이다. 대시보드에서 잰 조합 가운데 못 미친 것은 둘이다. 라이트의 아바타 글자(강조 글자 `#0A52EE` 가 아바타 바탕 `#C6D9FD` 위에 놓여 4.27:1)와 다크의 강조 면 위 흰 글자(월 선택 알약, 메뉴 강조 행, 3.84:1)다. 선택 면 위 강조 글자는 라이트 5.08:1, 다크 6.24:1이다.
- 흰 글자를 올리는 면은 `accent` 가 아니라 `accentFill` 로 칠한다. 다크에서 `#2F6BF2` 는 4.66:1이다. 시작 템플릿과 위젯의 `Theme.Accent` 는 `accentFill` 을 쓴다. 대시보드의 월 선택 알약과 메뉴 강조 행은 아직 `accent` 를 쓰므로 따라 하지 않는다.
- 팔레트는 역할별 구조체 하나에 라이트와 다크를 두고, 테마를 바꿀 때 모든 칸을 함께 섞는다(`palette.go`, [창 전체 색 전환](motion.md#창-전체-색-전환)). 이 파일은 대시보드와 시작 템플릿이 같은 것을 쓴다([코드로 가져다 쓰기](#코드로-가져다-쓰기)).

## 글자

| 쓰임 | 크기 | 굵기 |
|---|---|---|
| 큰 숫자 | 29.5 | Bold |
| 제목 | 20~21 | Bold |
| 내비 | 16 | 500, 선택 600 |
| 본문, 표, 입력 | 14.5~15.5 | 400~500 |
| 버튼 글자 | 14 | 600 |
| 축, 범례, 상태 | 13.5~14 | 400 |
| 보조, 배지, 말풍선 | 12.5 | 400~600 |

- 글꼴은 시스템 글꼴 하나다. 새 글꼴은 도입하지 않고 제안한다.
- 12.5보다 작은 글자를 쓰지 않는다.
- 굵기는 400, 500, 600, Bold만 쓴다. 선택되면 500에서 600으로 바꾼다.
- 위계는 크기, 굵기, 글자색 세 단계로 만든다.
- 바뀌는 숫자는 `ContentTransition` 으로 굴리고 `NoWrap` 을 건다.

## 간격과 모서리

- 사이드바 폭은 237이고, 앱 이름 줄과 상단 막대의 높이를 둘 다 69로 맞춘다.
- 본문 여백은 위 20, 오른쪽 18, 아래 34, 왼쪽 20이다. 카드 사이는 17~20이다.
- 카드 안 여백은 위 15~17, 오른쪽 14~21, 아래 0~20, 왼쪽 20~22이다. 아래가 0인 카드는 그래프나 도넛이 바닥까지 닿는다. 카드 안 간격은 13~26이다.
- 창 폭 1394 아래에서는 지표 카드 4열을 2열로, 그래프 두 장의 가로 배치를 세로로 바꾼다. 최소 창은 1024×680이다.

모서리는 요소 크기를 따라간다.

| 요소 | 모서리 |
|---|---|
| 카드 | 16 |
| 내비 항목, 계정 버튼 | 12 |
| 입력칸, 아이콘 버튼, 메뉴 패널 | 10 |
| 보조 버튼, 세그먼트 트랙, 메뉴 행 | 8 |
| 트랙 안 알약 | 7, 트랙보다 1 안쪽 |
| 표 머리글, 말풍선 | 6, 5 |
| 알약, 스위치, 아바타 | 높이의 절반 |

안쪽 요소의 모서리는 바깥보다 작게 잡는다. 테두리 1 안의 알약은 바깥보다 1 작고(트랙 8, 알약 7), 메뉴 패널 10 안의 행은 8이다.

## 누를 곳의 크기

작은 버튼은 쓰지 않는다.

- 누를 수 있는 것은 가로, 세로 모두 28 DIP 이상이다. WCAG 2.2의 최소 기준(24×24)보다 넉넉하게 잡는다.
- 아이콘만 있는 버튼은 32×32 이상이고, 안의 아이콘은 18~26이다.
- 자주 누르는 것(내비, 계정, 검색칸, 주요 버튼)은 높이 40 이상이다.
- 버튼 글자는 14 이상이고 좌우 여백은 12~18이다.
- 보이는 모양이 작아야 하면 누를 곳만 키운다. 아래 `clearButton` 처럼 바깥 버튼을 28×28로 잡고 그 안에 작은 모양을 놓는다.
- 표 안의 글자 링크는 셀 높이 전체를 누를 곳으로 잡는다.
- 붙어 있는 버튼 사이는 4 이상 띄우거나, 세그먼트처럼 한 트랙에 묶고 선택 면으로 구분한다.
- 그래프의 점(10×10) 같은 작은 표식을 누를 곳으로 쓰지 않는다. 그래프 영역 전체에서 가장 가까운 값을 고르고(`chart.Select`), 도넛은 조각과 범례 행을 함께 대상으로 삼는다.

대시보드에서 잰 값이다(1440×900, Tester의 `Find`).

| 요소 | 크기 |
|---|---|
| 내비 항목 | 210×45, 간격 4.5 |
| 사이드바 계정 버튼 | 198×46 |
| 검색칸 | 높이 40. 칸 어디를 눌러도 입력에 포커스가 간다 |
| 알림 버튼 | 34×34, 아이콘 26 |
| 다크 모드 스위치 | 58×30 |
| "View All" 버튼 | 높이 30, 좌우 여백 14 |
| 월 세그먼트 | 49.66×28, 트랙 300×30 |
| 지표 카드의 미니 그래프 | 85×45 |
| 매출 그래프 영역 | 640.6×176.5 |
| 도넛 | 182×182, 범례 행 157×20 |

대시보드에도 기준보다 작은 곳이 세 군데 있었고 모두 고쳤다. 예전 소스를 만나면 같은 방법으로 고친다.

| 요소 | 전 | 후 | 고친 법 |
|---|---|---|---|
| 드롭다운 메뉴 행 | 170×27.3 | 170×30 | `MinHeight(30)` |
| 검색칸 지우기 버튼 | 21×21 | 28×28, 보이는 원은 21 | 바깥 버튼을 28로 잡고 안에 원을 놓는다(`clearButton`) |
| 표 안 주문 번호 링크 | 85.9×17.1 | 85.9×37 | 셀 높이 전체를 버튼으로 잡고, 밑줄은 글자를 감싼 상자에 그린다 |

고친 뒤의 크기는 `main_test.go` 의 `TestEveryTargetIsLargeEnough` 가 지킨다. 메뉴를 열었을 때와 검색어를 넣었을 때도 잰다.

### 크기를 테스트로 지킨다

`Find` 는 그 글자를 보여 주거나 그 `Label` 을 단 첫 요소의 상자를 돌려준다. 부모가 자식보다 먼저 잡힌다. 그래서 이렇게 짓는다.

- 누를 수 있는 요소에는 안에 글자가 있어도 `Label` 을 준다. 없으면 안쪽 글자 상자가 잡힌다. 대시보드의 "View All" 과 월 세그먼트도 `Label` 을 달기 전에는 글자 53.6×16.5, 25.5×16.5가 잡혔다.
- 누를 것의 `Label` 은 화면 안에서 겹치지 않게 짓는다. 대시보드의 `Find("Dashboard")` 는 내비 항목이 아니라 앞에 있는 앱 이름 글자 104.1×23.8을 돌려준다.
- 바깥 줄을 눌러 입력칸에 포커스를 주는 구조라면 `Find` 는 안쪽 입력 상자를 준다. 대시보드 검색은 414×18.3이 잡힌다. 바깥 줄 높이는 따로 확인한다.

```go
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
```

```go test
func TestSmallTargetsNamesOnlyWhatIsTooSmall(t *testing.T) {
	query := "draft"
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Row(c).Padding(20).Gap(8).Children(func() {
			softButton(c, "Save")
			clearButton(c, &query)
			ui.ButtonBase(c).Label("Close").Size(21, 21)
		})
	}, 320, 200)
	tt.Advance(time.Second)
	small := smallTargets(tt, 28, "Save", "Clear search", "Close", "Missing")
	if len(small) != 2 || !strings.HasPrefix(small[0], "Close ") || !strings.HasPrefix(small[1], "Missing ") {
		t.Errorf("small targets %q", small)
	}
	if err := tt.Click("Clear search"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if query != "" {
		t.Errorf("the query stays: %q", query)
	}
}
```

앱에서는 화면마다 누를 수 있는 요소의 `Label` 을 모두 넘기고, 메뉴처럼 열려야 보이는 것은 연 뒤에 잰다. 새 앱에는 시작 템플릿의 `internal/view/app_test.go` 에 이 테스트가 들어 있다.

### 누를 곳만 키우기

```go
func clearButton(c *ui.Context, query *string) ui.Element {
	t := c.Theme()
	button := pressable(ui.ButtonBase(c), 0.86).Label("Clear search").Size(28, 28).Center().KeepFocus()
	hover := progress(button, "hover", button.Hovered(), fadeMotion)
	button.Children(func() {
		ui.Box(c).Size(21, 21).Radius(10.5).Background(t.TextMuted.Alpha(0.5 + 0.3*hover))
	})
	if button.Clicked() {
		*query = ""
	}
	return button
}
```

보이는 원은 21이고 누를 곳은 28이다. X 표시와 등장 움직임은 `examples/dashboard/components.go` 의 같은 이름 함수를 본다.

## 딱딱하지 않게

- 큰 면은 16, 버튼과 입력칸은 8~12로 둥글린다. 큰 면을 각지게 두지도, 모든 것을 알약으로 만들지도 않는다.
- 선은 1 DIP이고 바탕보다 조금만 진하다. 검은 선이나 굵은 테두리를 쓰지 않는다. 그래프 선(2~2.5), 그래프 점의 테두리(1.5~2), 포커스 링만 예외다.
- 카드 그림자는 y 1, 흐림 3, 불투명도 5%로 아주 옅다. 진한 그림자(y 6, 흐림 20, 18%)는 떠 있는 메뉴 패널에만 쓴다.
- 선택은 테두리가 아니라 옅은 강조 면, 강조 글자색, 굵기 600으로 보여 준다.
- 바탕은 옅은 회청색이고 카드는 거의 흰색이다. 카드가 바탕보다 밝아서 옅은 선만으로도 구분된다.
- 모든 상태 변화는 움직인다. 호버 면은 120ms에 번지고, 누르면 100ms 동안 줄어들며(내비와 메뉴 행 0.98, 버튼 0.96, 월 세그먼트 0.94, 주문 링크 0.97, 작은 원형 버튼 0.86), 선택 알약은 240ms 스프링으로 미끄러진다.
- 작은 반응을 군데군데 둔다. "View All" 화살표는 호버하면 오른쪽으로 3 움직이고, 메뉴 꺾쇠는 열리면 돌고, 알림 벨은 처음 열 때 한 번 흔들린다. 같은 정보를 두 번 알리지 않는 데서 멈춘다.

## 복잡하지 않게

- 화면 하나는 질문 하나에 답한다. 대시보드 첫 화면은 지표 4개, 월별 매출, 분류별 매출, 최근 주문 5개로 "이번 달은 어떤가"에 답한다.
- 내비 항목은 5개, 상단 막대는 검색, 알림, 테마 세 가지뿐이다.
- 같은 정보를 두 곳에 두지 않는다. 사용자 이름과 계정 메뉴는 사이드바 아래 한 곳에만 있다. 상단 막대에 한 번 더 두면 이름과 아바타만 겹치고 메뉴 내용도 같다.
- 카드마다 제목 하나에 동작 하나다("View All", 월 선택).
- 가끔 쓰는 동작은 메뉴 안에 둔다. 설정과 종료는 사이드바 아래의 계정 메뉴에 있다.
- 아이콘만 있는 버튼을 줄지어 놓지 않는다. 뜻이 바로 보이지 않으면 글자를 붙인다.
- 첫 화면의 표는 5행까지 보여 주고 나머지는 "View All" 로 넘긴다.
- 같은 종류의 요소는 크기와 여백이 같다. 카드 제목은 모두 20.5, 카드 왼쪽 여백은 20~22다.
- 최소 창 크기, 가장 긴 글자, 빈 데이터에서 먼저 확인한다. 깨지면 요소를 더하지 말고 덜어 낸다.

## 예시 지도

`examples/dashboard` 는 읽고 따라 하는 소스다. 돌려 보려면 임시 폴더에 복사한 뒤 `go tool mygo dev` 로 띄운다. 스킬 폴더 안에서 빌드하지 않는다. 화면을 이루는 부품은 앱 상태와 떨어져 있고, `shell.go` 와 `pages.go` 는 부품을 부르며 앱 상태를 잇는 얇은 층이다.

| 보고 싶은 것 | 파일과 함수 |
|---|---|
| 색 역할, 라이트와 다크, 테마 연결 | `palette.go` 의 `palette`, `fadedPalette`, `theme` |
| 모션 값, `progress`, `pressable`, 벨 키프레임 | `motion.go` |
| 사이드바, 미끄러지는 내비 알약, 아래쪽 계정 메뉴 | `frame.go` 의 `sidebar`, `navItem` |
| 상단 막대와 알림 벨 | `frame.go` 의 `topbar`, `notificationBell` |
| 사이드바, 상단 막대, 스크롤 영역을 한 틀로 | `frame.go` 의 `appFrame` |
| 검색칸, 스위치, 드롭다운 메뉴, 월 세그먼트 | `components.go` 의 `searchField`, `themeSwitch`, `menuButton`, `dropdown`, `segmented` |
| 카드, 표, 상태 알약 | `components.go` 의 `card`, `table`, `statusPill` |
| 지표 카드와 폭이 바뀌는 배지 | `components.go` 의 `metricCard`, `deltaPill` |
| "View All" 단추, 밑줄이 자라는 링크 | `components.go` 의 `moreButton`, `textLink` |
| 미니 그래프, 월별 영역 그래프, 도넛 | `charts.go` 의 `sparkline`, `areaChart`, `donut` |
| 그래프 말풍선 | `charts.go` 의 `pickPoint`, `chartTip` |
| 숫자 서식 | `formats.go` |
| 아이콘과 라이선스 | `icons.go`, `icons/` |
| 숨긴 타이틀바 창 | `window.go` |
| 앱 상태에 부품을 잇는 방법 | `shell.go`, `pages.go` |
| 부품만으로 화면을 조립하는 예 | `frame_test.go` 의 `TestAScreenBuiltFromPartsOnly` |
| 동작 테스트 | `main_test.go` |

## 코드로 가져다 쓰기

디자인 시스템의 코드는 값, 모션, 부품 세 층으로 나뉜다. 모두 대시보드의 파일이고, 시작 템플릿과 새 앱이 같은 파일을 쓴다.

| 파일 | 담은 것 | 새 앱에 |
|---|---|---|
| `palette.go` | 역할별 색(라이트, 다크), 테마를 바꿀 때 모든 칸을 섞는 `fadedPalette`, 위젯에 연결하는 `theme`, 위젯의 모든 색 칸까지 섞는 `fadedTheme` | 늘 복사 |
| `motion.go` | 모션 값 이름들, 벨 키프레임, `progress`, `pressable` | 늘 복사 |
| `palette_test.go`, `targets_test.go` | 팔레트와 `Theme` 의 모든 색 칸이 섞이는지, 글자 대비가 4.5:1 이상인지, 누를 곳이 28 DIP 이상인지 재는 `smallTargets` | 늘 복사 |
| `components.go` | 카드, 표, 메뉴, 세그먼트, 검색칸, 스위치, 지표 카드, 링크 | `--kit dashboard` |
| `charts.go` | 미니 그래프, 영역 그래프, 도넛, 말풍선 | `--kit dashboard` |
| `frame.go` | 사이드바, 상단 막대, 알림 벨, 페이지 틀 | `--kit dashboard` |
| `formats.go` | 금액, 비율, 천 단위 서식 | `--kit dashboard` |
| `icons.go`, `icons/` | Lucide 아이콘 15개를 읽는 코드와 SVG 파일, ISC 라이선스 | `--kit dashboard` |
| `window.go` | 타이틀바를 숨기고 기본 크기와 최소 크기를 정한 창 옵션 | `--kit dashboard`, 프로젝트 맨 위 |
| `components_test.go`, `charts_test.go`, `frame_test.go`, `formats_test.go`, `icons_test.go`, `window_test.go` | 부품만으로 만든 작은 화면에서 부품을 시험한다. 쓰는 예시도 된다 | `--kit dashboard` |

```sh
python3 $SKILL/scripts/new_app.py my-app --name "My App" --kit dashboard
```

부품은 앱 상태에 묶이지 않는다. 팔레트(`p`)를 받고, 앱에 묶인 것(아이콘, 상태 포인터, 콜백, 글자, 폭)만 매개변수로 받는다. 화면에 나오는 글자(제목, 이름, 비교 문구)는 모두 매개변수라서 부품 안에 박혀 있지 않다.

| 부품 | 하는 일 |
|---|---|
| `card(c, p)`, `cardTitle(c, title)` | 모서리 16, 1 DIP 선, 옅은 그림자의 면과 제목 20.5 |
| `statusPill(c, tone, label)` | 옅은 바탕, 진한 글자, 점 하나의 상태 알약 |
| `table(c, p, label, empty, columns, rows, cell)`, `column`, `cellText` | 머리글 행, 행 구분선, 빈 문구. 셀이 행 높이 전체를 차지해 링크를 크게 잡을 수 있다 |
| `menuButton(c, p, label, open)`, `chevron(c, p, icon, open)`, `dropdown(c, p, anchor, open, upward, entries)`, `menuEntry` | 키보드로 열고 닫는 메뉴. Esc로 닫으면 포커스가 버튼으로 돌아가고 행은 30 높이다 |
| `segmented(c, p, key, label, selected, names, width)` | 선택 알약이 스프링으로 미끄러지는 세그먼트 |
| `searchField(c, p, query, icon, placeholder)`, `clearButton` | 칸 어디를 눌러도 포커스가 가고, 지우기 버튼은 28×28이다 |
| `themeSwitch(c, p, on, label, sun, moon, changed)` | 58×30 스위치. 켜고 끄는 움직임과 아이콘이 있다 |
| `metricCard(c, p, metricView)`, `deltaPill(c, p, delta, up, down)` | 아이콘, 굴러서 바뀌는 큰 숫자, 색과 폭이 바뀌는 증감 배지, 미니 그래프를 한 카드에 담는다 |
| `moreButton(c, p, label, icon)`, `textLink(c, p, text)` | 호버하면 화살표가 3 움직이는 테두리 단추, 호버하면 밑줄이 자라는 표 안 링크. 둘 다 받은 쪽이 `Clicked()` 로 누름을 본다 |
| `sparkline(c, p, values, describe)` | 호버하면 선이 두꺼워지고 점과 말풍선이 뜨는 미니 그래프 |
| `areaChart(c, p, areaChartView)` | 눈금, 영역, 선, 점, 말풍선이 있는 그래프. `shaded` 번째 점까지 영역이 스프링으로 자란다 |
| `donut(c, p, donutView)` | 도넛과 범례. 조각이나 범례 행에 호버하면 조각이 자라고 나머지가 흐려지며, 키보드로도 훑는다 |
| `pickPoint(c, area, xs, describe)`, `chartTip(c, p, width, x, y, text)` | 그래프 영역 전체에서 가장 가까운 값을 고르고, 경계 안에 머무는 말풍선을 띄운다 |
| `sidebar(c, p, title, items, selected, choose, account)`, `navItem`, `logo`, `avatar` | 폭 237의 사이드바. 미끄러지는 선택 알약과 아래쪽 계정 메뉴가 있고, 숨긴 타이틀바면 창 버튼 자리를 비운다 |
| `topbar(c, p, content)`, `notificationBell(c, p, bell)` | 높이 69의 상단 막대(끌면 창이 움직이고 오른쪽 창 버튼 자리를 비운다). 벨은 열 때마다 한 번 흔들리고 점은 한 번 열어 보면 사라진다 |
| `appFrame(c, p, pageKey, rail, bar, body)` | 사이드바, 상단 막대, 페이지 영역을 한 틀로 묶는다. 페이지가 바뀌면 스크롤이 맨 위로 가고 새 페이지가 6 아래에서 들어온다 |
| `grouped`, `dollars`, `cents`, `thousands`, `rate`, `axisDollars` | `mygo/i18n` 으로 천 단위를 끊는 금액과 비율 서식 |
| `dashboardWindow(title, view)` | 1440×900, 최소 1024×680, 타이틀바 숨김, 크기를 기억하는 창 옵션 |
| `labelWeight(selected)` | 선택되면 500에서 600으로 |

- 새 앱은 `scripts/new_app.py` 가 파일을 `internal/view/` 로 복사한다. 패키지 이름만 `view` 로 바뀌고, `window.go` 와 그 테스트만 `main` 패키지 그대로 프로젝트 맨 위로 간다. 복사한 뒤에는 그 앱의 파일이니 색과 부품을 마음대로 고친다. 쓰지 않는 부품과 모션 값 줄은 지우고, 지운 부품의 테스트도 함께 지운다.
- 부품은 값이 바뀌고 비는 앱에서도 깨지지 않는다. 항목이 0개인 메뉴(`dropdown` 은 열리지 않고 닫힌 채로 둔다), 값이 하나뿐이거나 모두 같거나 모두 0인 그래프와 도넛, 지난달이 0이라 증감이 NaN이나 Inf인 배지(`deltaPill` 은 0으로 본다)가 그렇다. 사용자가 입력하는 긴 글자는 표 칸과 도넛 범례에서 한 줄로 끝나고 말줄임표가 붙는다. 시험은 `charts_test.go` 와 `components_test.go` 에 있다.
- 부품을 쓰는 방법은 `frame_test.go` 의 `TestAScreenBuiltFromPartsOnly` 가 보여 준다. 사이드바, 상단 막대, 지표 카드, 그래프, 도넛, 표를 부품만으로 한 화면에 조립하고 누르고 입력해 본다. 대시보드의 `shell.go` 와 `pages.go` 는 같은 부품을 부르며 앱 상태를 잇는 층이다.
- 아이콘은 `icons/` 폴더(ISC 라이선스 파일 포함)와 `icons.go` 가 함께 간다. 라이선스 파일을 지우지 않는다(`TestIconsParseAndTheirLicenseTravelsWithThem` 가 지킨다). 쓰지 않는 아이콘은 파일과 변수를 함께 지운다.
- `window.go` 는 `main.go` 에서 `mygo.NewWindow(dashboardWindow("My App", screen))` 로 쓴다. 숨긴 타이틀바는 `sidebar` 와 `topbar` 가 창 버튼 자리를 비우는 것과 짝이다. 일반 타이틀바 창을 쓸 앱은 이 파일을 쓰지 않아도 된다.
- 서식이 쓰는 `mygo/i18n` 은 시작 템플릿의 구조 검사(`architecture_test.go`)가 뷰 계층에서 허용한다. 뷰가 쓸 수 있는 MyGo 패키지는 `ui` 와 그 하위, `i18n` 이다.
- `dropdown` 의 강조 행과 `segmented` 의 선택 알약은 흰 글자를 `accent` 위에 올려 다크에서 3.84:1이다. 쓰는 앱은 `components.go` 의 두 곳의 `p.accent` 를 `p.accentFill` 로 바꾼다. 대시보드는 다크 화면이 달라져서 아직 바꾸지 않았다.
- 원본은 `examples/dashboard` 이고 `assets/` 의 사본과 같아야 한다. `motion.go` 는 값 목록이라 앱이 쓰지 않는 줄을 지워도 되고, 남긴 줄이 `assets/motion.go` 와 같은지 본다. `scripts/check_snippets.py` 가 이를 확인하고, 이 문서의 색이 `palette.go` 에 있는지도 확인한다. `new_app.py` 와 `new_app.py --kit dashboard` 가 만든 앱이 둘 다 도는지도 확인한다.
- 색을 더하면 `palette` 구조체, 두 팔레트 값, `channels()` 에 모두 넣는다. 하나라도 빠지면 `TestPaletteFadeCoversEveryColor` 가 실패한다. MyGo가 올라가 `Theme` 에 색 칸이 늘면 `themeColors` 에 넣으라고 `TestThemeFadeCoversEveryColorAndEndsExactly` 가 실패한다.

## 점검표

- 앱 화면이 `ui.View` 로만 그려지고 웹뷰, HTML, JavaScript가 없다.
- 강조색이 하나이고 글자색이 세 단계다. 글자 대비가 4.5:1 이상이다.
- 누를 수 있는 것이 모두 28×28 이상이고 그것을 지키는 테스트가 있다. 메뉴는 연 상태로 쟀다.
- 모서리, 선, 그림자가 위 값 안에 있다.
- 누를 수 있는 것마다 호버, 누름, 포커스 링이 있고, 선택은 면과 글자색으로 보인다.
- 라이트, 다크, 최소 창 크기에서 글자가 겹치거나 잘리지 않는다.
- 넣지 않는 것 목록에 걸리는 것이 없다.
