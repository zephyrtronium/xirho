package xi

import (
	"math"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xmath"
)

// Exblur applies a radial blur with strength depending on distance from a
// point. Exblur is intended to be used following other functions in a Then.
type Exblur struct {
	Str    xirho.Real `xirho:"strength"`
	Dist   xirho.Real `xirho:"dist"`
	Origin xirho.Vec3 `xirho:"origin"`
}

func (v *Exblur) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	ox, oy, oz := in.X-v.Origin[0], in.Y-v.Origin[1], in.Z-v.Origin[2]
	r := xmath.R3(ox, oy, oz)
	s := math.Pow(r*r, float64(v.Dist))*rng.Normal()*float64(v.Str) + r
	in.X = ox*s/r + v.Origin[0]
	in.Y = oy*s/r + v.Origin[1]
	in.Z = oz*s/r + v.Origin[2]
	return in
}

func (v *Exblur) Prep() {}

func init() {
	Register("exblur", func() xirho.F { return &Exblur{Dist: 0.5} })
}
