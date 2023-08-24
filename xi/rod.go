package xi

import (
	"math"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xmath"
)

// Rod creates a solid cylinder of a given radius, with circular cross-sections
// in the x/z plane.
type Rod struct {
	Radius float64 `xirho:"radius"`
}

func (v *Rod) Calc(in xirho.Pt, rng *xmath.RNG) xirho.Pt {
	s, c := math.Sincos(2 * math.Pi * rng.Uniform())
	return xirho.Pt{
		X: v.Radius * s,
		Y: in.Y + rng.Normal(),
		Z: v.Radius * c,
	}
}

func (v *Rod) Prep() {}

func init() {
	must("rod", func() xirho.Func { return &Rod{Radius: 0.1} })
}
