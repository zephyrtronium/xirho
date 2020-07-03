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
	"disc":        parseDisc,
	"flatten":     parseFlatten,
	"julia":       parseJulia,
	"julian":      parseJulian,
	"polar":       parsePolar,
	"spherical":   parseSpherical,
	"spherical3D": parseSpherical3D,
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

func parsePolar(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	in.Funcs = append(in.Funcs, maybeScaled(xi.Polar{}, attrs["polar"]))
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
