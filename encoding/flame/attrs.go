package flame

import (
	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xi"
)

// Funcs maps Flame xform attribute names to functions which parse attributes
// into function instances. Parameters of functions should not included here,
// e.g. julian is an entry but julian_power is not, as the parser should handle
// the parameters.
var Funcs = map[string]Parser{
	"linear":      parseLinear,
	"linear3D":    parseLinear,
	"blur":        parseBlur,
	"pre_blur":    parsePreblur,
	"bubble":      parseBubble,
	"elliptic":    parseElliptic,
	"disc":        parseDisc,
	"flatten":     parseFlatten,
	"julia":       parseJulia,
	"julian":      parseJulian,
	"mobius":      parseMobius,
	"mobiq":       parseMobiq,
	"polar":       parsePolar,
	"spherical":   parseSpherical,
	"spherical3D": parseSpherical3D,
	"splits":      parseSplits,
	"splits3D":    parseSplits3D,
}

// Parser is a function which parses a xirho function from XML attributes.
// Parsers should add the parsed function to one of pre, in, or post, and they
// should not modify attrs.
type Parser func(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax)

func parseLinear(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	v := attrs["linear"] + attrs["linear3D"]
	ax.Eye().Scale(v, v, v)
	in.Funcs = append(in.Funcs, &xi.Affine{Ax: ax})
}

func parseSpherical(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	f := xi.Then{
		Funcs: []xirho.F{
			xi.Spherical{},
			xi.Flatten{},
		},
	}
	in.Funcs = append(in.Funcs, maybeScaled(&f, attrs["spherical"]))
}

func parseSpherical3D(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	in.Funcs = append(in.Funcs, maybeScaled(xi.Spherical{}, attrs["spherical3D"]))
}

func parseBlur(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	in.Funcs = append(in.Funcs, maybeScaled(xi.Blur{}, attrs["blur"]))
}

func parsePreblur(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	pre.Funcs = append(pre.Funcs, maybeScaled(xi.Blur{}, attrs["pre_blur"]))
}

func parseBubble(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	in.Funcs = append(in.Funcs, maybeScaled(xi.Bubble{}, attrs["bubble"]))
}

func parseElliptic(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	in.Funcs = append(in.Funcs, maybeScaled(xi.CElliptic{}, attrs["elliptic"]))
}

func parseDisc(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	in.Funcs = append(in.Funcs, maybeScaled(xi.Disc{}, attrs["disc"]))
}

func parseFlatten(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	pre.Funcs = append(pre.Funcs, xi.Flatten{})
}

func parseJulia(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	f := xi.JuliaN{
		Power: 2,
		Dist:  1,
	}
	in.Funcs = append(in.Funcs, maybeScaled(&f, attrs["julia"]))
}

func parseJulian(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	f := xi.JuliaN{
		Power: xirho.Int(attrs["julian_power"]),
		Dist:  xirho.Real(attrs["julian_dist"]),
	}
	in.Funcs = append(in.Funcs, maybeScaled(&f, attrs["julian"]))
}

func parseMobius(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	f := xi.Mobius{
		Ar:     xirho.Real(attrs["Re_A"]),
		Avec:   xirho.Vec3{attrs["Im_A"], 0, 0},
		Br:     xirho.Real(attrs["Re_B"]),
		Bvec:   xirho.Vec3{attrs["Im_B"], 0, 0},
		Cr:     xirho.Real(attrs["Re_C"]),
		Cvec:   xirho.Vec3{attrs["Im_C"], 0, 0},
		Dr:     xirho.Real(attrs["Re_D"]),
		Dvec:   xirho.Vec3{attrs["Im_D"], 0, 0},
		InZero: 3,
	}
	t := xi.Then{
		Funcs: []xirho.F{
			xi.Flatten{},
			&f,
		},
	}
	in.Funcs = append(in.Funcs, &t)
}

func parseMobiq(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	f := xi.Mobius{
		Ar:     xirho.Real(attrs["mobiq_at"]),
		Avec:   xirho.Vec3{attrs["mobiq_ax"], attrs["mobiq_ay"], attrs["mobiq_az"]},
		Br:     xirho.Real(attrs["mobiq_bt"]),
		Bvec:   xirho.Vec3{attrs["mobiq_bx"], attrs["mobiq_by"], attrs["mobiq_bz"]},
		Cr:     xirho.Real(attrs["mobiq_ct"]),
		Cvec:   xirho.Vec3{attrs["mobiq_cx"], attrs["mobiq_cy"], attrs["mobiq_cz"]},
		Dr:     xirho.Real(attrs["mobiq_dt"]),
		Dvec:   xirho.Vec3{attrs["mobiq_dx"], attrs["mobiq_dy"], attrs["mobiq_dz"]},
		InZero: 3,
	}
	in.Funcs = append(in.Funcs, &f)
}

func parsePolar(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	in.Funcs = append(in.Funcs, maybeScaled(xi.Polar{}, attrs["polar"]))
}

func parseSplits(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	f := xi.Splits{
		X: xirho.Real(attrs["splits_x"]),
		Y: xirho.Real(attrs["splits_y"]),
	}
	t := xi.Then{
		Funcs: []xirho.F{
			xi.Flatten{},
			&f,
		},
	}
	in.Funcs = append(in.Funcs, &t)
}

func parseSplits3D(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	f := xi.Splits{
		X: xirho.Real(attrs["splits3D_x"]),
		Y: xirho.Real(attrs["splits3D_y"]),
		Z: xirho.Real(attrs["splits3D_z"]),
	}
	in.Funcs = append(in.Funcs, &f)
}

// maybeScaled returns f if v is 1 or a Then with f and Scale by v otherwise.
func maybeScaled(f xirho.F, v float64) xirho.F {
	if v == 1 {
		return f
	}
	return &xi.Then{
		Funcs: []xirho.F{
			f,
			&xi.Scale{Amount: xirho.Real(v)},
		},
	}
}
