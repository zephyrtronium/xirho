package xi_test

import (
	"testing"

	"github.com/zephyrtronium/xirho/fapi"
)

func TestExblurAPI(t *testing.T) {
	expect := map[string]fapi.Param{
		"strength": fapi.Real{},
		"dist":     fapi.Real{},
		"origin":   fapi.Vec3{},
	}
	ExpectAPI(t, expect, "exblur")
}
