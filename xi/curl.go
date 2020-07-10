package xi

import "github.com/zephyrtronium/xirho"

// Curl does curl
type Curl struct {
	C1 xirho.Real `xirho:"c1"`
	C2 xirho.Real `xirho:"c2"`
}

func (v *Curl) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	t1 := 1 + float64(v.C1)*in.X + float64(v.C2)*(in.X*in.X-in.Y*in.Y)
	t2 := float64(v.C1)*in.Y + 2*float64(v.C2)*in.X*in.Y
	r := t1*t1 + t2*t2
	in.X = (in.X*t1 + in.Y*t2) / r
	in.Y = (in.Y*t1 - in.X*t2) / r
	return in
}

func (v *Curl) Prep() {}

func init() {
	Register("curl", func() xirho.F { return &Curl{} })
}
