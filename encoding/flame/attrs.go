package flame

import (
	"math"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xi"
)

// Funcs maps Flame xform attribute names to functions which parse attributes
// into function instances. Parameters of functions should not included here,
// e.g. julian is an entry but julian_power is not, as the parser should handle
// the parameters.
var Funcs = map[string]Parser{
	"linear":        parseLinear,
	"linear3D":      parseLinear,
	"bipolar":       parseBipolar,
	"blur":          parseBlur,
	"pre_blur":      parsePreblur,
	"bubble":        parseBubble,
	"elliptic":      parseElliptic,
	"exp":           parseExp,
	"expo":          parseExpo,
	"curl":          parseCurl,
	"cylinder":      parseCylinder,
	"disc":          parseDisc,
	"flatten":       parseFlatten,
	"foci":          parseFoci,
	"gaussian_blur": parseGaussblur,
	"post_heat":     parsePostHeat,
	"julia":         parseJulia,
	"julian":        parseJulian,
	"lazysusan":     parseLazySusan,
	"log":           parseLog,
	"mobius":        parseMobius,
	"mobiq":         parseMobiq,
	"polar":         parsePolar,
	"rod":           parseRod,
	"scry":          parseScry,
	"spherical":     parseSpherical,
	"spherical3D":   parseSpherical3D,
	"pre_spherical": parsePrespherical,
	"splits":        parseSplits,
	"splits3D":      parseSplits3D,
	"unpolar":       parseUnpolar,
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

func parseBipolar(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	s := xirho.Angle(math.Mod(attrs["bipolar_shift"], 2))
	if s > 1 {
		s -= 2
	} else if s <= -1 {
		s += 2
	}
	in.Funcs = append(in.Funcs, maybeScaled(&xi.Bipolar{Shift: s * math.Pi}, attrs["bipolar"]))
}

func parseBlur(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	in.Funcs = append(in.Funcs, maybeScaled(xi.Blur{}, attrs["blur"]))
}

func parsePreblur(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	pre.Funcs = append(pre.Funcs, maybeScaled(xi.Gaussblur{}, attrs["pre_blur"]))
}

func parseBubble(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	in.Funcs = append(in.Funcs, maybeScaled(xi.Bubble{}, attrs["bubble"]))
}

func parseElliptic(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	in.Funcs = append(in.Funcs, maybeScaled(xi.CElliptic{}, attrs["elliptic"]))
}

func parseCurl(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	v := xi.Curl{
		C1: xirho.Real(attrs["curl_c1"]),
		C2: xirho.Real(attrs["curl_c2"]),
	}
	in.Funcs = append(in.Funcs, maybeScaled(&v, attrs["curl"]))
}

func parseCylinder(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	in.Funcs = append(in.Funcs, maybeScaled(xi.Cylinder{}, attrs["cylinder"]))
}

func parseDisc(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	in.Funcs = append(in.Funcs, maybeScaled(xi.Disc{}, attrs["disc"]))
}

func parseExp(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	in.Funcs = append(in.Funcs, maybeScaled(&xi.Exp{Base: math.E}, attrs["exp"]))
}

func parseExpo(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	z := complex(attrs["expo_real"], attrs["expo_imaginary"])
	in.Funcs = append(in.Funcs, maybeScaled(&xi.Exp{Base: xirho.Complex(z)}, attrs["expo"]))
}

func parseFlatten(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	pre.Funcs = append(pre.Funcs, xi.Flatten{})
}

func parseFoci(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	in.Funcs = append(in.Funcs, xi.Foci{})
}

func parseGaussblur(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	t := xi.Then{
		Funcs: xirho.FuncList{
			xi.Gaussblur{},
			xi.Flatten{},
		},
	}
	in.Funcs = append(in.Funcs, maybeScaled(&t, attrs["gaussian_blur"]))
}

func parsePostHeat(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	f := xi.Heat{
		ThetaT: xirho.Real(attrs["post_heat_theta_period"]),
		ThetaP: xirho.Angle(attrs["post_heat_theta_phase"]),
		ThetaA: xirho.Real(attrs["post_heat"] * attrs["post_heat_theta_amp"]),
		PhiT:   xirho.Real(attrs["post_heat_phi_period"]),
		PhiP:   xirho.Angle(attrs["post_heat_phi_phase"]),
		PhiA:   xirho.Real(attrs["post_heat"] * attrs["post_heat_phi_amp"]),
		RT:     xirho.Real(attrs["post_heat_r_period"]),
		RP:     xirho.Angle(attrs["post_heat_r_phase"]),
		RA:     xirho.Real(attrs["post_heat"] * attrs["post_heat_r_amp"]),
	}
	post.Funcs = append(post.Funcs, &f)
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

func parseLazySusan(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	f := xi.LazySusan{
		Inside:  xirho.Eye(),
		Outside: xirho.Eye(),
		Center:  xirho.Vec3{attrs["lazysusan_x"], -attrs["lazysusan_y"], 0},
		Radius:  xirho.Real(attrs["lazysusan"]),
		Spread:  xirho.Real(attrs["lazysusan_space"]),
		TwistZ:  xirho.Real(-attrs["lazysusan_twist"]),
	}
	f.Outside.Scale(1, 1, 0)
	f.Inside.RotZ(attrs["lazysusan_spin"]).Scale(1, 1, 0)
	in.Funcs = append(in.Funcs, &f)
}

func parseLog(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	f := xi.Log{Base: xirho.Complex(complex(attrs["log_base"], 0))}
	if f.Base == 0 {
		f.Base = math.E
	}
	in.Funcs = append(in.Funcs, maybeScaled(&f, attrs["log"]))
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
	in.Funcs = append(in.Funcs, maybeScaled(&t, attrs["mobius"]))
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
	in.Funcs = append(in.Funcs, maybeScaled(&f, attrs["mobiq"]))
}

func parsePolar(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	in.Funcs = append(in.Funcs, maybeScaled(xi.Polar{}, attrs["polar"]))
}

func parseRod(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	in.Funcs = append(in.Funcs, &xi.Rod{Radius: xirho.Real(attrs["rod"])})
}

func parseScry(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	f := xi.Scry{Radius: xirho.Real(attrs["scry"])}
	t := xi.Then{
		Funcs: xirho.FuncList{
			xi.Flatten{},
			&f,
		},
	}
	in.Funcs = append(in.Funcs, &t)
}

func parseSpherical(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	f := xi.Then{
		Funcs: []xirho.F{
			xi.Flatten{},
			xi.Spherical{},
		},
	}
	in.Funcs = append(in.Funcs, maybeScaled(&f, attrs["spherical"]))
}

func parseSpherical3D(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	in.Funcs = append(in.Funcs, maybeScaled(xi.Spherical{}, attrs["spherical3D"]))
}

func parsePrespherical(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	f := xi.Then{
		Funcs: []xirho.F{
			xi.Flatten{},
			xi.Spherical{},
		},
	}
	pre.Funcs = append(pre.Funcs, maybeScaled(&f, attrs["pre_spherical"]))
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
	in.Funcs = append(in.Funcs, maybeScaled(&t, attrs["splits"]))
}

func parseSplits3D(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	f := xi.Splits{
		X: xirho.Real(attrs["splits3D_x"]),
		Y: xirho.Real(attrs["splits3D_y"]),
		Z: xirho.Real(attrs["splits3D_z"]),
	}
	in.Funcs = append(in.Funcs, maybeScaled(&f, attrs["splits3D"]))
}

func parseUnpolar(attrs map[string]float64, pre, in, post *xi.Sum, ax xirho.Ax) {
	f := xi.Exp{Base: math.E}
	a := xi.Affine{}
	const sc = 1 / (2 * math.Pi)
	a.Ax.Eye().RotZ(math.Pi/2).Scale(-1, 1, 0)
	t := xi.Then{
		Funcs: xirho.FuncList{
			&a,
			&f,
		},
	}
	in.Funcs = append(in.Funcs, maybeScaled(&t, attrs["unpolar"]/(2*math.Pi)))
}

// maybeScaled returns f if v is 1 or a Then with f and Scale by v otherwise.
func maybeScaled(f xirho.F, v float64) xirho.F {
	if v == 1 {
		return f
	}
	if t, ok := f.(*xi.Then); ok {
		t.Funcs = append(t.Funcs, &xi.Scale{Amount: xirho.Real(v)})
		return t
	}
	return &xi.Then{
		Funcs: []xirho.F{
			f,
			&xi.Scale{Amount: xirho.Real(v)},
		},
	}
}
