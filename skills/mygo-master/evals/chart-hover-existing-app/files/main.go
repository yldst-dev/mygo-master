package main

import (
	"fmt"
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/chart"
)

var (
	regions = []string{"North", "South"}
	months  = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
	sales   = [2][6]float64{{42, 58, 76, 64, 91, 83}, {31, 47, 52, 70, 66, 88}}
)

type app struct {
	region int
}

func newApp() *app {
	return &app{}
}

func (a *app) view(c *ui.Context) {
	ui.Column(c).Fill().Padding(24).Gap(16).Children(func() {
		ui.Text(c, "Monthly sales").FontSize(22).Bold()
		ui.Row(c).Key("regions").Gap(8).Children(func() {
			for i, name := range regions {
				if ui.Button(c, name).Clicked() {
					a.region = i
				}
			}
		})
		a.chart(c)
	})
}

func (a *app) chart(c *ui.Context) {
	t := c.Theme()
	values := sales[a.region]
	x := chart.NewBand(months...).Padding(0.3, 0.15)
	y := chart.NewLinear(0, 100).Nice(4)
	chart.Plot(c, chart.Axes{Left: y.Ticks(4), Bottom: x.Ticks()}, func(area ui.Element) {
		for i, v := range values {
			chart.Rect(c, x.Band(i), y.Span(0, v)).Radius(4).Background(t.Accent).
				Tooltip(fmt.Sprintf("%s: %.0f", months[i], v))
		}
	}).Label("Sales by month").Height(280).Gap(6)
}

func main() {
	a := newApp()
	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{
			Title:     "Monthly sales",
			Width:     640,
			Height:    480,
			MinWidth:  480,
			MinHeight: 420,
			StateKey:  "main",
			Content:   ui.View(a.view),
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
