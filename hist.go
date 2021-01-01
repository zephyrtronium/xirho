package xirho

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"runtime"
	"sync/atomic"
	"unsafe"
)

// Hist is a uniform two-dimensional histogram.
type Hist struct {
	// bins is the histogram bins.
	counts []histBin
	// rows and cols are the histogram size.
	rows, cols int
	// lb is the logarithm of the average bin count.
	lb float64
	// exp is the reciprocal of the gamma factor applied to output pixels.
	exp float64
	// tr is the gamma threshold. The gamma factor is not applied to bins
	// with less than this fraction of the max counts.
	tr float64
	// br is brightness, a multiplier for alpha relative to color channels.
	br float64
	// stat is 1 if Stat has been called since the last use of Add.
	stat int32
}

type histBin struct {
	// r, g, b are red, green, and blue channels.
	r, g, b uint64
	// n is the bin count, which determines the alpha channel.
	n uint64
}

// HistMem estimates the memory usage in bytes of a histogram of a given size.
func HistMem(width, height int) int {
	return int(unsafe.Sizeof(Hist{})) + width*height*int(unsafe.Sizeof(histBin{}))
}

// NewHist allocates a histogram. An alternative to this function is to create
// a zero Hist value and call Reset and SetBrightness to initialize it.
func NewHist(width, height int) *Hist {
	return &Hist{
		counts: make([]histBin, width*height),
		rows:   height,
		cols:   width,
		exp:    1,
		br:     1,
	}
}

// Reset reinitializes the histogram counts. If the given width and height are
// not equal to the histogram's current respective sizes, then the histogram is
// completely reallocated.
func (h *Hist) Reset(width, height int) {
	if width != h.cols || height != h.rows {
		// Histograms can be very large, so we want to ensure the current
		// counts are collected before we attempt to allocate new ones.
		h.counts = nil
		runtime.GC()
		h.counts = make([]histBin, width*height)
		h.rows = height
		h.cols = width
		return
	}
	for i := range h.counts {
		h.counts[i] = histBin{}
	}
}

// Empty returns true if the histogram has zero size.
func (h *Hist) Empty() bool {
	return h.rows == 0 || h.cols == 0
}

// Add increments a histogram bucket by the given color. It is safe for
// multiple goroutines to call this concurrently.
func (h *Hist) Add(x, y int, c color.NRGBA64) {
	atomic.StoreInt32(&h.stat, 0)
	k := h.index(x, y)
	bin := &h.counts[k]
	atomic.AddUint64(&bin.r, uint64(c.R))
	atomic.AddUint64(&bin.g, uint64(c.G))
	atomic.AddUint64(&bin.b, uint64(c.B))
	atomic.AddUint64(&bin.n, uint64(c.A))
}

// index converts a coordinate to an index into the histogram counts. Panics if
// out of bounds in either dimension.
func (h *Hist) index(x, y int) int {
	if x < 0 || x >= h.cols {
		panic(fmt.Errorf("xirho: x=%d out of bounds (hist has %d cols)", x, h.cols))
	}
	if y < 0 || y >= h.rows {
		panic(fmt.Errorf("xirho: y=%d out of bounds (hist has %d rows)", y, h.rows))
	}
	return y*h.cols + x
}

// clscale is log10(0xffff). Histogram counts are in [0, 0xffff], but the flame
// algorithm is based on colors in [0, 1]. Rescaling final results to that
// range noticeably improves the brightness and coloration of images.
const clscale = 4.8164733037652496

// prep computes information needed to convert bins to colors.
func (h *Hist) prep() {
	var m uint64
	for i := range h.counts {
		b := &h.counts[i]
		m += atomic.LoadUint64(&b.n)
	}
	h.lb = math.Log10(float64(m) / float64(h.rows*h.cols*65536))
	atomic.StoreInt32(&h.stat, 1)
}

// SetBrightness sets the brightness parameters for the rendered image. br
// controls the brightness of color channels relative to alpha. gamma applies
// nonlinear brightness to the alpha channel to balance low-count bins, but
// only applies to bins exceeding a relative count of tr.
func (h *Hist) SetBrightness(br, gamma, tr float64) {
	h.br = br * 200 / 65536
	h.exp = 1 / gamma
	h.tr = tr
}

// Brightness returns the last brightness parameters passed to SetBrightness.
func (h *Hist) Brightness() (br, gamma, tr float64) {
	br = h.br * 65536 / 200
	gamma = 1 / h.exp
	tr = h.tr
	return br, gamma, tr
}

// Aspect returns the histogram's aspect ratio. If the histogram is empty, the
// result is 0.
func (h *Hist) Aspect() float64 {
	if h.Empty() {
		return 0
	}
	return float64(h.cols) / float64(h.rows)
}

// --- image.Image implementation for easy resizing ---

// ColorModel returns the histogram's internal color model.
func (h *Hist) ColorModel() color.Model {
	return color.NRGBA64Model
}

// Bounds returns the bounds of the histogram.
func (h *Hist) Bounds() image.Rectangle {
	return image.Rect(0, 0, h.cols, h.rows)
}

// At returns the color of a pixel in the histogram. Note that this is a fairly
// expensive operation.
func (h *Hist) At(x, y int) color.Color {
	if atomic.LoadInt32(&h.stat) == 0 {
		h.prep()
	}
	if x < 0 || x > h.cols || y < 0 || y > h.rows {
		return color.NRGBA64{}
	}
	bin := &h.counts[h.index(x, y)]
	r := atomic.LoadUint64(&bin.r)
	g := atomic.LoadUint64(&bin.g)
	b := atomic.LoadUint64(&bin.b)
	n := atomic.LoadUint64(&bin.n)
	if n == 0 {
		return color.NRGBA64{}
	}
	a := ascale(n, h.br, h.lb)
	ag := gamma(a, h.exp, h.tr)
	as := cscale(ag)
	if as <= 0 {
		return color.NRGBA64{}
	}
	p := color.NRGBA64{
		R: cscale(ag / a / float64(as) * float64(r)),
		G: cscale(ag / a / float64(as) * float64(g)),
		B: cscale(ag / a / float64(as) * float64(b)),
		A: as,
	}
	return p
}

func ascale(n uint64, br, lb float64) float64 {
	a := br * (math.Log10(float64(n+1)) - lb)
	return a
}

func cscale(c float64) uint16 {
	c *= 65536
	switch {
	case c < 0:
		return 0
	case c >= 65535:
		return 65535
	default:
		return uint16(c)
	}
}

func gamma(a, exp, tr float64) float64 {
	if a >= tr {
		return math.Pow(a, exp)
	}
	p := a / tr
	return p*math.Pow(a, exp) + (1-p)*a*math.Pow(tr, exp-1)
}
