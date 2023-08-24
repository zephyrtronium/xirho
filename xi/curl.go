package xi

import (
	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xmath"
)

// Curl does curl
type Curl struct {
	C1 float64 `xirho:"c1"`
	C2 float64 `xirho:"c2"`
}

func (v *Curl) Calc(in xirho.Pt, rng *xmath.RNG) xirho.Pt {
	t1 := 1 + v.C1*in.X + v.C2*(in.X*in.X-in.Y*in.Y)
	t2 := v.C1*in.Y + 2*v.C2*in.X*in.Y
	r := t1*t1 + t2*t2
	in.X = (in.X*t1 + in.Y*t2) / r
	in.Y = (in.Y*t1 - in.X*t2) / r
	return in
}

func (v *Curl) Prep() {}

func init() {
	must("curl", func() xirho.Func { return &Curl{} })
}
