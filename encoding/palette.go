package encoding

import (
	"bytes"
	"compress/lzw"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"image/color"
	"io"
)

// EncodePalette encodes a palette to the format used in xirho systems.
//
// The encoding may be lossy for some color models.
func EncodePalette(palette color.Palette) string {
	p := make([]byte, 2*4*len(palette))
	for i, c := range palette {
		r, g, b, a := c.RGBA()
		binary.BigEndian.PutUint16(p[2*(0*len(palette)+i):], uint16(a))
		binary.BigEndian.PutUint16(p[2*(1*len(palette)+i):], uint16(r))
		binary.BigEndian.PutUint16(p[2*(2*len(palette)+i):], uint16(g))
		binary.BigEndian.PutUint16(p[2*(3*len(palette)+i):], uint16(b))
	}
	var buf bytes.Buffer
	z := lzw.NewWriter(&buf, lzw.LSB, 8)
	if _, err := z.Write(p); err != nil {
		panic(err)
	}
	z.Close()
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

// DecodePalette decodes a palette from the format used in xirho systems.
func DecodePalette(s string) (color.Palette, error) {
	p, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("couldn't decode palette data as base64: %w", err)
	}
	z := lzw.NewReader(bytes.NewReader(p), lzw.LSB, 8).(*lzw.Reader)
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, z); err != nil {
		return nil, fmt.Errorf("couldn't decode palette data as lzw: %w", err)
	}
	p = buf.Bytes()
	palette := make(color.Palette, 0, len(p)/(2*4))
	for i := range palette {
		palette[i] = color.RGBA64{
			A: binary.BigEndian.Uint16(p[2*i:]),
			R: binary.BigEndian.Uint16(p[2*(i+len(palette)):]),
			G: binary.BigEndian.Uint16(p[2*(i+2*len(palette)):]),
			B: binary.BigEndian.Uint16(p[2*(i+3*len(palette)):]),
		}
	}
	return palette, nil
}
