package xi

import (
	"math"

	"github.com/zephyrtronium/xirho"
)

// Scry maps space onto a sphere.
type Scry struct {
	Radius xirho.Real `xirho:"radius"`
}

func (v *Scry) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	r := in.X*in.X + in.Y*in.Y + in.Z*in.Z
	s := math.Sqrt(r) * (r + 1/float64(v.Radius))
	in.X /= s
	in.Y /= s
	in.Z /= s
	return in
}

func (v *Scry) Prep() {}

func init() {
	Register("scry", func() xirho.F { return &Scry{Radius: 1} })
}
