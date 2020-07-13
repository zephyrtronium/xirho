package xi

import (
	"math"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xmath"
)

// Heat applies transverse and radial waves.
type Heat struct {
	ThetaT xirho.Real  `xirho:"planar wave period"`
	ThetaP xirho.Angle `xirho:"planar wave phase"`
	ThetaA xirho.Real  `xirho:"planar wave amp"`

	PhiT xirho.Real  `xirho:"axial wave period"`
	PhiP xirho.Angle `xirho:"axial wave phase"`
	PhiA xirho.Real  `xirho:"axial wave amp"`

	RT xirho.Real  `xirho:"radial wave period"`
	RP xirho.Angle `xirho:"radial wave phase"`
	RA xirho.Real  `xirho:"radial wave amp"`
}

func (v *Heat) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	r, theta, phi := xmath.Spherical(in.X, in.Y, in.Z)
	r += float64(v.RA) * math.Sin((2*math.Pi*r+float64(v.RP))/float64(v.RT))
	theta += float64(v.ThetaA) * math.Sin((2*math.Pi*r+float64(v.ThetaP))/float64(v.ThetaT))
	phi += float64(v.PhiA) * math.Sin((2*math.Pi*r+float64(v.PhiP))/float64(v.PhiT))
	in.X, in.Y, in.Z = xmath.FromSpherical(r, theta, phi)
	return in
}

func (v *Heat) Prep() {}

func newHeat() xirho.F {
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
