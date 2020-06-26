package xi

import (
	"math"

	"github.com/zephyrtronium/crazy"
	"github.com/zephyrtronium/xirho"
)

type Blur struct{}

func NewBlur() xirho.F {
	return Blur{}
}

func (Blur) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	d := crazy.Uniform0_1{rng}
	s, c := math.Sincos(2 * math.Pi * d.Next())
	r := d.Next()
	in.X = r * c
	in.Y = r * s
	return in
}

func (Blur) Params() []xirho.Param {
	return nil
}
