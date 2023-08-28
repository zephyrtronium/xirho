package xi_test

import (
	"testing"

	"github.com/zephyrtronium/xirho/fapi"
)

func TestHoleAPI(t *testing.T) {
	expect := map[string]fapi.Param{
		"amount": fapi.Real{},
		"origin": fapi.Vec3{},
	}
	ExpectAPI(t, expect, "hole", "spherivoid")
}
