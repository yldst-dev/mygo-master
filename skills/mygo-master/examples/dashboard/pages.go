package main

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/chart"
)

const (
	compactWidth = 1394
	monthWidth   = 49.66
	tipPadding   = 8
	tipFontSize  = 12.5
	donutRadius  = 91
	donutHole    = 54
	donutGrowth  = 5
	donutSink    = 3
	donutFade    = 0.8
)

func card(c *ui.Context, p *palette) ui.Element {
	return ui.Column(c).Radius(16).Background(p.card).Border(1, p.line).Shadow(0, 1, 3, 0, p.shadow)
}

func cardTitle(c *ui.Context, title string) {
	ui.Text(c, title).FontSize(20.5).Bold().Role(ui.RoleHeading).Margin(0, 0, 4, 0)
}

func (a *app) dashboardPage(c *ui.Context) {
	p := a.colors()
	width, _ := c.Size()
	columns := 4
	if width < compactWidth {
		columns = 2
	}
	ui.Grid(c).Columns(columns).GapX(17).GapY(17).Children(func() {
		for _, m := range metrics {
			a.metricCard(c, m)
		}
	})
	if width < compactWidth {
		ui.Column(c).Gap(18).Children(func() {
			a.revenueCard(c).Height(289)
			a.categoryCard(c).Height(289)
		})
	} else {
		ui.Row(c).Gap(18).Children(func() {
			a.revenueCard(c).Height(289).Grow(1).Basis(0).MinWidth(0)
			a.categoryCard(c).Height(289).Width(403)
		})
	}
	card(c, p).Padding(15, 16, 20, 22).Gap(13).Children(func() {
		ui.Row(c).Justify(ui.SpaceBetween).Children(func() {
			cardTitle(c, "Recent Orders")
			view := pressable(ui.ButtonBase(c), 0.96).Height(30).Padding(0, 14).Gap(5).Radius(8).Border(1, p.fieldLine)
			var nudge ui.Pose
			if view.Hovered() {
				nudge.X = 3
			}
			view.Background(p.field.Mix(p.navActive, progress(view, "hover", view.Hovered(), fadeMotion))).Children(func() {
				ui.Text(c, "View All").FontSize(14).FontWeight(600).TextColor(p.accentInk)
				ui.Icon(c, arrowIcon).Size(15, 15).TextColor(p.accentInk).AnimatePose(nudge, quickMotion)
			})
			if view.Clicked() {
				a.page, a.query = pageOrders, ""
			}
		})
		recent := filterOrders(a.query)
		if len(recent) > 5 {
			recent = recent[:5]
		}
		a.ordersTable(c, recent)
	})
}

func (a *app) metricCard(c *ui.Context, m metric) {
	p := a.colors()
	value, previous := m.at(a.month)
	card(c, p).Height(140).Padding(17, 14, 0, 22).Children(func() {
		ui.Row(c).Gap(18).Children(func() {
			ui.Icon(c, m.icon).Size(24, 24).TextColor(p.body)
			ui.Text(c, m.label).FontSize(14.5).FontWeight(500).TextColor(p.body).NoWrap()
		})
		ui.Box(c).ClipY().Margin(9, 0, 0, 0).Children(func() {
			ui.Text(c, m.format(value)).FontSize(29.5).Bold().NoWrap().ContentTransition(rollDigits)
		})
		ui.Row(c).Gap(7).Margin(11, 0, 0, 0).Children(func() {
			deltaPill(c, p, (value-previous)/previous*100)
			ui.Text(c, "vs. last month").FontSize(12.5).TextColor(p.muted).NoWrap().Transition(deltaFollow)
		})
		sparkline(c, p, m).Label(m.label+" trend").Absolute().Top(53).Right(14).Size(85, 45)
	})
}

func deltaPill(c *ui.Context, p *palette, delta float64) {
	icon, sign := arrowUpIcon, "+"
	if delta < 0 {
		icon, sign = arrowDownIcon, "−"
	}
	pill := ui.Row(c).Height(24).Padding(0, 9, 0, 7).Gap(4).Radius(12).ClipY().Transition(deltaResize)
	falling := progress(pill, "falling", delta < 0, fadeMotion)
	ink := p.good.ink.Mix(p.bad.ink, falling)
	pill.Background(p.good.background.Mix(p.bad.background, falling)).Children(func() {
		ui.Icon(c, icon).Size(14, 14).TextColor(ink).ContentTransition(ui.ContentTransition{})
		ui.Text(c, fmt.Sprintf("%s%.1f%%", sign, math.Abs(delta))).FontSize(12.5).FontWeight(600).TextColor(ink).NoWrap().ContentTransition(rollDigits)
	})
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
			Role(ui.RoleTooltip).ClipY().Transition(tipResize).Children(func() {
			ui.Text(c, text).FontSize(tipFontSize).TextColor(p.pageTop).NoWrap().ContentTransition(tipDigits)
		})
	}).IgnoreInput(true).Transition(markMotion)
}

func sparkline(c *ui.Context, p *palette, m metric) ui.Element {
	values := m.trend()
	lo, hi := chart.Extent(values...)
	scale := chart.NewLinear(lo, hi)
	xs, ys := make([]float32, len(values)), make([]float32, len(values))
	for i, v := range values {
		xs[i] = float32(i) / float32(len(values)-1)
		ys[i] = 0.2 + 0.75*scale.At(v)
	}
	spark := ui.Box(c).Radius(6)
	describe := func(i int) string { return trendMonths[i] + ": " + m.format(values[i]) }
	selected := pickPoint(c, spark, xs, describe)
	lit := progress(spark, "lit", selected >= 0, fadeMotion)
	return spark.Draw(func(pt *ui.Painter, r ui.Rect) {
		pt.FillPathGradient(chart.Area(r, xs, nil, ys, chart.CurveMonotone), ui.LinearGradient{From: p.accent.Alpha(0.14 + 0.1*lit), To: p.accent.Alpha(0), Angle: 180})
		pt.StrokePath(chart.Line(r, xs, ys, chart.CurveMonotone), 2+0.5*lit, p.accent)
	}).Children(func() {
		if selected >= 0 {
			chart.Point(c, xs[selected], ys[selected], 9, 9).Radius(4.5).Background(p.accent).Border(1.5, p.card).Transition(markMotion)
			chartTip(c, p, spark.Bounds().W, xs[selected], 1, describe(selected))
		}
	})
}

func (a *app) revenueCard(c *ui.Context) ui.Element {
	p := a.colors()
	return card(c, p).Padding(15, 19, 17, 22).Gap(26).Children(func() {
		ui.Row(c).Justify(ui.SpaceBetween).Children(func() {
			cardTitle(c, "Revenue (Monthly)")
			a.monthPicker(c)
		})
		a.revenueChart(c)
	})
}

func (a *app) monthPicker(c *ui.Context) {
	p := a.colors()
	picker := ui.SegmentedBase(c.Key("month"), &a.month, len(months))
	pillX := picker.Track.AnimateTo("pill", float32(a.month)*monthWidth, pillMotion)
	picker.Track.Label("Month").Height(30).Radius(8).Border(1, p.line).Draw(func(pt *ui.Painter, r ui.Rect) {
		pt.Fill(ui.Rect{X: r.X + 1 + pillX, Y: r.Y + 1, W: monthWidth, H: 28}, p.accent, 7)
	}).Children(func() {
		for i, name := range months {
			segment := pressable(picker.Segment(i), 0.94).Size(monthWidth, 28).Radius(7).Center()
			selected := i == a.month
			segment.Background(p.hover.Alpha(progress(segment, "hover", segment.Hovered() && !selected, fadeMotion)))
			offset := pillX - float32(i)*monthWidth
			covered := max(0, 1-max(offset, -offset)/monthWidth)
			segment.Children(func() {
				ui.Text(c, name).FontSize(14).FontWeight(labelWeight(selected)).TextColor(p.body.Mix(p.onAccent, covered))
			})
		}
	})
}

func (a *app) revenueChart(c *ui.Context) {
	p := a.colors()
	values := metrics[0].monthly[:]
	scale := chart.NewLinear(0, 150000)
	left := make([]chart.Tick, 7)
	for i := range left {
		v := 25000 * float64(i)
		label := "$0"
		if v > 0 {
			label = fmt.Sprintf("$%.0fK", v/1000)
		}
		left[i] = chart.Tick{At: scale.At(v), Value: v, Label: label}
	}
	xs, ys := make([]float32, len(values)), make([]float32, len(values))
	bottom := make([]chart.Tick, len(values))
	for i, v := range values {
		xs[i] = 0.025 + 0.1872*float32(i)
		ys[i] = scale.At(v)
		bottom[i] = chart.Tick{At: xs[i], Label: months[i]}
	}
	label := func(_ chart.Side, tick chart.Tick) {
		ui.Text(c, tick.Label).NoWrap().FontSize(14).TextColor(p.muted)
	}
	chart.Plot(c, chart.Axes{Left: left, Bottom: bottom, Label: label}, func(area ui.Element) {
		shadedXs, shadedYs := shadeUpTo(xs, ys, area.AnimateTo("shade", float32(a.month), pillMotion))
		area.Label("Monthly revenue").Draw(func(pt *ui.Painter, r ui.Rect) {
			for _, tick := range left {
				color := p.grid
				if tick.Value == 0 {
					color = p.axis
				}
				pt.Line(r.X, chart.Y(r, tick.At), r.X+r.W, chart.Y(r, tick.At), 1, color)
			}
			pt.Line(r.X, r.Y, r.X, r.Y+r.H, 1, p.grid)
			pt.FillPathGradient(chart.Area(r, shadedXs, nil, shadedYs, chart.CurveLinear), ui.LinearGradient{
				Angle: 180,
				Stops: []ui.GradientStop{{Color: p.areaStrong, At: 0}, {Color: p.areaStrong, At: 0.55}, {Color: p.areaFaint, At: 1}},
			})
			pt.StrokePath(chart.Line(r, xs, ys, chart.CurveLinear), 2.5, p.accent)
		})
		describe := func(i int) string { return months[i] + ": " + thousands(values[i]) }
		selected := pickPoint(c, area, xs, describe)
		if selected >= 0 {
			chart.RuleX(c.Key("guide"), xs[selected], 1).Background(p.axis).Transition(markMotion)
		}
		for i := range values {
			dot := chart.Point(c.Key(i), xs[i], ys[i], 10, 10).Radius(5).Border(2, p.accent)
			var pose ui.Pose
			if i == selected {
				pose.Scale = 1.3
			}
			dot.Background(p.card.Mix(p.accent, progress(dot, "filled", i == selected, quickMotion))).AnimatePose(pose, quickMotion)
		}
		if selected >= 0 {
			chartTip(c, p, area.Bounds().W, xs[selected], ys[selected], describe(selected))
		}
	}).Grow(1).Margin(0, 3, 0, 3).GapX(13).GapY(6)
}

func shadeUpTo(xs, ys []float32, at float32) ([]float32, []float32) {
	n := min(int(at), len(xs)-1)
	k := at - float32(n)
	if k <= 0 || n == len(xs)-1 {
		return xs[:n+1], ys[:n+1]
	}
	return append(xs[:n+1:n+1], xs[n]+(xs[n+1]-xs[n])*k), append(ys[:n+1:n+1], ys[n]+(ys[n+1]-ys[n])*k)
}

func (a *app) categoryCard(c *ui.Context) ui.Element {
	p := a.colors()
	shares := categoryShares[a.month]
	total, _ := metrics[0].at(a.month)
	return card(c, p).Padding(16, 19, 0, 20).Gap(23).Children(func() {
		cardTitle(c, "Sales by Category")
		ui.Row(c).Gap(23).Children(func() {
			donut := ui.Box(c).Size(2*donutRadius, 2*donutRadius).Shrink(0).Center().Label("Sales by category")
			shown := make([]float64, len(shares))
			for i, share := range shares {
				shown[i] = float64(donut.AnimateTo(i, float32(share), pillMotion))
			}
			slices := chart.Pie(shown, 0.03)
			active := -1
			if x, y, over := donut.PointerPosition(); over && donut.Hovered() {
				if i, ok := chart.SliceAt(slices, x-donutRadius, y-donutRadius, donutHole, donutRadius+donutGrowth); ok {
					active = i
				}
			}
			lit := make([]float32, len(slices))
			var dim float32
			ui.Column(c).Grow(1).Gap(15).Margin(0, 0, 10, 0).Children(func() {
				rows := make([]ui.Element, len(categories))
				for i, name := range categories {
					rows[i] = ui.Row(c).Height(20).Gap(10).ClipY().Focusable().Label(name)
					if rows[i].Hovered() || rows[i].FocusVisible() {
						active = i
					}
				}
				dim = progress(donut, "dim", active >= 0, fadeMotion)
				for i, name := range categories {
					lit[i] = progress(donut, name, i == active, fadeMotion)
					rows[i].TextColor(p.ink.Mix(p.muted, dim-lit[i])).Children(func() {
						ui.Box(c).Size(12, 12).Radius(6).Background(p.slices[i].Mix(p.card, donutFade*(dim-lit[i])))
						ui.Text(c, name).FontSize(13.5).Grow(1)
						ui.Text(c, fmt.Sprintf("%.0f%%", shares[i])).FontSize(13.5).ContentTransition(rollDigits)
					})
				}
			})
			donut.Draw(func(pt *ui.Painter, r ui.Rect) {
				cx, cy := r.X+r.W/2, r.Y+r.H/2
				for i, slice := range slices {
					pt.FillPath(chart.Sector(cx, cy, donutHole-donutSink*lit[i], donutRadius+donutGrowth*lit[i], slice), p.slices[i].Mix(p.card, donutFade*(dim-lit[i])))
				}
			}).Children(func() {
				ui.Box(c).ClipY().Children(func() {
					ui.Text(c, thousands(total)).FontSize(21).Bold().ContentTransition(rollDigits)
				})
				ui.Text(c, "Total").FontSize(14).TextColor(p.muted).Margin(3, 0, 0, 0)
			})
			ui.Box(c).Absolute().Left(0).Top(0).Size(2*donutRadius, 2*donutRadius).PassThrough().Children(func() {
				if active >= 0 {
					x, y := chart.Polar(0.5, 0.5, 0.4, slices[active].Mid())
					chartTip(c, p, 2*donutRadius, x, 1-y, categories[active]+": "+thousands(total*shares[active]/100))
				}
			})
		})
	})
}

type column struct {
	title  string
	weight float32
}

func (a *app) table(c *ui.Context, label string, columns []column, rows int, cell func(row, col int)) {
	p := a.colors()
	ui.Column(c).Label(label).Role(ui.RoleTable).Children(func() {
		ui.Row(c).Height(37).Padding(0, 12, 3, 12).Radius(6).Background(p.head).Role(ui.RoleRow).Children(func() {
			for _, col := range columns {
				ui.Text(c, col.title).FontSize(14.5).FontWeight(500).TextColor(p.body).Grow(col.weight).Basis(0).Role(ui.RoleColumnHeader)
			}
		})
		for row := range rows {
			ui.Row(c).Height(38).PaddingX(12).BorderWidth(0, 0, 1, 0).BorderColor(p.line).Role(ui.RoleRow).Children(func() {
				for i, col := range columns {
					ui.Row(c).Grow(col.weight).Basis(0).Role(ui.RoleCell).Children(func() {
						cell(row, i)
					})
				}
			})
		}
		if rows == 0 {
			ui.Text(c, "Nothing matches your search.").FontSize(14.5).TextColor(p.muted).Padding(12)
		}
	})
}

func cellText(c *ui.Context, text string) {
	ui.Text(c, text).FontSize(14.5).NoWrap()
}

func (a *app) ordersTable(c *ui.Context, rows []order) {
	p := a.colors()
	columns := []column{{"Order ID", 203}, {"Customer", 254}, {"Status", 289}, {"Amount", 163}, {"Date", 192}}
	a.table(c, "Orders", columns, len(rows), func(row, col int) {
		o := rows[row]
		switch col {
		case 0:
			open := pressable(ui.ButtonBase(c), 0.97).Label(o.id).Radius(4).Cursor(ui.CursorPointer)
			line := progress(open, "line", open.Hovered(), quickMotion)
			open.Draw(func(pt *ui.Painter, r ui.Rect) {
				pt.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - 1, W: r.W * line, H: 1}, p.accentInk, 0)
			}).Children(func() {
				ui.Text(c, o.id).FontSize(14.5).TextColor(p.accentInk)
			})
			if open.Clicked() {
				a.page, a.query = pageOrders, o.id
			}
		case 1:
			cellText(c, o.customer)
		case 2:
			statusPill(c, p.status(o.status), o.status)
		case 3:
			cellText(c, cents(o.amount))
		case 4:
			cellText(c, o.date)
		}
	})
}

func statusPill(c *ui.Context, colors tone, label string) {
	ui.Row(c).Height(24).Padding(0, 12, 0, 10).Gap(9).Radius(12).Background(colors.background).Children(func() {
		ui.Box(c).Size(8, 8).Radius(4).Background(colors.dot)
		ui.Text(c, label).FontSize(13.5).TextColor(colors.ink).NoWrap()
	})
}

func (a *app) analyticsPage(c *ui.Context) {
	p := a.colors()
	a.revenueCard(c).Height(360)
	card(c, p).Padding(15, 16, 20, 22).Gap(13).Children(func() {
		cardTitle(c, "Monthly Breakdown")
		columns := []column{{"Month", 1}, {"Total Revenue", 1}, {"Active Users", 1}, {"Orders", 1}, {"Conversion Rate", 1}}
		a.table(c, "Monthly breakdown", columns, len(months), func(row, col int) {
			if col == 0 {
				cellText(c, months[row])
				return
			}
			m := metrics[col-1]
			cellText(c, m.format(m.monthly[row]))
		})
	})
}

func (a *app) ordersPage(c *ui.Context) {
	p := a.colors()
	rows := filterOrders(a.query)
	card(c, p).Padding(15, 16, 20, 22).Gap(13).Children(func() {
		ui.Row(c).Justify(ui.SpaceBetween).Height(30).Children(func() {
			cardTitle(c, "All Orders")
			ui.Text(c, fmt.Sprintf("%d of %d", len(rows), len(orders))).FontSize(14).TextColor(p.muted)
		})
		a.ordersTable(c, rows)
	})
}

func (a *app) customersPage(c *ui.Context) {
	p := a.colors()
	rows := summarizeCustomers(filterOrders(a.query))
	card(c, p).Padding(15, 16, 20, 22).Gap(13).Children(func() {
		ui.Row(c).Height(30).Children(func() {
			cardTitle(c, "All Customers")
		})
		columns := []column{{"Customer", 2}, {"Orders", 1}, {"Total Spent", 1}, {"Last Order", 1}}
		a.table(c, "Customers", columns, len(rows), func(row, col int) {
			customer := rows[row]
			switch col {
			case 0:
				cellText(c, customer.name)
			case 1:
				cellText(c, fmt.Sprint(customer.orders))
			case 2:
				cellText(c, cents(customer.spent))
			case 3:
				cellText(c, customer.last)
			}
		})
	})
}

func (a *app) settingsPage(c *ui.Context) {
	p := a.colors()
	card(c, p).Padding(15, 21, 20, 22).Gap(13).Children(func() {
		cardTitle(c, "Appearance")
		ui.Row(c).Height(38).Justify(ui.SpaceBetween).Children(func() {
			ui.Text(c, "Dark mode").FontSize(14.5).TextColor(p.body)
			a.themeToggle(c, "Use dark mode")
		})
	})
}
