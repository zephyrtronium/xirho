package xi_test

import "testing"

func TestGaussblurAPI(t *testing.T) {
	ExpectAPI(t, nil, "gaussblur", "gaussian_blur")
}
