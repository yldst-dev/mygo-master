<div align="center">

# Dashboard

A native dashboard built with MyGo. It draws everything with `ui.View`, with no web view, HTML or JavaScript

[한국어](README.md) · English

</div>

Four metric cards, a monthly revenue chart, a sales-by-category donut and a recent orders table answer "how is this month going". It was built from a reference image, and every button, menu and chart answers the pointer with hover and motion. The orders and customers on screen are made-up sample data.

## Run it

You need Go. `go.mod` asks for Go 1.27.3; with the default `GOTOOLCHAIN=auto`, an older Go downloads that version for you.

```sh
go tool mygo dev
go test ./...
go tool mygo build -skip-dmg
```

`go tool mygo dev` rebuilds and restarts the app when a Go file changes. The build lands in `build/darwin-arm64/Dashboard.app`. The default window is 1440×900 and the smallest is 1024×680.

## Screens

| Screen | Contents |
|---|---|
| Dashboard | 4 metrics with sparklines, a revenue chart with a month pill, a donut, the 5 latest orders |
| Analytics | The revenue chart and a table of monthly figures |
| Orders | Every order, filtered by the search field |
| Customers | Orders summed up per customer |
| Settings | The dark mode switch |

The search field in the top bar filters by order ID, customer, status and date. The notification bell and the account menu open and close from the keyboard too.

## Files

| File | Role |
|---|---|
| `main.go` | Window options and the theme hookup |
| `shell.go` | App state, sidebar, top bar, notification bell |
| `pages.go` | Metrics and badges, charts, the donut, the 5 pages |
| `components.go` | Parts that are not tied to app state: card, status pill, table, dropdown menu, segmented control, search field, theme switch, chart tooltip. `new_app.py --kit dashboard` copies them into a new app |
| `data.go` | Sample data, number formats, order filtering and customer summaries |
| `palette.go` | Color roles (light and dark) and the theme fade. The same file as in the skill's starter |
| `motion.go` | Motion value names, `progress`, `pressable`. The lines this app uses from the starter's list of values |
| `icons.go`, `icons/` | [Lucide](https://lucide.dev) icons (ISC License) |
| `main_test.go` | Tests of the whole app, run without a window |
| `palette_test.go`, `components_test.go`, `targets_test.go` | Tests of the palette and the parts on their own. A new app gets them too |
| `mygo.json`, `resources/` | App name and identifier, icon |

It is a single package. The app predates the clean architecture rule, so it has no layers. For the structure of a new app, see the skill's starter.

## Design

The direction is a calm work tool. It looks refined by taking things away and lining them up.

- A bright blue-tinged gray background and one accent, blue (`#0A5FFE`). Dark mode is not an inversion; each role has its own value.
- Three text colors: title, label, muted. Cards are told apart by a corner radius of 16, a 1 DIP line and a very faint shadow.
- Everything clickable is at least 28×28 DIP, and `TestEveryTargetIsLargeEnough` keeps it that way. It measures with the menu open and with a search typed too.
- Every state change moves: hover 120ms, press 100ms, color change 200ms, chart marks a 180ms spring, pills and rolling figures a 240ms spring. Once motion ends, the frame requests stop.
- When the system asks for less motion, everything jumps to its end value.

The reasons behind the colors, type, spacing, corners, target sizes and motion values are in the `mygo-master` skill, in `references/design.md` and `references/motion.md`.

## Tests

`go test ./...` runs without a window, through `ui.Tester`. It checks the figures and texts on screen, the month pill, search, navigation, dark mode, the keys in menus, hover on the charts and the donut, that a tooltip stays inside its chart, that motion ends, less motion, target sizes, that the fade moves every palette color together, and text contrast. The parts (card, menu, segmented control, table, search field, switch, tooltip) are also tested on their own, in a tiny view made of nothing else.

## Known gaps

- In dark mode, white text on the month pill and on the highlighted menu row has a contrast of 3.84:1, short of 4.5:1. Painting them with `accentFill` from `palette.go` gives 4.66:1, but it changes the look a little, so it is not done yet.
- The avatar initials in light mode have a contrast of 4.27:1.
- The app icon is the default one that `mygo init` makes.
- It has only been run on macOS.
