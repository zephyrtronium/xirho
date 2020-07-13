package xi

import (
	"math"

	"github.com/zephyrtronium/xirho"
)

// Rod creates a solid cylinder of a given radius, with circular cross-sections
// in the x/z plane.
type Rod struct {
	Radius xirho.Real `xirho:"radius"`
}

func (v *Rod) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	s, c := math.Sincos(2 * math.Pi * rng.Uniform())
	return xirho.P{
		X: float64(v.Radius) * s,
		Y: in.Y + rng.Normal(),
		Z: float64(v.Radius) * c,
	}
}

func (v *Rod) Prep() {}

func init() {
	must("rod", func() xirho.F { return &Rod{Radius: 0.1} })
}
