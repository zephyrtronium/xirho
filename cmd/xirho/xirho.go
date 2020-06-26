package main

import (
	"context"
	"flag"
	"fmt"
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
	"github.com/zephyrtronium/xirho/xi"
)

func main() {
	var outname, profname string
	var sigint bool
	var timeout time.Duration
	var iters, hits int64
	var width, height int
	var osa int
	var gamma float64
	var resample string
	var procs int
	flag.StringVar(&outname, "png", "", "output filename (default stdout)")
	flag.StringVar(&profname, "prof", "", "CPU profile output (default no profiling)")
	flag.BoolVar(&sigint, "C", true, "save image on interrupt instead of exiting")
	flag.DurationVar(&timeout, "dur", 0, "max duration to render (default ignored)")
	flag.Int64Var(&iters, "iters", 0, "max iters (default ignored)")
	flag.Int64Var(&hits, "hits", 0, "max hits (default iters/5)")
	flag.IntVar(&width, "width", 1024, "output image width")
	flag.IntVar(&height, "height", 1024, "output image height")
	flag.IntVar(&osa, "osa", 1, "oversampling; histogram bins per pixel per axis")
	flag.Float64Var(&gamma, "gamma", 1, "gamma factor")
	flag.StringVar(&resample, "resample", "catmull-rom", "resampling method (catmull-rom, bilinear, approx-bilinear, or nearest)")
	flag.IntVar(&procs, "procs", 0, fmt.Sprintf("concurrent render routines (default %d)", runtime.GOMAXPROCS(0)-1))
	flag.Parse()
	resampler := resamplers[resample]
	if resampler == nil {
		log.Fatal("no resampler named", resample)
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
	out := os.Stdout
	if outname != "" {
		var err error
		out, err = os.Create(outname)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
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

	F := []xirho.F{
		first(),
		second(),
		third(),
		fourth(),
	}
	system := xirho.System{
		Funcs:   F,
		Final:   final(),
		Weights: []float64{6, 8, 1, 1},
		Graph:   defaultGraph(len(F)),
	}
	cam := xirho.Ax{}
	cam.Eye().Scale(0.5, 0.5, 0)
	log.Println("allocating histogram, estimated", xirho.HistMem(width*osa, height*osa)>>20, "MB")
	r := xirho.R{
		Hist:    xirho.NewHist(width*osa, height*osa, gamma),
		System:  system,
		Camera:  cam,
		Palette: mkpalette(),
		Procs:   procs,
		N:       iters,
		Q:       hits,
	}
	log.Println("rendering up to", r.N, "iters or", r.Q, "hits or", timeout)
	r.Render(ctx)
	log.Println("finished render with", r.Iters(), "iters,", r.Hits(), "hits")
	r.Hist.Stat()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.NRGBA64{A: 0xffff}), image.Point{}, draw.Src)
	log.Printf("drawing onto image of size %dx%d", width, height)
	resampler.Scale(img, img.Bounds(), r.Hist, r.Hist.Bounds(), draw.Over, nil)
	if outname != "" {
		log.Println("encoding to", outname)
	} else {
		log.Println("encoding to stdout")
	}
	err := png.Encode(out, img)
	if err != nil {
		panic(err)
	}
}

var resamplers = map[string]draw.Scaler{
	"catmull-rom":     draw.CatmullRom,
	"bilinear":        draw.BiLinear,
	"approx-bilinear": draw.ApproxBiLinear,
	"nearest":         draw.NearestNeighbor,
}

func mkpalette() []color.NRGBA64 {
	// return []color.NRGBA64{
	// 	{R: 0xffff, G: 0x0000, B: 0x0000, A: 0xffff},
	// 	{R: 0xffff, G: 0xffff, B: 0x0000, A: 0xffff},
	// 	{R: 0x0000, G: 0xffff, B: 0x0000, A: 0xffff},
	// 	{R: 0x0000, G: 0xffff, B: 0xffff, A: 0xffff},
	// 	{R: 0x0000, G: 0x0000, B: 0xffff, A: 0xffff},
	// 	{R: 0xffff, G: 0x0000, B: 0xffff, A: 0xffff},
	// 	{R: 0xffff, G: 0x0000, B: 0x0000, A: 0xffff},
	// 	{R: 0xffff, G: 0xffff, B: 0x0000, A: 0xffff},
	// 	{R: 0x0000, G: 0xffff, B: 0x0000, A: 0xffff},
	// 	{R: 0x0000, G: 0xffff, B: 0xffff, A: 0xffff},
	// 	{R: 0x0000, G: 0x0000, B: 0xffff, A: 0xffff},
	// 	{R: 0xffff, G: 0x0000, B: 0xffff, A: 0xffff},
	// 	{R: 0xffff, G: 0x0000, B: 0x0000, A: 0xffff},
	// 	{R: 0xffff, G: 0xffff, B: 0x0000, A: 0xffff},
	// 	{R: 0x0000, G: 0xffff, B: 0x0000, A: 0xffff},
	// 	{R: 0x0000, G: 0xffff, B: 0xffff, A: 0xffff},
	// 	{R: 0x0000, G: 0x0000, B: 0xffff, A: 0xffff},
	// 	{R: 0xffff, G: 0x0000, B: 0xffff, A: 0xffff},
	// }
	r := make([]color.NRGBA64, 256)
	for i := range r {
		r[i] = color.NRGBA64{R: uint16(i * i), G: uint16(i * i), B: uint16(i * i), A: 0xffff}
	}
	return r
}

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

func first() xirho.F {
	ax := xirho.Ax{}
	ax.Eye().RotZ(-math.Pi/2).Translate(1, 0, 0)
	return xi.NewThen(xi.NewAffine(ax, 0, 0.25), xi.NewSpherical())
}

func second() xirho.F {
	ax := xirho.Ax{}
	ax.Eye().RotZ(math.Pi / 2)
	return xi.NewThen(xi.NewAffine(ax, 1, 0.75), xi.NewSpherical())
}

func third() xirho.F {
	ax := xirho.Ax{}
	ax.Eye().Translate(3, 0, 0)
	return xi.NewAffine(ax, 0.5, 0.9)
}

func fourth() xirho.F {
	ax := xirho.Ax{}
	ax.Eye().Translate(-3, 0, 0)
	return xi.NewAffine(ax, 0.5, 0.9)
}

func final() xirho.F {
	ax := xirho.Ax{}
	ax.Eye().Translate(0.5, 0, 0)
	return xi.NewThen(xi.NewAffine(ax, 0, 1), xi.NewSpherical())
}
