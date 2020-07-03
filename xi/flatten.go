package xi

import "github.com/zephyrtronium/xirho"

// Flatten zeros the Z coordinate of the input.
type Flatten struct{}

func (Flatten) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	in.Z = 0
	return in
}

func (Flatten) Prep() {}

func init() {
	Register("flatten", func() xirho.F { return Flatten{} })
}
