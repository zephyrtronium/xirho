package xi

import (
	"math"

	"github.com/zephyrtronium/xirho"
)

type Polar struct{}

func NewPolar() xirho.F {
	return Polar{}
}

func (Polar) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	t := math.Atan2(in.X, in.Y) / math.Pi
	r := math.Hypot(in.X, in.Y) - 1
	in.X = t
	in.Y = r
	return in
}

func (Polar) Params() []xirho.Param {
	return nil
}
