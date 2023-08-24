package xi

import "github.com/zephyrtronium/xirho"

// Sum performs a list of functions, summing the spatial coordinates. An
// additional function controls the color coordinate.
type Sum struct {
	Funcs []xirho.Func `xirho:"funcs"`
	Color xirho.Func   `xirho:"color,optional"`
}

// newSum is a factory for Sum, defaulting to an empty function list.
func newSum() xirho.Func {
	return &Sum{}
}

func (f *Sum) Calc(in xirho.Pt, rng *xirho.RNG) (out xirho.Pt) {
	for _, v := range f.Funcs {
		p := v.Calc(in, rng)
		out.X += p.X
		out.Y += p.Y
		out.Z += p.Z
	}
	if f.Color != nil {
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
	if f.Color != nil {
		f.Color.Prep()
	}
}

func init() {
	must("sum", newSum)
}
