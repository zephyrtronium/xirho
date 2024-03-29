package xi

import (
	"math"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xmath"
)

// Bipolar does bipolar
type Bipolar struct {
	Shift float64 `xirho:"shift,angle"`
}

func (v *Bipolar) Calc(in xirho.Pt, rng *xmath.RNG) xirho.Pt {
	r := in.X*in.X + in.Y*in.Y
	y := math.Atan2(2*in.Y, r-1) - v.Shift
	// y is in (-2pi, 2pi]. Wrap to an angle.
	if y > math.Pi {
		y -= 2 * math.Pi
	} else if y < -math.Pi {
		y += 2 * math.Pi
	}
	in.X = math.Log((r+2*in.X+1)/(r-2*in.X+1)) / (2 * math.Pi)
	in.Y = y / math.Pi
	return in
}

func (v *Bipolar) Prep() {}

func init() {
	must("bipolar", func() xirho.Func { return &Bipolar{} })
}
