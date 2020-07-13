package xi

import "github.com/zephyrtronium/xirho"

// Scale applies linear scaling.
type Scale struct {
	Amount xirho.Real `xirho:"amount"`
}

func (v *Scale) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	in.X *= float64(v.Amount)
	in.Y *= float64(v.Amount)
	in.Z *= float64(v.Amount)
	return in
}

func (v *Scale) Prep() {}

func init() {
	must("scale", func() xirho.F { return &Scale{Amount: 1} })
}
