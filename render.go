package xirho

import (
	"context"
	"image/color"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/image/draw"

	"github.com/zephyrtronium/xirho/hist"
	"github.com/zephyrtronium/xirho/xmath"
)

// Render manages the rendering of a System onto a Hist.
type Render struct {
	// Hist is the target histogram.
	Hist *hist.Hist
	// Camera is the camera transform.
	Camera xmath.Affine
	// Palette is the colors used by the renderer.
	Palette color.Palette
	// n is the number of points calculated.
	n atomic.Int64
	// q is the number of points plotted.
	q atomic.Int64
}

// Render renders a System onto a Hist. Calculation is performed by procs
// goroutines, or by GOMAXPROCS goroutines if procs <= 0. Render returns after
// the context closes and after all its renderer goroutines finish. It is safe
// to call Render multiple times in succession to continue using the same
// histogram.
func (r *Render) Render(ctx context.Context, system System, procs int) {
	rng := xmath.NewRNG()
	if procs <= 0 {
		procs = runtime.GOMAXPROCS(0)
	}
	system.Prep()
	var wg sync.WaitGroup
	wg.Add(procs)
	for i := 0; i < procs; i++ {
		go func(rng xmath.RNG) {
			system.Iter(ctx, r, rng)
			wg.Done()
		}(rng)
		rng.Jump()
	}
	wg.Wait()
}

// RenderAsync manages asynchronous rendering. It is intended to be used in a
// go statement. The renderer does not begin work until receiving a System and
// other render settings over the change channel.
//
// RenderAsync is designed to allow a user interface to change rendering
// parameters and receive plots safely, without needing to explicitly
// synchronize worker goroutines. Whenever it receives items over the change or
// plot channels, RenderAsync handles pausing and resuming workers as needed to
// prevent data races. It also attempts to group together multiple changes and
// plot requests to reduce unnecessary work.
//
// Once the context closes, RenderAsync stops its workers, closes the imgs
// channel, and returns. If needed, other goroutines may join on RenderAsync by
// waiting for imgs to close. Until imgs closes, it is not safe to modify any
// of the renderer's fields.
func (r *Render) RenderAsync(ctx context.Context, change <-chan ChangeRender, plot <-chan PlotOnto, imgs chan<- draw.Image) {
	rng := xmath.NewRNG()
	//lint:ignore SA4006 indeed unused, but it's simpler to write this way
	rctx, cancel := context.WithCancel(ctx)
	defer close(imgs)
	var (
		wg     sync.WaitGroup
		procs  int
		system System
		out    chan<- draw.Image
		img    draw.Image
	)
	for {
		select {
		case <-ctx.Done():
			cancel()
			return
		case c := <-change:
			cancel()
			c = drainchg(c, change)
			rctx, cancel = context.WithCancel(ctx)
			x, y, osa := r.Hist.Width(), r.Hist.Height(), r.Hist.OSA()
			reset := false
			wg.Wait() // TODO: select with ctx.Done
			if !c.System.Empty() {
				system = c.System
				reset = true
			}
			if c.Size.Bins() != 0 {
				x, y, osa = c.Size.W, c.Size.H, c.Size.OSA
				reset = true
			}
			if c.Camera != nil {
				r.Camera = *c.Camera
				reset = true
			}
			if len(c.Palette) != 0 {
				r.Palette = append(color.Palette{}, c.Palette...)
				reset = true
			}
			if reset {
				r.Reset(x, y, osa)
			}
			procs = c.Procs
			r.start(rctx, &wg, procs, system, &rng)
		case work := <-plot:
			cancel()
			work = drainplot(work, plot)
			rctx, cancel = context.WithCancel(ctx)
			wg.Wait() // TODO: select with ctx.Done
			src := r.Hist.Image(work.ToneMap, r.Area(), r.Iters())
			work.Scale.Scale(work.Image, work.Image.Bounds(), src, src.Bounds(), draw.Over, nil)
			img = work.Image
			out = imgs
			r.start(rctx, &wg, procs, system, &rng)
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
func (r *Render) start(ctx context.Context, wg *sync.WaitGroup, procs int, system System, rng *xmath.RNG) {
	if system.Empty() {
		return
	}
	system.Prep()
	wg.Add(procs)
	for i := 0; i < procs; i++ {
		go func(rng xmath.RNG) {
			system.Iter(ctx, r, rng)
			wg.Done()
		}(*rng)
		rng.Jump()
	}
}

// plot plots a point.
func (r *Render) plot(x, y, z float64, c color.RGBA64, aspect float64) bool {
	x, y, _ = xmath.Tx(&r.Camera, x, y, z) // ignore z
	var col, row int
	if aspect >= 1 {
		y *= aspect
	} else {
		x /= aspect
	}
	// negated condition to catch nans
	if !(x >= -1 && x < 1 && y >= -1 && y < 1) {
		return false
	}
	col = int((x + 1) * 0.5 * float64(r.Hist.Cols()))
	row = int((y + 1) * 0.5 * float64(r.Hist.Rows()))
	r.Hist.Add(col, row, c)
	return true
}

// Area calculates the size in Cartesian units of the area viewed through the
// camera.
func (r *Render) Area() float64 {
	d := r.Camera.ProjArea()
	a := r.Hist.Aspect()
	if a > 1 {
		a = 1 / a
	}
	return a / d
}

// Iters returns the number of iterations the renderer has performed. It is
// safe to call this while the renderer is running.
func (r *Render) Iters() int64 {
	return r.n.Load()
}

// Hits returns the number of iterations the renderer has plotted. It is safe
// to call this while the renderer is running.
func (r *Render) Hits() int64 {
	return r.q.Load()
}

// ResetCounts resets the values returned by Iters and Hits to zero.
func (r *Render) ResetCounts() {
	r.n.Store(0)
	r.q.Store(0)
}

// Reset resets the histogram and the iteration counts. It is not safe to call
// this while the renderer is running.
func (r *Render) Reset(width, height, osa int) {
	r.ResetCounts()
	r.Hist.Reset(hist.Size{W: width, H: height, OSA: osa})
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
	// ToneMap is the tone mapping parameters for this render.
	ToneMap hist.ToneMap
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
	Size hist.Size
	// Camera is the new camera transform to use, if non-nil.
	Camera *xmath.Affine
	// Palette is the new palette to use, if it has nonzero length. The palette
	// is copied into the renderer.
	Palette color.Palette
	// Procs is the new number of worker goroutines to use. If this is zero,
	// then the renderer does no work until receiving a nonzero Procs.
	Procs int
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
