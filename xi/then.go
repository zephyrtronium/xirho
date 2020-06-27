package xi

import "github.com/zephyrtronium/xirho"

// Then performs a list of functions performed in a set order, without plotting
// intermediate results.
type Then struct {
	Funcs xirho.FuncList `xirho:"funcs"`
}

// newThen is a factory for Then, defaulting to an empty function list.
func newThen() xirho.F {
	return &Then{}
}

func (f *Then) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	for _, v := range f.Funcs {
		in = v.Calc(in, rng)
	}
	return in
}

func (f *Then) Prep() {
	for _, v := range f.Funcs {
		v.Prep()
	}
}

func init() {
	must("then", newThen)
}
