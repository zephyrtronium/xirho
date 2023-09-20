package xmath_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/zephyrtronium/xirho/xmath"
	"gonum.org/v1/gonum/mathext"
)

// TestGraphFill tests that the xmath PRNG can produce a random walk that
// follows every edge in fully connected graphs of various sizes. Passing this
// suggests that the generator should not systematically avoid any edge in the
// IFS "xaos" graph. Failure is indicated only by the test timing out. This is
// a randomized test; failures may be (but should not be) sporadic.
func TestGraphFill(t *testing.T) {
	rng := xmath.NewRNG()
	const alpha = 0.001
	for i := 2; i <= 128; i++ {
		p := graphFill(rng, i)
		if p < alpha {
			t.Logf("graph fill on %d nodes is not uniform at p=%f level (%f)", i, alpha, p)
		}
	}
	if testing.Short() {
		return
	}
	// graphFill is O(n²), assuming a uniform generator. Do these
	// longer tests concurrently.
	var wg sync.WaitGroup
	for i := 129; i <= 512; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			p := graphFill(rng, i)
			if p < alpha {
				t.Logf("graph fill on %d nodes is not uniform at p=%f level (%f)", i, alpha, p)
			}
		}(i)
	}
	wg.Wait()
}

func graphFill(rng xmath.RNG, degree int) float64 {
	edges := make([]int64, degree*degree)
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
	var n int64
	from := 0
	for k < degree*degree {
		to := 0
		x := rng.Uint64() & (1<<53 - 1)
		for next[to] < x {
			to++
		}
		if edges[from*degree+to] == 0 {
			k++
		}
		edges[from*degree+to]++
		n++
		from = to
	}
	// Walk extra steps to improve the sample.
	for i := n; i >= 0 || n < 10000; i-- {
		to := 0
		x := rng.Uint64() & (1<<53 - 1)
		for next[to] < x {
			to++
		}
		edges[from*degree+to]++
		n++
		from = to
	}
	// Calculate chi-squared statistic via Kahan sum.
	var x, c float64
	e := float64(n) / float64(degree*degree)
	for _, o := range edges {
		y := float64(o) - e
		y = y*y - c
		t := x + y
		c = t - x - y
		x = t
	}
	x /= e
	// Calculate p-value from chi-squared statistic.
	return mathext.GammaIncRegComp(float64(degree*degree-1)/2, x/2)
}

var doNotOptimize uint64

func BenchmarkUint64(b *testing.B) {
	rng := xmath.NewRNG()
	b.SetBytes(8)
	b.ResetTimer()
	var dnopt uint64
	for i := b.N; i != 0; i-- {
		dnopt = rng.Uint64()
	}
	doNotOptimize = dnopt
}
