package xi

import (
	"math"
	"math/cmplx"

	"github.com/zephyrtronium/xirho"
)

// Exp performs the complex exponential in a fixed base.
type Exp struct {
	Base complex128 `xirho:"base"`

	lb complex128
}

func (v *Exp) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt {
	z := cmplx.Exp(complex(in.X, in.Y) * v.lb)
	in.X = real(z)
	in.Y = imag(z)
	return in
}

func (v *Exp) Prep() {
	v.lb = cmplx.Log(v.Base)
}

func init() {
	must("exp", func() xirho.Func { return &Exp{Base: math.E} })
}
