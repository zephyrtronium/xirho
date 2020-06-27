package main

import (
	"context"
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

	"github.com/dim13/colormap"
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
	flag.IntVar(&procs, "procs", runtime.GOMAXPROCS(0), "concurrent render routines")
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

	system := params()
	cam := xirho.Ax{}
	cam.Eye().Scale(0.6, 0.6, 0)
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
	signal.Reset(os.Interrupt) // no rendering for ^C to interrupt
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

// ---- disc julian params ----

// func params() xirho.System {
// 	return xirho.System{
// 		Funcs:   []xirho.F{first(), second()},
// 		Final:   final(),
// 		Weights: []float64{30, 1},
// 		Graph:   defaultGraph(2),
// 	}
// }

// func first() xirho.F {
// 	ax := xirho.Ax{}
// 	ax.Eye().RotZ(-math.Pi/6).Scale(0.9, 0.9, 0)
// 	return &xi.Then{
// 		Funcs: xirho.FuncList{
// 			&xi.Affine{Ax: ax},
// 			&xi.ColorSpeed{Color: 0, Speed: 0.75},
// 			xi.Disc{},
// 		},
// 	}
// }

// func second() xirho.F {
// 	return &xi.Then{
// 		Funcs: xirho.FuncList{
// 			&xi.JuliaN{Power: 30, Dist: -1},
// 			&xi.ColorSpeed{Color: 1, Speed: 0.3},
// 		},
// 	}
// }

// func final() xirho.F {
// 	ax := xirho.Ax{}
// 	ax.Eye().Scale(3, 3, 0)
// 	ay := xirho.Ax{}
// 	ay.Eye().RotZ(math.Pi/4).Scale(2.5, 2.5, 0)
// 	return &xi.Then{
// 		Funcs: xirho.FuncList{
// 			&xi.Affine{Ax: ax},
// 			xi.Polar{},
// 			&xi.Affine{Ax: ay},
// 		},
// 	}
// }

// ---- spherical gasket params ----

// func params() xirho.System {
// 	return xirho.System{
// 		Funcs:   []xirho.F{first(), second(), third(), fourth()},
// 		Final:   final(),
// 		Weights: []float64{6, 8, 1, 1},
// 		Graph:   defaultGraph(4),
// 	}
// }

// func first() xirho.F {
// 	ax := xirho.Ax{}
// 	ax.Eye().RotZ(-math.Pi/2).Translate(1, 0, 0)
// 	return then(
// 		&xi.Affine{Ax: ax},
// 		&xi.ColorSpeed{Color: 0, Speed: 0.25},
// 		xi.Spherical{})
// }

// func second() xirho.F {
// 	ax := xirho.Ax{}
// 	ax.Eye().RotZ(math.Pi / 2)
// 	return then(
// 		&xi.Affine{Ax: ax},
// 		&xi.ColorSpeed{Color: 1, Speed: 0.75},
// 		xi.Spherical{})
// }

// func third() xirho.F {
// 	ax := xirho.Ax{}
// 	ax.Eye().Translate(3, 0, 0)
// 	return affspd(ax, 0.5, 0.9)
// }

// func fourth() xirho.F {
// 	ax := xirho.Ax{}
// 	ax.Eye().Translate(-3, 0, 0)
// 	return affspd(ax, 0.5, 0.9)
// }

// func final() xirho.F {
// 	ax := xirho.Ax{}
// 	ax.Eye().Translate(0.5, 0, 0)
// 	return then(&xi.Affine{Ax: ax}, xi.Spherical{})
// }

// ---- grand julian params ----

func params() xirho.System {
	return xirho.System{
		Funcs:   []xirho.F{first(), second(), third(), fourth()},
		Final:   final(),
		Weights: []float64{1, 12, 2, 2},
		Graph:   defaultGraph(4),
	}
}

func first() xirho.F {
	ax := xirho.Ax{}
	ax.Eye().Scale(10, 10, 10)
	ay := xirho.Ax{}
	ay.Eye().Scale(0.185, 0.185, 0.185)
	return then(xi.Blur{}, &xi.Affine{Ax: ax}, xi.Bubble{}, &xi.Affine{Ax: ay}, &xi.ColorSpeed{Color: 0.5, Speed: 0.5})
}

func second() xirho.F {
	ax := xirho.Ax{}
	ax.Eye().RotZ(math.Pi/4).Scale(1, 1, 0).Translate(0, 0.3, 0)
	return then(affspd(ax, 0, 0.75), &xi.JuliaN{Power: 2, Dist: -1})
}

func third() xirho.F {
	ax := xirho.Ax{}
	ax.Eye().RotZ(math.Pi / 4)
	ay := xirho.Ax{}
	ay.Eye().Scale(0.2, 0.2, 0)
	return then(&xi.Affine{Ax: ax}, &xi.JuliaN{Power: 15, Dist: -1}, affspd(ay, 1, 0.8))
}

func fourth() xirho.F {
	ax := xirho.Ax{}
	ax.Eye().RotZ(math.Pi / 4)
	ay := xirho.Ax{}
	ay.Eye().Scale(0.3, 0.3, 0)
	return then(&xi.Affine{Ax: ax}, &xi.JuliaN{Power: 8, Dist: -1}, affspd(ay, 1, 0.8))
}

func final() xirho.F {
	ax := xirho.Ax{}
	ax.Eye().Translate(0.5, 0, 0)
	return then(&xi.Affine{Ax: ax}, xi.Spherical{})
}
