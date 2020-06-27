package xi

import (
	"math"

	"github.com/zephyrtronium/crazy"
	"github.com/zephyrtronium/xirho"
)

// JuliaN does julian
type JuliaN struct {
	Power xirho.Int  `xirho:"power"`
	Dist  xirho.Real `xirho:"dist"`
}

// newJuliaN is a factory for JuliaN, defaulting Power to 3 and Dist to 1.
func newJuliaN() xirho.F {
	return &JuliaN{Power: 3, Dist: 1}
}

func (f *JuliaN) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	p1 := int64(f.Power)
	if p1 < 0 {
		p1 = -p1
	}
	p3 := float64(crazy.RNG{rng}.Uintn(uint(p1)))
	t := (math.Atan2(in.Y, in.X) + 2*math.Pi*p3) / float64(f.Power)
	r := math.Pow(math.Hypot(in.X, in.Y), float64(f.Dist)/float64(f.Power))
	s, c := math.Sincos(t)
	in.X = r * c
	in.Y = r * s
	return in
}

func (f *JuliaN) Prep() {}

func init() {
	must("julian", newJuliaN)
}
