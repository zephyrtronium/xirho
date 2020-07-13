package xi

import (
	"math"

	"github.com/zephyrtronium/crazy"
	"github.com/zephyrtronium/xirho"
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
	r := ox*ox + oy*oy + oz*oz
	// TODO: learn some geometry to do this without trig functions
	theta := math.Atan2(oy, ox)
	phi := math.Acos(oz / math.Sqrt(r))
	s := math.Pow(r, float64(v.Dist))*crazy.Normal{Source: rng, StdDev: float64(v.Str)}.Next() + math.Sqrt(r)
	st, ct := math.Sincos(theta)
	sp, cp := math.Sincos(phi)
	return xirho.P{
		X: s * sp * ct,
		Y: s * sp * st,
		Z: s * cp,
		C: in.C,
	}
}

func (v *Exblur) Prep() {}

func init() {
	Register("exblur", func() xirho.F { return &Exblur{Dist: 0.5} })
}
