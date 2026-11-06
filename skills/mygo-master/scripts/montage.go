//go:build ignore

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 8 {
		fmt.Fprintln(os.Stderr, "usage: go run montage.go <out.png> <x> <y> <w> <h> <columns> <frame.png>...")
		os.Exit(2)
	}
	number := func(i int) int {
		v, err := strconv.Atoi(os.Args[i])
		if err != nil {
			panic(err)
		}
		return v
	}
	x, y, w, h, columns := number(2), number(3), number(4), number(5), number(6)
	files := os.Args[7:]
	rows := (len(files) + columns - 1) / columns
	canvas := image.NewRGBA(image.Rect(0, 0, columns*(w+6)-6, rows*(h+6)-6))
	draw.Draw(canvas, canvas.Bounds(), image.NewUniform(color.RGBA{255, 0, 255, 255}), image.Point{}, draw.Src)
	for i, name := range files {
		f, err := os.Open(name)
		if err != nil {
			panic(err)
		}
		frame, err := png.Decode(f)
		f.Close()
		if err != nil {
			panic(err)
		}
		left, top := (i%columns)*(w+6), (i/columns)*(h+6)
		draw.Draw(canvas, image.Rect(left, top, left+w, top+h), frame, image.Pt(x, y), draw.Src)
	}
	out, err := os.Create(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer out.Close()
	if err := png.Encode(out, canvas); err != nil {
		panic(err)
	}
}
