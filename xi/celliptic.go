package xi

import (
	"math"
	"math/cmplx"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xmath"
)

// CElliptic is a simpler, conformal modification of elliptic.
type CElliptic struct{}

func newConfElliptic() xirho.Func {
	return CElliptic{}
}

func (CElliptic) Calc(in xirho.Pt, rng *xmath.RNG) xirho.Pt {
	z := 1 - cmplx.Acos(complex(in.X, in.Y))*(2/math.Pi)
	in.X = real(z)
	in.Y = imag(z)
	return in
}

func (CElliptic) Prep() {}

func init() {
	must("celliptic", newConfElliptic)
	// maybe elliptic as well?
}
