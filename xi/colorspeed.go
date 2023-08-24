package xi

import "github.com/zephyrtronium/xirho"

// ColorSpeed performs exponential smoothing on the input color coordinate
// toward a chosen color.
type ColorSpeed struct {
	// Color is the color coordinate toward which inputs move.
	Color float64 `xirho:"color,0,1"`
	// Speed is the smoothing rate. A value of 0 means the output color always
	// equals Color; a value of 1 means the output color always equals the
	// input color.
	Speed float64 `xirho:"speed,0,1"`
}

// newColorSpeed is a factory for ColorSpeed, defaulting Color to 0 and Speed
// to 1.
func newColorSpeed() xirho.Func {
	return &ColorSpeed{Speed: 1}
}

func (f *ColorSpeed) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt {
	in.C = in.C*f.Speed + (1-f.Speed)*f.Color
	return in
}

func (f *ColorSpeed) Prep() {}

func init() {
	must("colorspeed", newColorSpeed)
}
