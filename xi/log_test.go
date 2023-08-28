package xi_test

import (
	"testing"

	"github.com/zephyrtronium/xirho/fapi"
)

func TestLogAPI(t *testing.T) {
	expect := map[string]fapi.Param{
		"base": fapi.Complex{},
	}
	ExpectAPI(t, expect, "log")
}
