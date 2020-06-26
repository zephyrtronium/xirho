package xi

import (
	"math"

	"github.com/zephyrtronium/xirho"
)

// Disc does disc
type Disc struct{}

func NewDisc() xirho.F {
	return Disc{}
}

func (Disc) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	t := math.Atan2(in.X, in.Y) / math.Pi
	s, c := math.Sincos(math.Hypot(in.X, in.Y) * math.Pi)
	in.X = t * s
	in.Y = t * c
	return in
}

func (Disc) Params() []xirho.Param {
	return nil
}
