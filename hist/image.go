package hist

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

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

// TODO: saturation

// Note that this implementation is incorrect except on black backgrounds. It
// is legacy code. The API will change eventually.

// clscale is log10(0xffff). Histogram counts are in [0, 0xffff], but the flame
// algorithm is based on colors in [0, 1]. Subtracting this from log counts
// performs the conversion.
const clscale = 4.816473303765249707784354368778591143369496252776245939965515119387352293655218

// lwp is log10(200). Adding lwp performs a whitepoint adjustment.
const lwp = 2.301029995663981195213738

// Image creates a wrapper around the histogram that converts bins to pixels.
//
// The parameters are as follows. br is a multiplier for the alpha channel.
// gamma applies nonlinear brightness to brighten low-count bins. thresh
// specifies the minimum bin count to which to apply gamma as the ratio versus
// the number of iters per pixel. area is the area of the camera's visible
// plane in spatial units. iters is the total number of iterations run for the
// render.
//
// The histogram should not be modified while the wrapper is in use.
// Note that the wrapper holds a reference to the histogram's bins, so it
// should not be stored in any long-lived locations.
func (h *Hist) Image(tm ToneMap, area float64, iters int64) image.Image {
	// Convert to log early to avoid overflow and mitigate loss of precision.
	q := math.Log10(float64(len(h.counts))) - math.Log10(float64(iters))
	return &histImage{
		Hist: h,
		b:    tm.Brightness,
		g:    1 / tm.Gamma,
		t:    tm.GammaMin,
		lqa:  lwp - clscale + 4*math.Log10(float64(h.osa)) - math.Log10(area) + q,
	}
}

// histImage wraps a histogram with brightness parameters for rendering.
type histImage struct {
	*Hist
	b, g, t float64
	lqa     float64
}

func (h *histImage) ColorModel() color.Model {
	return color.RGBA64Model
}

func (h *histImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, h.cols, h.rows)
}

// At returns the color of a pixel in the histogram. Note that this is a fairly
// expensive operation.
func (h *histImage) At(x, y int) color.Color {
	return h.RGBA64At(x, y)
}

func (h *histImage) RGBA64At(x, y int) color.RGBA64 {
	if x < 0 || x >= h.cols || y < 0 || y >= h.rows {
		return color.RGBA64{}
	}
	bin := h.at(x, y)
	r := bin.r.Load()
	g := bin.g.Load()
	b := bin.b.Load()
	n := bin.n.Load()
	if n == 0 {
		return color.RGBA64{}
	}
	a := ascale(n, h.b, h.lqa)
	ag := gamma(aces(a), h.g, h.t)
	as := cscale(ag)
	if itdoesntworkatall {
		fmt.Printf("  at(%d,%d) h.b=%f h.g=%f h.t=%f h.lqa=%f rgbn=%d/%d/%d/%d a=%f ag=%f as=%d\n", x, y, h.b, h.g, h.t, h.lqa, r, g, b, n, a, ag, as)
	}
	if as <= 0 {
		return color.RGBA64{}
	}
	s := a / float64(n)
	rs := s * float64(r)
	gs := s * float64(g)
	bs := s * float64(b)
	p := color.RGBA64{
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
	a := br * (math.Log10(float64(n)) + lb)
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

func aces(x float64) float64 {
	// Approximate ACES filmic tone mapping curve by Krzysztof Narkowicz.
	// https://knarkowicz.wordpress.com/2016/01/06/aces-filmic-tone-mapping-curve/
	const (
		a = 2.51
		b = 0.03
		c = 2.43
		d = 0.59
		e = 0.14
	)
	return (x * (a*x + b)) / (x*(c*x+d) + e)
}
