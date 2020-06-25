package main

import (
	"context"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"golang.org/x/image/draw"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/variations"
)

func main() {
	F := []xirho.F{
		first(), // spam to simulate weights
		first(),
		first(),
		first(),
		first(),
		first(),
		second(),
		second(),
		second(),
		second(),
		second(),
		second(),
		second(),
		second(),
		third(),
		fourth(),
	}
	system := xirho.System{
		Funcs: F,
	}
	cam := xirho.Ax{}
	cam.Eye().Scale(0.25, 0.25, 0)
	r := xirho.R{
		Hist:    xirho.NewHist(4096, 4096, 4),
		System:  system,
		Camera:  cam,
		Palette: mkpalette(),
		N:       50e6,
	}
	r.Render(context.Background())
	r.Hist.Stat()
	img := image.NewRGBA(image.Rect(0, 0, 1024, 1024))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.NRGBA64{A: 0xffff}), image.Point{}, draw.Src)
	draw.CatmullRom.Scale(img, img.Bounds(), r.Hist, r.Hist.Bounds(), draw.Over, nil)
	err := png.Encode(os.Stdout, img)
	if err != nil {
		panic(err)
	}
}

func mkpalette() []color.NRGBA64 {
	return []color.NRGBA64{
		{R: 0xffff, G: 0x0000, B: 0x0000, A: 0xffff},
		{R: 0xffff, G: 0xffff, B: 0x0000, A: 0xffff},
		{R: 0x0000, G: 0xffff, B: 0x0000, A: 0xffff},
		{R: 0x0000, G: 0xffff, B: 0xffff, A: 0xffff},
		{R: 0x0000, G: 0x0000, B: 0xffff, A: 0xffff},
		{R: 0xffff, G: 0x0000, B: 0xffff, A: 0xffff},
		{R: 0xffff, G: 0x0000, B: 0x0000, A: 0xffff},
		{R: 0xffff, G: 0xffff, B: 0x0000, A: 0xffff},
		{R: 0x0000, G: 0xffff, B: 0x0000, A: 0xffff},
		{R: 0x0000, G: 0xffff, B: 0xffff, A: 0xffff},
		{R: 0x0000, G: 0x0000, B: 0xffff, A: 0xffff},
		{R: 0xffff, G: 0x0000, B: 0xffff, A: 0xffff},
		{R: 0xffff, G: 0x0000, B: 0x0000, A: 0xffff},
		{R: 0xffff, G: 0xffff, B: 0x0000, A: 0xffff},
		{R: 0x0000, G: 0xffff, B: 0x0000, A: 0xffff},
		{R: 0x0000, G: 0xffff, B: 0xffff, A: 0xffff},
		{R: 0x0000, G: 0x0000, B: 0xffff, A: 0xffff},
		{R: 0xffff, G: 0x0000, B: 0xffff, A: 0xffff},
	}
	// r := make([]color.NRGBA64, 256)
	// for i := range r {
	// 	r[i] = color.NRGBA64{R: uint16(i * i), G: uint16(i * i), B: uint16(i * i), A: 0xffff}
	// }
	// return r
}

func first() xirho.F {
	ax := xirho.Ax{}
	ax.Eye().RotZ(-math.Pi/2).Translate(1, 0, 0)
	return variations.NewThen(variations.NewAffine(ax, 0, 0.25), variations.NewSpherical())
}

func second() xirho.F {
	ax := xirho.Ax{}
	ax.Eye().RotZ(math.Pi / 2)
	return variations.NewThen(variations.NewAffine(ax, 1, 0.75), variations.NewSpherical())
}

func third() xirho.F {
	ax := xirho.Ax{}
	ax.Eye().Translate(3, 0, 0)
	return variations.NewAffine(ax, 0.5, 0.9)
}

func fourth() xirho.F {
	ax := xirho.Ax{}
	ax.Eye().Translate(-3, 0, 0)
	return variations.NewAffine(ax, 0.5, 0.9)
}
