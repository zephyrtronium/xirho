package xi

import (
	"math"
	"math/cmplx"

	"github.com/zephyrtronium/xirho"
)

// Exp performs the complex exponential in a fixed base.
type Exp struct {
	Base xirho.Complex `xirho:"base"`

	lb complex128
}

func (v *Exp) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	z := cmplx.Exp(complex(in.X, in.Y) * v.lb)
	in.X = real(z)
	in.Y = imag(z)
	return in
}

func (v *Exp) Prep() {
	v.lb = cmplx.Log(complex128(v.Base))
}

func init() {
	must("exp", func() xirho.F { return &Exp{Base: math.E} })
}
