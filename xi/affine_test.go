package xi_test

import (
	"testing"

	"github.com/zephyrtronium/xirho/fapi"
)

func TestAffineAPI(t *testing.T) {
	expect := map[string]fapi.Param{
		"transform": fapi.Affine{},
	}
	ExpectAPI(t, expect, "affine", "linear", "linear3D")
}
