package xmath_test

import (
	"fmt"
	"testing"

	"github.com/zephyrtronium/xirho/xmath"
)

// TestGraphFill tests that the xmath PRNG can produce a random walk that
// follows every edge in fully connected graphs of various sizes. Passing this
// suggests that the generator should not systematically avoid any edge in the
// IFS "xaos" graph. Failure is indicated only by the test timing out. This is
// a randomized test; failures may be (but should not be) sporadic.
func TestGraphFill(t *testing.T) {
	top := 512
	if testing.Short() {
		top = 128
	}
	rng := xmath.NewRNG()
	for i := 2; i <= top; i++ {
		t.Run(fmt.Sprintf("%d-nodes", i), func(t *testing.T) {
			graphFill(&rng, i)
		})
		t.Log("completed graph fill with", i, "nodes")
	}
}

func graphFill(rng *xmath.RNG, degree int) {
	edges := make([]bool, degree*degree)
	// Important that the edge choice algorithm mirror that used in xirho.
	const scale float64 = 1 << 53
	next := make([]uint64, degree)
	for i := range next {
		next[i] = uint64(float64(i+1) / float64(degree) * scale)
	}
	if next[len(next)-1] != 1<<53 {
		panic(fmt.Errorf("wrong maximum %064x, expected %064x", next[len(next)-1], uint64(1<<53)))
	}
	k := 0
	from := 0
	for k < degree*degree {
		to := 0
		x := rng.Uint64() & (1<<53 - 1)
		for next[to] < x {
			to++
		}
		if !edges[from*degree+to] {
			edges[from*degree+to] = true
			k++
		}
		from = to
	}
}
