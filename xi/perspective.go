package xi

import "github.com/zephyrtronium/xirho"

// Perspective applies a perspective transform to the 3D spatial coordinates.
type Perspective struct {
	Distance float64 `xirho:"distance"`
}

// newPerspective is a factory for Perspective, defaulting to a distance of 1.
func newPerspective() xirho.Func {
	return &Perspective{Distance: 1}
}

func (f *Perspective) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt {
	r := f.Distance / in.Z
	in.X *= r
	in.Y *= r
	return in
}

func (f *Perspective) Prep() {}

func init() {
	must("perspective", newPerspective)
}
