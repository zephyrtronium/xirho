package xi

import (
	"github.com/zephyrtronium/xirho"
)

// Gaussblur creates a spherical Gaussian blur.
type Gaussblur struct{}

func (Gaussblur) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt {
	return xirho.Pt{
		X: rng.Normal(),
		Y: rng.Normal(),
		Z: rng.Normal(),
		C: in.C,
	}
}

func (Gaussblur) Prep() {}

func init() {
	must("gaussblur", func() xirho.Func { return Gaussblur{} })
	must("gaussian_blur", func() xirho.Func { return Gaussblur{} })
}
