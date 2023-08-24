package xi

import "github.com/zephyrtronium/xirho"

// Splits spreads points away from the coordinate planes.
type Splits struct {
	X float64 `xirho:"x"`
	Y float64 `xirho:"y"`
	Z float64 `xirho:"z"`
}

func newSplits() xirho.Func {
	return &Splits{}
}

func (v *Splits) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt {
	if in.X >= 0 {
		in.X += v.X
	} else {
		in.X -= v.X
	}
	if in.Y >= 0 {
		in.Y += v.Y
	} else {
		in.Y -= v.Y
	}
	if in.Z >= 0 {
		in.Z += v.Z
	} else {
		in.Z -= v.Z
	}
	return in
}

func (v *Splits) Prep() {}

func init() {
	must("splits", newSplits)
	must("splits3D", newSplits)
}
