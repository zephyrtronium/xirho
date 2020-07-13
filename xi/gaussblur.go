package xi

import (
	"github.com/zephyrtronium/crazy"
	"github.com/zephyrtronium/xirho"
)

// Gaussblur creates a spherical Gaussian blur.
type Gaussblur struct{}

func (Gaussblur) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	r := crazy.NewNormal(rng, 0, 1)
	return xirho.P{
		X: r.Next(),
		Y: r.Next(),
		Z: r.Next(),
		C: in.C,
	}
}

func (Gaussblur) Prep() {}

func init() {
	Register("gaussblur", func() xirho.F { return Gaussblur{} })
	Register("gaussian_blur", func() xirho.F { return Gaussblur{} })
}
