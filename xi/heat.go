package xi

import (
	"math"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xmath"
)

// Heat applies transverse and radial waves.
type Heat struct {
	ThetaT float64 `xirho:"planar wave period"`
	ThetaP float64 `xirho:"planar wave phase,angle"`
	ThetaA float64 `xirho:"planar wave amp"`

	PhiT float64 `xirho:"axial wave period"`
	PhiP float64 `xirho:"axial wave phase,angle"`
	PhiA float64 `xirho:"axial wave amp"`

	RT float64 `xirho:"radial wave period"`
	RP float64 `xirho:"radial wave phase,angle"`
	RA float64 `xirho:"radial wave amp"`
}

func (v *Heat) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt {
	r, theta, phi := xmath.Spherical(in.X, in.Y, in.Z)
	r += v.RA * math.Sin((2*math.Pi*r+v.RP)/v.RT)
	theta += v.ThetaA * math.Sin((2*math.Pi*r+v.ThetaP)/v.ThetaT)
	phi += v.PhiA * math.Sin((2*math.Pi*r+v.PhiP)/v.PhiT)
	in.X, in.Y, in.Z = xmath.FromSpherical(r, theta, phi)
	return in
}

func (v *Heat) Prep() {}

func newHeat() xirho.Func {
	return &Heat{
		ThetaT: 1,
		ThetaP: 0,
		ThetaA: 0,
		PhiT:   1,
		PhiP:   0,
		PhiA:   0,
		RT:     1,
		RP:     0,
		RA:     0,
	}
}

func init() {
	must("heat", newHeat)
}
