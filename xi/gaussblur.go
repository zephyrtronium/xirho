package xi

import (
	"github.com/zephyrtronium/xirho"
)

// Gaussblur creates a spherical Gaussian blur.
type Gaussblur struct{}

func (Gaussblur) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	return xirho.P{
		X: rng.Normal(),
		Y: rng.Normal(),
		Z: rng.Normal(),
		C: in.C,
	}
}

func (Gaussblur) Prep() {}

func init() {
	Register("gaussblur", func() xirho.F { return Gaussblur{} })
	Register("gaussian_blur", func() xirho.F { return Gaussblur{} })
}
