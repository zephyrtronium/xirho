package xirho

import (
	"context"
	"image/color"
	"runtime"
	"sync"
)

// R manages the rendering of a System onto a Hist.
type R struct {
	// Hist is the target histogram.
	Hist *Hist
	// System is the system to render.
	System System
	// Camera is the camera transform.
	Camera Ax
	// Palette is the colors used by the renderer.
	Palette []color.NRGBA64
	// Procs is the number of goroutines to use in iterating the system. If
	// Procs <= 0, then Render instead uses max(1, GOMAXPROCS-1) goroutines.
	Procs int
	// N is the maximum number of iterations to perform. If N <= 0, then this
	// is not used as an exit condition.
	N int64
	// Q is the maximum number of times to plot, i.e. the maximum number of
	// iterations that produce points lying inside the histogram. If Q <= 0,
	// then this is not used as an exit condition.
	Q int64

	// n is the number of points calculated.
	n int64
	// q is the number of points plotted.
	q int64
	// aspect is the aspect ratio of the histogram.
	aspect float64
}

// Render renders a System onto a Hist. It returns after the context closes or
// after processing N points, and after all its renderer goroutines finish. It
// is safe to call Render multiple times in succession to continue using the
// same histogram, typically with increased N and Q. It is not safe to call
// Render multiple times concurrently, nor to modify any of r's fields
// concurrently.
func (r *R) Render(ctx context.Context) {
	procs := r.getProcs()
	rng := newRNG()
	r.aspect = float64(r.Hist.cols) / float64(r.Hist.rows)
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan P, procs)
	var wg sync.WaitGroup
	wg.Add(procs)
	for i := 0; i < procs; i++ {
		go func(rng RNG) {
			r.System.Iter(ctx, ch, rng)
			wg.Done()
		}(rng)
		rng.Jump()
	}
	for {
		select {
		case <-ctx.Done():
			// If our context is cancelled, then so are the workers', but vet
			// complains, and an extra cancel doesn't hurt anything.
			cancel()
			wg.Wait()
			return
		case p := <-ch:
			r.plot(p)
			if (r.N > 0 && r.n >= r.N) || (r.Q > 0 && r.q >= r.Q) {
				cancel()
				wg.Wait()
				return
			}
		}
	}
}

// plot plots a point.
func (r *R) plot(p P) {
	r.n++
	if !p.IsValid() {
		return
	}
	x, y, _ := Tx(&r.Camera, p.X, p.Y, p.Z) // ignore z
	var col, row int
	if r.aspect <= 1 {
		if x < -1 || x >= 1 || y < -r.aspect || y >= r.aspect {
			return
		}
		col = int((x + 1) * 0.5 * float64(r.Hist.cols))
		row = int((y + r.aspect) * 0.5 * float64(r.Hist.rows))
	} else {
		if x < -1/r.aspect || x >= 1/r.aspect || y < -1 || y >= 1 {
			return
		}
		col = int((x + 1/r.aspect) * 0.5 * float64(r.Hist.cols))
		row = int((y + 1) * 0.5 * float64(r.Hist.rows))
	}
	c := int(p.C * float64(len(r.Palette)))
	if c >= len(r.Palette) {
		// Since p.C can be 1.0, c can be out of bounds.
		c = len(r.Palette)
	}
	color := r.Palette[c]
	r.Hist.Add(col, row, color)
	r.q++
}

// getProcs gets the actual number of goroutines to spawn in Render.
func (r *R) getProcs() int {
	procs := r.Procs
	if procs <= 0 {
		procs = runtime.GOMAXPROCS(0) - 1
		if procs <= 0 {
			procs = 1
		}
	}
	return procs
}
