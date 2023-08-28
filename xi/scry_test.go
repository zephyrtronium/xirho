package xi_test

import (
	"testing"

	"github.com/zephyrtronium/xirho/fapi"
)

func TestScryAPI(t *testing.T) {
	expect := map[string]fapi.Param{
		"radius": fapi.Real{},
	}
	ExpectAPI(t, expect, "scry")
}
