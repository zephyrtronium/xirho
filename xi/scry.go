package xi

import (
	"math"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xmath"
)

// Scry maps space onto a sphere.
type Scry struct {
	Radius float64 `xirho:"radius"`
}

func (v *Scry) Calc(in xirho.Pt, rng *xmath.RNG) xirho.Pt {
	r := in.X*in.X + in.Y*in.Y + in.Z*in.Z
	s := math.Sqrt(r) * (r + 1/v.Radius)
	in.X /= s
	in.Y /= s
	in.Z /= s
	return in
}

func (v *Scry) Prep() {}

func init() {
	must("scry", func() xirho.Func { return &Scry{Radius: 1} })
}
