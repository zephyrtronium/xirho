package xi_test

import (
	"testing"

	"github.com/zephyrtronium/xirho/fapi"
)

func TestJuliaNAPI(t *testing.T) {
	expect := map[string]fapi.Param{
		"power": fapi.Int{},
		"dist":  fapi.Real{},
	}
	ExpectAPI(t, expect, "julian")
}
