package xi

import "github.com/zephyrtronium/xirho"

// Perspective applies a perspective transform to the 3D spatial coordinates.
type Perspective struct {
	Distance xirho.Real `xirho:"distance"`
}

// NewPerspective is a factory for Perspective, defaulting to a distance of 1.
func NewPerspective() xirho.F {
	return &Perspective{Distance: 1}
}

func (f *Perspective) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	r := float64(f.Distance) / in.Z
	in.X *= r
	in.Y *= r
	return in
}

func (f *Perspective) Prep() {}
