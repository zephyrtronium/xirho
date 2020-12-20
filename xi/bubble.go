package xi

import "github.com/zephyrtronium/xirho"

// Bubble maps the plane to a sphere.
type Bubble struct{}

// newBubble is a factory for Bubble.
func newBubble() xirho.Func {
	return Bubble{}
}

func (Bubble) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt {
	r := 4 / (in.X*in.X + in.Y*in.Y + in.Z*in.Z + 4)
	in.X *= r
	in.Y *= r
	in.Z *= r
	return in
}

func (Bubble) Prep() {}

func init() {
	must("bubble", newBubble)
}
