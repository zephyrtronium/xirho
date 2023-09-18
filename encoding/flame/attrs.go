package flame

import (
	"math"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xi"
	"github.com/zephyrtronium/xirho/xmath"
)

// Funcs maps Flame xform attribute names to functions which parse attributes
// into function instances. Parameters of functions should not included here,
// e.g. julian is an entry but julian_power is not; the latter is in
// KnownAttrs.
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
	"hemisphere":    parseHemisphere,
	"julia":         parseJulia,
	"julian":        parseJulian,
	"lazysusan":     parseLazySusan,
	"log":           parseLog,
	"mobius":        parseMobius,
	"mobiq":         parseMobiq,
	"noise":         parseNoise,
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

// KnownAttrs lists xform attributes which are not function names but should
// not be reported as unrecognized by the unmarshaler.
var KnownAttrs = map[string]bool{
	"weight":         true,
	"color":          true,
	"coefs":          true,
	"opacity":        true,
	"plotmode":       true,
	"pre_zscale":     true,
	"pre_ztranslate": true,
	"pre_rotate_x":   true,
	"pre_rotate_y":   true,

	"bipolar_shift":          true,
	"curl_c1":                true,
	"curl_c2":                true,
	"expo_real":              true,
	"expo_imaginary":         true,
	"julian_power":           true,
	"julian_dist":            true,
	"lazysusan_space":        true,
	"lazysusan_spin":         true,
	"lazysusan_twist":        true,
	"lazysusan_x":            true,
	"lazysusan_y":            true,
	"log_base":               true,
	"mobiq_at":               true,
	"mobiq_ax":               true,
	"mobiq_ay":               true,
	"mobiq_az":               true,
	"mobiq_bt":               true,
	"mobiq_bx":               true,
	"mobiq_by":               true,
	"mobiq_bz":               true,
	"mobiq_ct":               true,
	"mobiq_cx":               true,
	"mobiq_cy":               true,
	"mobiq_cz":               true,
	"mobiq_dt":               true,
	"mobiq_dx":               true,
	"mobiq_dy":               true,
	"mobiq_dz":               true,
	"post_heat_phi_amp":      true,
	"post_heat_phi_period":   true,
	"post_heat_phi_phase":    true,
	"post_heat_r_amp":        true,
	"post_heat_r_period":     true,
	"post_heat_r_phase":      true,
	"post_heat_theta_amp":    true,
	"post_heat_theta_period": true,
	"post_heat_theta_phase":  true,
	"splits_x":               true,
	"splits_y":               true,
	"splits3D_x":             true,
	"splits3D_y":             true,
	"splits3D_z":             true,

	// mobius variables
	"Re_A": true,
	"Im_A": true,
	"Re_B": true,
	"Im_B": true,
	"Re_C": true,
	"Im_C": true,
	"Re_D": true,
	"Im_D": true,
}

// Parser is a function which parses a xirho function from XML attributes.
// Parsers should add the parsed function to one of pre, in, or post, and they
// should not modify attrs.
type Parser func(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine)

func parseLinear(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	v := attrs["linear"] + attrs["linear3D"]
	ax.Eye().Scale(v, v, v)
	in.Funcs = append(in.Funcs, &xi.Affine{Ax: ax})
}

func parseBipolar(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	s := math.Mod(attrs["bipolar_shift"], 2)
	if s > 1 {
		s -= 2
	} else if s <= -1 {
		s += 2
	}
	in.Funcs = append(in.Funcs, maybeScaled(&xi.Bipolar{Shift: s * math.Pi}, attrs["bipolar"]))
}

func parseBlur(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	in.Funcs = append(in.Funcs, maybeScaled(xi.Blur{}, attrs["blur"]))
}

func parsePreblur(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	pre.Funcs = append(pre.Funcs, maybeScaled(xi.Gaussblur{}, attrs["pre_blur"]))
}

func parseBubble(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	in.Funcs = append(in.Funcs, maybeScaled(xi.Bubble{}, attrs["bubble"]))
}

func parseElliptic(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	in.Funcs = append(in.Funcs, maybeScaled(xi.CElliptic{}, attrs["elliptic"]))
}

func parseCurl(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	v := xi.Curl{
		C1: attrs["curl_c1"],
		C2: attrs["curl_c2"],
	}
	in.Funcs = append(in.Funcs, maybeScaled(&v, attrs["curl"]))
}

func parseCylinder(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	in.Funcs = append(in.Funcs, maybeScaled(xi.Cylinder{}, attrs["cylinder"]))
}

func parseDisc(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	in.Funcs = append(in.Funcs, maybeScaled(xi.Disc{}, attrs["disc"]))
}

func parseExp(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	in.Funcs = append(in.Funcs, maybeScaled(&xi.Exp{Base: math.E}, attrs["exp"]))
}

func parseExpo(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	z := complex(attrs["expo_real"], attrs["expo_imaginary"])
	in.Funcs = append(in.Funcs, maybeScaled(&xi.Exp{Base: z}, attrs["expo"]))
}

func parseFlatten(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	pre.Funcs = append(pre.Funcs, xi.Flatten{})
}

func parseFoci(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	in.Funcs = append(in.Funcs, xi.Foci{})
}

func parseGaussblur(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	t := xi.Then{
		Funcs: []xirho.Func{
			xi.Gaussblur{},
			xi.Flatten{},
		},
	}
	in.Funcs = append(in.Funcs, maybeScaled(&t, attrs["gaussian_blur"]))
}

func parsePostHeat(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	f := xi.Heat{
		ThetaT: attrs["post_heat_theta_period"],
		ThetaP: attrs["post_heat_theta_phase"],
		ThetaA: attrs["post_heat"] * attrs["post_heat_theta_amp"],
		PhiT:   attrs["post_heat_phi_period"],
		PhiP:   attrs["post_heat_phi_phase"],
		PhiA:   attrs["post_heat"] * attrs["post_heat_phi_amp"],
		RT:     attrs["post_heat_r_period"],
		RP:     attrs["post_heat_r_phase"],
		RA:     attrs["post_heat"] * attrs["post_heat_r_amp"],
	}
	post.Funcs = append(post.Funcs, &f)
}

func parseHemisphere(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	in.Funcs = append(in.Funcs, maybeScaled(xi.Hemisphere{}, attrs["hemisphere"]))
}

func parseJulia(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	f := xi.JuliaN{
		Power: 2,
		Dist:  1,
	}
	in.Funcs = append(in.Funcs, maybeScaled(&f, attrs["julia"]))
}

func parseJulian(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	f := xi.JuliaN{
		Power: int64(attrs["julian_power"]),
		Dist:  attrs["julian_dist"],
	}
	in.Funcs = append(in.Funcs, maybeScaled(&f, attrs["julian"]))
}

func parseLazySusan(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	f := xi.LazySusan{
		Inside:  xmath.Eye(),
		Outside: xmath.Eye(),
		Center:  [3]float64{attrs["lazysusan_x"], -attrs["lazysusan_y"], 0},
		Radius:  attrs["lazysusan"],
		Spread:  attrs["lazysusan_space"],
		TwistZ:  -attrs["lazysusan_twist"],
	}
	f.Outside.Scale(1, 1, 0)
	f.Inside.RotZ(attrs["lazysusan_spin"]).Scale(1, 1, 0)
	in.Funcs = append(in.Funcs, &f)
}

func parseLog(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	f := xi.Log{Base: complex(attrs["log_base"], 0)}
	if f.Base == 0 {
		f.Base = math.E
	}
	in.Funcs = append(in.Funcs, maybeScaled(&f, attrs["log"]))
}

func parseMobius(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	f := xi.Mobius{
		Ar:     attrs["Re_A"],
		Avec:   [3]float64{attrs["Im_A"], 0, 0},
		Br:     attrs["Re_B"],
		Bvec:   [3]float64{attrs["Im_B"], 0, 0},
		Cr:     attrs["Re_C"],
		Cvec:   [3]float64{attrs["Im_C"], 0, 0},
		Dr:     attrs["Re_D"],
		Dvec:   [3]float64{attrs["Im_D"], 0, 0},
		InZero: 3,
	}
	t := xi.Then{
		Funcs: []xirho.Func{
			xi.Flatten{},
			&f,
		},
	}
	in.Funcs = append(in.Funcs, maybeScaled(&t, attrs["mobius"]))
}

func parseMobiq(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	f := xi.Mobius{
		Ar:     attrs["mobiq_at"],
		Avec:   [3]float64{attrs["mobiq_ax"], attrs["mobiq_ay"], attrs["mobiq_az"]},
		Br:     attrs["mobiq_bt"],
		Bvec:   [3]float64{attrs["mobiq_bx"], attrs["mobiq_by"], attrs["mobiq_bz"]},
		Cr:     attrs["mobiq_ct"],
		Cvec:   [3]float64{attrs["mobiq_cx"], attrs["mobiq_cy"], attrs["mobiq_cz"]},
		Dr:     attrs["mobiq_dt"],
		Dvec:   [3]float64{attrs["mobiq_dx"], attrs["mobiq_dy"], attrs["mobiq_dz"]},
		InZero: 3,
	}
	in.Funcs = append(in.Funcs, maybeScaled(&f, attrs["mobiq"]))
}

func parseNoise(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	in.Funcs = append(in.Funcs, maybeScaled(xi.Noise{}, attrs["noise"]))
}

func parsePolar(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	in.Funcs = append(in.Funcs, maybeScaled(xi.Polar{}, attrs["polar"]))
}

func parseRod(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	in.Funcs = append(in.Funcs, &xi.Rod{Radius: attrs["rod"]})
}

func parseScry(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	f := xi.Scry{Radius: attrs["scry"]}
	t := xi.Then{
		Funcs: []xirho.Func{
			xi.Flatten{},
			&f,
		},
	}
	in.Funcs = append(in.Funcs, &t)
}

func parseSpherical(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	f := xi.Then{
		Funcs: []xirho.Func{
			xi.Flatten{},
			xi.Spherical{},
		},
	}
	in.Funcs = append(in.Funcs, maybeScaled(&f, attrs["spherical"]))
}

func parseSpherical3D(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	in.Funcs = append(in.Funcs, maybeScaled(xi.Spherical{}, attrs["spherical3D"]))
}

func parsePrespherical(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	f := xi.Then{
		Funcs: []xirho.Func{
			xi.Flatten{},
			xi.Spherical{},
		},
	}
	pre.Funcs = append(pre.Funcs, maybeScaled(&f, attrs["pre_spherical"]))
}

func parseSplits(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	f := xi.Splits{
		X: attrs["splits_x"],
		Y: attrs["splits_y"],
	}
	t := xi.Then{
		Funcs: []xirho.Func{
			xi.Flatten{},
			&f,
		},
	}
	in.Funcs = append(in.Funcs, maybeScaled(&t, attrs["splits"]))
}

func parseSplits3D(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	f := xi.Splits{
		X: attrs["splits3D_x"],
		Y: attrs["splits3D_y"],
		Z: attrs["splits3D_z"],
	}
	in.Funcs = append(in.Funcs, maybeScaled(&f, attrs["splits3D"]))
}

func parseUnpolar(attrs map[string]float64, pre, in, post *xi.Sum, ax xmath.Affine) {
	f := xi.Exp{Base: math.E}
	a := xi.Affine{}
	a.Ax.Eye().RotZ(math.Pi/2).Scale(-1, 1, 0)
	t := xi.Then{
		Funcs: []xirho.Func{
			&a,
			&f,
		},
	}
	in.Funcs = append(in.Funcs, maybeScaled(&t, attrs["unpolar"]/(2*math.Pi)))
}

// maybeScaled returns f if v is 1 or a Then with f and Scale by v otherwise.
func maybeScaled(f xirho.Func, v float64) xirho.Func {
	if v == 1 {
		return f
	}
	if t, ok := f.(*xi.Then); ok {
		t.Funcs = append(t.Funcs, &xi.Scale{Amount: v})
		return t
	}
	return &xi.Then{
		Funcs: []xirho.Func{
			f,
			&xi.Scale{Amount: v},
		},
	}
}
