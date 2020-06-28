package main

import (
	"context"
	"encoding/json"
	"flag"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"time"

	"golang.org/x/image/draw"

	"github.com/dim13/colormap"
	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/encoding"
	"github.com/zephyrtronium/xirho/xi"
)

func main() {
	var outname, profname, inname string
	var sigint bool
	var timeout time.Duration
	var iters, hits int64
	var width, height int
	var osa int
	var gamma float64
	var resample string
	var procs int
	var echo bool
	flag.StringVar(&outname, "png", "", "output filename (default stdout)")
	flag.StringVar(&profname, "prof", "", "CPU profile output (default no profiling)")
	flag.StringVar(&inname, "in", "", "input json filename (default stdin)")
	flag.BoolVar(&sigint, "C", true, "save image on interrupt instead of exiting")
	flag.DurationVar(&timeout, "dur", 0, "max duration to render (default ignored)")
	flag.Int64Var(&iters, "iters", 0, "max iters (default ignored)")
	flag.Int64Var(&hits, "hits", 0, "max hits (default iters/5)")
	flag.IntVar(&width, "width", 1024, "output image width")
	flag.IntVar(&height, "height", 1024, "output image height")
	flag.IntVar(&osa, "osa", 1, "oversampling; histogram bins per pixel per axis")
	flag.Float64Var(&gamma, "gamma", 1, "gamma factor")
	flag.StringVar(&resample, "resample", "catmull-rom", "resampling method (catmull-rom, bilinear, approx-bilinear, or nearest)")
	flag.IntVar(&procs, "procs", runtime.GOMAXPROCS(0), "concurrent render routines")
	flag.BoolVar(&echo, "echo", false, "print system encoding before rendering")
	flag.Parse()
	resampler := resamplers[resample]
	if resampler == nil {
		log.Fatalln("no resampler named", resample)
	}
	if profname != "" {
		prof, err := os.Create(profname)
		if err != nil {
			log.Fatal(err)
		}
		if err = pprof.StartCPUProfile(prof); err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}
	ctx := context.Background()
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	}
	if sigint {
		ctx, cancel = context.WithCancel(ctx)
		ch := make(chan os.Signal)
		signal.Notify(ch, os.Interrupt)
		go func() {
			<-ch
			signal.Reset(os.Interrupt)
			cancel()
		}()
	}
	if hits <= 0 {
		hits = iters / 5
	}

	var in io.Reader = os.Stdin
	if inname != "" {
		f, err := os.Open(inname)
		if err != nil {
			log.Fatalln("error opening input:", err)
		}
		in = f
	}
	d := json.NewDecoder(in)
	r, _, err := encoding.Unmarshal(d)
	if err != nil {
		log.Fatalln("error unmarshaling system:", err)
	}
	log.Println("allocating histogram, estimated", xirho.HistMem(width*osa, height*osa)>>20, "MB")
	r.Hist = xirho.NewHist(width*osa, height*osa, gamma)
	r.Procs = procs
	r.N = iters
	r.Q = hits
	if echo {
		m, err := encoding.Marshal(r)
		if err != nil {
			log.Fatalln("error reading system from input:", err)
		}
		log.Printf("system:\n%s\n", m)
	}
	log.Println("rendering up to", r.N, "iters or", r.Q, "hits or", timeout)
	r.Render(ctx)
	log.Println("finished render with", r.Iters(), "iters,", r.Hits(), "hits")
	signal.Reset(os.Interrupt) // no rendering for ^C to interrupt
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.NRGBA64{A: 0xffff}), image.Point{}, draw.Src)
	log.Printf("drawing onto image of size %dx%d", width, height)
	resampler.Scale(img, img.Bounds(), r.Hist, r.Hist.Bounds(), draw.Over, nil)
	out := os.Stdout
	if outname != "" {
		log.Println("encoding to", outname)
		var err error
		out, err = os.Create(outname)
		if err != nil {
			log.Fatalln("error creating output file:", err)
		}
		defer out.Close()
	} else {
		log.Println("encoding to stdout")
	}
	err = png.Encode(out, img)
	if err != nil {
		log.Fatalln("error encoding image:", err)
	}
}

var resamplers = map[string]draw.Scaler{
	"catmull-rom":     draw.CatmullRom,
	"bilinear":        draw.BiLinear,
	"approx-bilinear": draw.ApproxBiLinear,
	"nearest":         draw.NearestNeighbor,
}

func mkpalette() []color.NRGBA64 {
	// ---- matplotlib colormaps ----
	// cmap := colormap.Viridis
	// cmap := colormap.Inferno
	// cmap := colormap.Magma
	// cmap := colormap.Plasma
	// v := make([]color.NRGBA64, len(cmap))
	// for i, c := range cmap {
	// 	r, g, b, a := c.RGBA()
	// 	v[i] = color.NRGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a)}
	// }
	// return v
	// ---- rgb hexagon ----
	return []color.NRGBA64{
		{R: 0xffff, G: 0x0000, B: 0x0000, A: 0xffff},
		{R: 0xffff, G: 0xffff, B: 0x0000, A: 0xffff},
		{R: 0x0000, G: 0xffff, B: 0x0000, A: 0xffff},
		{R: 0x0000, G: 0xffff, B: 0xffff, A: 0xffff},
		{R: 0x0000, G: 0x0000, B: 0xffff, A: 0xffff},
		{R: 0xffff, G: 0x0000, B: 0xffff, A: 0xffff},
		{R: 0xffff, G: 0x0000, B: 0x0000, A: 0xffff},
		{R: 0xffff, G: 0xffff, B: 0x0000, A: 0xffff},
		{R: 0x0000, G: 0xffff, B: 0x0000, A: 0xffff},
		{R: 0x0000, G: 0xffff, B: 0xffff, A: 0xffff},
		{R: 0x0000, G: 0x0000, B: 0xffff, A: 0xffff},
		{R: 0xffff, G: 0x0000, B: 0xffff, A: 0xffff},
		{R: 0xffff, G: 0x0000, B: 0x0000, A: 0xffff},
		{R: 0xffff, G: 0xffff, B: 0x0000, A: 0xffff},
		{R: 0x0000, G: 0xffff, B: 0x0000, A: 0xffff},
		{R: 0x0000, G: 0xffff, B: 0xffff, A: 0xffff},
		{R: 0x0000, G: 0x0000, B: 0xffff, A: 0xffff},
		{R: 0xffff, G: 0x0000, B: 0xffff, A: 0xffff},
		{R: 0xffff, G: 0x0000, B: 0x0000, A: 0xffff},
		{R: 0xffff, G: 0xffff, B: 0x0000, A: 0xffff},
		{R: 0x0000, G: 0xffff, B: 0x0000, A: 0xffff},
		{R: 0x0000, G: 0xffff, B: 0xffff, A: 0xffff},
		{R: 0x0000, G: 0x0000, B: 0xffff, A: 0xffff},
		{R: 0xffff, G: 0x0000, B: 0xffff, A: 0xffff},
		{R: 0xffff, G: 0x0000, B: 0x0000, A: 0xffff},
		{R: 0xffff, G: 0xffff, B: 0x0000, A: 0xffff},
		{R: 0x0000, G: 0xffff, B: 0x0000, A: 0xffff},
		{R: 0x0000, G: 0xffff, B: 0xffff, A: 0xffff},
		{R: 0x0000, G: 0x0000, B: 0xffff, A: 0xffff},
		{R: 0xffff, G: 0x0000, B: 0xffff, A: 0xffff},
		{R: 0xffff, G: 0x0000, B: 0x0000, A: 0xffff},
		{R: 0xffff, G: 0xffff, B: 0x0000, A: 0xffff},
		{R: 0x0000, G: 0xffff, B: 0x0000, A: 0xffff},
		{R: 0x0000, G: 0xffff, B: 0xffff, A: 0xffff},
		{R: 0x0000, G: 0x0000, B: 0xffff, A: 0xffff},
		{R: 0xffff, G: 0x0000, B: 0xffff, A: 0xffff},
	}
	// ---- grayscale ----
	// r := make([]color.NRGBA64, 256)
	// for i := range r {
	// 	r[i] = color.NRGBA64{R: uint16(i * i), G: uint16(i * i), B: uint16(i * i), A: 0xffff}
	// }
	// return r
}

var _ = colormap.Inferno

func defaultGraph(n int) [][]float64 {
	r := make([][]float64, n)
	for i := range r {
		r[i] = make([]float64, n)
		for j := range r[i] {
			r[i][j] = 1
		}
	}
	return r
}

func then(fs ...xirho.F) xirho.F {
	return &xi.Then{Funcs: fs}
}

func affspd(ax xirho.Ax, color, speed float64) xirho.F {
	return &xi.Then{
		Funcs: xirho.FuncList{
			&xi.Affine{Ax: ax},
			&xi.ColorSpeed{Color: xirho.Real(color), Speed: xirho.Real(speed)},
		},
	}
}
