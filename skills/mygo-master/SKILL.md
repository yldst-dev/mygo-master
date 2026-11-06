---
name: mygo-master
description: >-
  MyGo(github.com/egoist/mygo)로 Go 데스크톱 앱을 만들거나 고칠 때 쓰는 종합 지침. 항상 최신
  릴리스의 MyGo로, 새 앱은 클린 아키텍처로 짠다. 프로젝트 시작, 뷰와 상태, 레이아웃, 위젯, 차트,
  메뉴, 테마, 빌드부터 부드러운 모션(호버, 누름, 메뉴, 페이지 전환, 그래프 말풍선, 숫자 굴림)의
  검증된 값과 레시피, Tester 테스트, 필름스트립, 변이 검사, 실제 창 자동화, 프레임 통계까지 담았다.
  go.mod에 github.com/egoist/mygo가 있는 프로젝트를 만지거나, "MyGo로 만들어줘", "이 앱에
  애니메이션 넣어줘", "부드럽게 해줘", "커서 올리면 표시되게", "go tool mygo", "ui.Tester", "Go
  네이티브 UI" 같은 말이 나오면 프레임워크 이름이 없어도 코드를 쓰기 전에 먼저 읽는다. Use for any
  native UI, motion, testing or build work on a Go app using MyGo, package ui or ui/chart.
---

# MyGo 마스터

MyGo 앱을 최신 버전과 클린 아키텍처로, 부드럽고 정확하고 꼼꼼하게 만드는 방법이다. 값과 절차는 실제 대시보드 앱을 만들고 요청자가 직접 써 본 뒤 "충분히 부드럽다"고 판단한 결과에서 왔다. MyGo v0.7.4에서 검증했다. 아래 명령의 `$SKILL` 은 이 SKILL.md가 있는 폴더다. Claude Code 개인 설치라면 `~/.claude/skills/mygo-master` 다.

## 다섯 가지 기준

- **최신 MyGo로 짠다.** 시작할 때 최신 릴리스를 확인하고, 그 버전의 소스와 문서로 API를 확인한다. 기억으로 쓰지 않는다. 뒤처진 프로젝트는 먼저 올린다.
- **클린 아키텍처로 짠다.** 새 앱은 도메인, 유스케이스, 어댑터, 뷰로 나누고 의존성은 안쪽으로만 향한다. 기존 앱은 그 앱의 구조를 따르고, 일회성 데모에는 층을 강요하지 않는다.
- **부드럽다.** 끊기는 곳이 없다. 커서가 닿으면 바로 반응하고, 중간에 방향이 바뀌어도 이어지며, 옆 요소까지 함께 움직인다. 튕기지 않고 지정 시간은 240ms를 넘지 않는다.
- **정확하다.** 쉬는 화면은 픽셀까지 맞고, 값은 경계 안에 있으며, 시안이 있으면 시안의 색과 크기를 잰다.
- **꼼꼼하다.** 포인터로 되는 일은 키보드로도 되고, 애니메이션이 끝나면 프레임을 요청하지 않는다. 테스트는 코드를 망가뜨려 힘을 확인한다. 돌리지 않은 검증은 돌렸다고 말하지 않는다.

## 어떤 문서를 읽을까

| 작업 | 읽을 것 |
|---|---|
| 새 앱, 버전 올리기, API 사용법, 위젯, 레이아웃, 차트, 빌드 | [references/framework.md](references/framework.md) |
| 코드를 어느 층에 둘지, 느린 작업, 구조 옮기기 | [references/architecture.md](references/architecture.md) |
| 움직임, 호버, 메뉴, 페이지 전환, 그래프 반응, 값 고르기 | [references/motion.md](references/motion.md) |
| 테스트, 필름스트립, 실제 창, 보고 | [references/verification.md](references/verification.md) |
| 이상한 증상, 작동하지 않는 전환, 멈추는 테스트 | [references/pitfalls.md](references/pitfalls.md) |

새 앱이면 `framework.md` 와 `architecture.md` 를 읽는다. 움직임을 넣는 작업이면 `motion.md` 와 `verification.md` 를 모두 읽고, 레시피를 고치기 전에 `pitfalls.md` 의 움직임 표를 훑는다.

## 작업 순서

### 1. 최신 버전과 프로젝트를 확인한다

```sh
go list -m -json github.com/egoist/mygo@latest
grep -n egoist/mygo go.mod
```

이 지시는 요청자가 직접 한 것이다. 그래서 MyGo 자체의 버전 변경은 요청받은 변경으로 본다.

- **새 앱:** `python3 $SKILL/scripts/new_app.py <폴더> --name "<영문 이름>"` 으로 시작한다. 최신 릴리스, 클린 아키텍처 구조, 모션 값, 테스트가 함께 들어온다. 한글 앱 이름은 macOS 빌드의 서명 단계에서 실패하므로 창 제목과 화면 글자에만 쓴다.
- **뒤처진 기존 앱:** 기능을 고치기 전에 [올린다](references/framework.md#기존-프로젝트를-최신으로-올리기). 요청자가 버전을 고정하라고 했거나 `replace` 가 있으면 그대로 두고 그 사실을 적는다.
- **최신이 v0.7.4보다 새로울 때:** 릴리스 노트를 읽고 이 스킬의 코드가 아직 도는지 확인한다. 깨진 곳은 새 API로 고쳐 쓰고 보고한다.

```sh
python3 $SKILL/scripts/check_snippets.py
```

기존 앱은 그 앱의 파일 구조, 모션 값, 팔레트를 먼저 따른다. 기존 화면을 건드린다면 지금 찍어 둔다([작업 전 기준선](references/verification.md#작업-전-기준선)).

### 2. 고치기 전에 짧게 계획한다

코드를 어느 층에 둘지 정한다. 움직임을 넣을 때는 코드 전에 짧은 표를 쓴다. 무엇이 어떤 상태에서 어떤 도구로 얼마나 움직이는지, 그리고 그때 무엇이 뛸 수 있는지 적는다. 값은 `motion.md` 의 값 표에서 고른다. 디자인 작업이면 요청자의 전역 디자인 지침도 따른다. 새 글꼴이나 외부 라이브러리는 제안만 한다.

### 3. 구현한다

항상 지키는 규칙은 이렇다.

- 쓰기 전에 그 버전의 문서와 `go doc` 으로 API를 확인하고, 새 버전에 이미 있는 기능을 다시 만들지 않는다.
- 요청하지 않은 기능은 넣지 않고 제안한다. 안정성에 필요한 처리(저장 실패, 깨진 파일, 포커스 복구)는 기능이 아니라 품질이므로 챙긴다.
- 업무 규칙은 `domain` 과 `usecase` 에 두고 MyGo를 가져오지 않는다. `ui.Element` 와 `*ui.Context` 는 뷰 밖으로 나가지 않는다. 인터페이스는 바깥 세계에 닿는 포트에만 둔다.
- 뷰는 상태의 함수이고 프레임마다 돈다. 안에서 작업을 시작하지 않고, 요소를 저장하지 않는다. 다른 고루틴은 `Window.Update` 로 넘긴다.
- 메인 스레드로 넘어가는 호출(`mygo.Theme.SetSource` 등)은 `main.go` 에서 함수로 주입한다. 그래야 테스트가 멈추지 않는다.
- 조건부 형제와 움직이는 목록에는 키를 주고, 글자가 없는 요소에는 `Label` 을 준다.
- 내장 위젯에 움직임을 넣을 자리가 없으면 베이스 위젯으로 짠다. 그때 잃는 것을 보고한다.
- 색은 스칼라로 애니메이션하고 팔레트에서 섞는다. `assets/motion.go` 의 값과 `progress`, `pressable` 을 쓴다. 전환에는 `Position`, `Size` 처럼 움직일 항목을 명시한다.
- 포즈로 되는 움직임은 포즈로 한다. 뷰를 다시 돌리지 않아 가장 싸다.
- 포인터 입력을 받지 않을 말풍선에는 `IgnoreInput(true)` 를 건다.
- 쉬는 화면은 바뀌지 않아야 한다. 호버 표면 때문에 여백을 바꿨다면 픽셀로 확인한다.
- Enter로 제출하는 입력칸은 한글 조합 중 Enter 문제를 확인한다([함정 모음](references/pitfalls.md#입력과-포커스)).

### 4. 검증한다

변경 크기에 맞게 고른다. 전체 표는 [verification.md](references/verification.md#어느-층을-돌릴지)에 있다. 움직임을 새로 넣었다면 전부 돌린다.

```sh
gofmt -l .
go vet ./...
go tool mygo vet .
go tool mygo migrate-ui .
CGO_ENABLED=0 go test -count=1 -timeout 120s ./...
```

`gofmt -l` 과 `migrate-ui` 미리 보기는 출력이 없어야 한다.

| 도구 | 용도 |
|---|---|
| `scripts/new_app.py` | 최신 MyGo와 클린 아키텍처로 새 앱을 만든다 |
| `assets/starter/architecture_test.go` | 의존성 방향을 테스트로 지킨다 |
| `assets/film_test.go` + `scripts/overlay.py` | 저장소를 건드리지 않고 필름스트립과 쉬는 화면을 찍는다 |
| `scripts/pixdiff.go` | 쉬는 화면과 시안을 픽셀로 비교한다 |
| `scripts/mutate.py` | 코드를 망가뜨려 테스트의 힘을 확인한다 |
| `scripts/drive.py` | 자동화 플러그인을 넣은 시험 빌드를 띄우고 합성 입력을 보낸다 |
| `scripts/frame_stats.py` | 프레임 로그를 요약한다 |
| `scripts/montage.go` | 실제 창 캡처를 한 장으로 모은다 |
| `scripts/check_snippets.py` | 이 스킬의 코드와 API 이름이 최신 MyGo에 맞는지 확인한다. `--project` 로 프로젝트의 고정 버전도 볼 수 있다 |

필름스트립은 Read 도구로 직접 본다. 움직임은 숫자로 확인되지 않는다.

### 5. 끝내고 보고한다

끝내기 전에 `git status` 로 실행 파일 같은 부산물을 찾고, `pgrep -fl` 로 시험 프로세스가 남지 않았는지 본다. 보고에는 쓴 MyGo 버전, 달라진 것, 요청 밖으로 바뀐 것, 고친 문제, 숫자로 된 검증 결과, 확인하지 못한 것을 적는다. 합성 이벤트로 본 것과 실제 마우스로 본 것을 구분한다. 요청받지 않으면 커밋하지 않는다.

## 끝났다고 말하기 전에

- 프로젝트가 최신 MyGo 위에 있다. 아니면 그 이유를 적었다.
- 새 앱은 `architecture_test.go` 가 통과하고, 업무 규칙이 창 없이 테스트된다.
- 정적 검사와 테스트를 지금 코드에서 돌렸다.
- 쉬는 상태, 호버 중, 호버 뒤에 `FramePending` 이 false다.
- 새 테스트를 변이 검사로 확인했다.
- 새 움직임을 필름스트립으로 봤다. 최소 창 크기와 다크 모드도 봤다.
- 쉬는 화면 차이를 모두 설명할 수 있다.
- 키보드로도 쓸 수 있고, 줄인 움직임에서 끝값으로 간다.
- 창에 보이는 변경은 실제 창에서 돌렸거나, 돌리지 못했다고 적었다.
- 저장소에 부산물이 없고 시험 프로세스가 없다.
