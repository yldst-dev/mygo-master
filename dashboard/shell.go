package main

import (
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

var navigation = []navEntry{
	{homeIcon, "Dashboard"},
	{chartIcon, "Analytics"},
	{cartIcon, "Orders"},
	{usersIcon, "Customers"},
	{settingsIcon, "Settings"},
}

type app struct {
	page          page
	month         int
	query         string
	dark          bool
	alertsSeen    bool
	alertsOpen    bool
	bellRings     uint64
	accountOpen   bool
	shown         *palette
	setAppearance func(mygo.ThemeSource)
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
	appFrame(c, p, a.page, func() { a.rail(c) }, func() { a.bar(c) }, func() {
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
}

func (a *app) rail(c *ui.Context) {
	sidebar(c, a.colors(), "Dashboard", navigation, int(a.page), func(i int) { a.page = page(i) }, accountView{
		label:    "Account menu",
		initials: "JD",
		name:     "John Doe",
		role:     "Admin",
		chevron:  chevronIcon,
		open:     &a.accountOpen,
		entries:  a.accountEntries(),
	})
}

func (a *app) bar(c *ui.Context) {
	p := a.colors()
	topbar(c, p, func() {
		searchField(c, p, &a.query, searchIcon, "Search for anything...")
		ui.Spacer(c)
		notificationBell(c, p, bellView{
			label:   "Notifications",
			icon:    bellIcon,
			open:    &a.alertsOpen,
			seen:    &a.alertsSeen,
			rings:   &a.bellRings,
			entries: a.alertEntries(),
		})
		themeSwitch(c, p, &a.dark, "Dark mode", sunIcon, moonIcon, a.applyTheme).Margin(0, 2, 0, 29)
	})
}

func (a *app) alertEntries() []menuEntry {
	entries := make([]menuEntry, 3)
	for i, o := range orders[:len(entries)] {
		entries[i] = menuEntry{label: o.id + "  " + o.status, choose: func() { a.page, a.query = pageOrders, o.id }}
	}
	return entries
}

func (a *app) accountEntries() []menuEntry {
	return []menuEntry{
		{label: "Settings", choose: func() { a.page = pageSettings }},
		{label: "Quit", separated: true, choose: func() { mygo.App.Quit() }},
	}
}
