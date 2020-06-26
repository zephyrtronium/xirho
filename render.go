package xirho

import (
	"context"
	"image/color"
	"sync"
	"sync/atomic"
	"time"
)

// R manages the rendering of a System onto a Hist.
type R struct {
	// These fields must be first on 32-bit platforms because they are updated
	// atomically.
	// n is the number of points calculated.
	n int64
	// q is the number of points plotted.
	q int64

	// Hist is the target histogram.
	Hist *Hist
	// System is the system to render.
	System System
	// Camera is the camera transform.
	Camera Ax
	// Palette is the colors used by the renderer.
	Palette []color.NRGBA64
	// Procs is the number of goroutines to use in iterating the system. If
	// Procs <= 0, then Render instead uses GOMAXPROCS goroutines.
	Procs int
	// N is the maximum number of iterations to perform. If N <= 0, then this
	// is not used as an exit condition.
	N int64
	// Q is the maximum number of times to plot, i.e. the maximum number of
	// iterations that produce points lying inside the histogram. If Q <= 0,
	// then this is not used as an exit condition.
	Q int64

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
	rng := newRNG()
	r.aspect = float64(r.Hist.rows) / float64(r.Hist.cols)
	ctx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	wg.Add(r.Procs)
	for i := 0; i < r.Procs; i++ {
		go func(rng RNG) {
			r.System.Iter(ctx, r, rng)
			wg.Done()
		}(rng)
		rng.Jump()
	}
	ticker := time.NewTicker(15 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			// If our context is cancelled, then so are the workers', but vet
			// complains, and an extra cancel doesn't hurt anything.
			cancel()
			wg.Wait()
			return
		case <-ticker.C:
			if (r.N > 0 && atomic.LoadInt64(&r.n) >= r.N) || (r.Q > 0 && atomic.LoadInt64(&r.q) >= r.Q) {
				cancel()
				wg.Wait()
				return
			}
		}
	}
}

// plot plots a point.
func (r *R) plot(p P) bool {
	atomic.AddInt64(&r.n, 1)
	if !p.IsValid() {
		return false
	}
	x, y, _ := Tx(&r.Camera, p.X, p.Y, p.Z) // ignore z
	var col, row int
	if r.aspect <= 1 {
		if x < -1 || x >= 1 || y < -r.aspect || y >= r.aspect {
			return false
		}
		col = int((x + 1) * 0.5 * float64(r.Hist.cols))
		row = int((y/r.aspect + 1) * 0.5 * float64(r.Hist.rows))
	} else {
		if x < -1/r.aspect || x >= 1/r.aspect || y < -1 || y >= 1 {
			return false
		}
		col = int((x*r.aspect + 1) * 0.5 * float64(r.Hist.cols))
		row = int((y + 1) * 0.5 * float64(r.Hist.rows))
	}
	c := int(p.C * float64(len(r.Palette)))
	if c >= len(r.Palette) {
		// Since p.C can be 1.0, c can be out of bounds.
		c = len(r.Palette)
	}
	color := r.Palette[c]
	r.Hist.Add(col, row, color)
	return true
}

// Iters returns the number of iterations the renderer has performed. It is
// safe to call this while the renderer is running.
func (r *R) Iters() int64 {
	return atomic.LoadInt64(&r.n)
}

// Hits returns the number of iterations the renderer has plotted. It is safe
// to call this while the renderer is running.
func (r *R) Hits() int64 {
	return atomic.LoadInt64(&r.q)
}
