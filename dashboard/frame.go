package main

import "github.com/egoist/mygo/ui"

const (
	navItemHeight = 45
	navItemGap    = 4.5
)

type navEntry struct {
	icon  *ui.SVG
	label string
}

type accountView struct {
	label    string
	initials string
	name     string
	role     string
	chevron  *ui.SVG
	open     *bool
	entries  []menuEntry
}

type bellView struct {
	label   string
	icon    *ui.SVG
	open    *bool
	seen    *bool
	rings   *uint64
	entries []menuEntry
}

func logo(c *ui.Context, p *palette) {
	ui.Box(c).Size(30, 30).Center().Children(func() {
		ui.Box(c).Size(24, 24).Radius(7.5).Background(p.accent).Rotate(45).Center().Children(func() {
			ui.Box(c).Size(13, 13).Radius(3).Background(p.onAccent)
		})
	})
}

func avatar(c *ui.Context, p *palette, initials string) {
	ui.Box(c).Size(38, 38).Radius(19).Background(p.avatar).Center().Children(func() {
		ui.Text(c, initials).FontSize(15).FontWeight(600).TextColor(p.accentInk)
	})
}

func navItem(c *ui.Context, p *palette, selected bool, icon *ui.SVG, label string) bool {
	item := pressable(ui.ButtonBase(c), 0.98).Height(navItemHeight).Padding(0, 14, 0, 18).Gap(18).Radius(12).Justify(ui.Start).Label(label)
	item.Background(p.hover.Alpha(progress(item, "hover", item.Hovered() && !selected, fadeMotion)))
	color := p.body.Mix(p.accentInk, progress(item, "selected", selected, fadeMotion))
	item.Children(func() {
		ui.Icon(c, icon).Size(22, 22).TextColor(color)
		ui.Text(c, label).FontSize(16).FontWeight(labelWeight(selected)).TextColor(color)
	})
	return item.Clicked()
}

func sidebar(c *ui.Context, p *palette, title string, items []navEntry, selected int, choose func(int), who accountView) ui.Element {
	return ui.Column(c).Width(237).Gradient(p.railTop, p.railBottom, 180).BorderWidth(0, 1, 0, 0).BorderColor(p.line).Children(func() {
		ui.Row(c).Height(69).Padding(c.TitleBar().Height/2, 0, 0, 23).Gap(7).DragWindow().Children(func() {
			logo(c, p)
			ui.Text(c, title).FontSize(20.2).Bold()
		})
		nav := ui.Column(c).Margin(20, 13, 0, 13).Gap(navItemGap)
		pillY := nav.AnimateTo("pill", float32(selected)*(navItemHeight+navItemGap), pillMotion)
		nav.Draw(func(pt *ui.Painter, r ui.Rect) {
			pt.Fill(ui.Rect{X: r.X, Y: r.Y + pillY, W: r.W, H: navItemHeight}, p.navActive, 12)
		}).Children(func() {
			for i, item := range items {
				if navItem(c, p, i == selected, item.icon, item.label) {
					choose(i)
				}
			}
		})
		ui.Spacer(c)
		ui.Box(c).Height(1).Margin(0, 18, 0, 22).Background(p.line)
		account := menuButton(c, p, who.label, who.open).Height(46).Margin(16, 16, 34, 22).Padding(0, 6).Gap(16).Radius(12).Justify(ui.Start)
		account.Children(func() {
			avatar(c, p, who.initials)
			ui.Column(c).Grow(1).Gap(2).Margin(0, 0, 4, 0).Children(func() {
				ui.Text(c, who.name).FontSize(14.5).FontWeight(600)
				ui.Text(c, who.role).FontSize(12.5).TextColor(p.muted)
			})
			chevron(c, p, who.chevron, *who.open)
		})
		dropdown(c, p, account, who.open, true, who.entries)
	})
}

func topbar(c *ui.Context, p *palette, content func()) ui.Element {
	return ui.Row(c).Height(69).Padding(0, 16+c.TitleBar().Right, 0, 28).Background(p.bar).BorderWidth(0, 0, 1, 0).BorderColor(p.line).DragWindow().Children(content)
}

func notificationBell(c *ui.Context, p *palette, b bellView) ui.Element {
	wasOpen := *b.open
	bell := menuButton(c, p, b.label, b.open).Size(34, 34).Radius(10)
	if *b.open && !wasOpen {
		*b.rings++
	}
	*b.seen = *b.seen || *b.open
	bell.Children(func() {
		icon := ui.Icon(c, b.icon).Size(26, 26).TextColor(p.body).TransformOrigin(0.5, 0.1)
		if *b.rings > 0 {
			icon.Transform(ui.Rotation(icon.Keyframes("ring", bellRing, ui.KeyframeOptions{Run: *b.rings})))
		}
		if !*b.seen {
			ui.Box(c).Absolute().Top(2).Right(0).Size(10, 10).Radius(5).Background(p.accent).Transition(dotExit)
		}
	})
	dropdown(c, p, bell, b.open, false, b.entries)
	return bell
}

func appFrame(c *ui.Context, p *palette, pageKey any, rail, bar, body func()) ui.Element {
	return ui.Row(c).Fill().AlignItems(ui.Stretch).Children(func() {
		rail()
		ui.Column(c).Grow(1).Basis(0).MinWidth(0).LinearGradient(ui.LinearGradient{From: p.pageTop, To: p.pageBottom, Angle: 180, Start: 0.5, End: 1}).Children(func() {
			bar()
			ui.Scroll(c.Key(pageKey)).Grow(1).Padding(20, 18, 34, 20).Transition(pageEnter).Children(func() {
				ui.Column(c).Gap(20).Children(body)
			})
		})
	})
}
