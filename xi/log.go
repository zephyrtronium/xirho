package xi

import (
	"math"
	"math/cmplx"

	"github.com/zephyrtronium/xirho"
)

// Log computes the complex logarithm of the input point treated as x+iy.
type Log struct {
	Base xirho.Complex `xirho:"base"`

	lb complex128
}

func (v *Log) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	z := cmplx.Log(complex(in.X, in.Y)) * v.lb
	in.X = real(z)
	in.Y = imag(z)
	return in
}

func (v *Log) Prep() {
	v.lb = 1 / cmplx.Log(complex128(v.Base))
}

func init() {
	Register("log", func() xirho.F { return &Log{Base: math.E} })
}
