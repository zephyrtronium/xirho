package xi

import (
	"math"

	"github.com/zephyrtronium/crazy"
	"github.com/zephyrtronium/xirho"
)

// Blur produces a noisy solid circle with unit radius.
type Blur struct{}

// newBlur is a factory for Blur.
func newBlur() xirho.F {
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

func (Blur) Prep() {}

func init() {
	must("blur", newBlur)
}
