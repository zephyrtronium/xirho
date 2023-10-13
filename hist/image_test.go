package hist_test

import (
	"image/color"
	"testing"

	"github.com/zephyrtronium/xirho/hist"
)

func TestAtOOB(t *testing.T) {
	z := []int{0, 1, 10}
	for _, sz := range z {
		h := hist.New(hist.Size{W: sz, H: sz, OSA: 1})
		src := h.Image(hist.ToneMap{Brightness: 1, Gamma: 1}, 1, 1)
		// top edge
		for i := -1; i <= sz; i++ {
			c := src.At(i, -1)
			if c != (color.RGBA64{}) {
				t.Errorf("wrong color at %d,-1: want color.RGBA64{}, got %#v", i, c)
			}
		}
		// left edge
		for i := -1; i <= sz; i++ {
			c := src.At(-1, i)
			if c != (color.RGBA64{}) {
				t.Errorf("wrong color at -1,%d: want color.RGBA64{}, got %#v", i, c)
			}
		}
		// bottom edge
		for i := -1; i <= sz; i++ {
			c := src.At(i, sz)
			if c != (color.RGBA64{}) {
				t.Errorf("wrong color at %d,%d: want color.RGBA64{}, got %#v", i, sz, c)
			}
		}
		// right edge
		for i := -1; i <= sz; i++ {
			c := src.At(sz, i)
			if c != (color.RGBA64{}) {
				t.Errorf("wrong color at %d,%d: want color.RGBA64{}, got %#v", sz, i, c)
			}
		}
	}
}

func TestAt(t *testing.T) {
	clr := color.RGBA64{R: 0xffff, G: 0xffff, B: 0xffff, A: 0xffff}
	z := []int{1, 10}
	for _, width := range z {
		for _, height := range z {
			h := hist.New(hist.Size{W: width, H: height, OSA: 1})
			src := h.Image(hist.ToneMap{Brightness: 1, Contrast: 1, Gamma: 1}, 1, 1)
			for x := 0; x < width; x++ {
				for y := 0; y < height; y++ {
					c := src.At(x, y)
					if c != (color.RGBA64{}) {
						t.Errorf("wrong color at %d,%d: want color.RGBA64{}, got %#v", x, y, c)
					}
				}
			}
			src = nil
			for x := 0; x < width; x++ {
				for y := 0; y < height; y++ {
					h.Add(x, y, clr)
				}
			}
			src = h.Image(hist.ToneMap{Brightness: 1, Contrast: 4, Gamma: 1}, 1, 1)
			for x := 0; x < width; x++ {
				for y := 0; y < height; y++ {
					c := src.At(x, y)
					if c != clr {
						t.Errorf("wrong color at %d,%d: want %#v, got %#v", x, y, clr, c)
					}
				}
			}
		}
	}
}
