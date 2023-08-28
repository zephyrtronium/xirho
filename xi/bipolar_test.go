package xi_test

import (
	"testing"

	"github.com/zephyrtronium/xirho/fapi"
)

func TestBipolarAPI(t *testing.T) {
	expect := map[string]fapi.Param{
		"shift": fapi.Angle{},
	}
	ExpectAPI(t, expect, "bipolar")
}
