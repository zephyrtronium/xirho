package xi

import (
	"math"

	"github.com/zephyrtronium/crazy"
	"github.com/zephyrtronium/xirho"
)

// Farblur applies an affine-transformed Gaussian blur with strength varying
// according to distance from a point. Farblur is intended to be used following
// other functions in a Then.
type Farblur struct {
	Origin xirho.Vec3   `xirho:"origin"`
	Ax     xirho.Affine `xirho:"affine"`
	Dist   xirho.Real   `xirho:"dist"`
}

func (v *Farblur) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	ox, oy, oz := in.X-v.Origin[0], in.Y-v.Origin[1], in.Z-v.Origin[2]
	r := ox*ox + oy*oy + oz*oz
	g := crazy.NewNormal(rng, 0, math.Pow(r, float64(v.Dist)))
	x, y, z := xirho.Tx(&v.Ax, g.Next(), g.Next(), g.Next())
	in.X += x
	in.Y += y
	in.Z += z
	return in
}

func (v *Farblur) Prep() {}

func newFarblur() xirho.F {
	return &Farblur{
		Ax:   xirho.Eye(),
		Dist: 0.5,
	}
}

func init() {
	Register("farblur", newFarblur)
}
