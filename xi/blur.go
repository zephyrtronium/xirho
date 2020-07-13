package xi

import (
	"math"

	"github.com/zephyrtronium/xirho"
)

// Blur produces a noisy solid circle with unit radius.
type Blur struct{}

// newBlur is a factory for Blur.
func newBlur() xirho.F {
	return Blur{}
}

func (Blur) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	s, c := math.Sincos(2 * math.Pi * rng.Uniform())
	r := rng.Uniform()
	in.X = r * c
	in.Y = r * s
	return in
}

func (Blur) Prep() {}

func init() {
	must("blur", newBlur)
}
