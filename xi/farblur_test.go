package xi_test

import (
	"testing"

	"github.com/zephyrtronium/xirho/fapi"
)

func TestFarblurAPI(t *testing.T) {
	expect := map[string]fapi.Param{
		"origin": fapi.Vec3{},
		"affine": fapi.Affine{},
		"dist":   fapi.Real{},
	}
	ExpectAPI(t, expect, "farblur")
}
