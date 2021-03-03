package main

import (
	"bufio"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/image/draw"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/encoding"
	"github.com/zephyrtronium/xirho/encoding/flame"
	"github.com/zephyrtronium/xirho/xmath"
)

func interactive(ctx context.Context, s *encoding.System, w, h int, res draw.Scaler, tm xirho.ToneMap, bg color.Color, osa, procs int) {
	r := &xirho.Render{Hist: xirho.NewHist(w*osa, h*osa)}
	status := status{
		r:      r,
		change: make(chan xirho.ChangeRender, 1),
		plot:   make(chan xirho.PlotOnto),
		imgs:   make(chan draw.Image),
		onto: xirho.PlotOnto{
			Scale:   res,
			ToneMap: tm,
		},
		bg:    image.Uniform{C: bg},
		sz:    image.Pt(w, h),
		osa:   osa,
		procs: procs,
	}
	if s != nil {
		cam := s.Camera
		c := xirho.ChangeRender{
			System:  s.System,
			Size:    image.Pt(w*status.osa, h*status.osa),
			Camera:  &cam,
			Palette: s.Palette,
			Procs:   status.procs,
		}
		status.change <- c
	}
	ctx, status.cancel = context.WithCancel(ctx)
	defer status.cancel()
	in := bufio.NewScanner(os.Stdin)
	go r.RenderAsync(ctx, status.change, status.plot, status.imgs)
	for ctx.Err() == nil {
		loop(ctx, &status, in)
	}
}

func loop(ctx context.Context, status *status, in *bufio.Scanner) {
	fmt.Print("] ")
	if !in.Scan() {
		status.cancel()
		return
	}
	line := strings.TrimSpace(in.Text())
	k := strings.IndexFunc(line, unicode.IsSpace)
	if k < 0 {
		k = len(line)
	}
	cmd := line[:k]
	line = strings.TrimSpace(line[k:])
	for _, c := range commands {
		if c.is(cmd) {
			if c.exec == nil {
				fmt.Println("not implemented yet :(")
				return
			}
			c.exec(ctx, status, line)
			return
		}
	}
	fmt.Println("no such command; try help")
}

type status struct {
	r      *xirho.Render
	cancel context.CancelFunc
	change chan xirho.ChangeRender
	plot   chan xirho.PlotOnto
	imgs   chan draw.Image
	onto   xirho.PlotOnto
	bg     image.Uniform
	sz     image.Point
	osa    int
	procs  int
	tick   *time.Ticker
}

type command struct {
	name []string
	desc string
	exec func(ctx context.Context, status *status, line string)
}

func (c *command) is(cmd string) bool {
	for _, n := range c.name {
		if strings.EqualFold(n, cmd) {
			return true
		}
	}
	return false
}

var commands = []*command{
	{
		name: []string{"help", "?"},
		desc: `show commands list & render status`,
		exec: nil, // set in init to avoid initialization cycle
	},
	{
		name: []string{"open", "load", "o"},
		desc: `open xirho json file & begin render`,
		exec: open,
	},
	{
		name: []string{"flame", "flam3"},
		desc: `open flame xml file & begin render`,
		exec: flam3,
	},
	{
		name: []string{"width", "w"},
		desc: `set image width, preserving aspect ratio`,
		exec: width,
	},
	{
		name: []string{"height", "h"},
		desc: `set image height, preserving aspect ratio`,
		exec: height,
	},
	{
		name: []string{"resize", "size", "rsz", "sz"},
		desc: `set image size`,
		exec: size,
	},
	{
		name: []string{"oversample", "osa"},
		desc: `set histogram size multiplier`,
		exec: oversample,
	},
	{
		name: []string{"pause", "stop"},
		desc: `pause render`,
		exec: pause,
	},
	{
		name: []string{"unpause", "go"},
		desc: `unpause render`,
		exec: unpause,
	},
	{
		name: []string{"bg"},
		desc: `set render background color`,
		exec: background,
	},
	{
		name: []string{"brightness", "br", "b"},
		desc: `set render brightness`,
		exec: brightness,
	},
	{
		name: []string{"gamma", "g"},
		desc: `set render gamma factor`,
		exec: gamma,
	},
	{
		name: []string{"threshold", "gt", "t"},
		desc: `set render gamma threshold`,
		exec: threshold,
	},
	{
		name: []string{"scaler", "scale", "resample"},
		desc: `set resampling method for rendering`,
		exec: scaler,
	},
	{
		name: []string{"procs", "goroutines", "workers"},
		desc: `set number of worker processes`,
		exec: procs,
	},
	{
		name: []string{"dx"},
		desc: `translate camera horizontally`,
		exec: camdx,
	},
	{
		name: []string{"dy"},
		desc: `translate camera vertically`,
		exec: camdy,
	},
	{
		name: []string{"dz"},
		desc: `translate camera forward/backward`,
		exec: camdz,
	},
	{
		name: []string{"roll"},
		desc: `rotate camera about z axis`,
		exec: camroll,
	},
	{
		name: []string{"pitch"},
		desc: `rotate camera about x axis`,
		exec: campitch,
	},
	{
		name: []string{"yaw"},
		desc: `rotate camera about y axis`,
		exec: camyaw,
	},
	{
		name: []string{"zoom", "z"},
		desc: `zoom camera in or out`,
		exec: camzoom,
	},
	{
		name: []string{"eye"},
		desc: `reset camera to a reasonable default`,
		exec: cameye,
	},
	{
		name: []string{"render", "r"},
		desc: `plot render and encode to a png file`,
		exec: render,
	},
	{
		name: []string{"preview", "quick", "p"},
		desc: `plot render using a fast resampling method`,
		exec: preview,
	},
	{
		name: []string{"quit", "exit", "q", "x"},
		desc: `exit interactive mode`,
		exec: quit,
	},
}

func init() {
	for i, c := range commands {
		if c.name[0] == "help" {
			commands[i].exec = help
			return
		}
	}
	panic("didn't set up help")
}

func help(ctx context.Context, status *status, line string) {
	if line == "?" {
		fmt.Println("Cheeky...")
	}
	fmt.Println("xirho interactive mode")
	fmt.Println("Commands (enter ? after a command for its help):")
	for _, c := range commands {
		fmt.Printf("\t%s - %s\n", strings.Join(c.name, " "), c.desc)
	}
	fmt.Printf("Rendering with %d procs:\n", status.procs)
	n, q := status.r.Iters(), status.r.Hits()
	fmt.Printf("Ran %d iters, plotted %d points, hit ratio %f\n", n, q, float64(q)/float64(n))
	fmt.Printf("Output image size %dx%d\n", status.sz.X, status.sz.Y)
	cols, rows := status.r.Hist.Cols(), status.r.Hist.Rows()
	fmt.Printf("Histogram oversampled %dx, size %dx%d (%d MB)\n", status.osa, cols, rows, xirho.HistMem(cols, rows)>>20)
	fmt.Printf("Plotting brightness %f, gamma %f, gamma threshold %f\n", status.onto.ToneMap.Brightness, status.onto.ToneMap.Gamma, status.onto.ToneMap.GammaMin)
	r, g, b, a := status.bg.C.RGBA()
	fmt.Printf("Plot background RGBA: #%02x%02x%02x%02x\n", r>>8, g>>8, b>>8, a>>8)
}

func open(ctx context.Context, status *status, line string) {
	const usage = `open <file.json>
	Decode a xirho system from the given file path and begin rendering it.
	The width or height of the current render is preserved depending on
	the aspect ratio of the new system, width if landscape or height if
	portrait. If an error occurs, the current render remains but is
	paused.`
	if line == "" || line == "?" {
		fmt.Println(usage)
		return
	}
	select {
	case <-ctx.Done():
		return
	case status.change <- xirho.ChangeRender{}:
		// pause render since we're going to discard it
	}
	f, err := os.Open(line)
	if err != nil {
		fmt.Println(err)
		return
	}
	d := json.NewDecoder(f)
	s, err := encoding.Unmarshal(d)
	f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	if s.Meta != nil {
		fmt.Printf("%s (%s)\n", s.Meta.Title, s.Meta.Date.Format(time.Stamp))
		fmt.Println("Author(s):", strings.Join(s.Meta.Authors, ", "))
		fmt.Println("Licensed under", s.Meta.License)
	}
	w, h := xmath.Fit(status.sz.X, status.sz.Y, s.Aspect)
	status.onto.ToneMap = s.ToneMap
	status.sz = image.Pt(w, h)
	cam := s.Camera
	c := xirho.ChangeRender{
		System:  s.System,
		Size:    image.Pt(w*status.osa, h*status.osa),
		Camera:  &cam,
		Palette: s.Palette,
		Procs:   status.procs,
	}
	select {
	case <-ctx.Done():
		return
	case status.change <- c:
		// do nothing
	}
}

func flam3(ctx context.Context, status *status, line string) {
	const usage = `flame <file.xml>
	Decode a flam3 system from the given file path and begin rendering it.
	The width or height of the current render is preserved depending on
	the aspect ratio of the new system, width if landscape or height if
	portrait. If an error occurs, the current render remains but is
	paused.`
	if line == "" || line == "?" {
		fmt.Println(usage)
		return
	}
	select {
	case <-ctx.Done():
		return
	case status.change <- xirho.ChangeRender{}:
		// pause render since we're going to discard it
	}
	f, err := os.Open(line)
	if err != nil {
		fmt.Println(err)
		return
	}
	d := xml.NewDecoder(f)
	s, err := flame.Unmarshal(d)
	f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(s.Meta.Title) // always present
	if len(s.Unrecognized) != 0 {
		fmt.Printf("Unrecognized attributes: %q\n", s.Unrecognized)
	}
	w, h := xmath.Fit(status.sz.X, status.sz.Y, s.Aspect)
	status.onto.ToneMap = s.ToneMap
	status.bg = image.Uniform{C: s.BG}
	status.sz = image.Pt(w, h)
	cam := s.Camera
	c := xirho.ChangeRender{
		System:  s.System,
		Size:    image.Pt(w*status.osa, h*status.osa),
		Camera:  &cam,
		Palette: s.Palette,
		Procs:   status.procs,
	}
	select {
	case <-ctx.Done():
		return
	case status.change <- c:
		// do nothing
	}
}

func width(ctx context.Context, status *status, line string) {
	const usage = `width <px>
	Set the output image width in pixels. The current aspect ratio is
	preserved. The existing oversampling value is applied.
	Resets rendering progress.`
	if line == "" || line == "?" {
		fmt.Println(usage)
		return
	}
	w, err := strconv.Atoi(line)
	if err != nil {
		fmt.Println(err)
		return
	}
	h := int(float64(w)/status.r.Hist.Aspect() + 0.5)
	status.sz = image.Pt(w, h)
	c := xirho.ChangeRender{Size: image.Pt(w*status.osa, h*status.osa), Procs: status.procs}
	select {
	case <-ctx.Done():
		return
	case status.change <- c:
		// do nothing
	}
}

func height(ctx context.Context, status *status, line string) {
	const usage = `height <px>
	Set the output image height in pixels. The current aspect ratio is
	preserved. The existing oversampling value is applied.
	Resets rendering progress.`
	if line == "" || line == "?" {
		fmt.Println(usage)
		return
	}
	h, err := strconv.Atoi(line)
	if err != nil {
		fmt.Println(err)
		return
	}
	w := int(float64(h)*status.r.Hist.Aspect() + 0.5)
	status.sz = image.Pt(w, h)
	c := xirho.ChangeRender{Size: image.Pt(w*status.osa, h*status.osa), Procs: status.procs}
	select {
	case <-ctx.Done():
		return
	case status.change <- c:
		// do nothing
	}
}

func size(ctx context.Context, status *status, line string) {
	const usage = `size <w>x<h>
	Set the output image size in pixels.
	Resets rendering progress.`
	if line == "" || line == "?" {
		fmt.Println(usage)
		return
	}
	line = strings.ToLower(line)
	f := strings.Split(line, "x")
	if len(f) != 2 {
		fmt.Println(`size <w>x<h>, e.g. size 1024x1024`)
		return
	}
	w, err := strconv.Atoi(f[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	if w <= 0 {
		fmt.Println("can't set width to", w)
		return
	}
	h, err := strconv.Atoi(f[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	if h <= 0 {
		fmt.Println("can't set height to", h)
		return
	}
	status.sz = image.Pt(w, h)
	c := xirho.ChangeRender{Size: image.Pt(w*status.osa, h*status.osa), Procs: status.procs}
	select {
	case <-ctx.Done():
		return
	case status.change <- c:
		// do nothing
	}
}

func oversample(ctx context.Context, status *status, line string) {
	const usage = `oversample <n>
	Set histogram bins per pixel per axis. Be careful, memory usage grows
	very rapidly as this increases!
	Resets rendering progress.`
	if line == "" || line == "?" {
		fmt.Println(usage)
		return
	}
	osa, err := strconv.Atoi(line)
	if err != nil {
		fmt.Println(err)
		return
	}
	if osa <= 0 {
		fmt.Println("can't set oversampling to", osa)
		return
	}
	w, h := status.sz.X*osa, status.sz.Y*osa
	fmt.Printf("new histogram memory usage will be %d MB\n", xirho.HistMem(w, h)>>20)
	status.osa = osa
	c := xirho.ChangeRender{Size: image.Pt(w, h), Procs: status.procs}
	select {
	case <-ctx.Done():
		return
	case status.change <- c:
		// do nothing
	}
}

func pause(ctx context.Context, status *status, line string) {
	const usage = `pause
	Pause rendering until unpause or any command that resets rendering
	progress.`
	if line == "?" {
		fmt.Println(usage)
		return
	}
	c := xirho.ChangeRender{}
	select {
	case <-ctx.Done():
		return
	case status.change <- c:
		// do nothing
	}
}

func unpause(ctx context.Context, status *status, line string) {
	const usage = `unpause
	Unpause rendering after a previous pause.`
	if line == "?" {
		fmt.Println(usage)
		return
	}
	c := xirho.ChangeRender{Procs: status.procs}
	select {
	case <-ctx.Done():
		return
	case status.change <- c:
		// do nothing
	}
}

func background(ctx context.Context, status *status, line string) {
	const usage = `bg <#rrggbbaa>
	Set background color for rendered images. The color may use one, two,
	or four hexadecimal digits per channel, laid out as RGB or RGBA.`
	if line == "" || line == "?" {
		fmt.Println(usage)
		return
	}
	line = strings.TrimPrefix(line, "#")
	var r, g, b, a uint16
	a = 0xffff
	switch len(line) {
	case 3: // rgb
		c, err := strconv.ParseUint(line[0:1], 16, 4)
		if err != nil {
			fmt.Println(err)
			return
		}
		r = uint16(c) * 0x1111
		c, err = strconv.ParseUint(line[1:2], 16, 4)
		if err != nil {
			fmt.Println(err)
			return
		}
		g = uint16(c) * 0x1111
		c, err = strconv.ParseUint(line[2:3], 16, 4)
		if err != nil {
			fmt.Println(err)
			return
		}
		b = uint16(c) * 0x1111
	case 4: // rgba
		c, err := strconv.ParseUint(line[0:1], 16, 4)
		if err != nil {
			fmt.Println(err)
			return
		}
		r = uint16(c) * 0x1111
		c, err = strconv.ParseUint(line[1:2], 16, 4)
		if err != nil {
			fmt.Println(err)
			return
		}
		g = uint16(c) * 0x1111
		c, err = strconv.ParseUint(line[2:3], 16, 4)
		if err != nil {
			fmt.Println(err)
			return
		}
		b = uint16(c) * 0x1111
		c, err = strconv.ParseUint(line[3:4], 16, 4)
		if err != nil {
			fmt.Println(err)
			return
		}
		a = uint16(c) * 0x1111
	case 6: // rrggbb
		c, err := strconv.ParseUint(line[0:2], 16, 8)
		if err != nil {
			fmt.Println(err)
			return
		}
		r = uint16(c) * 0x0101
		c, err = strconv.ParseUint(line[2:4], 16, 8)
		if err != nil {
			fmt.Println(err)
			return
		}
		g = uint16(c) * 0x0101
		c, err = strconv.ParseUint(line[4:6], 16, 8)
		if err != nil {
			fmt.Println(err)
			return
		}
		b = uint16(c) * 0x0101
	case 8: // rrggbbaa
		c, err := strconv.ParseUint(line[0:2], 16, 8)
		if err != nil {
			fmt.Println(err)
			return
		}
		r = uint16(c) * 0x0101
		c, err = strconv.ParseUint(line[2:4], 16, 8)
		if err != nil {
			fmt.Println(err)
			return
		}
		g = uint16(c) * 0x0101
		c, err = strconv.ParseUint(line[4:6], 16, 8)
		if err != nil {
			fmt.Println(err)
			return
		}
		b = uint16(c) * 0x0101
		c, err = strconv.ParseUint(line[6:8], 16, 8)
		if err != nil {
			fmt.Println(err)
			return
		}
		a = uint16(c) * 0x0101
	case 12: // rrrrggggbbbb
		c, err := strconv.ParseUint(line[0:4], 16, 16)
		if err != nil {
			fmt.Println(err)
			return
		}
		r = uint16(c)
		c, err = strconv.ParseUint(line[4:8], 16, 16)
		if err != nil {
			fmt.Println(err)
			return
		}
		g = uint16(c)
		c, err = strconv.ParseUint(line[8:12], 16, 16)
		if err != nil {
			fmt.Println(err)
			return
		}
		b = uint16(c)
	case 16: // rrrrggggbbbbaaaa
		c, err := strconv.ParseUint(line[0:4], 16, 16)
		if err != nil {
			fmt.Println(err)
			return
		}
		r = uint16(c)
		c, err = strconv.ParseUint(line[4:8], 16, 16)
		if err != nil {
			fmt.Println(err)
			return
		}
		g = uint16(c)
		c, err = strconv.ParseUint(line[8:12], 16, 16)
		if err != nil {
			fmt.Println(err)
			return
		}
		b = uint16(c)
		c, err = strconv.ParseUint(line[12:16], 16, 16)
		if err != nil {
			fmt.Println(err)
			return
		}
		a = uint16(c)
	}
	status.bg = image.Uniform{C: color.NRGBA64{R: r, G: g, B: b, A: a}}
}

func brightness(ctx context.Context, status *status, line string) {
	const usage = `brightness <x>
	Set brightness, the scaling of alpha relative to color channels for
	rendered images. x may be any finite number greater than 0.`
	if line == "" || line == "?" {
		fmt.Println(usage)
		return
	}
	x, err := strconv.ParseFloat(line, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	if x <= 0 || !xmath.IsFinite(x) {
		fmt.Println("can't set brightness to", x)
		return
	}
	status.onto.ToneMap.Brightness = x
}

func gamma(ctx context.Context, status *status, line string) {
	const usage = `gamma <x>
	Set gamma factor, which controls nonlinear scaling according to
	per-bin brightness. x may be any finite number greater than 0.`
	if line == "" || line == "?" {
		fmt.Println(usage)
		return
	}
	x, err := strconv.ParseFloat(line, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	if x <= 0 || !xmath.IsFinite(x) {
		fmt.Println("can't set gamma to", x)
		return
	}
	status.onto.ToneMap.Gamma = x
}

func threshold(ctx context.Context, status *status, line string) {
	const usage = `threshold <x>
	Set gamma threshold, which controls the minimum bin brightness to
	gamma scaling is applied, relative to the expected bin count. x may be
	in the interval [0, 1].`
	if line == "" || line == "?" {
		fmt.Println(usage)
		return
	}
	x, err := strconv.ParseFloat(line, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !(x >= 0 && x <= 1) { // condition negated to detect nan
		fmt.Println("can't set gamma threshold to", x)
		return
	}
	status.onto.ToneMap.GammaMin = x
}

func scaler(ctx context.Context, status *status, line string) {
	const usage = `scaler <name>
	Set the resampling method used to scale down oversampled histograms to
	output images when using the render command. Note that methods other
	than nearest and approx-bilinear may be very slow with large
	oversampling values.`
	if line == "" || line == "?" {
		fmt.Println(usage)
		fmt.Println("Available resampling methods:")
		names := make([]string, 0, len(resamplers))
		for name := range resamplers {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			fmt.Printf("\t%s\n", name)
		}
		return
	}
	line = strings.ToLower(line)
	r := resamplers[line]
	if r == nil {
		fmt.Println("no such resampler")
		return
	}
	status.onto.Scale = r
}

func procs(ctx context.Context, status *status, line string) {
	const usage = `procs <n>
	Set the number of worker threads used to iterate the system. n must be
	greater than 0. Also unpauses if currently paused.`
	if line == "" || line == "?" {
		fmt.Println(usage)
		return
	}
	n, err := strconv.Atoi(line)
	if err != nil {
		fmt.Println(err)
		return
	}
	if n <= 0 {
		fmt.Println("can't use", n, "procs")
		return
	}
	status.procs = n
	c := xirho.ChangeRender{Procs: n}
	select {
	case <-ctx.Done():
		return
	case status.change <- c:
		// do nothing
	}
}

func camdx(ctx context.Context, status *status, line string) {
	const usage = `dx <d>
	Translate the camera along the horizontal axis. Positive moves right.`
	if line == "" || line == "?" {
		fmt.Println(usage)
		return
	}
	d, err := strconv.ParseFloat(line, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	cam := status.r.Camera
	cam.Translate(d, 0, 0)
	c := xirho.ChangeRender{Camera: &cam, Procs: status.procs}
	select {
	case <-ctx.Done():
		return
	case status.change <- c:
		// do nothing
	}
}

func camdy(ctx context.Context, status *status, line string) {
	const usage = `dy <d>
	Translate the camera along the vertical axis. Positive moves down.`
	if line == "" || line == "?" {
		fmt.Println(usage)
		return
	}
	d, err := strconv.ParseFloat(line, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	cam := status.r.Camera
	cam.Translate(0, d, 0)
	c := xirho.ChangeRender{Camera: &cam, Procs: status.procs}
	select {
	case <-ctx.Done():
		return
	case status.change <- c:
		// do nothing
	}
}

func camdz(ctx context.Context, status *status, line string) {
	const usage = `dz <d>
	Translate the camera along the depth axis. Positive moves forward.`
	if line == "" || line == "?" {
		fmt.Println(usage)
		return
	}
	d, err := strconv.ParseFloat(line, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	cam := status.r.Camera
	cam.Translate(0, 0, d)
	c := xirho.ChangeRender{Camera: &cam, Procs: status.procs}
	select {
	case <-ctx.Done():
		return
	case status.change <- c:
		// do nothing
	}
}

func camroll(ctx context.Context, status *status, line string) {
	const usage = `roll <d>
	Rotate the camera about the depth axis by an angle in clockwise
	degrees.`
	if line == "" || line == "?" {
		fmt.Println(usage)
		return
	}
	d, err := strconv.ParseFloat(line, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	cam := status.r.Camera
	cam.Roll(d * -math.Pi / 180)
	c := xirho.ChangeRender{Camera: &cam, Procs: status.procs}
	select {
	case <-ctx.Done():
		return
	case status.change <- c:
		// do nothing
	}
}

func campitch(ctx context.Context, status *status, line string) {
	const usage = `pitch <d>
	Rotate the camera about the horizontal axis by an angle in clockwise
	degrees.`
	if line == "" || line == "?" {
		fmt.Println(usage)
		return
	}
	d, err := strconv.ParseFloat(line, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	cam := status.r.Camera
	cam.Pitch(d * -math.Pi / 180)
	c := xirho.ChangeRender{Camera: &cam, Procs: status.procs}
	select {
	case <-ctx.Done():
		return
	case status.change <- c:
		// do nothing
	}
}

func camyaw(ctx context.Context, status *status, line string) {
	const usage = `yaw <d>
	Rotate the camera about the vertical axis by an angle in clockwise
	degrees.`
	if line == "" || line == "?" {
		fmt.Println(usage)
		return
	}
	d, err := strconv.ParseFloat(line, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	cam := status.r.Camera
	cam.Yaw(d * -math.Pi / 180)
	c := xirho.ChangeRender{Camera: &cam, Procs: status.procs}
	select {
	case <-ctx.Done():
		return
	case status.change <- c:
		// do nothing
	}
}

func camzoom(ctx context.Context, status *status, line string) {
	const usage = `zoom <d>
	Zoom the camera in or out by a multiplicative factor. Values greater
	than 1 zoom in and between 0 and 1 zoom out.`
	if line == "" || line == "?" {
		fmt.Println(usage)
		return
	}
	d, err := strconv.ParseFloat(line, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	cam := status.r.Camera
	cam.Zoom(d)
	c := xirho.ChangeRender{Camera: &cam, Procs: status.procs}
	select {
	case <-ctx.Done():
		return
	case status.change <- c:
		// do nothing
	}
}

func cameye(ctx context.Context, status *status, line string) {
	const usage = `eye
	Set the camera to the identity transform.`
	if line == "?" {
		fmt.Println(usage)
		return
	}
	cam := status.r.Camera
	cam.Eye()
	c := xirho.ChangeRender{Camera: &cam, Procs: status.procs}
	select {
	case <-ctx.Done():
		return
	case status.change <- c:
		// do nothing
	}
}

func render(ctx context.Context, status *status, line string) {
	const usage = `render <output.png>
	Render the current histogram to a PNG file. If not paused, rendering
	automatically resumes afterward.`
	if line == "" || line == "?" {
		fmt.Println(usage)
		return
	}
	onto := status.onto
	onto.Image = image.NewNRGBA64(image.Rect(0, 0, status.sz.X, status.sz.Y))
	draw.Draw(onto.Image, image.Rect(0, 0, status.sz.X, status.sz.Y), &status.bg, image.Point{}, draw.Src)
	var t time.Time
	select {
	case <-ctx.Done():
		return
	case status.plot <- onto:
		t = time.Now()
	}
	select {
	case <-ctx.Done():
		return
	case img := <-status.imgs:
		if img == nil {
			return
		}
		d := time.Since(t)
		fmt.Printf("rendered in %v, now encoding to %s\n", d, line)
		f, err := os.Create(line)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Println("error closing output file:", err)
			}
		}()
		if err := png.Encode(f, img); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func preview(ctx context.Context, status *status, line string) {
	const usage = `preview [output.png]
	Render the current histogram to a PNG file (defaulting to
	xirho-preview.png) using approx-bilinear resampling. If not paused,
	rendering automatically resumes afterward.`
	if line == "?" {
		fmt.Println(usage)
		return
	}
	if line == "" {
		line = "xirho-preview.png"
	}
	onto := status.onto
	onto.Image = image.NewNRGBA64(image.Rect(0, 0, status.sz.X, status.sz.Y))
	draw.Draw(onto.Image, image.Rect(0, 0, status.sz.X, status.sz.Y), &status.bg, image.Point{}, draw.Src)
	onto.Scale = draw.ApproxBiLinear
	var t time.Time
	select {
	case <-ctx.Done():
		return
	case status.plot <- onto:
		t = time.Now()
	}
	select {
	case <-ctx.Done():
		return
	case img := <-status.imgs:
		if img == nil {
			return
		}
		d := time.Since(t)
		fmt.Printf("rendered in %v, now encoding to %s\n", d, line)
		f, err := os.Create(line)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Println("error closing output file:", err)
			}
		}()
		if err := png.Encode(f, img); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func quit(ctx context.Context, status *status, line string) {
	const usage = `quit
	Stop iterating and exit the program.`
	if line == "?" {
		fmt.Println(usage)
		return
	}
	status.cancel()
}
