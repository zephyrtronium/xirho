// The xirho command implements a basic renderer using xirho.
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"encoding/xml"
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"time"

	"golang.org/x/image/draw"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/encoding"
	"github.com/zephyrtronium/xirho/encoding/flame"
	"github.com/zephyrtronium/xirho/hist"
)

func main() {
	var intr bool
	var outname, profname, inname, flamename, dumpname string
	var sigint bool
	var timeout time.Duration
	var width, height int
	var osa int
	var tm xirho.ToneMap
	var resample string
	var procs int
	var echo bool
	var bgr, bgg, bgb, bga int
	flag.BoolVar(&intr, "i", false, "interactive mode")
	flag.StringVar(&outname, "png", "", "output filename (default stdout)")
	flag.StringVar(&profname, "prof", "", "CPU profile output (default no profiling)")
	flag.StringVar(&inname, "in", "", "input json filename (default stdin)")
	flag.StringVar(&flamename, "flame", "", "input flame filename")
	flag.BoolVar(&sigint, "C", true, "save image on interrupt instead of exiting (ignored when interactive)")
	flag.DurationVar(&timeout, "dur", 0, "max duration to render (default ignored; always ignored when interactive)")
	flag.IntVar(&width, "width", 1024, "output image width")
	flag.IntVar(&height, "height", 1024, "output image height")
	flag.IntVar(&osa, "osa", 1, "oversampling; histogram bins per pixel per axis")
	flag.Float64Var(&tm.Gamma, "gamma", 0, "gamma factor")
	flag.Float64Var(&tm.GammaMin, "thresh", 0, "gamma threshold")
	flag.Float64Var(&tm.Brightness, "bright", 0, "brightness")
	flag.StringVar(&resample, "resample", "catmull-rom", "resampling method (catmull-rom, bilinear, approx-bilinear, or nearest)")
	flag.IntVar(&procs, "procs", runtime.GOMAXPROCS(0), "concurrent render routines")
	flag.BoolVar(&echo, "echo", false, "print system encoding before rendering")
	flag.IntVar(&bgr, "bg.r", 0, "background red (0-255)")
	flag.IntVar(&bgg, "bg.g", 0, "background green (0-255)")
	flag.IntVar(&bgb, "bg.b", 0, "background blue (0-255)")
	flag.IntVar(&bga, "bg.a", 255, "background alpha (0-255)")
	flag.StringVar(&dumpname, "raw-histogram-dump", "", "dump raw histogram data to file")
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
	if timeout > 0 && !intr {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	}
	if sigint && !intr {
		ctx, cancel = context.WithCancel(ctx)
		ch := make(chan os.Signal)
		signal.Notify(ch, os.Interrupt)
		go func() {
			<-ch
			signal.Reset(os.Interrupt)
			cancel()
		}()
	}

	var s *encoding.System
	var err error
	u := color.NRGBA64{
		R: uint16(bgr * 0x0101),
		G: uint16(bgg * 0x0101),
		B: uint16(bgb * 0x0101),
		A: uint16(bga * 0x0101),
	}
	switch {
	case inname != "":
		f, err := os.Open(inname)
		if err != nil {
			log.Fatalln("error opening input:", err)
		}
		d := json.NewDecoder(f)
		d.UseNumber()
		s, err = encoding.Unmarshal(d)
		if err != nil {
			log.Fatalln("error unmarshaling system:", err)
		}
	case flamename != "":
		f, err := os.Open(flamename)
		if err != nil {
			log.Fatalln("error opening input:", err)
		}
		d := xml.NewDecoder(f)
		s, err = flame.Unmarshal(d)
		if err != nil {
			log.Fatalln("error unmarshaling system:", err)
		}
	}
	if tm != (xirho.ToneMap{}) {
		s.ToneMap = tm
	}
	if intr {
		interactive(ctx, s, width, height, resampler, tm, u, osa, procs)
		return
	}
	if s == nil {
		log.Fatal("no system to render")
	}
	log.Println("allocating histogram, estimated", hist.MemFor(width*osa, height*osa)>>20, "MB")
	r := &xirho.Render{
		Hist:    xirho.NewHist(width*osa, height*osa),
		Camera:  s.Camera,
		Palette: s.Palette,
	}
	if echo {
		m, err := encoding.Marshal(s.System, r, s.ToneMap, nil, s.Meta)
		if err != nil {
			log.Fatalln("error reading system from input:", err)
		}
		log.Printf("system:\n%s\n", m)
	}
	log.Println("rendering for", timeout, "or until ^C")
	r.Render(ctx, s.System, procs)
	log.Println("finished render with", r.Iters(), "iters,", r.Hits(), "hits")
	signal.Reset(os.Interrupt) // no rendering for ^C to interrupt

	if dumpname != "" {
		dumpto(dumpname, r.Hist)
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), image.NewUniform(u), image.Point{}, draw.Src)
	log.Printf("drawing onto image of size %dx%d", width, height)
	src := r.Hist.Image(tm, r.Area(), r.Iters(), osa)
	resampler.Scale(img, img.Bounds(), src, src.Bounds(), draw.Over, nil)
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
	"lanczos1":        lanczos(1),
	"lanczos3":        lanczos(3),
	"lanczos5":        lanczos(5),
}

func lanczos(a float64) *draw.Kernel {
	return &draw.Kernel{
		Support: a,
		At: func(x float64) float64 {
			if x == 0 {
				return 1
			}
			return a * math.Sin(math.Pi*x) * math.Sin(math.Pi*x/a) / (math.Pi * math.Pi * x * x)
		},
	}
}

func dumpto(fn string, h *xirho.Hist) {
	f, err := os.Create(fn)
	if err != nil {
		log.Println("couldn't dump histogram:", err)
		return
	}
	w := bufio.NewWriter(f)
	log.Println("dumping histogram to", fn)
	n, err := h.WriteTo(w)
	if err != nil {
		log.Println("error after writing", n, "bytes:", err)
		return
	}
	if err := w.Flush(); err != nil {
		log.Println("error flushing buffer after writing", n, "bytes:", err)
		return
	}
	if err := f.Close(); err != nil {
		log.Println("error closing dump after writing", n, "bytes:", err)
	}
	log.Println("dumped", n, "bytes")
}
