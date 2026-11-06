# 클린 아키텍처로 짜기

요청자의 기본 지시다. 새 MyGo 앱은 클린 아키텍처로 짠다. 층을 나누는 목적은 형식이 아니라 세 가지다. 업무 규칙을 창 없이 테스트하고, 저장소나 네트워크를 바꿔도 안쪽을 건드리지 않으며, 뷰가 표시에만 집중하게 한다.

목차: [언제 적용하는가](#언제-적용하는가) · [의존성 규칙](#의존성-규칙) · [어디에 둘까](#어디에-둘까) · [MyGo에 맞춘 규칙](#mygo에-맞춘-규칙) · [느린 유스케이스](#느린-유스케이스) · [층별 테스트](#층별-테스트) · [경계 검사](#경계-검사) · [피할 것](#피할-것) · [기존 앱 옮기기](#기존-앱-옮기기)

## 언제 적용하는가

| 상황 | 할 일 |
|---|---|
| 새 앱 | `scripts/new_app.py` 로 시작한다. 이 구조가 기본이다 |
| 기존 앱이 이미 이 구조 | 그대로 따른다 |
| 기존 앱이 다른 구조 | 그 프로젝트의 구조를 따른다. 요청 없이 옮기지 않고, 옮기는 편이 낫다면 제안만 한다. 새 업무 규칙은 그 구조 안에서라도 뷰 밖의 함수로 빼서 창 없이 테스트한다 |
| 구조를 옮겨 달라는 요청 | [기존 앱 옮기기](#기존-앱-옮기기)를 따른다 |
| 일회성 데모, 한 파일 실험, 작은 스크립트 | 계층과 인터페이스를 강요하지 않는다 |

## 의존성 규칙

소스의 의존성은 안쪽으로만 향한다. 안쪽은 바깥을 모른다.

```text
main.go ──> internal/view ─────┐
   │                           ├──> internal/usecase ──> internal/domain
   └──────> internal/adapter ──┘
```

| 층 | 폴더 | 둘 것 | 가져올 수 있는 것 |
|---|---|---|---|
| 도메인 | `internal/domain` | 엔티티, 값 객체, 업무 규칙 | 표준 라이브러리 |
| 유스케이스 | `internal/usecase` | 앱이 하는 일, 필요한 포트 | `domain` |
| 어댑터 | `internal/adapter` | 포트 구현: 파일, DB, HTTP, 시스템 기능 | `domain`, `usecase`, `github.com/egoist/mygo` 전체 |
| 뷰 | `internal/view` | 화면 상태, 뷰 함수, 팔레트, 모션, 아이콘 | `domain`, `usecase`, `github.com/egoist/mygo/ui` 와 그 하위 패키지 |
| 조립 | `main.go` | 객체 조립, 창 옵션, `App.Run` | 전부 |

뷰 패키지 이름을 `ui` 로 하면 `github.com/egoist/mygo/ui` 와 헷갈리므로 `view` 로 한다. `main.go` 는 `mygo build` 의 기본 위치인 프로젝트 맨 위에 둔다.

## 어디에 둘까

| 코드 | 위치 |
|---|---|
| 합계, 검증, 상태 전이처럼 화면이 없어도 성립하는 규칙 | `domain` |
| 저장과 불러오기의 흐름, 여러 엔티티에 걸친 동작, 검색과 집계 | `usecase` |
| JSON 변환, 파일 경로, SQL, HTTP 요청 | `adapter` |
| 금액 표시 형식, 날짜 문자열, 색, 아이콘, 간격 | `view` |
| 열림 여부, 현재 페이지, 입력 중인 글자, 호버, 선택, 애니메이션 값 | `view` |
| 저장 경로 결정, 외관 설정, 메뉴와 트레이 연결 | `main.go` 또는 `adapter` |

판단이 어려우면 이렇게 묻는다. 창이 없는 명령줄 도구로 앱을 다시 만들어도 필요한가? 그렇다면 `domain` 이나 `usecase` 다. 데스크톱이 아니라 서버로 옮겨도 그대로인가? 그렇지 않은 입출력이면 `adapter` 다.

## MyGo에 맞춘 규칙

- **`ui.Element` 와 `*ui.Context` 는 뷰 밖으로 나가지 않는다.** 유스케이스의 인자나 반환값에 넣지 않는다.
- **뷰는 `github.com/egoist/mygo` 를 가져오지 않는다.** 뷰에 필요한 시스템 기능은 `main.go` 에서 함수로 넣어 준다. 메인 스레드로 넘어가는 호출이 테스트를 멈추는 문제가 구조적으로 사라진다.
- **유스케이스는 구체 타입으로 쓴다.** 뷰는 `*usecase.Tasks` 를 직접 받는다. 구현이 하나인 인터페이스를 만들지 않는다.
- **포트만 인터페이스로 둔다.** 바깥 세계에 닿는 곳(`TaskStore`)을 쓰는 쪽인 `usecase` 에 선언한다. 실제 구현과 테스트용 가짜, 두 구현이 있으므로 값어치를 한다. 메서드가 하나뿐이면 함수 값으로 충분할 수 있다.
- **표시 변환은 뷰에서 한다.** 유스케이스는 값과 오류를 돌려주고, 글자로 바꾸는 일은 뷰가 한다.
- **화면 상태와 앱 상태를 나눈다.** 껐다 켜도 남아야 하면 앱 상태이고 유스케이스를 거친다. 그 밖에는 뷰 구조체에 둔다.
- **뷰 안에서 부르는 유스케이스는 빠른 조회만 한다.** 뷰는 프레임마다 돈다. 변경은 `Clicked`, `OnClick`, `OnChange` 같은 처리된 동작에서 부른다.
- **파일은 크기대로 나눈다.** 뷰 패키지 안을 `shell.go`, `pages.go`, `palette.go`, `motion.go`, `icons.go` 처럼 나눈다. 패키지를 더 쪼개는 것은 화면끼리 정말 독립적일 때만 한다.

뷰와 조립은 시작 템플릿에 이미 이 형태로 있다. `assets/starter/` 의 `main.go`, `internal/view/app.go`, `internal/usecase/tasks.go`, `internal/adapter/filestore.go` 를 먼저 읽는다. `internal/view/motion.go` 는 `new_app.py` 가 `assets/motion.go` 에서 만든다. 할 일 도메인은 예시이므로 실제 도메인으로 바꾸고 남기지 않는다.

웹 창을 쓰는 앱도 같다. `mygo.Bind` 로 넘기는 서비스는 바깥층이고 유스케이스를 부른다. 이 줄은 문서 기준이다.

## 느린 유스케이스

파일, 네트워크처럼 오래 걸릴 수 있는 작업은 메인 스레드에서 돌리지 않는다. 고루틴에서 값을 받아 온 뒤 상태 변경만 메인 스레드로 넘긴다. 유스케이스는 공유 상태 없이 값을 돌려주게 짜면 잠금이 필요 없다. 넘기는 함수는 `main.go` 에서 `win.Update` 를 넣어 준다.

```go
type reportScreen struct {
	fetch   func() ([]string, error)
	publish func(func())
	rows    []string
	problem error
	loading bool
}

func (s *reportScreen) refresh() {
	if s.loading {
		return
	}
	s.loading = true
	go func() {
		rows, err := s.fetch()
		s.publish(func() {
			s.rows, s.problem, s.loading = rows, err, false
		})
	}()
}
```

```go test
func TestRefreshPublishesOnTheMainThreadOnce(t *testing.T) {
	calls := 0
	done := make(chan func())
	s := &reportScreen{
		fetch: func() ([]string, error) {
			calls++
			return []string{"row"}, nil
		},
		publish: func(update func()) { done <- update },
	}
	s.refresh()
	s.refresh()
	update := <-done
	if len(s.rows) != 0 || !s.loading {
		t.Error("the state changes before the main thread takes it")
	}
	update()
	if calls != 1 || len(s.rows) != 1 || s.loading {
		t.Errorf("calls %d, rows %v, loading %v", calls, s.rows, s.loading)
	}
}
```

`win` 은 `Content` 를 받아 만들어지므로 두 단계로 조립한다. 화면 구조체를 만들고, `win := mygo.NewWindow(...)` 를 부른 뒤 `screen.publish = win.Update` 를 넣는다. 결과가 낡았을 수 있으면(검색어가 바뀐 뒤 돌아온 응답) 요청 번호를 비교해 버린다. 여러 고루틴이 같은 값을 건드리면 `go test -race` 로 확인한다.

## 층별 테스트

| 층 | 테스트 | 도구 |
|---|---|---|
| `domain` | 규칙과 경계값 | 평범한 Go 테스트 |
| `usecase` | 흐름, 저장 실패 시 상태 유지 | 메모리 가짜 포트 |
| `adapter` | 왕복 저장, 없는 파일, 깨진 파일 | `t.TempDir` |
| `view` | 표시, 입력, 움직임, 쉬는 상태 | `ui.Tester` + 실제 유스케이스 + 메모리 포트 |
| 전체 | 의존성 방향 | `architecture_test.go` |
| `main.go` | 조립이 실제로 도는지 | 시험 빌드의 실제 창 |

뷰 테스트에서 유스케이스를 가짜로 바꾸지 않는다. 포트만 가짜로 바꾸면 업무 규칙까지 함께 검증된다.

## 경계 검사

`assets/starter/architecture_test.go` 는 `internal/` 의 모든 Go 파일을 읽어 금지된 가져오기가 있으면 실패한다. 새 층을 만들면 이 파일에 먼저 선언해야 한다. 그래서 층이 말없이 늘지 않는다. `go test ./...` 에 포함되므로 평소에는 따로 신경 쓰지 않아도 된다.

이 테스트는 파일을 직접 읽으므로 `-overlay` 변이로는 검증되지 않는다. 임시 복사본의 파일을 실제로 고쳐 확인한다.

## 피할 것

| 하지 말 것 | 이유 | 대신 |
|---|---|---|
| 구현이 하나뿐인 유스케이스 인터페이스 | 값어치 없는 간접 호출 | 구체 타입 |
| 층마다 같은 필드의 구조체를 복사 | 변환 코드만 는다 | 실제로 모양이 달라질 때만 나눈다 |
| `domain` 에 표시 문자열이나 색 | 안쪽이 화면을 알게 된다 | `view` 에서 변환 |
| 뷰에 저장 로직이나 파일 경로 | 창 없이 테스트할 수 없다 | `usecase` 와 `adapter` |
| 뷰에서 `mygo.Theme.SetSource` 를 직접 호출 | 테스트가 멈춘다 | 함수로 주입 |
| 유스케이스에서 `Window.Update` 호출 | 안쪽이 프레임워크를 알게 된다 | 뷰에서 넘긴 함수로 처리 |
| 뷰 패키지를 파일마다 패키지로 분리 | 경계 없이 가져오기만 는다 | 한 패키지 안의 파일로 나눈다 |
| 팩토리, 전역 컨테이너, DI 프레임워크 | `main.go` 의 조립 몇 줄로 끝난다 | 생성자에 값을 넘긴다 |

## 기존 앱 옮기기

요청받았을 때만 한다. 한 번에 한 단계씩 옮기고 단계마다 테스트와 쉬는 화면 비교를 돌린다.

1. 쉬는 화면과 테스트 결과를 기준선으로 남긴다.
2. 순수한 규칙과 타입을 `domain` 으로 옮긴다. 뷰는 그대로 둔다.
3. 필터, 집계, 저장 흐름을 `usecase` 로 옮기고 입출력을 포트로 뺀다.
4. 입출력 구현을 `adapter` 로 옮긴다.
5. 뷰를 `internal/view` 로 옮기고 시스템 기능 호출을 주입으로 바꾼다.
6. `main.go` 에 조립만 남기고 `architecture_test.go` 를 넣는다.
7. 쉬는 화면이 모두 같은지 확인한다. 구조만 옮겼으므로 픽셀이 달라지면 안 된다.
