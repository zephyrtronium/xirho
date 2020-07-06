package xi

import (
	"math"

	"github.com/zephyrtronium/crazy"
	"github.com/zephyrtronium/xirho"
)

// Rod creates a solid cylinder of a given radius, with circular cross-sections
// in the x/z plane.
type Rod struct {
	Radius xirho.Real `xirho:"radius"`
}

func (v *Rod) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	u := crazy.Uniform0_1{Source: rng}
	s, c := math.Sincos(2 * math.Pi * u.Next())
	return xirho.P{
		X: float64(v.Radius) * s,
		Y: in.Y + u.Next() + u.Next() + u.Next() + u.Next() - 2,
		Z: float64(v.Radius) * c,
	}
}

func (v *Rod) Prep() {}

func init() {
	Register("rod", func() xirho.F { return &Rod{Radius: 0.1} })
}
