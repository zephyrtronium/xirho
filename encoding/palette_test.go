package encoding_test

import (
	"image/color"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/zephyrtronium/xirho/encoding"
)

func TestPaletteRoundTrip(t *testing.T) {
	cases := []struct {
		name    string
		palette color.Palette
	}{
		{
			name: "single",
			palette: color.Palette{
				color.Transparent,
			},
		},
		{
			name: "hex",
			palette: color.Palette{
				color.RGBA{R: 0xff, A: 0xff},
				color.RGBA{R: 0xff, G: 0xff, A: 0xff},
				color.RGBA{G: 0xff, A: 0xff},
				color.RGBA{G: 0xff, B: 0xff, A: 0xff},
				color.RGBA{B: 0xff, A: 0xff},
				color.RGBA{R: 0xff, B: 0xff, A: 0xff},
			},
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			m := encoding.EncodePalette(c.palette)
			d, err := encoding.DecodePalette(m)
			if err != nil {
				t.Errorf("failed to decode palette: %v", err)
			}
			want := rgbaPalette(c.palette)
			got := rgbaPalette(d)
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("failed to round-trip (+got/-want):\n%s", diff)
			}
		})
	}
}

func rgbaPalette(p color.Palette) []color.RGBA64 {
	o := make([]color.RGBA64, len(p))
	for i, c := range p {
		r, g, b, a := c.RGBA()
		o[i] = color.RGBA64{
			R: uint16(r),
			G: uint16(g),
			B: uint16(b),
			A: uint16(a),
		}
	}
	return o
}
