package main

import (
	"context"
	"image"
	"image/color"
	"image/png"
	"os"

	"golang.org/x/image/draw"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/variations"
)

func main() {
	F := []xirho.F{
		variations.NewAffine(),
		variations.NewAffine(),
		variations.NewAffine(),
		// variations.NewAffine(),
	}
	a := [3][3]float64{
		{0.5, 0, 0},
		{0, 0.5, 0},
		{0, 0, 0},
	}
	*F[0].Params()[0].(xirho.Affine).V = xirho.Ax{A: a, B: [3]float64{-1, -1, 0}}
	*F[1].Params()[0].(xirho.Affine).V = xirho.Ax{A: a, B: [3]float64{-1, 1, 0}}
	*F[2].Params()[0].(xirho.Affine).V = xirho.Ax{A: a, B: [3]float64{1, -1, 0}}
	// *F[3].Params()[0].(xirho.Affine).V = xirho.Ax{A: a, B: [3]float64{1, 1, 0}}
	*F[0].Params()[1].(xirho.Real).V = 0
	*F[1].Params()[1].(xirho.Real).V = 0.25
	*F[2].Params()[1].(xirho.Real).V = 0.75
	// *F[3].Params()[1].(xirho.Real).V = 1
	*F[0].Params()[2].(xirho.Real).V = 0.5
	*F[1].Params()[2].(xirho.Real).V = 0.5
	*F[2].Params()[2].(xirho.Real).V = 0.5
	// *F[3].Params()[2].(xirho.Real).V = 0.5
	system := xirho.System{
		Funcs: F,
	}
	palette := mkpalette()
	r := xirho.R{
		Hist:    xirho.NewHist(256, 256, 1),
		System:  system,
		Camera:  *xirho.Eye(),
		Palette: palette,
		Q:       1e7,
	}
	r.Render(context.Background())
	r.Hist.Stat()
	img := image.NewNRGBA(r.Hist.Bounds())
	draw.Draw(img, img.Bounds(), image.NewUniform(color.NRGBA64{A: 0xffff}), image.ZP, draw.Src)
	draw.Draw(img, img.Bounds(), r.Hist, image.ZP, draw.Over)
	err := png.Encode(os.Stdout, img)
	if err != nil {
		panic(err)
	}
}

func mkpalette() []color.NRGBA64 {
	// return []color.NRGBA64{
	// 	{R: 0xffff, G: 0x0000, B: 0x0000, A: 0xffff},
	// 	{R: 0xffff, G: 0xffff, B: 0x0000, A: 0xffff},
	// 	{R: 0x0000, G: 0xffff, B: 0x0000, A: 0xffff},
	// 	{R: 0x0000, G: 0xffff, B: 0xffff, A: 0xffff},
	// 	{R: 0x0000, G: 0x0000, B: 0xffff, A: 0xffff},
	// 	{R: 0xffff, G: 0x0000, B: 0xffff, A: 0xffff},
	// }
	r := make([]color.NRGBA64, 256)
	for i := range r {
		r[i] = color.NRGBA64{R: uint16(i * i), G: uint16(i) << 8, B: uint16(i * i), A: 0xffff}
	}
	return r
}
