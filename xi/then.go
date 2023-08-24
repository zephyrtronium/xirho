package xi

import (
	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xmath"
)

// Then performs a list of functions performed in a set order, without plotting
// intermediate results.
type Then struct {
	Funcs []xirho.Func `xirho:"funcs"`
}

// newThen is a factory for Then, defaulting to an empty function list.
func newThen() xirho.Func {
	return &Then{}
}

func (f *Then) Calc(in xirho.Pt, rng *xmath.RNG) xirho.Pt {
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
