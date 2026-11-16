package main

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/chart"
)

const (
	tipPadding  = 8
	tipFontSize = 12.5
)

func card(c *ui.Context, p *palette) ui.Element {
	return ui.Column(c).Radius(16).Background(p.card).Border(1, p.line).Shadow(0, 1, 3, 0, p.shadow)
}

func cardTitle(c *ui.Context, title string) {
	ui.Text(c, title).FontSize(20.5).Bold().Role(ui.RoleHeading).Margin(0, 0, 4, 0)
}

func statusPill(c *ui.Context, colors tone, label string) {
	ui.Row(c).Height(24).Padding(0, 12, 0, 10).Gap(9).Radius(12).Background(colors.background).Children(func() {
		ui.Box(c).Size(8, 8).Radius(4).Background(colors.dot)
		ui.Text(c, label).FontSize(13.5).TextColor(colors.ink).NoWrap()
	})
}

type column struct {
	title  string
	weight float32
}

func table(c *ui.Context, p *palette, label, empty string, columns []column, rows int, cell func(row, col int)) {
	ui.Column(c).Label(label).Role(ui.RoleTable).Children(func() {
		ui.Row(c).Height(37).Padding(0, 12, 3, 12).Radius(6).Background(p.head).Role(ui.RoleRow).Children(func() {
			for _, col := range columns {
				ui.Text(c, col.title).FontSize(14.5).FontWeight(500).TextColor(p.body).Grow(col.weight).Basis(0).Role(ui.RoleColumnHeader)
			}
		})
		for row := range rows {
			ui.Row(c).Height(38).PaddingX(12).BorderWidth(0, 0, 1, 0).BorderColor(p.line).Role(ui.RoleRow).Children(func() {
				for i, col := range columns {
					ui.Row(c).Grow(col.weight).Basis(0).FillHeight().Role(ui.RoleCell).Children(func() {
						cell(row, i)
					})
				}
			})
		}
		if rows == 0 {
			ui.Text(c, empty).FontSize(14.5).TextColor(p.muted).Padding(12)
		}
	})
}

func cellText(c *ui.Context, text string) {
	ui.Text(c, text).FontSize(14.5).NoWrap()
}

type menuEntry struct {
	label     string
	separated bool
	choose    func()
}

func menuButton(c *ui.Context, p *palette, label string, open *bool) ui.Element {
	button := pressable(ui.ButtonBase(c), 0.96).Label(label).Role(ui.RoleMenuButton)
	if button.Clicked() || !*open && button.Shortcut(0, ui.KeyDown) {
		*open = !*open
	}
	return button.Background(p.hover.Alpha(progress(button, "hover", *open || button.Hovered(), fadeMotion))).Expanded(*open)
}

func chevron(c *ui.Context, p *palette, icon *ui.SVG, open bool) ui.Element {
	var pose ui.Pose
	if open {
		pose.Rotate = 180
	}
	return ui.Icon(c, icon).Size(18, 18).TextColor(p.ink).AnimatePose(pose, quickMotion)
}

func dropdown(c *ui.Context, p *palette, anchor ui.Element, open *bool, upward bool, entries []menuEntry) {
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
				rows[i] = pressable(ui.ButtonBase(c), 0.98).Label(entry.label).Role(ui.RoleMenuItem).FocusRing(false).MinHeight(30).Padding(6, 8).Radius(8).Justify(ui.Start).Shrink(0)
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

func segmented(c *ui.Context, p *palette, key, label string, selected *int, names []string, width float32) {
	picker := ui.SegmentedBase(c.Key(key), selected, len(names))
	pillX := picker.Track.AnimateTo("pill", float32(*selected)*width, pillMotion)
	picker.Track.Label(label).Height(30).Radius(8).Border(1, p.line).Draw(func(pt *ui.Painter, r ui.Rect) {
		pt.Fill(ui.Rect{X: r.X + 1 + pillX, Y: r.Y + 1, W: width, H: 28}, p.accent, 7)
	}).Children(func() {
		for i, name := range names {
			segment := pressable(picker.Segment(i), 0.94).Label(name).Size(width, 28).Radius(7).Center()
			chosen := i == *selected
			segment.Background(p.hover.Alpha(progress(segment, "hover", segment.Hovered() && !chosen, fadeMotion)))
			offset := pillX - float32(i)*width
			covered := max(0, 1-max(offset, -offset)/width)
			segment.Children(func() {
				ui.Text(c, name).FontSize(14).FontWeight(labelWeight(chosen)).TextColor(p.body.Mix(p.onAccent, covered))
			})
		}
	})
}

func searchField(c *ui.Context, p *palette, query *string, icon *ui.SVG, placeholder string) ui.Element {
	var input ui.Element
	field := ui.Row(c).Width(480).MinWidth(200).Margin(0, 16, 0, 0).Height(40).Padding(0, 14).Gap(16).Radius(10).
		Background(p.field).FontSize(15.5).Cursor(ui.CursorText).Children(func() {
		ui.Icon(c, icon).Size(20, 20).TextColor(p.muted)
		input = ui.TextInputBase(c, query).Label("Search").Placeholder(placeholder).Grow(1).MinWidth(80)
		if *query != "" {
			clearButton(c, p, query, input)
		}
	})
	line := p.fieldLine.Mix(p.line.Mix(p.ink, 0.25), progress(field, "hover", field.Hovered(), fadeMotion))
	field.Border(1, line.Mix(p.accent, progress(field, "focus", input.Focused(), fadeMotion)))
	if field.Clicked() {
		input.Focus()
	}
	return field
}

func clearButton(c *ui.Context, p *palette, query *string, input ui.Element) {
	clear := pressable(ui.ButtonBase(c), 0.86).Label("Clear search").Size(28, 28).Center().KeepFocus().Cursor(ui.CursorDefault).Transition(popMotion)
	hover := progress(clear, "hover", clear.Hovered(), fadeMotion)
	clear.Children(func() {
		ui.Box(c).Size(21, 21).Radius(10.5).Background(p.muted.Alpha(0.5 + 0.3*hover)).Draw(func(pt *ui.Painter, r ui.Rect) {
			k := r.W * 0.3
			var x ui.Path
			x.MoveTo(r.X+k, r.Y+k).LineTo(r.X+r.W-k, r.Y+r.H-k).MoveTo(r.X+r.W-k, r.Y+k).LineTo(r.X+k, r.Y+r.H-k)
			pt.StrokePath(&x, 1.5, p.pageTop)
		})
	})
	if clear.Clicked() || input.Shortcut(0, ui.KeyEscape) {
		*query = ""
		input.Focus()
	}
}

func themeSwitch(c *ui.Context, p *palette, on *bool, label string, sun, moon *ui.SVG, changed func()) ui.Element {
	toggle := ui.SwitchBase(c, on).Label(label).Size(58, 30).Radius(15).OnChange(changed)
	icon := sun
	if *on {
		icon = moon
	}
	side := progress(toggle, "knob", *on, pillMotion)
	width := 26 + 5*progress(toggle, "hold", toggle.Pressed(), pressMotion)
	track := p.toggle.Mix(p.accent, 0.14*progress(toggle, "hover", toggle.Hovered(), fadeMotion))
	return toggle.Background(track).Children(func() {
		ui.Box(c).Absolute().Top(2).Left(2+side*(54-width)).Size(width, 26).Radius(13).Background(p.knob).
			Shadow(0, 1, 3, 0, p.shadow).Center().Children(func() {
			ui.Icon(c, icon).Size(17, 17).TextColor(p.ink).ContentTransition(ui.ContentTransition{})
		})
	})
}

func labelWeight(selected bool) int {
	if selected {
		return 600
	}
	return 500
}

func pickPoint(c *ui.Context, area ui.Element, xs []float32, describe func(int) string) int {
	selected := ui.Local(area, "selected", func() int { return -1 })
	chart.Select(c, area, selected, xs, describe)
	if !area.Hovered() && !area.Focused() {
		*selected = -1
	}
	return *selected
}

func chartTip(c *ui.Context, p *palette, width, x, y float32, text string) {
	textWidth, _ := c.MeasureText(0, ui.Span{Text: text, Size: tipFontSize})
	reach := (tipPadding + textWidth/2) / width
	chart.Anchor(c.Key("tip"), min(max(x, reach), 1-reach), y, ui.Center, ui.End, func() {
		ui.Row(c).Height(26).PaddingX(tipPadding).Margin(0, 0, 8, 0).Radius(5).Background(p.ink).Shadow(0, 2, 8, 0, ui.RGBA(0, 0, 0, 0.2)).
			Role(ui.RoleTooltip).ClipY().Transition(tipFit).Children(func() {
			ui.Text(c, text).FontSize(tipFontSize).TextColor(p.pageTop).NoWrap().ContentTransition(tipDigits)
		})
	}).IgnoreInput(true).Transition(markMotion)
}
