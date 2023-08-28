package xi_test

import (
	"testing"

	"github.com/zephyrtronium/xirho/fapi"
)

func TestHeatAPI(t *testing.T) {
	expect := map[string]fapi.Param{
		"planar wave period": fapi.Real{},
		"planar wave phase":  fapi.Angle{},
		"planar wave amp":    fapi.Real{},
		"axial wave period":  fapi.Real{},
		"axial wave phase":   fapi.Angle{},
		"axial wave amp":     fapi.Real{},
		"radial wave period": fapi.Real{},
		"radial wave phase":  fapi.Angle{},
		"radial wave amp":    fapi.Real{},
	}
	ExpectAPI(t, expect, "heat")
}
