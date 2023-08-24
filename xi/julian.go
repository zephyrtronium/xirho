package xi

import (
	"math"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xmath"
)

// JuliaN does julian
type JuliaN struct {
	Power int64   `xirho:"power"`
	Dist  float64 `xirho:"dist"`
}

// newJuliaN is a factory for JuliaN, defaulting Power to 3 and Dist to 1.
func newJuliaN() xirho.Func {
	return &JuliaN{Power: 3, Dist: 1}
}

func (f *JuliaN) Calc(in xirho.Pt, rng *xmath.RNG) xirho.Pt {
	p1 := int(f.Power)
	if p1 < 0 {
		p1 = -p1
	}
	p3 := float64(rng.Intn(p1))
	t := (math.Atan2(in.Y, in.X) + 2*math.Pi*p3) / float64(f.Power)
	r := math.Pow(math.Hypot(in.X, in.Y), f.Dist/float64(f.Power))
	s, c := math.Sincos(t)
	in.X = r * c
	in.Y = r * s
	return in
}

func (f *JuliaN) Prep() {}

func init() {
	must("julian", newJuliaN)
}
