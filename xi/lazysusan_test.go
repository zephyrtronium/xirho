package xi_test

import (
	"testing"

	"github.com/zephyrtronium/xirho/fapi"
)

func TestLazySusanAPI(t *testing.T) {
	expect := map[string]fapi.Param{
		"inside":  fapi.Affine{},
		"outside": fapi.Affine{},
		"center":  fapi.Vec3{},
		"radius":  fapi.Real{},
		"spread":  fapi.Real{},
		"twistZ":  fapi.Real{},
	}
	ExpectAPI(t, expect, "lazysusan")
}
