package xi

import (
	"math"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xmath"
)

// Noise produces a noisy blotch thing.
type Noise struct{}

func (Noise) Calc(in xirho.Pt, rng *xmath.RNG) xirho.Pt {
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
