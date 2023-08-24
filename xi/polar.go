package xi

import (
	"math"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xmath"
)

// Polar maps the x/y rectangular coordinates of the input to polar.
type Polar struct{}

// newPolar is a factory for Polar.
func newPolar() xirho.Func {
	return Polar{}
}

func (Polar) Calc(in xirho.Pt, rng *xmath.RNG) xirho.Pt {
	t := math.Atan2(in.X, in.Y) / math.Pi
	r := math.Hypot(in.X, in.Y) - 1
	in.X = t
	in.Y = r
	return in
}

func (Polar) Prep() {}

func init() {
	must("polar", newPolar)
}
