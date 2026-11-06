package main

import (
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

type page int

const (
	pageDashboard page = iota
	pageAnalytics
	pageOrders
	pageCustomers
	pageSettings
)

const (
	navItemHeight = 45
	navItemGap    = 4.5
)

var navigation = []struct {
	page  page
	icon  *ui.SVG
	label string
}{
	{pageDashboard, homeIcon, "Dashboard"},
	{pageAnalytics, chartIcon, "Analytics"},
	{pageOrders, cartIcon, "Orders"},
	{pageCustomers, usersIcon, "Customers"},
	{pageSettings, settingsIcon, "Settings"},
}

var (
	pillMotion  = ui.Animation{Spring: ui.SpringOptions{Duration: 240 * time.Millisecond}}
	themeMotion = ui.Animation{Duration: 200 * time.Millisecond}
	quickMotion = ui.Animation{Duration: 150 * time.Millisecond}
	fadeMotion  = ui.Animation{Duration: 120 * time.Millisecond}
	pressMotion = ui.Animation{Duration: 100 * time.Millisecond}
	pageEnter   = ui.ElementTransition{Colors: true, Duration: 180 * time.Millisecond, Enter: &ui.Motion{Y: 6, Opacity: 0.35}}
	dotExit     = ui.ElementTransition{Colors: true, Duration: 150 * time.Millisecond, Exit: &ui.Motion{Scale: 0.4}}
	clearPop    = ui.ElementTransition{Position: true, Duration: 120 * time.Millisecond, Enter: &ui.Motion{Scale: 0.6}, Exit: &ui.Motion{Scale: 0.6}}
	rollDigits  = ui.ContentTransition{Numeric: true, Animation: &pillMotion}
	deltaResize = ui.ElementTransition{Size: true, Animation: &pillMotion}
	deltaFollow = ui.ElementTransition{Position: true, Animation: &pillMotion}
	markSpring  = ui.Animation{Spring: ui.SpringOptions{Duration: 180 * time.Millisecond}}
	markMotion  = ui.ElementTransition{Position: true, Animation: &markSpring, Enter: &ui.Motion{Scale: 0.9}, Exit: &ui.Motion{Scale: 0.9}}
	tipResize   = ui.ElementTransition{Size: true, Animation: &markSpring}
	tipDigits   = ui.ContentTransition{Numeric: true, Animation: &fadeMotion}
	bellRing    = []ui.Keyframe{
		{Value: 0},
		{At: 70 * time.Millisecond, Value: 14},
		{At: 150 * time.Millisecond, Value: -10},
		{At: 230 * time.Millisecond, Value: 6},
		{At: 300 * time.Millisecond, Value: -3},
		{At: 360 * time.Millisecond, Value: 0},
	}
)

type menuEntry struct {
	label     string
	separated bool
	choose    func()
}

type app struct {
	page            page
	month           int
	query           string
	dark            bool
	alertsSeen      bool
	alertsOpen      bool
	bellRings       uint64
	accountOpen     bool
	sideAccountOpen bool
	shown           *palette
	setAppearance   func(mygo.ThemeSource)
}

func newApp(setAppearance func(mygo.ThemeSource)) *app {
	return &app{month: len(months) - 1, setAppearance: setAppearance}
}

func (a *app) colors() *palette {
	return a.shown
}

func (a *app) applyTheme() {
	if a.dark {
		a.setAppearance(mygo.ThemeDark)
		return
	}
	a.setAppearance(mygo.ThemeLight)
}

func (a *app) view(c *ui.Context) {
	a.shown = fadedPalette(progress(c.Root(), "theme", a.dark, themeMotion))
	p := a.colors()
	c.SetTheme(p.theme(a.dark))
	ui.Row(c).Fill().AlignItems(ui.Stretch).Children(func() {
		a.sidebar(c)
		ui.Column(c).Grow(1).Basis(0).MinWidth(0).LinearGradient(ui.LinearGradient{From: p.pageTop, To: p.pageBottom, Angle: 180, Start: 0.5, End: 1}).Children(func() {
			a.topbar(c)
			ui.Scroll(c.Key(a.page)).Grow(1).Padding(20, 18, 34, 20).Transition(pageEnter).Children(func() {
				ui.Column(c).Gap(20).Children(func() {
					switch a.page {
					case pageAnalytics:
						a.analyticsPage(c)
					case pageOrders:
						a.ordersPage(c)
					case pageCustomers:
						a.customersPage(c)
					case pageSettings:
						a.settingsPage(c)
					default:
						a.dashboardPage(c)
					}
				})
			})
		})
	})
}

func (a *app) sidebar(c *ui.Context) {
	p := a.colors()
	ui.Column(c).Width(237).Gradient(p.railTop, p.railBottom, 180).BorderWidth(0, 1, 0, 0).BorderColor(p.line).Children(func() {
		ui.Row(c).Height(69).Padding(c.TitleBar().Height/2, 0, 0, 23).Gap(7).DragWindow().Children(func() {
			logo(c, p)
			ui.Text(c, "Dashboard").FontSize(20.2).Bold()
		})
		nav := ui.Column(c).Margin(20, 13, 0, 13).Gap(navItemGap)
		pillY := nav.AnimateTo("pill", float32(a.page)*(navItemHeight+navItemGap), pillMotion)
		nav.Draw(func(pt *ui.Painter, r ui.Rect) {
			pt.Fill(ui.Rect{X: r.X, Y: r.Y + pillY, W: r.W, H: navItemHeight}, p.navActive, 12)
		}).Children(func() {
			for _, item := range navigation {
				a.navItem(c, item.page, item.icon, item.label)
			}
		})
		ui.Spacer(c)
		ui.Box(c).Height(1).Margin(0, 18, 0, 22).Background(p.line)
		account := menuButton(c, p, "Account menu", &a.sideAccountOpen).Height(46).Margin(16, 16, 34, 22).Padding(0, 6).Gap(16).Radius(12).Justify(ui.Start)
		account.Children(func() {
			avatar(c, p)
			ui.Column(c).Grow(1).Gap(2).Margin(0, 0, 4, 0).Children(func() {
				ui.Text(c, "John Doe").FontSize(14.5).FontWeight(600)
				ui.Text(c, "Admin").FontSize(12.5).TextColor(p.muted)
			})
			chevron(c, p, a.sideAccountOpen)
		})
		a.dropdown(c, account, &a.sideAccountOpen, true, a.accountEntries())
	})
}

func logo(c *ui.Context, p *palette) {
	ui.Box(c).Size(30, 30).Center().Children(func() {
		ui.Box(c).Size(24, 24).Radius(7.5).Background(p.accent).Rotate(45).Center().Children(func() {
			ui.Box(c).Size(13, 13).Radius(3).Background(p.onAccent)
		})
	})
}

func progress(e ui.Element, key string, on bool, motion ui.Animation) float32 {
	target := float32(0)
	if on {
		target = 1
	}
	return e.AnimateTo(key, target, motion)
}

func pressable(e ui.Element, scale float32) ui.Element {
	var pose ui.Pose
	if e.Pressed() {
		pose.Scale = scale
	}
	return e.AnimatePose(pose, pressMotion)
}

func labelWeight(selected bool) int {
	if selected {
		return 600
	}
	return 500
}

func (a *app) navItem(c *ui.Context, target page, icon *ui.SVG, label string) {
	p := a.colors()
	item := pressable(ui.ButtonBase(c), 0.98).Height(navItemHeight).Padding(0, 14, 0, 18).Gap(18).Radius(12).Justify(ui.Start).Label(label)
	selected := a.page == target
	item.Background(p.hover.Alpha(progress(item, "hover", item.Hovered() && !selected, fadeMotion)))
	color := p.body.Mix(p.accentInk, progress(item, "selected", selected, fadeMotion))
	item.Children(func() {
		ui.Icon(c, icon).Size(22, 22).TextColor(color)
		ui.Text(c, label).FontSize(16).FontWeight(labelWeight(selected)).TextColor(color)
	})
	if item.Clicked() {
		a.page = target
	}
}

func avatar(c *ui.Context, p *palette) {
	ui.Box(c).Size(38, 38).Radius(19).Background(p.avatar).Center().Children(func() {
		ui.Text(c, "JD").FontSize(15).FontWeight(600).TextColor(p.accentInk)
	})
}

func (a *app) topbar(c *ui.Context) {
	p := a.colors()
	ui.Row(c).Height(69).Padding(0, 16+c.TitleBar().Right, 0, 28).Background(p.bar).BorderWidth(0, 0, 1, 0).BorderColor(p.line).DragWindow().Children(func() {
		a.searchField(c)
		ui.Spacer(c)
		a.notifications(c)
		a.themeToggle(c, "Dark mode").Margin(0, 17, 0, 29)
		ui.Box(c).Size(1, 34).Background(p.line)
		account := menuButton(c, p, "Account", &a.accountOpen).Height(46).Margin(0, 0, 0, 13).Padding(0, 6).Gap(13).Radius(12)
		account.Children(func() {
			avatar(c, p)
			ui.Text(c, "John Doe").FontSize(15.5).FontWeight(500)
			chevron(c, p, a.accountOpen).Margin(0, 0, 0, 1)
		})
		a.dropdown(c, account, &a.accountOpen, false, a.accountEntries())
	})
}

func (a *app) searchField(c *ui.Context) {
	p := a.colors()
	var input ui.Element
	field := ui.Row(c).Width(480).MinWidth(200).Margin(0, 16, 0, 0).Height(40).Padding(0, 14).Gap(16).Radius(10).
		Background(p.field).FontSize(15.5).Cursor(ui.CursorText).Children(func() {
		ui.Icon(c, searchIcon).Size(20, 20).TextColor(p.muted)
		input = ui.TextInputBase(c, &a.query).Label("Search").Placeholder("Search for anything...").Grow(1).MinWidth(80)
		if a.query != "" {
			a.clearButton(c, input)
		}
	})
	line := p.fieldLine.Mix(p.line.Mix(p.ink, 0.25), progress(field, "hover", field.Hovered(), fadeMotion))
	field.Border(1, line.Mix(p.accent, progress(field, "focus", input.Focused(), fadeMotion)))
	if field.Clicked() {
		input.Focus()
	}
}

func (a *app) clearButton(c *ui.Context, input ui.Element) {
	p := a.colors()
	clear := pressable(ui.Box(c), 0.86).Size(21, 21).Radius(10.5).Label("Clear search").Role(ui.RoleButton).KeepFocus().Cursor(ui.CursorDefault).Transition(clearPop)
	clear.Background(p.muted.Alpha(0.5 + 0.3*progress(clear, "hover", clear.Hovered(), fadeMotion))).Draw(func(pt *ui.Painter, r ui.Rect) {
		k := r.W * 0.3
		var x ui.Path
		x.MoveTo(r.X+k, r.Y+k).LineTo(r.X+r.W-k, r.Y+r.H-k).MoveTo(r.X+r.W-k, r.Y+k).LineTo(r.X+k, r.Y+r.H-k)
		pt.StrokePath(&x, 1.5, p.pageTop)
	})
	if clear.Clicked() || input.Shortcut(0, ui.KeyEscape) {
		a.query = ""
		input.Focus()
	}
}

func (a *app) notifications(c *ui.Context) {
	p := a.colors()
	wasOpen := a.alertsOpen
	bell := menuButton(c, p, "Notifications", &a.alertsOpen).Size(34, 34).Radius(10)
	if a.alertsOpen && !wasOpen {
		a.bellRings++
	}
	a.alertsSeen = a.alertsSeen || a.alertsOpen
	bell.Children(func() {
		icon := ui.Icon(c, bellIcon).Size(26, 26).TextColor(p.body).TransformOrigin(0.5, 0.1)
		if a.bellRings > 0 {
			icon.Transform(ui.Rotation(icon.Keyframes("ring", bellRing, ui.KeyframeOptions{Run: a.bellRings})))
		}
		if !a.alertsSeen {
			ui.Box(c).Absolute().Top(2).Right(0).Size(10, 10).Radius(5).Background(p.accent).Transition(dotExit)
		}
	})
	entries := make([]menuEntry, 3)
	for i, o := range orders[:len(entries)] {
		entries[i] = menuEntry{label: o.id + "  " + o.status, choose: func() { a.page, a.query = pageOrders, o.id }}
	}
	a.dropdown(c, bell, &a.alertsOpen, false, entries)
}

func (a *app) themeToggle(c *ui.Context, label string) ui.Element {
	p := a.colors()
	toggle := ui.SwitchBase(c, &a.dark).Label(label).Size(58, 30).Radius(15).OnChange(a.applyTheme)
	icon := sunIcon
	if a.dark {
		icon = moonIcon
	}
	side := progress(toggle, "knob", a.dark, pillMotion)
	width := 26 + 5*progress(toggle, "hold", toggle.Pressed(), pressMotion)
	track := p.toggle.Mix(p.accent, 0.14*progress(toggle, "hover", toggle.Hovered(), fadeMotion))
	return toggle.Background(track).Children(func() {
		ui.Box(c).Absolute().Top(2).Left(2+side*(54-width)).Size(width, 26).Radius(13).Background(p.knob).
			Shadow(0, 1, 3, 0, p.shadow).Center().Children(func() {
			ui.Icon(c, icon).Size(17, 17).TextColor(p.ink).ContentTransition(ui.ContentTransition{})
		})
	})
}

func (a *app) accountEntries() []menuEntry {
	return []menuEntry{
		{label: "Settings", choose: func() { a.page = pageSettings }},
		{label: "Quit", separated: true, choose: func() { mygo.App.Quit() }},
	}
}

func menuButton(c *ui.Context, p *palette, label string, open *bool) ui.Element {
	button := pressable(ui.ButtonBase(c), 0.96).Label(label).Role(ui.RoleMenuButton)
	if button.Clicked() || !*open && button.Shortcut(0, ui.KeyDown) {
		*open = !*open
	}
	return button.Background(p.hover.Alpha(progress(button, "hover", *open || button.Hovered(), fadeMotion))).Expanded(*open)
}

func chevron(c *ui.Context, p *palette, open bool) ui.Element {
	var pose ui.Pose
	if open {
		pose.Rotate = 180
	}
	return ui.Icon(c, chevronIcon).Size(18, 18).TextColor(p.ink).AnimatePose(pose, quickMotion)
}

func (a *app) dropdown(c *ui.Context, anchor ui.Element, open *bool, upward bool, entries []menuEntry) {
	p := a.colors()
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
		ui.PopoverBase(c.Key(open), anchor, open, func(panel ui.Element) {
			panel.Role(ui.RoleMenu).FocusGroup(ui.Vertical).Margin(4, 0, 0, 0).MinWidth(180).Padding(4).Radius(10).
				Background(p.pageTop).Border(1, p.line).Shadow(0, 6, 20, 0, ui.RGBA(0, 0, 0, 0.18)).
				TransformOrigin(0, origin).Transition(motion)
			if panel.Shortcut(0, ui.KeyTab) || panel.Shortcut(ui.Shift, ui.KeyTab) {
				*open = false
			}
			rows := make([]ui.Element, len(entries))
			lastHovered := ui.Local(panel, "hovered", func() int { return -1 })
			hovered := -1
			for i, entry := range entries {
				if entry.separated {
					ui.Box(c).Height(1).Margin(4, 8).Background(p.line).Shrink(0)
				}
				rows[i] = pressable(ui.ButtonBase(c), 0.98).Label(entry.label).Role(ui.RoleMenuItem).FocusRing(false).Padding(6, 8).Radius(8).Justify(ui.Start).Shrink(0)
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
				row.Background(p.accent.Alpha(lit))
				color := p.ink.Mix(p.onAccent, lit)
				row.Children(func() {
					ui.Text(c, entries[i].label).SingleLine().TextColor(color)
				})
				if row.Clicked() {
					*open = false
					entries[i].choose()
				}
			}
		})
	})
	first.AutoFocus()
}
