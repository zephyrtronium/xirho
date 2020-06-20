package variations

import "github.com/zephyrtronium/xirho"

type Affine struct {
	ax    xirho.Ax
	c, sp float64

	p []xirho.Param
}

func NewAffine() xirho.F {
	ax := &Affine{
		ax: *xirho.Eye(),
	}
	ax.p = []xirho.Param{
		xirho.AffineParam(&ax.ax, "transform"),
		xirho.RealParam(&ax.c, "color", true, 0, 1),
		xirho.RealParam(&ax.sp, "color weight", true, 0, 1),
	}
	return ax
}

func (v *Affine) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	x, y, z := xirho.Tx(&v.ax, in.X, in.Y, in.Z)
	return xirho.P{
		X: x,
		Y: y,
		Z: z,
		C: v.sp*in.C + (1-v.sp)*v.c,
	}
}

func (v *Affine) Params() []xirho.Param {
	return v.p
}
