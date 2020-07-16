package xirho

import (
	"context"
	"image"
	"image/color"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/image/draw"

	"github.com/zephyrtronium/xirho/xmath"
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

	// Meta contains metadata about the fractal.
	Meta *Metadata
}

// Render renders a System onto a Hist. It returns after the context closes or
// after processing N points, and after all its renderer goroutines finish. It
// is safe to call Render multiple times in succession to continue using the
// same histogram, typically with increased N and Q. It is not safe to call
// Render multiple times concurrently, nor to modify any of r's fields
// concurrently.
func (r *R) Render(ctx context.Context) {
	rng := xmath.NewRNG()
	r.aspect = float64(r.Hist.rows) / float64(r.Hist.cols)
	ctx, cancel := context.WithCancel(ctx)
	procs := r.Procs
	if procs <= 0 {
		procs = runtime.GOMAXPROCS(0)
	}
	r.System.Prep()
	var wg sync.WaitGroup
	wg.Add(procs)
	for i := 0; i < procs; i++ {
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

// RenderAsync manages asynchronous rendering. It is intended to be used
// in a go statement.
//
// RenderAsync is designed to allow a user interface to change rendering
// parameters and receive plots safely, without needing to explicitly
// synchronize worker goroutines. Whenever it receives items over the change or
// plot channels, RenderAsync handles pausing and resuming workers as needed to
// prevent data races. It also attempts to group together
//
// If the renderer contains a non-empty system and histogram upon calling this
// method, RenderAsync immediately begins rendering using them, using Render's
// rules for the number of worker goroutines. Otherwise, it waits for a
// ChangeRender value to initialize the render parameters.
//
// Once the context closes, RenderAsync stops its workers, closes the imgs
// channel, and returns.
func (r *R) RenderAsync(ctx context.Context, change <-chan ChangeRender, plot <-chan PlotOnto, imgs chan<- draw.Image) {
	rng := xmath.NewRNG()
	rctx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	defer close(imgs)
	procs := 0
	if !r.System.Empty() && !r.Hist.Empty() {
		procs = r.Procs
		if procs <= 0 {
			procs = runtime.GOMAXPROCS(0)
		}
		r.start(rctx, &wg, procs, &rng)
	}
	var out chan<- draw.Image
	var img draw.Image
	for {
		select {
		case <-ctx.Done():
			cancel()
			return
		case c := <-change:
			cancel()
			c = drainchg(c, change)
			rctx, cancel = context.WithCancel(ctx)
			x, y := r.Hist.cols, r.Hist.rows
			wg.Wait() // TODO: select with ctx.Done
			if !c.System.Empty() {
				if !c.emptysz() {
					x, y = c.Size.X, c.Size.Y
				}
				r.Hist.Reset(x, y)
			} else if !c.emptysz() {
				x, y = c.Size.X, c.Size.Y
				r.Hist.Reset(x, y)
			}
			r.aspect = float64(y) / float64(x)
			procs = c.Procs
			r.start(rctx, &wg, procs, &rng)
		case work := <-plot:
			cancel()
			work = drainplot(work, plot)
			rctx, cancel = context.WithCancel(ctx)
			wg.Wait() // TODO: select with ctx.Done
			r.Hist.SetBrightness(work.Bright, work.Gamma, work.Thresh)
			work.Scale.Scale(work.Image, work.Image.Bounds(), r.Hist, r.Hist.Bounds(), draw.Over, nil)
			img = work.Image
			out = imgs
			r.start(rctx, &wg, procs, &rng)
		case out <- img:
			// out is normally nil, so this case will not be selected. It is
			// set to imgs when we have an image to send; once we send the
			// image, we can set out back to nil. This way, we automatically
			// consolidate a proportion of rapid draw requests.
			out = nil
		}
	}
}

// start starts worker goroutines with the given context.
func (r *R) start(ctx context.Context, wg *sync.WaitGroup, procs int, rng *xmath.RNG) {
	wg.Add(procs)
	for i := 0; i < procs; i++ {
		go func(rng xmath.RNG) {
			r.System.Iter(ctx, r, rng)
			wg.Done()
		}(*rng)
		rng.Jump()
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
		c = len(r.Palette) - 1
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

// drainchg pulls items from a ChangeRender channel until doing so would block,
// returning the last item obtained.
func drainchg(c ChangeRender, change <-chan ChangeRender) ChangeRender {
	runtime.Gosched()
	for {
		select {
		case c = <-change: // do nothing
		default:
			return c
		}
	}
}

// drainplot pulls items from a PlotOnto channel until doing so would block,
// returning the last item obtained.
func drainplot(work PlotOnto, plot <-chan PlotOnto) PlotOnto {
	runtime.Gosched()
	for {
		select {
		case work = <-plot: // do nothing
		default:
			return work
		}
	}
}

// PlotOnto is a work item for RenderAsync to plot onto.
type PlotOnto struct {
	// Image is the image to plot onto. The histogram is plotted using the Over
	// Porter-Duff operator.
	Image draw.Image
	// Scale is the resampling method to use to resample the histogram to the
	// size of Image.
	Scale draw.Scaler
	// Bright, Gamma, and Thresh correspond to the same parameters to
	// Hist.SetBrightness.
	Bright, Gamma, Thresh float64
}

// ChangeRender signals to RenderAsync to modify its system, histogram, or
// number of workers. RenderAsync can be paused without discarding render
// progress by sending this type's zero value.
type ChangeRender struct {
	// System is the new system to render. If the system is empty, then the
	// renderer continues using its previous non-empty system.
	System System
	// Size is the new histogram size to render. If this is the zero value,
	// then the histogram is neither resized nor reset. If this is equal to the
	// histogram's current size, then all plotting progress is cleared.
	Size image.Point
	// Procs is the new number of worker goroutines to use. If this is zero,
	// then the renderer does no work until receiving a nonzero Procs.
	Procs int
}

// emptysz returns true if the change's size is empty.
func (c ChangeRender) emptysz() bool {
	return c.Size.X == 0 || c.Size.Y == 0
}

// Metadata holds metadata about a fractal.
type Metadata struct {
	// Title is the name of the fractal.
	Title string `json:"title"`
	// Authors is the list of people who created the fractal.
	Authors []string `json:"authors"`
	// Date is the time the fractal was last modified.
	Date time.Time `json:"date"`
	// License is the license under which the fractal parameters are shared.
	// Typically this would be the title of the license, e.g. "CC4-BY-SA",
	// rather than the full license text.
	License string `json:"license"`
}
