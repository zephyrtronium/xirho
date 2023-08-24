package xi

import (
	"math"
	"math/cmplx"

	"github.com/zephyrtronium/xirho"
)

// Log computes the complex logarithm of the input point treated as x+iy.
type Log struct {
	Base complex128 `xirho:"base"`

	lb complex128
}

func (v *Log) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt {
	z := cmplx.Log(complex(in.X, in.Y)) * v.lb
	in.X = real(z)
	in.Y = imag(z)
	return in
}

func (v *Log) Prep() {
	v.lb = 1 / cmplx.Log(v.Base)
}

func init() {
	must("log", func() xirho.Func { return &Log{Base: math.E} })
}
