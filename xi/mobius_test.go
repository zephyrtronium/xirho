package xi_test

import (
	"testing"

	"github.com/zephyrtronium/xirho/fapi"
)

func TestMobiusAPI(t *testing.T) {
	expect := map[string]fapi.Param{
		"A.scalar":    fapi.Real{},
		"A.vector":    fapi.Vec3{},
		"B.scalar":    fapi.Real{},
		"B.vector":    fapi.Vec3{},
		"C.scalar":    fapi.Real{},
		"C.vector":    fapi.Vec3{},
		"D.scalar":    fapi.Real{},
		"D.vector":    fapi.Vec3{},
		"input blank": fapi.List{},
	}
	ExpectAPI(t, expect, "mobius", "mobiq")
}
