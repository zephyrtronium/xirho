package xi

import "github.com/zephyrtronium/xirho"

type ColorSpeed struct {
	color, speed float64

	p []xirho.Param
}

func NewColorSpeed(color, speed float64) xirho.F {
	c := &ColorSpeed{color: color, speed: speed}
	c.p = []xirho.Param{
		xirho.RealParam(&c.color, "color", true, 0, 1),
		xirho.RealParam(&c.speed, "speed", true, 0, 1),
	}
	return c
}

func (f *ColorSpeed) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	in.C = in.C*f.speed + (1-f.speed)*f.color
	return in
}

func (f *ColorSpeed) Params() []xirho.Param {
	return f.p
}
