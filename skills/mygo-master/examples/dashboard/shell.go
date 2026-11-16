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

var bellRing = []ui.Keyframe{
	{Value: 0},
	{At: 70 * time.Millisecond, Value: 14},
	{At: 150 * time.Millisecond, Value: -10},
	{At: 230 * time.Millisecond, Value: 6},
	{At: 300 * time.Millisecond, Value: -3},
	{At: 360 * time.Millisecond, Value: 0},
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
	darkness := progress(c.Root(), "theme", a.dark, themeMotion)
	a.shown = fadedPalette(darkness)
	p := a.colors()
	c.SetTheme(fadedTheme(darkness))
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
			chevron(c, p, chevronIcon, a.sideAccountOpen)
		})
		dropdown(c, p, account, &a.sideAccountOpen, true, a.accountEntries())
	})
}

func logo(c *ui.Context, p *palette) {
	ui.Box(c).Size(30, 30).Center().Children(func() {
		ui.Box(c).Size(24, 24).Radius(7.5).Background(p.accent).Rotate(45).Center().Children(func() {
			ui.Box(c).Size(13, 13).Radius(3).Background(p.onAccent)
		})
	})
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
		searchField(c, p, &a.query, searchIcon, "Search for anything...")
		ui.Spacer(c)
		a.notifications(c)
		themeSwitch(c, p, &a.dark, "Dark mode", sunIcon, moonIcon, a.applyTheme).Margin(0, 17, 0, 29)
		ui.Box(c).Size(1, 34).Background(p.line)
		account := menuButton(c, p, "Account", &a.accountOpen).Height(46).Margin(0, 0, 0, 13).Padding(0, 6).Gap(13).Radius(12)
		account.Children(func() {
			avatar(c, p)
			ui.Text(c, "John Doe").FontSize(15.5).FontWeight(500)
			chevron(c, p, chevronIcon, a.accountOpen).Margin(0, 0, 0, 1)
		})
		dropdown(c, p, account, &a.accountOpen, false, a.accountEntries())
	})
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
	dropdown(c, p, bell, &a.alertsOpen, false, entries)
}

func (a *app) accountEntries() []menuEntry {
	return []menuEntry{
		{label: "Settings", choose: func() { a.page = pageSettings }},
		{label: "Quit", separated: true, choose: func() { mygo.App.Quit() }},
	}
}
