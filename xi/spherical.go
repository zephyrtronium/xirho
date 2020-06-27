package xi

import (
	"github.com/zephyrtronium/xirho"
)

// Spherical calculates the conjugate of the complex reciprocal of the x and y
// of the input point treated as x+iy. The z and c coordinates are unchanged.
type Spherical struct{}

// NewSpherical is a factory for Spherical.
func NewSpherical() xirho.F {
	return Spherical{}
}

func (Spherical) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	r := in.X*in.X + in.Y*in.Y + in.Z*in.Z
	in.X /= r
	in.Y /= r
	in.Z /= r
	return in
}

func (Spherical) Prep() {}
