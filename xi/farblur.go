package xi

import (
	"math"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xmath"
)

// Farblur applies an affine-transformed Gaussian blur with strength varying
// according to distance from a point. Farblur is intended to be used following
// other functions in a Then.
type Farblur struct {
	Origin [3]float64   `xirho:"origin"`
	Ax     xmath.Affine `xirho:"affine"`
	Dist   float64      `xirho:"dist"`
}

func (v *Farblur) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt {
	ox, oy, oz := in.X-v.Origin[0], in.Y-v.Origin[1], in.Z-v.Origin[2]
	r := ox*ox + oy*oy + oz*oz
	s := math.Pow(r, v.Dist)
	x, y, z := xmath.Tx(&v.Ax, rng.Normal()*s, rng.Normal()*s, rng.Normal()*s)
	in.X += x
	in.Y += y
	in.Z += z
	return in
}

func (v *Farblur) Prep() {}

func newFarblur() xirho.Func {
	return &Farblur{
		Ax:   xmath.Eye(),
		Dist: 0.5,
	}
}

func init() {
	must("farblur", newFarblur)
}
