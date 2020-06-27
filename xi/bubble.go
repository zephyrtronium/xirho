package xi

import "github.com/zephyrtronium/xirho"

// Bubble maps the plane to a sphere.
type Bubble struct{}

// NewBubble is a factory for Bubble.
func NewBubble() xirho.F {
	return Bubble{}
}

func (Bubble) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	r := 4 / (in.X*in.X + in.Y*in.Y + in.Z*in.Z + 4)
	in.X *= r
	in.Y *= r
	in.Z *= r
	return in
}

func (Bubble) Prep() {}
