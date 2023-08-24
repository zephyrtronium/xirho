package xirho

import (
	"context"
	"fmt"
	"image/color"
	"sync/atomic"
	"unsafe"

	"github.com/zephyrtronium/xirho/xmath"
)

// System is a generalized iterated function system.
type System struct {
	// Nodes is the system's node list.
	Nodes []Node
	// Final is an additional function applied after each function, if it is
	// non-nil. The result from Final is used only for plotting; the input to
	// it is the same as the input to the next iteration's function.
	Final Func
}

// Node describes the properties of a single node within a system.
type Node struct {
	// Func is the function which transforms points.
	Func Func
	// Opacity scales the alpha channel of points plotted by the node. It
	// must be in the interval [0, 1].
	Opacity float64
	// Weight controls the proportion of iterations which map to this node.
	// It must be a finite, nonnegative number.
	Weight float64
	// Graph is the weights from this node to each other node in the
	// system. If the graph is shorter than the number of nodes in the
	// system, then the missing values are treated as being 1.
	Graph []float64

	// Label is the label for this node.
	Label string
}

// iterator manages the iterations of a System by a single goroutine.
type iterator struct {
	// n is the number of nodes in the system.
	n int
	// nodes is the system node list.
	nodes unsafe.Pointer // *[n]Func
	// final is the system final.
	final Func
	// nclrs is the number of colors in the palette.
	nclrs int
	// palette is the renderer's palette converted to RGBA.
	palette unsafe.Pointer // *[nclrs]color.RGBA64
	// rng is the iterator's source of randomness.
	rng xmath.RNG
	// op is the pre-multiplied opacities of each function in the system.
	op unsafe.Pointer // *[n]uint64
	// w is the pre-multiplied weights of each edge in the directed graph.
	w unsafe.Pointer // *[n][n]uint64
}

// nodeat gets the nth node in the system. This does not perform bounds checks.
func (it *iterator) nodeat(n int) Func {
	return *(*Func)(unsafe.Add(it.nodes, uintptr(n)*unsafe.Sizeof(Func(nil))))
}

// colorat gets the nth color in the palette. This does not perform bounds
// checks.
func (it *iterator) colorat(n int) color.RGBA64 {
	return *(*color.RGBA64)(unsafe.Add(it.palette, uintptr(n)*unsafe.Sizeof(color.RGBA64{})))
}

// opat gets the pre-multiplied opacity of the nth node in the system. This
// does not perform bounds checks.
func (it *iterator) opat(n int) uint64 {
	return *(*uint64)(unsafe.Add(it.op, uintptr(n)*unsafe.Sizeof(uint64(0))))
}

// wrow gets a pointer to the pre-multiplied edge weights from node n in the
// system. This does not perform bounds checks.
func (it *iterator) wrow(n int) *uint64 {
	return (*uint64)(unsafe.Add(it.w, uintptr(it.n*n)*unsafe.Sizeof(uint64(0))))
}

// nextw gets the next edge weight from an array returned by it.wrow.
func nextw(w *uint64) *uint64 {
	return (*uint64)(unsafe.Add(unsafe.Pointer(w), unsafe.Sizeof(uint64(0))))
}

// Prep calls the Prep method of each function in the system. It should be
// called once before any call to Iter.
func (s System) Prep() {
	for _, f := range s.Nodes {
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
func (s System) Iter(ctx context.Context, r *Render, rng xmath.RNG) {
	if err := s.Check(); err != nil {
		panic(err)
	}
	it := iterator{rng: rng}
	it.prep(s, r.Palette)
	aspect := r.Hist.Aspect()
	p, k := it.fuse() // p may not be valid!
	done := ctx.Done()
	var n, q int
	for {
		p = it.nodeat(k).Calc(p, &it.rng)
		n++
		// If a function has opacity α, that means we plot its points with
		// probability α. If we don't plot a point, then there's no reason
		// to apply the final, since that is only a nonlinear camera.
		if op := it.opat(k); op >= 1<<53 || (op > 0 && it.rng.Uint64()%(1<<53) < op) {
			fp := it.doFinal(p)
			if !fp.IsValid() {
				p, k = it.fuse()
				continue
			}
			i := int(fp.C * float64(it.nclrs))
			if i >= it.nclrs {
				// Since fp.C can be 1.0, i can be out of bounds.
				i = it.nclrs - 1
			}
			if r.plot(fp.X, fp.Y, fp.Z, it.colorat(i), aspect) {
				q++
			}
		}
		k = it.next(k)
		if n == 25000 {
			atomic.AddInt64(&r.n, int64(n))
			t := atomic.AddInt64(&r.q, int64(q))
			n, q = 0, 0
			// Some random-ish condition that's fast to check to decide
			// whether to re-fuse. 0x8 is the lowest bit set in 25000, so
			// this will be every other group if the hit ratio is 1.0.
			if t&0x8 == 0 {
				p, k = it.fuse()
			}
			select {
			case <-done:
				return
			default:
				// continue on
			}
		}
	}
}

// Check verifies that the system is properly configured: it contains at least
// one node, no opacities are outside [0, 1], and no weight is negative or
// non-finite. If any of these conditions is false, then the returned error
// describes the problem.
func (s System) Check() error {
	if s.Empty() {
		return fmt.Errorf("xirho: cannot render an empty system")
	}
	for i, f := range s.Nodes {
		if !xmath.IsFinite(f.Opacity) {
			return fmt.Errorf("xirho: non-finite opacity %v for func %d", f.Opacity, i)
		}
		if f.Opacity < 0 || f.Opacity > 1 {
			return fmt.Errorf("xirho: out of bounds opacity %v for func %d", f.Opacity, i)
		}
		if !xmath.IsFinite(f.Weight) {
			return fmt.Errorf("xirho: non-finite weight %v for func %d", f.Weight, i)
		}
		if f.Weight < 0 {
			return fmt.Errorf("xirho: negative weight %v for func %d", f.Weight, i)
		}
		for j, x := range f.Graph {
			if !xmath.IsFinite(x) {
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
	return len(s.Nodes) == 0
}

// doFinal applies the system's Final function to the point, if present.
func (it *iterator) doFinal(p Pt) Pt {
	if it.final != nil {
		p = it.final.Calc(p, &it.rng)
	}
	return p
}

// fuseLen is the number of iterations to perform before beginning to plot.
const fuseLen = 30

// fuse obtains initial conditions to plot points from the system.
func (it *iterator) fuse() (Pt, int) {
	p := Pt{
		X: it.rng.Uniform()*2 - 1,
		Y: it.rng.Uniform()*2 - 1,
		Z: it.rng.Uniform()*2 - 1,
		C: it.rng.Uniform(),
	}
	k := it.next(it.rng.Intn(it.n))
	for i := 0; i < fuseLen; i++ {
		p = it.nodeat(k).Calc(p, &it.rng)
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
	w := it.wrow(k)
	i := 0
	for v >= *w {
		w = nextw(w)
		i++
	}
	return i
}

// prep sets up the iterator's weighted directed graph, which controls the
// probability of each function being chosen based on the current one, and
// pre-multiplies brightnesses
func (it *iterator) prep(s System, p color.Palette) {
	it.final = s.Final
	var nodes []Func
	var w []uint64
	var op []uint64
	var palette []color.RGBA64
	switch l := len(s.Nodes); l {
	case 0:
		panic("xirho: iterator prep on empty system (unreachable)")
	case 1:
		it.n = 1
		nodes = []Func{s.Nodes[0].Func}
		w = []uint64{^uint64(0)} // even if the weight is 0
	default:
		nodes = make([]Func, len(s.Nodes))
		for i, f := range s.Nodes {
			nodes[i] = f.Func
		}
		w = make([]uint64, len(s.Nodes)*len(s.Nodes))
		// Let F denote the set of nodes in the system. Let f denote the
		// current node.
		// Each node in F has its own weight, and f has a weight to each node
		// in F (including itself). The probability of choosing g in F as the
		// next node is then the product of g's weight and the weight from f to
		// g, divided by the total weight of all nodes. Then we have a
		// probability distribution.
		// Numerical stability is important here, and this is called only once
		// per proc per render, so we can afford relatively expensive
		// algorithms like Kahan summation. We also scale to 2^53 instead of
		// 2^64-1 so that float64 doesn't lose precision over integers.
		// Furthermore, since we take 53-bit random numbers in it.next, scaling
		// 1.0 by 2^53 means the last element will always be greater than any
		// variate, which simplifies the loop.
		const scale float64 = 1 << 53
		wb := make([]float64, len(s.Nodes))
		for i, f := range s.Nodes {
			for j := copy(wb, f.Graph); j < len(s.Nodes); j++ {
				// Fill in missing values with 1.
				wb[j] = 1
			}
			for j, g := range s.Nodes {
				wb[j] *= g.Weight
			}
			sum := cumsum(wb)
			if sum == 0 {
				// 0 sum would give nan for every element. Avoid nan.
				w := w[i*len(s.Nodes) : (i+1)*len(s.Nodes)]
				for j := range w {
					w[j] = ^uint64(0)
				}
				continue
			}
			for j, x := range wb {
				w[i*len(s.Nodes)+j] = uint64(x / sum * scale)
			}
		}
	}
	// Calculate opacity probabilities. The idea here is essentially the same
	// as in fixed-point weights.
	op = make([]uint64, len(s.Nodes))
	for i, f := range s.Nodes {
		op[i] = uint64(f.Opacity * (1 << 53))
	}
	// Pre-multiply palette.
	if len(p) == 0 {
		p = color.Palette{color.RGBA64{}}
	}
	palette = make([]color.RGBA64, len(p))
	for i, c := range p {
		r, g, b, a := c.RGBA()
		palette[i] = color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a)}
	}
	it.n = len(s.Nodes)
	it.nodes = unsafe.Pointer(&nodes[0])
	it.w = unsafe.Pointer(&w[0])
	it.op = unsafe.Pointer(&op[0])
	it.nclrs = len(p)
	it.palette = unsafe.Pointer(&palette[0])
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
