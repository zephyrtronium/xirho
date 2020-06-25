package xi

import (
	"github.com/zephyrtronium/xirho"
)

// Spherical calculates the conjugate of the complex reciprocal of the x and y
// of the input point treated as x+iy. The z and c coordinates are unchanged.
type Spherical struct{}

func NewSpherical() xirho.F {
	return Spherical{}
}

func (Spherical) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	r := in.X*in.X + in.Y*in.Y
	in.X /= r
	in.Y /= r
	return in
}

func (Spherical) Params() []xirho.Param {
	return nil
}
