package xi

import (
	"math"

	"github.com/zephyrtronium/xirho"
)

// Foci does foci
type Foci struct{}

func (Foci) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	ex := math.Exp(in.X) / 2
	sy, cy := math.Sincos(in.Y)
	d := (ex + 1/(4*ex) + cy)
	in.X = (ex - 1/(4*ex)) / d
	in.Y = sy / d
	return in
}

func (Foci) Prep() {}

func init() {
	must("foci", func() xirho.F { return Foci{} })
}
