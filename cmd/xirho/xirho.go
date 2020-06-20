package main

import (
	"context"
	"image/color"
	"image/png"
	"os"

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
		Hist:    xirho.NewHist(64, 64, 1),
		System:  system,
		Camera:  *xirho.Eye(),
		Palette: palette,
		N:       1e6,
	}
	r.Render(context.Background())
	r.Hist.Stat()
	err := png.Encode(os.Stdout, r.Hist)
	if err != nil {
		panic(err)
	}
}

func mkpalette() []color.NRGBA64 {
	r := make([]color.NRGBA64, 256)
	for i := range r {
		r[i] = color.NRGBA64{R: uint16(i), G: uint16(i), B: uint16(i), A: 0xffff}
	}
	return r
}
