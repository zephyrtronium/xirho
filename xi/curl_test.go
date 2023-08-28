package xi_test

import (
	"testing"

	"github.com/zephyrtronium/xirho/fapi"
)

func TestCurlAPI(t *testing.T) {
	expect := map[string]fapi.Param{
		"c1": fapi.Real{},
		"c2": fapi.Real{},
	}
	ExpectAPI(t, expect, "curl")
}
