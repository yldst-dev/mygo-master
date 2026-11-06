//go:build ignore

package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
)

func load(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return png.Decode(f)
}

func distance(a, b uint32) uint32 {
	if a > b {
		return a - b
	}
	return b - a
}

func compare(name, left, right string) bool {
	a, err := load(left)
	if err != nil {
		fmt.Printf("%-24s %v\n", name, err)
		return false
	}
	b, err := load(right)
	if err != nil {
		fmt.Printf("%-24s %v\n", name, err)
		return false
	}
	r := a.Bounds()
	if r != b.Bounds() {
		fmt.Printf("%-24s size differs: %v, %v\n", name, r, b.Bounds())
		return false
	}
	count, worst, total := 0, uint32(0), uint64(0)
	minX, minY, maxX, maxY := r.Max.X, r.Max.Y, -1, -1
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			ar, ag, ab, _ := a.At(x, y).RGBA()
			br, bg, bb, _ := b.At(x, y).RGBA()
			dr, dg, db := distance(ar, br)>>8, distance(ag, bg)>>8, distance(ab, bb)>>8
			total += uint64(dr + dg + db)
			if d := max(dr, dg, db); d > 0 {
				count++
				worst = max(worst, d)
				minX, minY, maxX, maxY = min(minX, x), min(minY, y), max(maxX, x), max(maxY, y)
			}
		}
	}
	if count == 0 {
		fmt.Printf("%-24s identical\n", name)
		return true
	}
	mean := float64(total) / float64(r.Dx()*r.Dy()*3) / 255 * 100
	fmt.Printf("%-24s %d px differ, worst %d/255, mean %.2f%%, box x %d-%d y %d-%d\n", name, count, worst, mean, minX, maxX, minY, maxY)
	return false
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: go run pixdiff.go <before.png|dir> <after.png|dir>")
		os.Exit(2)
	}
	left, right := os.Args[1], os.Args[2]
	same := true
	if info, err := os.Stat(left); err == nil && info.IsDir() {
		names, _ := filepath.Glob(filepath.Join(left, "*.png"))
		if len(names) == 0 {
			fmt.Fprintln(os.Stderr, "no PNG files in", left)
			os.Exit(2)
		}
		for _, path := range names {
			name := filepath.Base(path)
			same = compare(name, path, filepath.Join(right, name)) && same
		}
	} else {
		same = compare(filepath.Base(left), left, right)
	}
	if !same {
		os.Exit(1)
	}
}
