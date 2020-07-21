package xirho

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/zephyrtronium/xirho/xmath"
)

// RNG is the randomness source type.
type RNG = xmath.RNG

// System is a generalized iterated function system.
type System struct {
	// Funcs is the system's function list.
	Funcs []SysFunc
	// Final is an additional function applied after each function, if it is
	// non-nil. The result from Final is used only for plotting; the input to
	// it is the same as the input to the next iteration's function.
	Final F
}

// SysFunc describes the properties of a single function within a system.
type SysFunc struct {
	// Func is the function which transforms points.
	Func F
	// Opacity scales the alpha channel of points plotted by the function. It
	// must be in the interval [0, 1].
	Opacity float64
	// Weight controls the proportion of iterations which map to this function.
	// It must be a finite, nonnegative number.
	Weight float64
	// Graph is the weights from this function to each other function in the
	// system. If the graph is shorter than the number of functions in the
	// system, then the missing values are treated as being 1.
	Graph []float64

	// Label is the label for this function.
	Label string
}

// iterator manages the iterations of a System by a single goroutine.
type iterator struct {
	// funcs is the system function list.
	funcs []F
	// final is the system final.
	final F
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
		f.Func.Prep()
	}
	if s.Final != nil {
		s.Final.Prep()
	}
}

// Iter iterates the function system and plots points onto r. It continues
// iterating until the context's Done channel is closed. rng should be seeded
// to a distinct state for each call to this method. Iter panics if Check
// returns an error.
func (s System) Iter(ctx context.Context, r *R, rng RNG) {
	if err := s.Check(); err != nil {
		panic(err)
	}
	it := iterator{rng: rng}
	it.prep(s)
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
			if !p.IsValid() {
				p, k = it.fuse()
				continue
			}
			p = it.funcs[k].Calc(p, &it.rng)
			// If a function has opacity α, that means we plot its points with
			// probability α. If we don't plot a point, then there's no reason
			// to apply the final, since that is only a nonlinear camera.
			if it.op[k] >= 1<<53 || (it.op[k] > 0 && it.rng.Uint64()%(1<<53) < it.op[k]) {
				fp := it.doFinal(p)
				if r.plot(fp) {
					q++
				}
			}
			k = it.next(k)
			n++
		}
	}
}

// Check verifies that the system is properly configured: it contains at least
// one function, no opacities are outside [0, 1], and no weight is negative or
// non-finite. If any of these conditions is false, then the returned error
// describes the problem.
func (s System) Check() error {
	if s.Empty() {
		return fmt.Errorf("xirho: cannot render an empty system")
	}
	for i, f := range s.Funcs {
		if f.Opacity-f.Opacity != 0 {
			return fmt.Errorf("xirho: non-finite opacity %v for func %d", f.Opacity, i)
		}
		if f.Opacity < 0 || f.Opacity > 1 {
			return fmt.Errorf("xirho: out of bounds opacity %v for func %d", f.Opacity, i)
		}
		if f.Weight-f.Weight != 0 {
			return fmt.Errorf("xirho: non-finite weight %v for func %d", f.Weight, i)
		}
		if f.Weight < 0 {
			return fmt.Errorf("xirho: negative weight %v for func %d", f.Weight, i)
		}
		for j, x := range f.Graph {
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

// Empty returns whether the system contains no functions.
func (s System) Empty() bool {
	return len(s.Funcs) == 0
}

// doFinal applies the system's Final function to the point, if present.
func (it *iterator) doFinal(p P) P {
	if it.final != nil {
		p = it.final.Calc(p, &it.rng)
	}
	return p
}

// fuseLen is the number of iterations to perform before beginning to plot.
const fuseLen = 30

// fuse obtains initial conditions to plot points from the system.
func (it *iterator) fuse() (P, int) {
	p := P{
		X: it.rng.Uniform()*2 - 1,
		Y: it.rng.Uniform()*2 - 1,
		Z: it.rng.Uniform()*2 - 1,
		C: it.rng.Uniform(),
	}
	k := it.next(it.rng.Intn(len(it.funcs)))
	for i := 0; i < fuseLen; i++ {
		p = it.funcs[k].Calc(p, &it.rng)
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
	w := it.w[k*len(it.funcs) : (k+1)*len(it.funcs)]
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
func (it *iterator) prep(s System) {
	it.final = s.Final
	switch l := len(s.Funcs); l {
	case 0:
		it.funcs = nil
		it.w = nil
	case 1:
		it.funcs = []F{s.Funcs[0].Func}
		it.w = []uint64{^uint64(0)} // even if the weight is 0
	default:
		it.funcs = make([]F, len(s.Funcs))
		for i, f := range s.Funcs {
			it.funcs[i] = f.Func
		}
		it.w = make([]uint64, len(s.Funcs)*len(s.Funcs))
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
		wb := make([]float64, len(s.Funcs))
		for i, f := range s.Funcs {
			for j := copy(wb, f.Graph); j < len(s.Funcs); j++ {
				// Fill in missing values with 1.
				wb[j] = 1
			}
			for j, g := range s.Funcs {
				wb[j] *= g.Weight
			}
			sum := cumsum(wb)
			if sum == 0 {
				// 0 sum would give nan for every element. Avoid nan.
				w := it.w[i*len(s.Funcs) : (i+1)*len(s.Funcs)]
				for j := range w {
					w[j] = ^uint64(0)
				}
				continue
			}
			for j, x := range wb {
				it.w[i*len(s.Funcs)+j] = uint64(x / sum * scale)
			}
		}
	}
	// Calculate opacity probabilities. The idea here is essentially the same
	// as in fixed-point weights.
	it.op = make([]uint64, len(s.Funcs))
	for i, f := range s.Funcs {
		it.op[i] = uint64(f.Opacity * (1 << 53))
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
