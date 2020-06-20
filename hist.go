package xirho

import (
	"fmt"
	"image"
	"image/color"
	"math"
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
	// stat is true if Stat has been called since the last use of Add.
	stat bool
}

type histBin struct {
	// r, g, b are average red, green, and blue channels as 16.16 fixed point.
	r, g, b uint32
	// n is the bin count, which determines the alpha channel.
	n uint32
}

// Add increments a histogram bucket by the given color. The color's alpha
// channel is ignored.
func (h *Hist) Add(x, y int, c color.NRGBA64) {
	k := h.index(x, y)
	bin := h.counts[k] // TODO: benchmark this vs. taking pointer
	bin.n++
	if bin.n == 0 {
		// Overflow. Just return; the bin in memory is unmodified.
		return
	}
	h.stat = false
	if bin.n >= 1<<31-1 {
		// We've hit this bin two billion times. Its color is decided.
		// Note that this also prevents overflow in averaging calculations
		// below, since the intermediate product can never have its MSB set.
		// We still want to save the new count to preserve dynamic range.
		h.counts[k].n = bin.n
		return
	}
	// Average each component with the new channels using Knuth's online mean
	// algorithm: m_n = (m_{n-1}*(n-1) + x_n)/n. Each channel is a 16.16
	// fixed-point number; m_{n-1}*(n-1) multiplies a 16.16 by a 32.0, giving
	// a 48.16; adding can't overflow because of above, so still 48.16; and
	// dividing by a 32.0 brings us to a 16.48; shift back to 16.16.
	bin.r = uint32((uint64(bin.r)*uint64(bin.n-1) + uint64(c.R)<<16) / uint64(bin.n) >> 32)
	bin.g = uint32((uint64(bin.g)*uint64(bin.n-1) + uint64(c.G)<<16) / uint64(bin.n) >> 32)
	bin.b = uint32((uint64(bin.b)*uint64(bin.n-1) + uint64(c.B)<<16) / uint64(bin.n) >> 32)
	h.counts[k] = bin
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

// Stat computes information about the histogram needed to convert bins to
// colors. It must be called prior to any call to At.
func (h *Hist) Stat() {
	// Find the maximum and compute lb.
	var m uint32
	for _, b := range h.counts {
		if b.n > m {
			m = b.n
		}
	}
	h.lb = math.Log(float64(m))
	h.stat = true
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
	if !h.stat {
		panic(fmt.Errorf("xirho: no call to Stat, cannot compute pixels"))
	}
	if x < 0 || x > h.cols || y < 0 || y > h.rows {
		return color.RGBA64{}
	}
	bin := h.counts[h.index(x, y)]
	if bin.n == 0 {
		return color.RGBA64{}
	}
	alpha := math.Log(float64(bin.n)) / h.lb // logarithmic tone mapping
	alpha = math.Pow(alpha, h.exp)           // gamma
	alpha *= 65536                           // scale to uint16
	if alpha < 0 {                           // clip to uint16
		alpha = 0
	} else if alpha > 65535 {
		alpha = 65535
	}
	return color.RGBA64{
		R: uint16(bin.r >> 16),
		G: uint16(bin.g >> 16),
		B: uint16(bin.b >> 16),
		A: uint16(alpha),
	}
}
