package xi

import (
	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xmath"
)

// Hole translates points radially away from a point.
type Hole struct {
	Amount xirho.Real `xirho:"amount"`
	Origin xirho.Vec3 `xirho:"origin"`
}

func (v *Hole) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	ox, oy, oz := in.X-v.Origin[0], in.Y-v.Origin[1], in.Z-v.Origin[2]
	r := xmath.R3(ox, oy, oz)
	s := 1 + float64(v.Amount)/r
	in.X = ox*s + v.Origin[0]
	in.Y = oy*s + v.Origin[1]
	in.Z = oz*s + v.Origin[2]
	return in
}

func (v *Hole) Prep() {}

func init() {
	must("hole", func() xirho.F { return &Hole{} })
	must("spherivoid", func() xirho.F { return &Hole{} })
}
