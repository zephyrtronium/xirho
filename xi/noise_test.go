package xi_test

import "testing"

func TestNoiseAPI(t *testing.T) {
	ExpectAPI(t, nil, "noise")
}
