package main

import (
	"context"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"

	"golang.org/x/image/draw"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xi"
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
	cam.Eye().Scale(0.5, 0.5, 0)
	const width, height = 16384, 16384
	log.Println("allocating histogram, estimated", xirho.HistMem(width, height)>>20, "MB")
	const iters = width * height * 5
	r := xirho.R{
		Hist:    xirho.NewHist(width, height, 4),
		System:  system,
		Camera:  cam,
		Palette: mkpalette(),
		N:       iters,
		Q:       iters / 5,
	}
	log.Println("rendering up to", r.N, "iters or", r.Q, "hits")
	r.Render(context.Background())
	log.Println("finished render with", r.Iters(), "iters,", r.Hits(), "hits")
	r.Hist.Stat()
	const outWidth, outHeight = 1024, 1024
	img := image.NewRGBA(image.Rect(0, 0, outWidth, outHeight))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.NRGBA64{A: 0xffff}), image.Point{}, draw.Src)
	log.Printf("drawing onto image of size %dx%d", outWidth, outHeight)
	draw.CatmullRom.Scale(img, img.Bounds(), r.Hist, r.Hist.Bounds(), draw.Over, nil)
	log.Println("encoding to stdout")
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
	return xi.NewThen(xi.NewAffine(ax, 0, 0.25), xi.NewSpherical())
}

func second() xirho.F {
	ax := xirho.Ax{}
	ax.Eye().RotZ(math.Pi / 2)
	return xi.NewThen(xi.NewAffine(ax, 1, 0.75), xi.NewSpherical())
}

func third() xirho.F {
	ax := xirho.Ax{}
	ax.Eye().Translate(3, 0, 0)
	return xi.NewAffine(ax, 0.5, 0.9)
}

func fourth() xirho.F {
	ax := xirho.Ax{}
	ax.Eye().Translate(-3, 0, 0)
	return xi.NewAffine(ax, 0.5, 0.9)
}
