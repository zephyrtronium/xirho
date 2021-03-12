// Package hist defines an IFS histogram and related image processing routines.
package hist

import (
	"encoding/binary"
	"image/color"
	"io"
	"math/bits"
	"runtime"
	"sync/atomic"
	"unsafe"
)

// Hist is a uniform two-dimensional histogram.
type Hist struct {
	// arr is a pointer to first element of counts.
	arr unsafe.Pointer
	// rows and cols are the raw histogram size, after multiplication by osa.
	rows, cols int
	// osa is the oversampling factor.
	osa int
	// counts is the slice backed by arr. It is kept as a separate field for
	// the convenience of less performance-sensitive algorithms.
	counts []bin
}

type bin struct {
	// r, g, b are the red, green, and blue channels.
	r, g, b uint64
	// n is the bin count.
	n uint64
}

// Size describes a histogram size, including width, height, and oversampling.
type Size struct {
	// W and H are the number of columns and rows, respectively.
	W, H int
	// OSA is the oversampling factor, the number of bins per axis to be
	// resampled into each pixel. A value of 0 results in an empty histogram.
	OSA int
}

// Bins computes the number of bins in a histogram of this size. If Overflows
// returns true, then the result is zero.
func (sz Size) Bins() int {
	if sz.Overflows() {
		return 0
	}
	return sz.W * sz.H * sz.OSA * sz.OSA
}

// Mem estimates the memory usage in bytes of a histogram (and not accumulator)
// of this size. If Overflows returns true, then the result is zero.
func (sz Size) Mem() uintptr {
	if sz.Overflows() {
		return 0
	}
	return unsafe.Sizeof(Hist{}) + uintptr(sz.W*sz.H*sz.OSA*sz.OSA)*unsafe.Sizeof(bin{})
}

// Overflows returns true when the memory required by a histogram of this
// size would overflow the size of an integer. This is always true if either
// dimension or the oversampling factor is negative.
func (sz Size) Overflows() bool {
	nw := bits.Len64(uint64(sz.W))
	nh := bits.Len64(uint64(sz.H))
	no := bits.Len64(uint64(sz.OSA))
	return nw+nh+2*no >= bits.UintSize
}

// MemFor estimates the memory usage in bytes of a histogram (and not
// accumulator) of a given size. It assumes Overflows returns false for
// the given width and height.
func MemFor(width, height int) int {
	return int(unsafe.Sizeof(Hist{})) + width*height*int(unsafe.Sizeof(bin{}))
}

// Overflows returns true when the memory required by a histogram of the
// given size would overflow the size of an integer.
func Overflows(width, height int) bool {
	// Lazy approach: do the multiplication.
	// mask is the bits that are allowed to be set in the low word of the
	// multiplication result.
	mask := ^uint64(0) >> uint64(bits.Len64(uint64(unsafe.Sizeof(bin{}))))
	// Convert to int64 first to sign extend so we return true for negatives.
	hi, lo := bits.Mul64(uint64(int64(width)), uint64(int64(height)))
	return hi != 0 || lo&^mask != 0
}

// check panics if the size cannot be used to allocate a histogram, e.g. due to
// negative size or integer overflow.
func (sz Size) check() {
	if sz.W < 0 || sz.H < 0 || sz.OSA < 0 {
		panic("xirho: cannot make negative size histogram")
	}
	if sz.Overflows() {
		panic("xirho: histogram size overflows")
	}
}

// New allocates a new histogram.
func New(sz Size) *Hist {
	sz.check()
	h := Hist{
		cols:   sz.W * sz.OSA,
		rows:   sz.H * sz.OSA,
		osa:    sz.OSA,
		counts: make([]bin, sz.Bins()),
	}
	if len(h.counts) > 0 {
		h.arr = unsafe.Pointer(&h.counts[0])
	}
	return &h
}

// Reset reinitializes the histogram counts. If the given size requires a
// different number of bins from the current one, then the entire histogram is
// reallocated.
func (h *Hist) Reset(sz Size) {
	sz.check()
	// Might be different sizes with the same number of bins.
	h.cols, h.rows, h.osa = sz.W*sz.OSA, sz.H*sz.OSA, sz.OSA
	if sz.Bins() != len(h.counts) {
		// Histograms can be very large, so we want to ensure the current
		// counts are collected before we attempt to allocate new ones.
		h.arr = nil
		h.counts = nil
		runtime.GC()
		h.counts = make([]bin, sz.Bins())
		if len(h.counts) > 0 {
			h.arr = unsafe.Pointer(&h.counts[0])
		}
		return
	}
	for i := range h.counts {
		h.counts[i] = bin{}
	}
}

// IsEmpty returns true if the histogram has zero bins.
func (h *Hist) IsEmpty() bool {
	return len(h.counts) == 0
}

// checkBounds controls whether at checks histogram bounds. Only disable this
// if everything that calls at or Add is thoroughly tested!
const checkBounds = true

// at gets the bin at a given x and y. May panics if either dimension is
// out of bounds.
func (h *Hist) at(x, y int) *bin {
	if checkBounds {
		if x < 0 || y < 0 || x >= h.cols || y >= h.rows {
			panic("xirho: histogram position out of bounds")
		}
	}
	return (*bin)(unsafe.Pointer(uintptr(h.arr) + uintptr(y*h.cols+x)*unsafe.Sizeof(bin{})))
}

// Add increments a histogram bucket by a color. It is safe for multiple
// goroutines to call this concurrently.
func (h *Hist) Add(x, y int, c color.RGBA64) {
	bin := h.at(x, y)
	atomic.AddUint64(&bin.r, uint64(c.R))
	atomic.AddUint64(&bin.g, uint64(c.G))
	atomic.AddUint64(&bin.b, uint64(c.B))
	atomic.AddUint64(&bin.n, uint64(c.A))
}

// Width returns the horizontal size of the histogram in pixels.
func (h *Hist) Width() int {
	if h.osa == 0 {
		return 0
	}
	return h.cols / h.osa
}

// Height returns the vertical size of the histogram in pixels.
func (h *Hist) Height() int {
	if h.osa == 0 {
		return 0
	}
	return h.rows / h.osa
}

// OSA returns the histogram oversampling factor.
func (h *Hist) OSA() int {
	return h.osa
}

// Size returns the histogram's size.
func (h *Hist) Size() Size {
	return Size{
		W:   h.Width(),
		H:   h.Height(),
		OSA: h.OSA(),
	}
}

// Cols returns the horizontal size of the histogram in bins. Use Width to get
// the size in pixels.
func (h *Hist) Cols() int {
	return h.cols
}

// Rows returns the vertical size of the histogram in bins. Use Height to get
// the size in pixels.
func (h *Hist) Rows() int {
	return h.rows
}

// Aspect returns the histogram's aspect ratio. If the histogram is empty, the
// result is 0.
func (h *Hist) Aspect() float64 {
	if h.IsEmpty() {
		return 0
	}
	return float64(h.cols) / float64(h.rows)
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
