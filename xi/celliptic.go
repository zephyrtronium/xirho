package xi

import (
	"math"
	"math/cmplx"

	"github.com/zephyrtronium/xirho"
)

// CElliptic is a simpler, conformal modification of elliptic.
type CElliptic struct{}

func newConfElliptic() xirho.F {
	return CElliptic{}
}

func (CElliptic) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	z := 1 - cmplx.Acos(complex(in.X, in.Y))*(2/math.Pi)
	in.X = real(z)
	in.Y = imag(z)
	return in
}

func (CElliptic) Prep() {}

func init() {
	Register("celliptic", newConfElliptic)
	// maybe elliptic as well?
}
