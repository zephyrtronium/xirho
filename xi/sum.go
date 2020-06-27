package xi

import "github.com/zephyrtronium/xirho"

// Sum performs a list of functions, summing the spatial coordinates and
// averaging the color coordinate.
type Sum struct {
	Funcs xirho.FuncList `xirho:"funcs"`
}

// NewSum is a factory for Sum, defaulting to an empty function list.
func NewSum() xirho.F {
	return &Sum{}
}

func (f *Sum) Calc(in xirho.P, rng *xirho.RNG) (out xirho.P) {
	for _, v := range f.Funcs {
		p := v.Calc(in, rng)
		out.X += p.X
		out.Y += p.Y
		out.Z += p.Z
		out.C += p.C
	}
	out.C /= float64(len(f.Funcs))
	return out
}

func (f *Sum) Prep() {}
