package xi

import (
	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xmath"
)

// Hole translates points radially away from a point.
type Hole struct {
	Amount float64    `xirho:"amount"`
	Origin [3]float64 `xirho:"origin"`
}

func (v *Hole) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt {
	ox, oy, oz := in.X-v.Origin[0], in.Y-v.Origin[1], in.Z-v.Origin[2]
	r := xmath.R3(ox, oy, oz)
	s := 1 + v.Amount/r
	in.X = ox*s + v.Origin[0]
	in.Y = oy*s + v.Origin[1]
	in.Z = oz*s + v.Origin[2]
	return in
}

func (v *Hole) Prep() {}

func init() {
	must("hole", func() xirho.Func { return &Hole{} })
	must("spherivoid", func() xirho.Func { return &Hole{} })
}
