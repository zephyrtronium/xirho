package xi_test

import (
	"testing"

	"github.com/zephyrtronium/xirho/fapi"
)

func TestSumAPI(t *testing.T) {
	expect := map[string]fapi.Param{
		"funcs": fapi.FuncList{},
		"color": fapi.Func{},
	}
	ExpectAPI(t, expect, "sum")
}
