package view

import "github.com/egoist/mygo/ui"

func softButton(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	button := pressable(ui.ButtonBase(c), 0.96).Height(30).Padding(0, 14).Radius(8).Label(label)
	hover := progress(button, "hover", button.Hovered(), fadeMotion)
	return button.Background(t.SurfaceHover.Alpha(hover)).Children(func() {
		ui.Text(c, label).FontSize(14).TextColor(t.TextMuted.Mix(t.Text, hover))
	})
}

func primaryButton(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	button := pressable(ui.ButtonBase(c), 0.96).Height(40).Padding(0, 18).Radius(10).Label(label)
	hover := progress(button, "hover", button.Hovered(), fadeMotion)
	return button.Background(t.Accent.Mix(t.Text, 0.12*hover)).Children(func() {
		ui.Text(c, label).FontSize(14.5).FontWeight(600).TextColor(t.AccentText)
	})
}

func toggle(c *ui.Context, on *bool, label string, changed func()) ui.Element {
	t := c.Theme()
	s := ui.SwitchBase(c, on).Label(label).Size(52, 30).Radius(15).OnChange(changed)
	side := progress(s, "knob", *on, pillMotion)
	width := 26 + 5*progress(s, "hold", s.Pressed(), pressMotion)
	hover := progress(s, "hover", s.Hovered(), fadeMotion)
	return s.Background(t.SurfacePressed.Mix(t.Accent, side).Mix(t.Text, 0.06*hover)).Children(func() {
		ui.Box(c).Absolute().Top(2).Left(2+side*(48-width)).Size(width, 26).Radius(13).
			Background(ui.Hex("#FFFFFF")).Shadow(0, 1, 3, 0, ui.RGBA(0, 0, 0, 0.2))
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

func lightTheme() *ui.Theme {
	t := *ui.LightTheme()
	t.Background, t.Surface, t.SurfaceHover, t.SurfacePressed = ui.Hex("#F6F7FC"), ui.Hex("#FEFEFF"), ui.Hex("#EEF1F9"), ui.Hex("#DAE2F5")
	t.Border, t.Text, t.TextMuted = ui.Hex("#E6E9F0"), ui.Hex("#0B1220"), ui.Hex("#545B75")
	accent := ui.Hex("#0A5FFE")
	t.Accent, t.AccentHover, t.AccentPressed, t.AccentText, t.Focus, t.Radius = accent, accent, accent, ui.Hex("#FFFFFF"), accent, 8
	return &t
}

func darkTheme() *ui.Theme {
	t := *ui.DarkTheme()
	t.Background, t.Surface, t.SurfaceHover, t.SurfacePressed = ui.Hex("#0B0F1A"), ui.Hex("#141B2C"), ui.Hex("#18213A"), ui.Hex("#232D45")
	t.Border, t.Text, t.TextMuted = ui.Hex("#232C40"), ui.Hex("#F3F6FC"), ui.Hex("#8E98AE")
	accent := ui.Hex("#2F6BF2")
	t.Accent, t.AccentHover, t.AccentPressed, t.AccentText, t.Focus, t.Radius = accent, accent, accent, ui.Hex("#FFFFFF"), accent, 8
	return &t
}

func themeColors(t *ui.Theme) []*ui.Color {
	return []*ui.Color{
		&t.Background, &t.Surface, &t.SurfaceHover, &t.SurfacePressed, &t.Border, &t.Text, &t.TextMuted,
		&t.Accent, &t.AccentHover, &t.AccentPressed, &t.AccentText,
		&t.Danger, &t.Warning, &t.Success, &t.Selection, &t.Focus, &t.Scrollbar,
	}
}

func fadedTheme(darkness float32) *ui.Theme {
	light, dark := lightTheme(), darkTheme()
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
