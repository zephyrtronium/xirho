package xirho

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/zephyrtronium/crazy"
)

// RNG is the randomness source type.
type RNG = crazy.Xoshiro

// MaxFuncs is the maximum number of unique functions that a system may hold.
const MaxFuncs = 65536

// System is a generalized iterated function system.
type System struct {
	// Funcs is the system's function list.
	Funcs []F
	// Final is an additional function applied after each function if non-nil.
	Final F
	// Opacity scales the alpha channel of points plotted by each function. It
	// must be the same length as Funcs, and each element must be in the
	// interval [0, 1].
	Opacity []float64
	// Weights controls the proportion of iterations which map to each func. It
	// must be the same length as Funcs, and each element must be a
	// finite, nonnegative number.
	Weights []float64
	// Graph is the weights from the row function to each column function. It
	// must be of size len(Funcs) × len(Funcs), and each element must be a
	// finite, nonnegative number.
	Graph [][]float64

	// Labels gives the labels for each non-final function in the system, if it
	// is not nil.
	Labels []string
}

// iterator manages the iterations of a System by a single goroutine.
type iterator struct {
	System
	// rng is the iterator's source of randomness.
	rng RNG
	// op is the pre-multiplied opacities of each function in the system.
	op []uint64
	// w is the pre-multiplied weights of each edge in the directed graph.
	w []uint64
}

// Prep calls the Prep method of each function in the system. It should be
// called once before any call to Iter.
func (s System) Prep() {
	for _, f := range s.Funcs {
		f.Prep()
	}
	if s.Final != nil {
		s.Final.Prep()
	}
}

// Iter iterates the function system and sends output points over results. It
// continues iterating until the context's Done channel is closed. rng should
// be seeded to a distinct state for each call to this method. Iter panics if
// Check returns an error.
func (s System) Iter(ctx context.Context, r *R, rng RNG) {
	if err := s.Check(); err != nil {
		panic(err)
	}
	it := iterator{System: s, rng: rng}
	it.prep()
	p, k := it.fuse() // p may not be valid!
	done := ctx.Done()
	var n, q int64
	for {
		select {
		case <-done:
			atomic.AddInt64(&r.n, n)
			atomic.AddInt64(&r.q, q)
			return
		default:
			n++
			if n == 0x1000 {
				atomic.AddInt64(&r.n, n)
				n = 0
			}
			if !p.IsValid() {
				p, k = it.fuse()
				continue
			}
			fp := it.final(p)
			if r.plot(fp, it.op[k]) {
				q++
				if q == 0x1000 {
					atomic.AddInt64(&r.q, q)
					q = 0
				}
			}
			p = it.Funcs[k].Calc(p, &it.rng)
			k = it.next(k)
		}
	}
}

// Check verifies that the system is properly configured: it has as many
// opacities and weights as functions, the directed graph links to every
// function, no opacities are outside [0, 1], and neither the weights nor the
// directed graph contain a negative or non-finite element. If any of these
// conditions is false, then the returned error describes the problem.
func (s System) Check() error {
	if len(s.Funcs) != len(s.Opacity) {
		return fmt.Errorf("xirho: size mismatch, have %d funcs and %d opacities", len(s.Funcs), len(s.Opacity))
	}
	for i, x := range s.Opacity {
		if x-x != 0 {
			return fmt.Errorf("xirho: non-finite opacity %v for func %d", x, i)
		}
		if x < 0 || x > 1 {
			return fmt.Errorf("xirho: out of bounds opacity %v for func %d", x, i)
		}
	}
	if len(s.Funcs) != len(s.Weights) {
		return fmt.Errorf("xirho: size mismatch, have %d funcs and %d weights", len(s.Funcs), len(s.Weights))
	}
	for i, x := range s.Weights {
		if x-x != 0 {
			return fmt.Errorf("xirho: non-finite weight %v for func %d", x, i)
		}
		if x < 0 {
			return fmt.Errorf("xirho: negative weight %v for func %d", x, i)
		}
	}
	for i, g := range s.Graph {
		if len(s.Funcs) != len(g) {
			return fmt.Errorf("xirho: size mismatch, have %d funcs but graph node %d has %d weights", len(s.Funcs), i, len(g))
		}
		for j, x := range g {
			if x-x != 0 {
				return fmt.Errorf("xirho: non-finite weight %v for func %d to %d", x, i, j)
			}
			if x < 0 {
				return fmt.Errorf("xirho: negative weight %v for func %d to %d", x, i, j)
			}
		}
	}
	return nil
}

// final applies the system's Final function to the point, if present.
func (it *iterator) final(p P) P {
	if it.Final != nil {
		p = it.Final.Calc(p, &it.rng)
	}
	return p
}

// fuseLen is the number of iterations to perform before beginning to plot.
const fuseLen = 30

// fuse obtains initial conditions to plot points from the system.
func (it *iterator) fuse() (P, int) {
	d := crazy.Uniform{Source: &it.rng, Low: -1, High: 1}
	p := P{
		X: d.Next(),
		Y: d.Next(),
		Z: d.Next(),
		C: crazy.Uniform0_1{Source: &it.rng}.Next(),
	}
	k := it.next(crazy.RNG{Source: &it.rng}.Intn(len(it.Funcs)))
	for i := 0; i < fuseLen+1; i++ {
		p = it.Funcs[k].Calc(p, &it.rng)
		if !p.IsValid() {
			break
		}
		k = it.next(k)
	}
	return p, k
}

// next obtains the next function to use from the current one.
func (it *iterator) next(k int) int {
	v := it.rng.Uint64() & (1<<53 - 1)
	w := it.w[k*len(it.Funcs) : (k+1)*len(it.Funcs)]
	for i, x := range w {
		if v < x {
			return i
		}
	}
	panic("unreachable")
}

// prep sets up the iterator's weighted directed graph, which controls the
// probability of each function being chosen based on the current one, and
// pre-multiplies brightnesses
func (it *iterator) prep() {
	switch l := len(it.Funcs); l {
	case 0:
		it.w = nil
	case 1:
		it.w = []uint64{^uint64(0)} // even if the weight is 0
	default:
		it.w = make([]uint64, len(it.Funcs)*len(it.Funcs))
		// Let F denote the set of functions in the system. Let f denote the
		// current function.
		// Each function in F has its own weight, and f has a weight to each
		// function in F (including itself). The probability of choosing g in F
		// as the next function is then the product of g's weight and the
		// weight from f to g, divided by the total weight of all functions.
		// Then we have a probability distribution.
		// Numerical stability is important here, and this is called only once
		// per proc per render, so we can afford relatively expensive
		// algorithms like Kahan summation. We also scale to 2^53 instead of
		// 2^64-1 so that float64 doesn't lose precision over integers.
		// Furthermore, since we take 53-bit random numbers in it.next, scaling
		// 1.0 by 2^53 means the last element will always be greater than any
		// variate, which simplifies the loop.
		const scale float64 = 1 << 53
		wb := make([]float64, len(it.Funcs))
		for i, g := range it.Graph {
			copy(wb, g)
			for j, x := range it.Weights {
				wb[j] *= x
			}
			sum := cumsum(wb)
			if sum == 0 {
				// 0 sum would give nan for every element. Avoid nan.
				w := it.w[i*len(it.Funcs) : (i+1)*len(it.Funcs)]
				for j := range w {
					w[j] = ^uint64(0)
				}
				continue
			}
			for j, x := range wb {
				it.w[i*len(it.Funcs)+j] = uint64(x / sum * scale)
			}
		}
	}
	// Calculate opacity multipliers. Opacity is a 16.0 by 1.47 fixed-point
	// multiplication.
	it.op = make([]uint64, len(it.Funcs))
	for i, x := range it.Opacity {
		it.op[i] = uint64(x * (1 << 47))
	}
}

// cumsum computes the cumulative sum of float64s without loss of precision
// and returns the sum.
func cumsum(f []float64) float64 {
	var sum, c float64
	for i, x := range f {
		y := x - c
		f[i] = sum + y
		c = f[i] - sum - y
		sum = f[i]
	}
	return sum
}

// newRNG creates a new seeded RNG instance.
func newRNG() RNG {
	rng := RNG{}
	crazy.CryptoSeeded(&rng, 8)
	return rng
}
