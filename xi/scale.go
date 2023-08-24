package xi

import "github.com/zephyrtronium/xirho"

// Scale applies linear scaling.
type Scale struct {
	Amount float64 `xirho:"amount"`
}

func (v *Scale) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt {
	in.X *= v.Amount
	in.Y *= v.Amount
	in.Z *= v.Amount
	return in
}

func (v *Scale) Prep() {}

func init() {
	must("scale", func() xirho.Func { return &Scale{Amount: 1} })
}
