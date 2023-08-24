package xi

import (
	"math"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xmath"
)

// Hemisphere projects the plane onto a half-sphere.
type Hemisphere struct{}

func (Hemisphere) Calc(in xirho.Pt, rng *xmath.RNG) xirho.Pt {
	r := 1 / math.Sqrt(in.X*in.X+in.Y*in.Y+1)
	return xirho.Pt{X: r * in.X, Y: r * in.Y, Z: r, C: in.C}
}

func (Hemisphere) Prep() {}

func init() {
	must("hemisphere", func() xirho.Func { return Hemisphere{} })
}
