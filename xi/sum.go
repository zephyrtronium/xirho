package xi

import "github.com/zephyrtronium/xirho"

// Sum performs a list of functions, summing the spatial coordinates. An
// additional function controls the color coordinate.
type Sum struct {
	Funcs xirho.FuncList `xirho:"funcs"`
	Color xirho.Func     `xirho:"color"`
}

// newSum is a factory for Sum, defaulting to an empty function list.
func newSum() xirho.F {
	return &Sum{}
}

func (f *Sum) Calc(in xirho.P, rng *xirho.RNG) (out xirho.P) {
	for _, v := range f.Funcs {
		p := v.Calc(in, rng)
		out.X += p.X
		out.Y += p.Y
		out.Z += p.Z
	}
	if f.Color.F != nil {
		out.C = f.Color.Calc(in, rng).C
	} else {
		out.C = in.C
	}
	return out
}

func (f *Sum) Prep() {
	for _, v := range f.Funcs {
		v.Prep()
	}
	if f.Color.F != nil {
		f.Color.Prep()
	}
}

func init() {
	must("sum", newSum)
}
