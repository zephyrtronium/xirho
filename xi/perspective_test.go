package xi_test

import (
	"testing"

	"github.com/zephyrtronium/xirho/fapi"
)

func TestPerspectiveAPI(t *testing.T) {
	expect := map[string]fapi.Param{
		"distance": fapi.Real{},
	}
	ExpectAPI(t, expect, "perspective")
}
