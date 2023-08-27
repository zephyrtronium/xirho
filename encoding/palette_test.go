package encoding_test

import (
	"encoding/json"
	"image/color"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/encoding"
	"github.com/zephyrtronium/xirho/xi"
)

func TestPalette(t *testing.T) {
	// This test is to verify that the extracted helpers produce the same
	// output as the existing encoding, before we refactor to just use
	// EncodePalette everywhere.
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
			// Create a valid system to encode, marshal it, then unmarshal just
			// the palette and compare.
			s := encoding.System{
				System: xirho.System{
					Nodes: []xirho.Node{
						{
							Func: xi.Blur{},
						},
					},
				},
				Palette: c.palette,
			}
			m, err := s.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}
			var d struct {
				Palette string `json:"palette"`
			}
			if err := json.Unmarshal(m, &d); err != nil {
				t.Fatal(err)
			}
			{
				// Test that we encode the same way.
				got := encoding.EncodePalette(c.palette)
				if got != d.Palette {
					t.Errorf("wrong encoded palette: want %q, got %q", d.Palette, got)
				}
			}
			{
				// Test that we decode the same way.
				want := rgbaPalette(c.palette)
				p, err := encoding.DecodePalette(d.Palette)
				if err != nil {
					t.Errorf("couldn't decode palette: %v", err)
				}
				got := rgbaPalette(p)
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("wrong result (+got/-want):\n%s", diff)
				}
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
