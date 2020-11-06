package xmath_test

import (
	"math"
	"testing"

	"github.com/zephyrtronium/xirho/xmath"
)

func TestIsFinite(t *testing.T) {
	cases := []struct {
		x  float64
		ok bool
	}{
		{0, true},
		{math.MaxFloat64, true},
		{-math.MaxFloat64, true},
		{math.Inf(0), false},
		{math.Inf(-1), false},
		{math.NaN(), false},
	}
	for _, c := range cases {
		if xmath.IsFinite(c.x) != c.ok {
			t.Errorf("wrong finitude for %v: wanted %t, got %t", c.x, c.ok, xmath.IsFinite(c.x))
		}
	}
}

func TestAngle(t *testing.T) {
	cases := []float64{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
		1e1, 2e2, 3e3, 4e4, 5e5, 6e6, 7e7, 8e8, 9e9,
		-1, -2, -3, -4, -5, -6, -7, -8, -9,
		-1e1, -2e2, -3e3, -4e4, -5e5, -6e6, -7e7, -8e8, -9e9,
		math.Pi, -math.Pi,
		math.Nextafter(math.Pi, 100), math.Nextafter(math.Pi, 0),
		math.Nextafter(-math.Pi, -100), math.Nextafter(-math.Pi, 0),
	}
	for _, c := range cases {
		x := xmath.Angle(c)
		if !(-math.Pi < x && x <= math.Pi) {
			t.Errorf("%.17g to angle %.17g outside (%.17g, %.17g]", c, x, -math.Pi, math.Pi)
		}
		cs, cc := math.Sincos(c)
		xs, xc := math.Sincos(x)
		// optimistic
		if math.Abs((xs-cs)/cs) > 1e-6 && math.Abs(xs-cs) > 1e-12 {
			t.Errorf("sines of original and wrapped angles %g and %g are far apart: %.17g != %.17g", c, x, cs, xs)
		}
		if math.Abs((xc-cc)/cc) > 1e-6 && math.Abs(xs-cs) > 1e-12 {
			t.Errorf("cosines of original and wrapped angles %g and %g are far apart: %.17g != %.17g", c, x, cc, xc)
		}
	}
}
