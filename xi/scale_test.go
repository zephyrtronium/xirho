package xi_test

import (
	"testing"

	"github.com/zephyrtronium/xirho/fapi"
)

func TestScaleAPI(t *testing.T) {
	expect := map[string]fapi.Param{
		"amount": fapi.Real{},
	}
	ExpectAPI(t, expect, "scale")
}
