// Package encoding implements marshaling and unmarshaling function systems.
package encoding

import (
	"bytes"
	"compress/lzw"
	"encoding/binary"
	"encoding/json"
	"image/color"
	"io"

	"github.com/zephyrtronium/xirho"
)

// marshaler controls the marshaling and unmarshaling of a full function system
// and renderer.
//
// NOTE: The definition of marshaler must be kept up to date with xirho.System
// and xirho.R.
type marshaler struct {
	Meta *xirho.Metadata `json:"meta,omitempty"`
	// system params
	Funcs []*funcm `json:"funcs"`
	Final *funcm   `json:"final,omitempty"`
	// renderer params
	Aspect float64  `json:"aspect"`
	Camera xirho.Ax `json:"camera"`
	// brightness params
	Bright float64 `json:"bright"`
	Gamma  float64 `json:"gamma"`
	Thresh float64 `json:"thresh"`
	// Palette is formed by concatenating each channel of the NRGBA64 palette
	// in ARGB order as big-endian, then LZW-encoding the result.
	Palette []byte `json:"palette"`
}

// Marshal creates a JSON encoding of the renderer and system information
// needed to serialize the system. If system.Check returns a non-nil error,
// then that error is returned instead.
func Marshal(system xirho.System, r *xirho.Render, tm xirho.ToneMap) ([]byte, error) {
	// TODO: wrap errors
	if err := system.Check(); err != nil {
		return nil, err
	}
	m := marshaler{
		Meta:   r.Meta,
		Funcs:  make([]*funcm, len(system.Nodes)),
		Camera: r.Camera,
		Bright: tm.Brightness,
		Gamma:  tm.Gamma,
		Thresh: tm.GammaMin,
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
	m.Aspect = float64(r.Hist.Cols()) / float64(r.Hist.Rows())
	palette := make([]byte, 0, 2*4*len(r.Palette))
	for _, c := range r.Palette {
		palette = append(palette, byte(c.A>>8), byte(c.A))
	}
	for _, c := range r.Palette {
		palette = append(palette, byte(c.R>>8), byte(c.R))
	}
	for _, c := range r.Palette {
		palette = append(palette, byte(c.G>>8), byte(c.G))
	}
	for _, c := range r.Palette {
		palette = append(palette, byte(c.B>>8), byte(c.B))
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

// Unmarshal decodes a xirho renderer from serialized JSON. The returned aspect
// is the number of columns per row in the histogram. The histogram should have
// its Reset method called before use. The Procs, N, and Q fields are left 0.
// Calling UseNumber on the decoder allows Unmarshal to guarantee full
// precision for xirho.Int function parameters.
func Unmarshal(d *json.Decoder) (system xirho.System, render *xirho.Render, tm xirho.ToneMap, aspect float64, err error) {
	// TODO: wrap errors
	m := marshaler{}
	if err = d.Decode(&m); err != nil {
		return
	}
	render = &xirho.Render{
		Hist:   &xirho.Hist{},
		Camera: m.Camera,
	}
	system = xirho.System{
		Nodes: make([]xirho.Node, len(m.Funcs)),
	}
	aspect = m.Aspect
	for i, f := range m.Funcs {
		system.Nodes[i].Func, err = unf(f)
		if err != nil {
			return
		}
		system.Nodes[i].Opacity = f.Opacity
		system.Nodes[i].Weight = f.Weight
		system.Nodes[i].Graph = f.Graph
		system.Nodes[i].Label = f.Label
	}
	if m.Final != nil {
		system.Final, err = unf(m.Final)
		if err != nil {
			return
		}
	}
	tm = xirho.ToneMap{Brightness: m.Bright, Gamma: m.Gamma, GammaMin: m.Thresh}
	z := lzw.NewReader(bytes.NewReader(m.Palette), lzw.LSB, 8)
	var buf bytes.Buffer
	if _, err = io.Copy(&buf, z); err != nil {
		return
	}
	palette := buf.Bytes()
	render.Palette = make([]color.NRGBA64, len(palette)/(2*4))
	for i := range render.Palette {
		render.Palette[i] = color.NRGBA64{
			A: binary.BigEndian.Uint16(palette[2*i:]),
			R: binary.BigEndian.Uint16(palette[2*(i+len(render.Palette)):]),
			G: binary.BigEndian.Uint16(palette[2*(i+2*len(render.Palette)):]),
			B: binary.BigEndian.Uint16(palette[2*(i+3*len(render.Palette)):]),
		}
	}
	return
}
