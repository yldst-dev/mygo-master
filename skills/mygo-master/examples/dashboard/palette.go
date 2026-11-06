package main

import "github.com/egoist/mygo/ui"

type tone struct {
	background, ink, dot ui.Color
}

type palette struct {
	railTop, railBottom, pageTop, pageBottom, bar, card ui.Color
	line, grid, axis, field, fieldLine, head            ui.Color
	ink, body, muted                                    ui.Color
	accent, accentInk, onAccent                         ui.Color
	navActive, hover, avatar, toggle, knob              ui.Color
	shadow, areaStrong, areaFaint                       ui.Color
	good, info, bad                                     tone
	slices                                              [4]ui.Color
}

var lightPalette = palette{
	railTop: ui.Hex("#FBFBFE"), railBottom: ui.Hex("#F7F9FE"),
	pageTop: ui.Hex("#F6F7FC"), pageBottom: ui.Hex("#EEF1F8"),
	bar: ui.Hex("#F7F8FD"), card: ui.Hex("#FEFEFF"),
	line: ui.Hex("#E6E9F0"), grid: ui.Hex("#E7E9F0"), axis: ui.Hex("#BCC5D5"),
	field: ui.Hex("#F4F5FA"), fieldLine: ui.Hex("#DCE1EB"), head: ui.Hex("#F2F4FA"),
	ink: ui.Hex("#0B1220"), body: ui.Hex("#434A62"), muted: ui.Hex("#545B75"),
	accent: ui.Hex("#0A5FFE"), accentInk: ui.Hex("#0A52EE"), onAccent: ui.Hex("#FFFFFF"),
	navActive: ui.Hex("#E3EBFC"), hover: ui.Hex("#EEF1F9"), avatar: ui.Hex("#C6D9FD"),
	toggle: ui.Hex("#DAE2F5"), knob: ui.Hex("#FFFFFF"),
	shadow:     ui.RGBA(15, 23, 42, 0.05),
	areaStrong: ui.Hex("#0A5FFE").Alpha(0.2), areaFaint: ui.Hex("#0A5FFE").Alpha(0.05),
	good:   tone{ui.Hex("#D9F5DB"), ui.Hex("#146C1E"), ui.Hex("#0E7A12")},
	info:   tone{ui.Hex("#DBEAFE"), ui.Hex("#12308F"), ui.Hex("#0A5FFE")},
	bad:    tone{ui.Hex("#FEE5E5"), ui.Hex("#A31818"), ui.Hex("#E5191D")},
	slices: [4]ui.Color{ui.Hex("#0A5FFE"), ui.Hex("#70A9FE"), ui.Hex("#C4DAFD"), ui.Hex("#E7E8EA")},
}

var darkPalette = palette{
	railTop: ui.Hex("#0F1524"), railBottom: ui.Hex("#0F1524"),
	pageTop: ui.Hex("#0B0F1A"), pageBottom: ui.Hex("#0B0F1A"),
	bar: ui.Hex("#0F1524"), card: ui.Hex("#141B2C"),
	line: ui.Hex("#232C40"), grid: ui.Hex("#232C40"), axis: ui.Hex("#3B465E"),
	field: ui.Hex("#182034"), fieldLine: ui.Hex("#2A3450"), head: ui.Hex("#1A2236"),
	ink: ui.Hex("#F3F6FC"), body: ui.Hex("#C4CCDD"), muted: ui.Hex("#8E98AE"),
	accent: ui.Hex("#3D7BFF"), accentInk: ui.Hex("#7CAEFF"), onAccent: ui.Hex("#FFFFFF"),
	navActive: ui.Hex("#1B2B4D"), hover: ui.Hex("#18213A"), avatar: ui.Hex("#1F3566"),
	toggle: ui.Hex("#232D45"), knob: ui.Hex("#0F1524"),
	shadow:     ui.RGBA(0, 0, 0, 0.25),
	areaStrong: ui.Hex("#3D7BFF").Alpha(0.28), areaFaint: ui.Hex("#3D7BFF").Alpha(0.04),
	good:   tone{ui.Hex("#12331F"), ui.Hex("#7CE2A0"), ui.Hex("#45C46B")},
	info:   tone{ui.Hex("#15294F"), ui.Hex("#9CC2FF"), ui.Hex("#3D7BFF")},
	bad:    tone{ui.Hex("#3A1618"), ui.Hex("#F7A9A9"), ui.Hex("#F0545A")},
	slices: [4]ui.Color{ui.Hex("#3D7BFF"), ui.Hex("#6BA5FA"), ui.Hex("#2C4B80"), ui.Hex("#3A4254")},
}

func (p *palette) channels() []*ui.Color {
	return []*ui.Color{
		&p.railTop, &p.railBottom, &p.pageTop, &p.pageBottom, &p.bar, &p.card,
		&p.line, &p.grid, &p.axis, &p.field, &p.fieldLine, &p.head,
		&p.ink, &p.body, &p.muted,
		&p.accent, &p.accentInk, &p.onAccent,
		&p.navActive, &p.hover, &p.avatar, &p.toggle, &p.knob,
		&p.shadow, &p.areaStrong, &p.areaFaint,
		&p.good.background, &p.good.ink, &p.good.dot,
		&p.info.background, &p.info.ink, &p.info.dot,
		&p.bad.background, &p.bad.ink, &p.bad.dot,
		&p.slices[0], &p.slices[1], &p.slices[2], &p.slices[3],
	}
}

func fadedPalette(darkness float32) *palette {
	switch darkness {
	case 0:
		return &lightPalette
	case 1:
		return &darkPalette
	}
	faded := lightPalette
	from, to := faded.channels(), darkPalette.channels()
	for i, color := range from {
		*color = color.Mix(*to[i], darkness)
	}
	return &faded
}

func (p *palette) status(status string) tone {
	switch status {
	case statusCompleted:
		return p.good
	case statusCancelled:
		return p.bad
	}
	return p.info
}

func (p *palette) theme(dark bool) *ui.Theme {
	t := *ui.LightTheme()
	if dark {
		t = *ui.DarkTheme()
	}
	t.Background, t.Surface, t.SurfaceHover = p.pageTop, p.card, p.hover
	t.Border, t.Text, t.TextMuted = p.line, p.ink, p.muted
	t.Accent, t.AccentHover, t.AccentPressed, t.AccentText = p.accent, p.accent, p.accent, p.onAccent
	t.Focus, t.Radius = p.accent, 8
	return &t
}
