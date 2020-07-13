package xi

import (
	"math"

	"github.com/zephyrtronium/xirho"
)

// Cylinder creates a cylinder with circular cross-sections of radius 1 in the
// x/z plane.
type Cylinder struct{}

func (Cylinder) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	s, c := math.Sincos(in.X)
	in.X = s
	in.Z = c
	return in
}

func (Cylinder) Prep() {}

func init() {
	must("cylinder", func() xirho.F { return Cylinder{} })
}
