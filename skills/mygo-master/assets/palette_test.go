package main

import (
	"math"
	"reflect"
	"testing"

	"github.com/egoist/mygo/ui"
)

func luminance(c ui.Color) float64 {
	channel := func(v uint8) float64 {
		s := float64(v) / 255
		if s <= 0.03928 {
			return s / 12.92
		}
		return math.Pow((s+0.055)/1.055, 2.4)
	}
	return 0.2126*channel(c.R) + 0.7152*channel(c.G) + 0.0722*channel(c.B)
}

func contrast(a, b ui.Color) float64 {
	hi, lo := luminance(a), luminance(b)
	if lo > hi {
		hi, lo = lo, hi
	}
	return (hi + 0.05) / (lo + 0.05)
}

func TestPaletteFadeCoversEveryColor(t *testing.T) {
	var count func(reflect.Type) int
	count = func(typ reflect.Type) int {
		switch {
		case typ == reflect.TypeOf(ui.Color{}):
			return 1
		case typ.Kind() == reflect.Array:
			return typ.Len() * count(typ.Elem())
		case typ.Kind() == reflect.Struct:
			total := 0
			for i := range typ.NumField() {
				total += count(typ.Field(i).Type)
			}
			return total
		}
		return 0
	}
	if want, got := count(reflect.TypeOf(palette{})), len(lightPalette.channels()); want != got {
		t.Errorf("the palette has %d colors, the fade moves %d", want, got)
	}
	if fadedPalette(0) != &lightPalette || fadedPalette(1) != &darkPalette {
		t.Error("the ends of the fade are not the palettes themselves")
	}
}

func TestPaletteTextIsReadable(t *testing.T) {
	for name, p := range map[string]*palette{"light": &lightPalette, "dark": &darkPalette} {
		for _, pair := range []struct {
			what       string
			text, base ui.Color
		}{
			{"ink on the card", p.ink, p.card},
			{"body on the card", p.body, p.card},
			{"muted on the card", p.muted, p.card},
			{"ink on the page", p.ink, p.pageTop},
			{"muted on the page", p.muted, p.pageTop},
			{"text on the accent fill", p.onAccent, p.accentFill},
		} {
			if got := contrast(pair.text, pair.base); got < 4.5 {
				t.Errorf("%s, %s: contrast %.2f is under 4.5", name, pair.what, got)
			}
		}
	}
}

func TestThemeCarriesThePaletteIntoTheWidgets(t *testing.T) {
	for _, dark := range []bool{false, true} {
		p := &lightPalette
		if dark {
			p = &darkPalette
		}
		theme := p.theme(dark)
		if theme.Accent != p.accentFill || theme.InactiveAccent != theme.Accent || theme.Focus != p.accent {
			t.Errorf("dark %v: accent %v, inactive %v, focus %v", dark, theme.Accent, theme.InactiveAccent, theme.Focus)
		}
		if theme.Surface != p.card || theme.SurfacePressed != p.toggle || theme.Text != p.ink {
			t.Errorf("dark %v: the theme does not carry the palette", dark)
		}
	}
}

func TestThemeFadeCoversEveryColorAndEndsExactly(t *testing.T) {
	want := 0
	kind := reflect.TypeOf(ui.Theme{})
	for i := range kind.NumField() {
		if kind.Field(i).Type == reflect.TypeOf(ui.Color{}) && kind.Field(i).IsExported() {
			want++
		}
	}
	if got := len(themeColors(&ui.Theme{})); got != want {
		t.Errorf("the theme has %d colors, the fade moves %d", want, got)
	}
	if fadedTheme(0) != restThemes[0] || fadedTheme(1) != restThemes[1] {
		t.Error("the ends of the fade are not the themes themselves")
	}
	middle := fadedTheme(0.5)
	for i, color := range themeColors(middle) {
		from, to := *themeColors(restThemes[0])[i], *themeColors(restThemes[1])[i]
		if from != to && *color == from || from != to && *color == to {
			t.Errorf("color %d has not moved to the middle", i)
		}
	}
}
