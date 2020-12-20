package xi

import "github.com/zephyrtronium/xirho"

// Splits spreads points away from the coordinate planes.
type Splits struct {
	X xirho.Real `xirho:"x"`
	Y xirho.Real `xirho:"y"`
	Z xirho.Real `xirho:"z"`
}

func newSplits() xirho.Func {
	return &Splits{}
}

func (v *Splits) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt {
	if in.X >= 0 {
		in.X += float64(v.X)
	} else {
		in.X -= float64(v.X)
	}
	if in.Y >= 0 {
		in.Y += float64(v.Y)
	} else {
		in.Y -= float64(v.Y)
	}
	if in.Z >= 0 {
		in.Z += float64(v.Z)
	} else {
		in.Z -= float64(v.Z)
	}
	return in
}

func (v *Splits) Prep() {}

func init() {
	must("splits", newSplits)
	must("splits3D", newSplits)
}
