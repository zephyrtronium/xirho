package xirho

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"sync/atomic"
	"unsafe"
)

// Hist is a uniform two-dimensional histogram.
type Hist struct {
	// TODO: benchmark SoA versus AoS
	// bins is the histogram bins.
	counts []histBin
	// rows and cols are the histogram size.
	rows, cols int
	// lb is the logarithm of the maximum bin count.
	lb float64
	// exp is the reciprocal of the gamma factor applied to output pixels.
	exp float64
	// br is brightness, a 32.16 fixed-point multiplier for color channels
	// relative to alpha.
	br uint64
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

// NewHist allocates a histogram.
func NewHist(width, height int) *Hist {
	return &Hist{
		counts: make([]histBin, width*height),
		rows:   height,
		cols:   width,
		exp:    1,
		br:     1 << 16,
	}
}

// Add increments a histogram bucket by the given color. It is safe for
// multiple goroutines to call this concurrently.
func (h *Hist) Add(x, y int, c color.NRGBA64) {
	atomic.StoreInt32(&h.stat, 0)
	k := h.index(x, y)
	bin := &h.counts[k]
	atomic.AddUint64(&bin.r, uint64(c.R)*h.br>>16)
	atomic.AddUint64(&bin.g, uint64(c.G)*h.br>>16)
	atomic.AddUint64(&bin.b, uint64(c.B)*h.br>>16)
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
	// Find the maximum and compute lb.
	var m uint64
	for _, b := range h.counts {
		if n := atomic.LoadUint64(&b.n); n > m {
			m = n
		}
	}
	h.lb = math.Log10(float64(m)) - clscale
	atomic.StoreInt32(&h.stat, 1)
}

// SetGamma sets the gamma factor for output pixels.
func (h *Hist) SetGamma(gamma float64) {
	h.exp = 1 / gamma
}

// SetBrightness sets the brightness, which is a multiplier for the color
// channels (i.e. excluding alpha). Brightness is logarithmic to account for
// the logarithmic tone mapping; values larger than ln(2^32) (about 22.18) are
// limited to that value. NaN values result in no change.
func (h *Hist) SetBrightness(b float64) {
	b = math.Exp(b)
	if math.IsNaN(b) {
		return
	}
	const up = float64((1<<32)<<16-1) / (1 << 16)
	if b < 0 {
		b = 0
	} else if b > up {
		b = up
	}
	h.br = uint64(b * (1 << 16))
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
		return color.RGBA64{}
	}
	bin := &h.counts[h.index(x, y)]
	r := atomic.LoadUint64(&bin.r)
	g := atomic.LoadUint64(&bin.g)
	b := atomic.LoadUint64(&bin.b)
	n := atomic.LoadUint64(&bin.n)
	if n == 0 {
		return color.RGBA64{}
	}
	return color.RGBA64{
		R: cscale(r, h.lb),
		G: cscale(g, h.lb),
		B: cscale(b, h.lb),
		A: cscaleg(n, h.lb, h.exp),
	}
}

// cscale scales a bin count to a color component.
func cscale(n uint64, lb float64) uint16 {
	a := (math.Log10(float64(n)) - clscale) / lb // logarithmic tone mapping
	a *= 65536                                   // scale to uint16
	if a < 0 {                                   // clip to uint16
		a = 0
	} else if a > 65535 {
		a = 65535
	}
	return uint16(a)
}

// cscaleg scales a bin count to a color component with gamma.
func cscaleg(n uint64, lb, exp float64) uint16 {
	a := (math.Log10(float64(n)) - clscale) / lb // logarithmic tone mapping
	a = math.Pow(a, exp)                         // gamma
	a *= 65536                                   // scale to uint16
	if a < 0 {                                   // clip to uint16
		a = 0
	} else if a > 65535 {
		a = 65535
	}
	return uint16(a)
}
