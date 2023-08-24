package xi

import (
	"math"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xmath"
)

// Disc does disc.
type Disc struct{}

// newDisc is a factory for Disc.
func newDisc() xirho.Func {
	return Disc{}
}

func (Disc) Calc(in xirho.Pt, rng *xmath.RNG) xirho.Pt {
	t := math.Atan2(in.X, in.Y) / math.Pi
	s, c := math.Sincos(math.Hypot(in.X, in.Y) * math.Pi)
	in.X = t * s
	in.Y = t * c
	return in
}

func (Disc) Prep() {}

func init() {
	must("disc", newDisc)
}
