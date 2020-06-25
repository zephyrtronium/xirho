package xi

import (
	"math"

	"github.com/zephyrtronium/crazy"
	"github.com/zephyrtronium/xirho"
)

// JuliaN does julian
type JuliaN struct {
	power int64
	dist  float64

	p []xirho.Param
}

func NewJuliaN() xirho.F {
	f := &JuliaN{
		power: 3,
		dist:  1,
	}
	f.p = []xirho.Param{
		xirho.IntParam(&f.power, "power", false, 0, 0),
		xirho.RealParam(&f.dist, "dist", false, 0, 0),
	}
	return f
}

func (f *JuliaN) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	p1 := f.power
	if p1 < 0 {
		p1 = -p1
	}
	p3 := float64(crazy.RNG{rng}.Uintn(uint(p1)))
	t := (math.Atan(in.Y/in.X) + 2*math.Pi*p3) / float64(f.power)
	r := math.Pow(math.Hypot(in.X, in.Y), f.dist/float64(f.power))
	s, c := math.Sincos(t)
	in.X = r * c
	in.Y = r * s
	return in
}

func (f *JuliaN) Params() []xirho.Param {
	return f.p
}
