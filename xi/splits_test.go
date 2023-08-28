package xi_test

import (
	"testing"

	"github.com/zephyrtronium/xirho/fapi"
)

func TestSplitsAPI(t *testing.T) {
	expect := map[string]fapi.Param{
		"x": fapi.Real{},
		"y": fapi.Real{},
		"z": fapi.Real{},
	}
	ExpectAPI(t, expect, "splits", "splits3D")
}
