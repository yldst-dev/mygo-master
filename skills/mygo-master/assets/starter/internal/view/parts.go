package view

import "github.com/egoist/mygo/ui"

func softButton(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	button := pressable(ui.ButtonBase(c), 0.96).Height(30).Padding(0, 14).Radius(8).Label(label)
	hover := progress(button, "hover", button.Hovered(), fadeMotion)
	return button.Background(t.SurfaceHover.Alpha(hover)).Children(func() {
		ui.Text(c, label).FontSize(14)
	})
}

func steadyScroll(c *ui.Context, build func()) {
	ui.Scroll(c).Grow(1).Children(func() {
		outer := ui.Column(c)
		outer.Children(func() {
			inner := ui.Column(c).Children(build)
			if height := inner.Bounds().H; height > 0 {
				outer.MinHeight(outer.AnimateTo("height", height, pillMotion))
			}
		})
	})
}

func themeColors(t *ui.Theme) []*ui.Color {
	return []*ui.Color{
		&t.Background, &t.Surface, &t.SurfaceHover, &t.SurfacePressed, &t.Border, &t.Text, &t.TextMuted,
		&t.Accent, &t.AccentHover, &t.AccentPressed, &t.AccentText,
		&t.Danger, &t.Warning, &t.Success, &t.Selection, &t.Focus, &t.Scrollbar,
	}
}

func fadedTheme(darkness float32) *ui.Theme {
	light, dark := ui.LightTheme(), ui.DarkTheme()
	switch darkness {
	case 0:
		return light
	case 1:
		return dark
	}
	faded := *light
	from, to := themeColors(&faded), themeColors(dark)
	for i, color := range from {
		*color = color.Mix(*to[i], darkness)
	}
	return &faded
}
