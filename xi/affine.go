package xi

import "github.com/zephyrtronium/xirho"

type Affine struct {
	ax    xirho.Ax
	c, sp float64

	p []xirho.Param
}

func NewAffine(ax xirho.Ax, color, speed float64) xirho.F {
	tx := &Affine{
		ax: ax,
		c:  color,
		sp: speed,
	}
	tx.p = []xirho.Param{
		xirho.AffineParam(&tx.ax, "transform"),
		xirho.RealParam(&tx.c, "color", true, 0, 1),
		xirho.RealParam(&tx.sp, "color weight", true, 0, 1),
	}
	return tx
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
