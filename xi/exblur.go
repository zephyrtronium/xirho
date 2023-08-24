package xi

import (
	"math"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xmath"
)

// Exblur applies a radial blur with strength depending on distance from a
// point. Exblur is intended to be used following other functions in a Then.
type Exblur struct {
	Str    float64    `xirho:"strength"`
	Dist   float64    `xirho:"dist"`
	Origin [3]float64 `xirho:"origin"`
}

func (v *Exblur) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt {
	ox, oy, oz := in.X-v.Origin[0], in.Y-v.Origin[1], in.Z-v.Origin[2]
	r := xmath.R3(ox, oy, oz)
	s := math.Pow(r*r, v.Dist)*rng.Normal()*v.Str + r
	in.X = ox*s/r + v.Origin[0]
	in.Y = oy*s/r + v.Origin[1]
	in.Z = oz*s/r + v.Origin[2]
	return in
}

func (v *Exblur) Prep() {}

func init() {
	must("exblur", func() xirho.Func { return &Exblur{Dist: 0.5} })
}
