package xi_test

import (
	"testing"

	"github.com/zephyrtronium/xirho/fapi"
)

func TestThenAPI(t *testing.T) {
	expect := map[string]fapi.Param{
		"funcs": fapi.FuncList{},
	}
	ExpectAPI(t, expect, "then")
}
