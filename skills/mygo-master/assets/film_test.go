package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

type film struct {
	t      *testing.T
	tt     *ui.Tester
	rect   image.Rectangle
	frames []*image.RGBA
}

func newFilm(t *testing.T, tt *ui.Tester, scale float32, r ui.Rect, pad float32) *film {
	t.Helper()
	tt.SetScale(scale)
	return &film{t: t, tt: tt, rect: image.Rect(
		int((r.X-pad)*scale), int((r.Y-pad)*scale),
		int((r.X+r.W+pad)*scale), int((r.Y+r.H+pad)*scale),
	)}
}

func (f *film) grab() {
	source := f.tt.Image()
	frame := image.NewRGBA(image.Rect(0, 0, f.rect.Dx(), f.rect.Dy()))
	draw.Draw(frame, frame.Bounds(), source, f.rect.Min, draw.Src)
	f.frames = append(f.frames, frame)
}

func (f *film) play(step time.Duration, count int) {
	for range count {
		f.tt.Advance(step)
		f.grab()
	}
}

func (f *film) save(path string) {
	f.t.Helper()
	w, h := f.rect.Dx(), f.rect.Dy()
	columns := max(1, min(len(f.frames), 1600/(w+6)))
	rows := (len(f.frames) + columns - 1) / columns
	sheet := image.NewRGBA(image.Rect(0, 0, columns*(w+6)-6, rows*(h+6)-6))
	draw.Draw(sheet, sheet.Bounds(), image.NewUniform(color.RGBA{255, 0, 255, 255}), image.Point{}, draw.Src)
	for i, frame := range f.frames {
		x, y := (i%columns)*(w+6), (i/columns)*(h+6)
		draw.Draw(sheet, image.Rect(x, y, x+w, y+h), frame, image.Point{}, draw.Src)
	}
	writePNG(f.t, path, sheet)
}

func saveRest(t *testing.T, directory, name string, view func(*ui.Context), width, height int) {
	t.Helper()
	writePNG(t, filepath.Join(directory, name+".png"), ui.Render(view, width, height, 2))
}

func writePNG(t *testing.T, path string, picture image.Image) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if err := png.Encode(file, picture); err != nil {
		t.Fatal(err)
	}
}
