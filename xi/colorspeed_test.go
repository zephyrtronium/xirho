package xi_test

import (
	"testing"

	"github.com/zephyrtronium/xirho/fapi"
)

func TestColorSpeedAPI(t *testing.T) {
	expect := map[string]fapi.Param{
		"color": fapi.Real{},
		"speed": fapi.Real{},
	}
	ExpectAPI(t, expect, "colorspeed")
}
