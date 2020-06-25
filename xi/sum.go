package xi

import "github.com/zephyrtronium/xirho"

// Sum performs a list of functions, summing the spatial coordinates and
// averaging the color coordinate.
type Sum struct {
	funcs []xirho.F

	p []xirho.Param
}

// NewSum creates a Sum function.
func NewSum(funcs ...xirho.F) xirho.F {
	f := &Sum{
		funcs: append([]xirho.F{}, funcs...), // copy
	}
	f.p = []xirho.Param{
		xirho.FuncListParam(&f.funcs, "funcs"),
	}
	return f
}

func (f *Sum) Calc(in xirho.P, rng *xirho.RNG) (out xirho.P) {
	for _, v := range f.funcs {
		p := v.Calc(in, rng)
		out.X += p.X
		out.Y += p.Y
		out.Z += p.Z
		out.C += p.C
	}
	out.C /= float64(len(f.funcs))
	return out
}

func (f *Sum) Params() []xirho.Param {
	return f.p
}
