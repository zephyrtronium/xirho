package xirho

import (
	"github.com/zephyrtronium/xirho/hist"
)

// Hist is a uniform two-dimensional histogram.
type Hist = hist.Hist

// NewHist allocates a new histogram.
func NewHist(x, y int) *Hist {
	return hist.New(x, y)
}

// ToneMap holds the parameters describing conversion from histogram bin counts
// to color and alpha channels.
type ToneMap = hist.ToneMap
