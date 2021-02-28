package xi

import (
	"math"

	"github.com/zephyrtronium/xirho"
)

// Noise produces a noisy blotch thing.
type Noise struct{}

func (Noise) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt {
	s, c := math.Sincos(2 * math.Pi * rng.Uniform())
	r := rng.Uniform()
	in.X *= r * c
	in.Y *= r * s
	return in
}

func (Noise) Prep() {}

func init() {
	must("noise", func() xirho.Func { return Noise{} })
}
