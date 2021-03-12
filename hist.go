package xirho

import (
	"github.com/zephyrtronium/xirho/hist"
)

// Hist is a uniform two-dimensional histogram.
type Hist = hist.Hist

// HistSize is a histogram size including oversampling.
type HistSize = hist.Size

// NewHist allocates a new histogram. osa is the oversampling factor, the
// number of histogram bins per axis in a pixel.
func NewHist(sz HistSize) *Hist {
	return hist.New(sz)
}

// ToneMap holds the parameters describing conversion from histogram bin counts
// to color and alpha channels.
type ToneMap = hist.ToneMap
