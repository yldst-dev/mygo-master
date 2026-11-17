package main

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/chart"
)

const (
	tipPadding  = 8
	tipFontSize = 12.5

	areaInset = 0.025
	areaSpan  = 0.936

	donutRadius = 91
	donutHole   = 54
	donutGrowth = 5
	donutSink   = 3
	donutFade   = 0.8
)

type areaChartView struct {
	label     string
	names     []string
	values    []float64
	ceiling   float64
	step      float64
	axisLabel func(v float64) string
	shaded    int
	describe  func(i int) string
}

type donutView struct {
	label    string
	names    []string
	shares   []float64
	center   string
	caption  string
	describe func(i int) string
}

func sparkline(c *ui.Context, p *palette, values []float64, describe func(i int) string) ui.Element {
	lo, hi := chart.Extent(values...)
	scale := chart.NewLinear(lo, hi)
	xs, ys := make([]float32, len(values)), make([]float32, len(values))
	for i, v := range values {
		xs[i] = float32(i) / float32(len(values)-1)
		ys[i] = 0.2 + 0.75*scale.At(v)
	}
	spark := ui.Box(c).Radius(6)
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

func areaChart(c *ui.Context, p *palette, v areaChartView) ui.Element {
	scale := chart.NewLinear(0, v.ceiling)
	left := make([]chart.Tick, int(v.ceiling/v.step)+1)
	for i := range left {
		value := v.step * float64(i)
		left[i] = chart.Tick{At: scale.At(value), Value: value, Label: v.axisLabel(value)}
	}
	xs, ys := make([]float32, len(v.values)), make([]float32, len(v.values))
	bottom := make([]chart.Tick, len(v.values))
	var stride float32
	if len(v.values) > 1 {
		stride = areaSpan / float32(len(v.values)-1)
	}
	for i, value := range v.values {
		xs[i] = areaInset + stride*float32(i)
		ys[i] = scale.At(value)
		bottom[i] = chart.Tick{At: xs[i], Label: v.names[i]}
	}
	label := func(_ chart.Side, tick chart.Tick) {
		ui.Text(c, tick.Label).NoWrap().FontSize(14).TextColor(p.muted)
	}
	return chart.Plot(c, chart.Axes{Left: left, Bottom: bottom, Label: label}, func(area ui.Element) {
		shadedXs, shadedYs := shadeUpTo(xs, ys, area.AnimateTo("shade", float32(v.shaded), pillMotion))
		area.Label(v.label).Draw(func(pt *ui.Painter, r ui.Rect) {
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
		selected := pickPoint(c, area, xs, v.describe)
		if selected >= 0 {
			chart.RuleX(c.Key("guide"), xs[selected], 1).Background(p.axis).Transition(markMotion)
		}
		for i := range v.values {
			dot := chart.Point(c.Key(i), xs[i], ys[i], 10, 10).Radius(5).Border(2, p.accent)
			var pose ui.Pose
			if i == selected {
				pose.Scale = 1.3
			}
			dot.Background(p.card.Mix(p.accent, progress(dot, "filled", i == selected, quickMotion))).AnimatePose(pose, quickMotion)
		}
		if selected >= 0 {
			chartTip(c, p, area.Bounds().W, xs[selected], ys[selected], v.describe(selected))
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

func donut(c *ui.Context, p *palette, v donutView) ui.Element {
	slice := func(i int) ui.Color { return p.slices[i%len(p.slices)] }
	return ui.Row(c).Gap(23).Children(func() {
		ring := ui.Box(c).Size(2*donutRadius, 2*donutRadius).Shrink(0).Center().Label(v.label)
		shown := make([]float64, len(v.shares))
		for i, share := range v.shares {
			shown[i] = float64(ring.AnimateTo(i, float32(share), pillMotion))
		}
		slices := chart.Pie(shown, 0.03)
		active := -1
		if x, y, over := ring.PointerPosition(); over && ring.Hovered() {
			if i, ok := chart.SliceAt(slices, x-donutRadius, y-donutRadius, donutHole, donutRadius+donutGrowth); ok {
				active = i
			}
		}
		lit := make([]float32, len(slices))
		var dim float32
		ui.Column(c).Grow(1).Gap(15).Margin(0, 0, 10, 0).Children(func() {
			rows := make([]ui.Element, len(v.names))
			for i, name := range v.names {
				rows[i] = ui.Row(c).Height(20).Gap(10).ClipY().Focusable().Label(name)
				if rows[i].Hovered() || rows[i].FocusVisible() {
					active = i
				}
			}
			dim = progress(ring, "dim", active >= 0, fadeMotion)
			for i, name := range v.names {
				lit[i] = progress(ring, name, i == active, fadeMotion)
				rows[i].TextColor(p.ink.Mix(p.muted, dim-lit[i])).Children(func() {
					ui.Box(c).Size(12, 12).Radius(6).Background(slice(i).Mix(p.card, donutFade*(dim-lit[i])))
					ui.Text(c, name).FontSize(13.5).Grow(1)
					ui.Text(c, fmt.Sprintf("%.0f%%", v.shares[i])).FontSize(13.5).ContentTransition(rollDigits)
				})
			}
		})
		ring.Draw(func(pt *ui.Painter, r ui.Rect) {
			cx, cy := r.X+r.W/2, r.Y+r.H/2
			for i, part := range slices {
				pt.FillPath(chart.Sector(cx, cy, donutHole-donutSink*lit[i], donutRadius+donutGrowth*lit[i], part), slice(i).Mix(p.card, donutFade*(dim-lit[i])))
			}
		}).Children(func() {
			ui.Box(c).ClipY().Children(func() {
				ui.Text(c, v.center).FontSize(21).Bold().ContentTransition(rollDigits)
			})
			ui.Text(c, v.caption).FontSize(14).TextColor(p.muted).Margin(3, 0, 0, 0)
		})
		ui.Box(c).Absolute().Left(0).Top(0).Size(2*donutRadius, 2*donutRadius).PassThrough().Children(func() {
			if active >= 0 {
				x, y := chart.Polar(0.5, 0.5, 0.4, slices[active].Mid())
				chartTip(c, p, 2*donutRadius, x, 1-y, v.describe(active))
			}
		})
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
			Role(ui.RoleTooltip).ClipY().Transition(tipFit).Children(func() {
			ui.Text(c, text).FontSize(tipFontSize).TextColor(p.pageTop).NoWrap().ContentTransition(tipDigits)
		})
	}).IgnoreInput(true).Transition(markMotion)
}
