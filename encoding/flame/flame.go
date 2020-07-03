// Package flame implements parsing the XML-based Flame format.
package flame

import (
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xi"
)

// Unmarshal decodes a renderer from Flame XML.
func Unmarshal(d *xml.Decoder) (r *xirho.R, aspect float64, bg color.NRGBA64, err error) {
	var flm flame
	if err = d.Decode(&flm); err != nil {
		return
	}
	sz, err := nums(flm.Size)
	if err != nil {
		return
	}
	aspect = sz[0] / sz[1]
	tr, err := nums(flm.Center)
	if err != nil {
		return
	}
	var cam xirho.Ax
	cam.Eye()
	msz := sz[0]
	if sz[1] > msz { // <?
		msz = sz[1]
	}
	cam.Scale(4*flm.Scale/msz, 4*flm.Scale/msz, 4*flm.Scale/msz) // ?
	cam.RotX(-flm.Yaw)
	cam.RotY(-flm.Pitch)
	cam.RotZ(-flm.Angle)
	cam.Translate(-tr[0], -tr[1], flm.Zpos)
	bgc, err := nums(flm.Background)
	if err != nil {
		return
	}
	bg = color.NRGBA64{
		R: uint16(bgc[0] * 0xffff),
		G: uint16(bgc[1] * 0xffff),
		B: uint16(bgc[2] * 0xffff),
		A: 0xffff,
	}
	system := xirho.System{
		Funcs:   make([]xirho.F, len(flm.Xforms)),
		Weights: make([]float64, len(flm.Xforms)),
		Graph:   make([][]float64, len(flm.Xforms)),
	}
	var df decoded
	for i, xf := range flm.Xforms {
		df, err = decodexf(xf, false)
		system.Funcs[i] = df.f
		system.Weights[i] = df.weight
		if df.graph == nil {
			df.graph = make([]float64, len(flm.Xforms))
			for j := range df.graph {
				df.graph[j] = 1
			}
		}
		system.Graph[i] = df.graph
	}
	if flm.Final.XMLName.Local != "" {
		// A finalxform is an xform with a different name and some missing fields,
		// so we can decode it easily by making an xform out of the final and then
		// grabbing the information we care about.
		xf := xform{
			Color:    flm.Final.Color,
			Symmetry: flm.Final.Symmetry,
			Coefs:    flm.Final.Coefs,
			Post:     flm.Final.Post,
			Opacity:  flm.Final.Opacity,
			Attrs:    flm.Final.Attrs,
		}
		df, err = decodexf(xf, true)
		if err != nil {
			return
		}
		system.Final = df.f
	}
	r = &xirho.R{
		Hist:    &xirho.Hist{},
		System:  system,
		Camera:  cam,
		Palette: parsepalette(flm.Palette),
	}
	r.Hist.SetBrightness(flm.Brightness, flm.Gamma, flm.Thresh)
	// TODO: perspective
	return
}

// decoded is a decoded transform.
type decoded struct {
	f      xirho.F
	weight float64
	graph  []float64
}

// decodexf decodes an xform.
func decodexf(xf xform, final bool) (d decoded, err error) {
	d.weight = xf.Weight
	if xf.Chaos != "" {
		d.graph, err = nums(xf.Chaos)
		if err != nil {
			return
		}
	}
	// Collect variations first so we can check them for basic ones.
	vars := make(map[string]float64, len(xf.Attrs))
	for _, attr := range xf.Attrs {
		var v float64
		v, err = strconv.ParseFloat(attr.Value, 64)
		if err != nil {
			return
		}
		vars[attr.Name.Local] = v
	}
	// Decode affine transform.
	var ax xirho.Ax
	if ax, err = decodetx(xf.Coefs, "transform"); err != nil {
		return
	}
	// Check for variations that are really part of the transform.
	if v, ok := vars["linear"]; ok {
		ax.Scale(v, v, v)
	}
	if v, ok := vars["linear3D"]; ok {
		ax.Scale(v, v, v)
	}
	if _, ok := vars["flatten"]; ok {
		ax.Scale(1, 1, 0)
	}
	if v, ok := vars["pre_zscale"]; ok {
		ax.Scale(1, 1, v)
	}
	if v, ok := vars["pre_ztranslate"]; ok {
		ax.Translate(0, 0, v)
	}
	if v, ok := vars["pre_rotate_x"]; ok { // was that the name?
		ax.RotX(v)
	}
	if v, ok := vars["pre_rotate_y"]; ok {
		ax.RotY(v)
	}
	if v, ok := vars["pre_rotate_z"]; ok { // I think this one exists too
		ax.RotZ(v)
	}
	// Decode post-transform, if it exists.
	px := xirho.Eye()
	if xf.Post != "" {
		px, err = decodetx(xf.Post, "post-transform")
		if err != nil {
			return
		}
	}
	// Decode other variations. In the fractal flame algorithm, they are always
	// summed (except maybe pre/post are thend in a fixed order?).
	pre, in, post := xi.Sum{}, xi.Sum{}, xi.Sum{}
	for name := range vars {
		if parse, ok := Funcs[name]; ok {
			parse(vars, &pre, &in, &post, ax)
		}
	}
	// Then everything together. Only save the parts that do something.
	f := xi.Then{}
	if ax != xirho.Eye() {
		f.Funcs = append(f.Funcs, &xi.Affine{Ax: ax})
	}
	if s := sumdefault(pre); s != nil {
		f.Funcs = append(f.Funcs, s)
	}
	if s := sumdefault(in); s != nil {
		f.Funcs = append(f.Funcs, s)
	}
	if xf.Symmetry != 1 {
		cs := xi.ColorSpeed{
			Color: xirho.Real(xf.Color),
			Speed: xirho.Real(xf.Symmetry+1) / 2,
		}
		f.Funcs = append(f.Funcs, &cs)
	}
	if s := sumdefault(post); s != nil {
		f.Funcs = append(f.Funcs, s)
	}
	if px != xirho.Eye() {
		f.Funcs = append(f.Funcs, &xi.Affine{Ax: px})
	}
	switch len(f.Funcs) {
	case 0:
		// No functions. We skipped ax because it was an affine eye, so use it,
		// unless we're decoding a final, since finals are allowed to be nil.
		if !final {
			d.f = &xi.Affine{Ax: ax}
		}
	case 1:
		// No reason to Then a single function.
		d.f = f.Funcs[0]
	default:
		d.f = &f
	}
	return
}

// decodetx decodes an affine transform.
func decodetx(coefs, name string) (ax xirho.Ax, err error) {
	var a []float64
	a, err = nums(coefs)
	if err != nil {
		return
	}
	if len(a) != 6 {
		err = fmt.Errorf("invalid %s coefs: expected list of 6 numbers, got %v", name, a)
		return
	}
	ax = aff2to3(a)
	return
}

// sumdefault returns a function encapsulating the behavior of a Sum based on
// its function list and color function. The returned value is nil if the Sum
// contains no functions.
func sumdefault(f xi.Sum) xirho.F {
	switch {
	case len(f.Funcs) == 0 && f.Color.F == nil:
		return nil
	case len(f.Funcs) == 1 && f.Color.F == nil:
		return f.Funcs[0]
	case len(f.Funcs) == 0 && f.Color.F != nil:
		return f.Color.F
	default:
		return &f
	}
}

// aff2to3 converts a Flame 2D affine matrix to a xirho.Ax transform.
func aff2to3(a []float64) (ax xirho.Ax) {
	ax.Eye()
	ax[0] = a[0]
	ax[1] = a[2]
	ax[3] = a[4]
	ax[4] = a[1]
	ax[5] = a[3]
	ax[7] = a[5]
	return ax
}

// nums parses a space-separated list of numbers into a slice of float64s.
func nums(s string) ([]float64, error) {
	words := strings.Fields(s)
	r := make([]float64, 0, len(words))
	for _, word := range words {
		if word == "" {
			continue
		}
		v, err := strconv.ParseFloat(word, 64)
		if err != nil {
			return nil, err
		}
		r = append(r, v)
	}
	return r, nil
}

// parsepalette parses a flame palette.
func parsepalette(p palette) []color.NRGBA64 {
	r := make([]color.NRGBA64, 0, p.Count)
	for _, line := range strings.Fields(p.Data) {
		// The error from DecodeString is not checked because DecodeString
		// returns the decoded bytes before the error and we have no resolution
		// strategy other than to continue.
		v, _ := hex.DecodeString(line)
		for len(v) >= 3 {
			c := color.NRGBA64{
				R: uint16(v[0]) * 0x0101,
				G: uint16(v[1]) * 0x0101,
				B: uint16(v[2]) * 0x0101,
				A: 0xffff,
			}
			r = append(r, c)
			v = v[3:]
		}
	}
	return r
}

type flame struct {
	XMLName    xml.Name   `xml:"flame"`
	Name       string     `xml:"name,attr"`
	Size       string     `xml:"size,attr"`
	Center     string     `xml:"center,attr"`
	Scale      float64    `xml:"scale,attr"`
	Angle      float64    `xml:"angle,attr"`
	Pitch      float64    `xml:"cam_pitch,attr"`
	Yaw        float64    `xml:"cam_yaw,attr"`
	Zpos       float64    `xml:"cam_zpos,attr"`
	Background string     `xml:"background,attr"`
	Brightness float64    `xml:"brightness,attr"`
	Gamma      float64    `xml:"gamma,attr"`
	Thresh     float64    `xml:"gamma_threshold,attr"`
	Xforms     []xform    `xml:"xform"`
	Final      finalxform `xml:"finalxform"`
	Palette    palette    `xml:"palette"`
}

type xform struct {
	XMLName  xml.Name   `xml:"xform"`
	Weight   float64    `xml:"weight,attr"`
	Color    float64    `xml:"color,attr"`
	Symmetry float64    `xml:"symmetry,attr"`
	Coefs    string     `xml:"coefs,attr"`
	Post     string     `xml:"post,attr"`
	Chaos    string     `xml:"chaos,attr"`
	Opacity  float64    `xml:"opacity,attr"`
	Attrs    []xml.Attr `xml:",any,attr"`
}

type finalxform struct {
	XMLName  xml.Name   `xml:"finalxform"`
	Color    float64    `xml:"color,attr"`
	Symmetry float64    `xml:"symmetry,attr"`
	Coefs    string     `xml:"coefs,attr"`
	Post     string     `xml:"post,attr"`
	Opacity  float64    `xml:"opacity,attr"`
	Attrs    []xml.Attr `xml:",any,attr"`
}

type palette struct {
	XMLName xml.Name `xml:"palette"`
	Count   int      `xml:"count,attr"`
	Format  string   `xml:"format,attr"`
	Data    string   `xml:",chardata"`
}
