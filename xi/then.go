package xi

import "github.com/zephyrtronium/xirho"

// Then performs a list of functions performed in a set order, without plotting
// intermediate results.
type Then struct {
	funcs []xirho.F

	p []xirho.Param
}

// NewThen creates a Then function.
func NewThen(funcs ...xirho.F) xirho.F {
	f := &Then{
		funcs: append([]xirho.F{}, funcs...), // copy
	}
	f.p = []xirho.Param{
		xirho.FuncListParam(&f.funcs, "funcs"),
	}
	return f
}

func (f *Then) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	for _, v := range f.funcs {
		in = v.Calc(in, rng)
	}
	return in
}

func (f *Then) Params() []xirho.Param {
	return f.p
}
