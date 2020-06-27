package xi

import "github.com/zephyrtronium/xirho"

// Affine performs an affine transform.
type Affine struct {
	Ax xirho.Affine `xirho:"transform"`
}

// newAffine is a factory for Affine, defaulting to an identity transform.
func newAffine() xirho.F {
	tx := &Affine{}
	tx.Ax.Eye()
	return tx
}

func (v *Affine) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	x, y, z := xirho.Tx(&v.Ax, in.X, in.Y, in.Z)
	return xirho.P{
		X: x,
		Y: y,
		Z: z,
		C: in.C,
	}
}

func (v *Affine) Prep() {}

func init() {
	must("affine", newAffine)
	must("linear", newAffine)
	must("linear3D", newAffine)
}
