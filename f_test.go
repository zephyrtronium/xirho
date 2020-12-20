package xirho_test

import (
	"math"
	"testing"

	"github.com/zephyrtronium/xirho"
)

func TestPtIsValid(t *testing.T) {
	succeed := []xirho.Pt{
		{X: math.Nextafter(math.Inf(0), 0)},
		{Y: math.Nextafter(math.Inf(0), 0)},
		{Z: math.Nextafter(math.Inf(0), 0)},
		{C: 0},
		{C: math.Nextafter(0, 1)},
		{C: 1},
		{C: math.Nextafter(1, 0)},
		{X: math.Nextafter(math.Inf(0), 0), Y: math.Nextafter(math.Inf(-1), 0)},
	}
	fail := []xirho.Pt{
		{C: math.Nextafter(0, -1)},
		{C: math.Nextafter(1, 2)},
	}
	// Check non-finite values exhaustively.
	vs := []float64{0, math.NaN(), math.Inf(0), math.Inf(-1)}
	for _, x := range vs {
		for _, y := range vs {
			for _, z := range vs {
				for _, c := range vs {
					if x != 0 || y != 0 || z != 0 || c != 0 {
						fail = append(fail, xirho.Pt{X: x, Y: y, Z: z, C: c})
					}
				}
			}
		}
	}
	for _, p := range succeed {
		if !p.IsValid() {
			t.Error(p, "not reported valid")
		}
	}
	for _, p := range fail {
		if p.IsValid() {
			t.Error(p, "reported valid")
		}
	}
}
