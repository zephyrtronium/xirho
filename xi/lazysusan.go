package xi

import (
	"math"

	"github.com/zephyrtronium/xirho"
)

// LazySusan transforms points by different affine transforms depending on
// whether the input is inside a selection sphere.
type LazySusan struct {
	Inside  xirho.Affine `xirho:"inside"`
	Outside xirho.Affine `xirho:"outside"`
	Center  [3]float64   `xirho:"center"`
	Radius  float64      `xirho:"radius"`
	Spread  float64      `xirho:"spread"`
	TwistZ  float64      `xirho:"twistZ"`
}

func (v *LazySusan) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt {
	x, y, z := in.X-v.Center[0], in.Y-v.Center[1], in.Z-v.Center[2]
	r := math.Sqrt(x*x + y*y + z*z)
	var ax xirho.Affine
	if r < v.Radius {
		ax = v.Inside
		if v.TwistZ != 0 {
			ax.RotZ(v.TwistZ * (v.Radius - r))
		}
	} else {
		ax = v.Outside
		if v.Spread != 0 {
			sc := v.Radius * (1 + v.Spread/r)
			ax.Scale(sc, sc, sc)
		}
	}
	x, y, z = xirho.Tx(&ax, x, y, z)
	return xirho.Pt{
		X: x + v.Center[0],
		Y: y + v.Center[1],
		Z: z + v.Center[2],
		C: in.C,
	}
}

func (v *LazySusan) Prep() {}

func newLazySusan() xirho.Func {
	return &LazySusan{
		Inside:  xirho.Eye(),
		Outside: xirho.Eye(),
		Radius:  1,
	}
}

func init() {
	must("lazysusan", newLazySusan)
}
