package xirho

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"io"
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
func (h *Hist) Add(x, y int, c color.RGBA64) {
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

// Cols returns the horizontal size of the histogram in bins.
func (h *Hist) Cols() int {
	return h.cols
}

// Rows returns the vertical size of the histogram in bins.
func (h *Hist) Rows() int {
	return h.rows
}

// Aspect returns the histogram's aspect ratio. If the histogram is empty, the
// result is 0.
func (h *Hist) Aspect() float64 {
	if h.Empty() {
		return 0
	}
	return float64(h.cols) / float64(h.rows)
}

// ToneMap holds the parameters describing conversion from histogram bin counts
// to color and alpha channels.
type ToneMap struct {
	// Brightness is a multiplier for the log-alpha channel.
	Brightness float64
	// Gamma is a nonlinear scaler that boosts low- and high-count
	// bins differently.
	Gamma float64
	// GammaMin is the minimum log-alpha value to which to apply gamma scaling
	// as a ratio versus the number of iterations per output pixel. Should
	// generally be between 0 and 1, inclusive.
	GammaMin float64
}

// Image creates a wrapper around the histogram that converts bins to pixels.
//
// The parameters are as follows. br is a multiplier for the alpha channel.
// gamma applies nonlinear brightness to brighten low-count bins. thresh
// specifies the minimum bin count to which to apply gamma as the ratio versus
// the number of iters per pixel. area is the area of the camera's visible
// plane in spatial units. iters is the total number of iterations run for the
// render. osa is the oversampling, the expected number of histogram bins per
// pixel per axis (although the image may be rescaled to any size).
//
// The histogram should not be modified while the wrapper is in use.
// Note that the wrapper holds a reference to the histogram's bins, so it
// should not be stored in any long-lived locations.
func (h *Hist) Image(tm ToneMap, area float64, iters int64, osa int) image.Image {
	if area <= 0 {
		area = 1
	}
	if osa <= 0 {
		osa = 1
	}
	// Convert to log early to avoid overflow and mitigate loss of precision.
	q := math.Log10(float64(len(h.counts))) - math.Log10(float64(iters))
	return &histImage{
		Hist: *h,
		b:    tm.Brightness * 0xffff,
		g:    1 / tm.Gamma,
		t:    tm.GammaMin,
		lqa:  4*math.Log10(float64(osa)) - math.Log10(area) + q - 2*clscale,
	}
}

// histImage wraps a histogram with brightness parameters for rendering.
type histImage struct {
	Hist
	b, g, t float64
	lqa     float64
}

func (h *histImage) ColorModel() color.Model {
	return color.NRGBA64Model
}

func (h *histImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, h.cols, h.rows)
}

// At returns the color of a pixel in the histogram. Note that this is a fairly
// expensive operation.
func (h *histImage) At(x, y int) color.Color {
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
	a := ascale(n, h.b, h.lqa)
	ag := gamma(a, h.g, h.t)
	as := cscale(ag)
	if itdoesntworkatall {
		fmt.Printf("  at(%d,%d) h.b=%f h.g=%f h.t=%f h.lqa=%f rgbn=%d/%d/%d/%d a=%f ag=%f as=%d\n", x, y, h.b, h.g, h.t, h.lqa, r, g, b, n, a, ag, as)
	}
	if as <= 0 {
		return color.NRGBA64{}
	}
	s := a / float64(n)
	rs := s * float64(r)
	gs := s * float64(g)
	bs := s * float64(b)
	p := color.NRGBA64{
		R: cscale(rs),
		G: cscale(gs),
		B: cscale(bs),
		A: as,
	}
	if itdoesntworkatall {
		fmt.Printf("at(%d,%d) p=%v s=%g rgb=%g/%g/%g\n", x, y, p, s, rs, gs, bs)
	}
	return p
}

const itdoesntworkatall = false

func ascale(n uint64, br, lb float64) float64 {
	a := br * (math.Log10(float64(n)) - lb)
	return a / float64(n)
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

// WriteTo dumps the histogram contents. The first two 8-byte words are the
// size in columns and rows, respectively. Then, each bin's red count is
// written in row-major order, then each blue, green, and alpha count. Each
// value written is an 8-byte little-endian integer.
//
// It is not safe to call WriteTo while the histogram may be plotted onto.
func (h *Hist) WriteTo(w io.Writer) (n int64, err error) {
	b := make([]byte, 16)
	binary.LittleEndian.PutUint64(b[0:8], uint64(h.cols))
	binary.LittleEndian.PutUint64(b[8:16], uint64(h.rows))
	k, err := w.Write(b[:16])
	n += int64(k)
	if err != nil {
		return n, err
	}
	b = b[:8]
	for _, bin := range h.counts {
		binary.LittleEndian.PutUint64(b, bin.r)
		k, err = w.Write(b)
		n += int64(k)
		if err != nil {
			return n, err
		}
	}
	for _, bin := range h.counts {
		binary.LittleEndian.PutUint64(b, bin.g)
		k, err = w.Write(b)
		n += int64(k)
		if err != nil {
			return n, err
		}
	}
	for _, bin := range h.counts {
		binary.LittleEndian.PutUint64(b, bin.b)
		k, err = w.Write(b)
		n += int64(k)
		if err != nil {
			return n, err
		}
	}
	for _, bin := range h.counts {
		binary.LittleEndian.PutUint64(b, bin.n)
		k, err = w.Write(b)
		n += int64(k)
		if err != nil {
			return n, err
		}
	}
	return n, nil
}
