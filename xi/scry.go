package xi

import (
	"math"

	"github.com/zephyrtronium/xirho"
)

// Scry maps space onto a sphere.
type Scry struct {
	Radius xirho.Real `xirho:"radius"`
}

func (v *Scry) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt {
	r := in.X*in.X + in.Y*in.Y + in.Z*in.Z
	s := math.Sqrt(r) * (r + 1/float64(v.Radius))
	in.X /= s
	in.Y /= s
	in.Z /= s
	return in
}

func (v *Scry) Prep() {}

func init() {
	must("scry", func() xirho.Func { return &Scry{Radius: 1} })
}
