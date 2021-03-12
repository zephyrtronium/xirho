// Package encoding implements marshaling and unmarshaling function systems.
package encoding

import (
	"bytes"
	"compress/lzw"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"image"
	"image/color"
	"io"
	"strconv"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xmath"
)

// System holds a xirho system and its rendering parameters for marshaling or
// unmarshaling.
type System struct {
	System  xirho.System
	ToneMap xirho.ToneMap
	Aspect  float64
	Camera  xirho.Affine
	BG      color.NRGBA64
	Palette color.Palette

	Meta *xirho.Metadata

	// Unrecognized is the list of unrecognized function type names following
	// unmarshaling a system. It may contain duplicates.
	Unrecognized []string
	// Err contains any error that occurred while decoding this system.
	Err error
}

// Wrap wraps a xirho system, renderer, tone mapping, and optionally a
// background color and metadata into a serializable system.
func Wrap(system xirho.System, r *xirho.Render, tm xirho.ToneMap, bg *color.NRGBA64, meta *xirho.Metadata) *System {
	s := System{
		System:  system,
		ToneMap: tm,
		Aspect:  r.Hist.Aspect(),
		Camera:  r.Camera,
		Palette: r.Palette,
		Meta:    meta,
	}
	if bg != nil {
		s.BG = *bg
	}
	return &s
}

// Render creates a new renderer with a histogram scaled to the given size,
// preserving the encoded aspect ratio.
func (s *System) Render(sz image.Point, osa int) *xirho.Render {
	w, h := xmath.Fit(sz.X, sz.Y, s.Aspect)
	return &xirho.Render{
		Hist:    xirho.NewHist(xirho.HistSize{W: w, H: h, OSA: osa}),
		Camera:  s.Camera,
		Palette: s.Palette,
	}
}

// MarshalJSON marshals the system as a JSON object. If the xirho system in s
// produces an error from Check(), then that error is returned instead.
func (s *System) MarshalJSON() ([]byte, error) {
	system := s.System
	// TODO: wrap errors
	if err := system.Check(); err != nil {
		return nil, err
	}
	m := marshaler{
		Funcs:  make([]*funcm, len(system.Nodes)),
		Camera: s.Camera,
		Bright: s.ToneMap.Brightness,
		Gamma:  s.ToneMap.Gamma,
		Thresh: s.ToneMap.GammaMin,
		Aspect: s.Aspect,
		Meta:   s.Meta,
	}
	for i, f := range system.Nodes {
		e, err := newFuncm(f.Func)
		e.Opacity = f.Opacity
		e.Weight = f.Weight
		e.Graph = f.Graph
		e.Label = f.Label
		if err != nil {
			return nil, err
		}
		m.Funcs[i] = e
	}
	if system.Final != nil {
		e, err := newFuncm(system.Final)
		if err != nil {
			return nil, err
		}
		m.Final = e
	}
	if s.BG != (color.NRGBA64{}) {
		m.BG = (*bgcolor)(&s.BG)
	}
	palette := make([]byte, 2*4*len(s.Palette))
	for i, c := range s.Palette {
		r, g, b, a := c.RGBA()
		binary.BigEndian.PutUint16(palette[2*(0*len(s.Palette)+i):], uint16(a))
		binary.BigEndian.PutUint16(palette[2*(1*len(s.Palette)+i):], uint16(r))
		binary.BigEndian.PutUint16(palette[2*(2*len(s.Palette)+i):], uint16(g))
		binary.BigEndian.PutUint16(palette[2*(3*len(s.Palette)+i):], uint16(b))
	}
	var buf bytes.Buffer
	z := lzw.NewWriter(&buf, lzw.LSB, 8)
	if _, err := z.Write(palette); err != nil {
		panic(err)
	}
	z.Close()
	m.Palette = buf.Bytes()
	return json.Marshal(m)
}

// UnmarshalJSON decodes a system from a JSON object.
func (s *System) UnmarshalJSON(b []byte) (err error) {
	*s = System{} // Ensure all fields are cleared on error.
	// TODO: wrap errors
	defer func() {
		s.Err = err
	}()
	m := marshaler{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return err
	}
	s.Camera = m.Camera
	s.System = xirho.System{
		Nodes: make([]xirho.Node, len(m.Funcs)),
	}
	s.Aspect = m.Aspect
	for i, f := range m.Funcs {
		s.System.Nodes[i].Func, err = unf(f)
		if err != nil {
			return err
		}
		s.System.Nodes[i].Opacity = f.Opacity
		s.System.Nodes[i].Weight = f.Weight
		s.System.Nodes[i].Graph = f.Graph
		s.System.Nodes[i].Label = f.Label
	}
	if m.Final != nil {
		s.System.Final, err = unf(m.Final)
		if err != nil {
			return err
		}
	}
	s.ToneMap = xirho.ToneMap{Brightness: m.Bright, Gamma: m.Gamma, GammaMin: m.Thresh}
	if m.BG != nil {
		s.BG = color.NRGBA64(*m.BG)
	}
	z := lzw.NewReader(bytes.NewReader(m.Palette), lzw.LSB, 8)
	var buf bytes.Buffer
	if _, err = io.Copy(&buf, z); err != nil {
		return err
	}
	palette := buf.Bytes()
	s.Palette = make(color.Palette, len(palette)/(2*4))
	for i := range s.Palette {
		s.Palette[i] = color.NRGBA64{
			A: binary.BigEndian.Uint16(palette[2*i:]),
			R: binary.BigEndian.Uint16(palette[2*(i+len(s.Palette)):]),
			G: binary.BigEndian.Uint16(palette[2*(i+2*len(s.Palette)):]),
			B: binary.BigEndian.Uint16(palette[2*(i+3*len(s.Palette)):]),
		}
	}
	s.Meta = m.Meta
	return nil
}

// marshaler controls the marshaling and unmarshaling of a full function system
// and renderer.
type marshaler struct {
	Meta *xirho.Metadata `json:"meta,omitempty"`
	// system params
	Funcs []*funcm `json:"funcs"`
	Final *funcm   `json:"final,omitempty"`
	// renderer params
	Aspect float64      `json:"aspect"`
	Camera xirho.Affine `json:"camera"`
	// brightness params
	Bright float64 `json:"bright"`
	Gamma  float64 `json:"gamma"`
	Thresh float64 `json:"thresh"`
	// bg color, if any
	BG *bgcolor `json:"bg,omitempty"`
	// Palette is formed by concatenating each channel of the NRGBA64 palette
	// in ARGB order as big-endian, then LZW-encoding the result.
	Palette []byte `json:"palette"`
}

// bgcolor serializes an NRGBA64 color in a friendlier format.
type bgcolor color.NRGBA64

func (c *bgcolor) MarshalText() ([]byte, error) {
	if c == nil {
		return nil, nil
	}
	b := make([]byte, 1, 17)
	b[0] = '#'
	hex.Encode(b[1:], []byte{byte(c.R >> 8), byte(c.R)})
	hex.Encode(b[5:], []byte{byte(c.G >> 8), byte(c.G)})
	hex.Encode(b[9:], []byte{byte(c.B >> 8), byte(c.B)})
	hex.Encode(b[13:], []byte{byte(c.A >> 8), byte(c.A)})
	return b, nil
}

func (c *bgcolor) UnmarshalText(text []byte) error {
	var r, g, b, a uint16
	switch len(text) {
	case 3: // rgb
		c, err := strconv.ParseUint(string(text[0:1]), 16, 4)
		if err != nil {
			return err
		}
		r = uint16(c) * 0x1111
		c, err = strconv.ParseUint(string(text[1:2]), 16, 4)
		if err != nil {
			return err
		}
		g = uint16(c) * 0x1111
		c, err = strconv.ParseUint(string(text[2:3]), 16, 4)
		if err != nil {
			return err
		}
		b = uint16(c) * 0x1111
	case 4: // rgba
		c, err := strconv.ParseUint(string(text[0:1]), 16, 4)
		if err != nil {
			return err
		}
		r = uint16(c) * 0x1111
		c, err = strconv.ParseUint(string(text[1:2]), 16, 4)
		if err != nil {
			return err
		}
		g = uint16(c) * 0x1111
		c, err = strconv.ParseUint(string(text[2:3]), 16, 4)
		if err != nil {
			return err
		}
		b = uint16(c) * 0x1111
		c, err = strconv.ParseUint(string(text[3:4]), 16, 4)
		if err != nil {
			return err
		}
		a = uint16(c) * 0x1111
	case 6: // rrggbb
		c, err := strconv.ParseUint(string(text[0:2]), 16, 8)
		if err != nil {
			return err
		}
		r = uint16(c) * 0x0101
		c, err = strconv.ParseUint(string(text[2:4]), 16, 8)
		if err != nil {
			return err
		}
		g = uint16(c) * 0x0101
		c, err = strconv.ParseUint(string(text[4:6]), 16, 8)
		if err != nil {
			return err
		}
		b = uint16(c) * 0x0101
	case 8: // rrggbbaa
		c, err := strconv.ParseUint(string(text[0:2]), 16, 8)
		if err != nil {
			return err
		}
		r = uint16(c) * 0x0101
		c, err = strconv.ParseUint(string(text[2:4]), 16, 8)
		if err != nil {
			return err
		}
		g = uint16(c) * 0x0101
		c, err = strconv.ParseUint(string(text[4:6]), 16, 8)
		if err != nil {
			return err
		}
		b = uint16(c) * 0x0101
		c, err = strconv.ParseUint(string(text[6:8]), 16, 8)
		if err != nil {
			return err
		}
		a = uint16(c) * 0x0101
	case 12: // rrrrggggbbbb
		c, err := strconv.ParseUint(string(text[0:4]), 16, 16)
		if err != nil {
			return err
		}
		r = uint16(c)
		c, err = strconv.ParseUint(string(text[4:8]), 16, 16)
		if err != nil {
			return err
		}
		g = uint16(c)
		c, err = strconv.ParseUint(string(text[8:12]), 16, 16)
		if err != nil {
			return err
		}
		b = uint16(c)
	case 16: // rrrrggggbbbbaaaa
		c, err := strconv.ParseUint(string(text[0:4]), 16, 16)
		if err != nil {
			return err
		}
		r = uint16(c)
		c, err = strconv.ParseUint(string(text[4:8]), 16, 16)
		if err != nil {
			return err
		}
		g = uint16(c)
		c, err = strconv.ParseUint(string(text[8:12]), 16, 16)
		if err != nil {
			return err
		}
		b = uint16(c)
		c, err = strconv.ParseUint(string(text[12:16]), 16, 16)
		if err != nil {
			return err
		}
		a = uint16(c)
	}
	*c = bgcolor{R: r, G: g, B: b, A: a}
	return nil
}

// Marshal creates a JSON encoding of the renderer and system information
// needed to serialize the system. If system.Check returns a non-nil error,
// then that error is returned instead.
func Marshal(system xirho.System, r *xirho.Render, tm xirho.ToneMap, bg *color.NRGBA64, meta *xirho.Metadata) ([]byte, error) {
	return Wrap(system, r, tm, bg, meta).MarshalJSON()
}

// Unmarshal decodes a xirho renderer from serialized JSON.
// Calling UseNumber on the decoder allows Unmarshal to guarantee full
// precision for xirho.Int function parameters.
func Unmarshal(d *json.Decoder) (*System, error) {
	var s System
	err := d.Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
