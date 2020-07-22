package xirho_test

import (
	"context"
	"image/color"
	"math"
	"testing"
	"time"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xmath"
)

type prepf bool

func (v *prepf) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	return in
}

func (v *prepf) Prep() {
	*v = true
}

func TestSystemPrep(t *testing.T) {
	s := xirho.System{
		Funcs: []xirho.SysFunc{
			{Func: new(prepf), Weight: 1},
			{Func: new(prepf), Weight: 1},
			{Func: new(prepf), Weight: 1},
		},
		Final: new(prepf),
	}
	s.Prep()
	for i, f := range s.Funcs {
		if !*f.Func.(*prepf) {
			t.Error("function", i, "not prepped")
		}
	}
	if !*s.Final.(*prepf) {
		t.Error("final not prepped")
	}
}

func TestSystemCheck(t *testing.T) {
	// First, check that a well-defined system passes.
	s := xirho.System{
		Funcs: []xirho.SysFunc{
			{Func: new(prepf), Weight: 1, Graph: []float64{1}},
		},
	}
	if err := s.Check(); err != nil {
		t.Error("system", s, "gave unexpected check error", err)
	}
	cases := map[string]xirho.System{
		"empty":          {},
		"emptyWithFinal": {Final: new(prepf)},
		"opacityNegative": {
			Funcs: []xirho.SysFunc{
				{Func: new(prepf), Opacity: math.Nextafter(0, -1), Weight: 1},
			},
		},
		"opacityExcess": {
			Funcs: []xirho.SysFunc{
				{Func: new(prepf), Opacity: math.Nextafter(1, 2), Weight: 1},
			},
		},
		"opacityNan": {
			Funcs: []xirho.SysFunc{
				{Func: new(prepf), Opacity: math.NaN(), Weight: 1},
			},
		},
		"weightNegative": {
			Funcs: []xirho.SysFunc{
				{Func: new(prepf), Weight: math.Nextafter(0, -1)},
			},
		},
		"weightInf": {
			Funcs: []xirho.SysFunc{
				{Func: new(prepf), Weight: math.Inf(0)},
			},
		},
		"weightNan": {
			Funcs: []xirho.SysFunc{
				{Func: new(prepf), Weight: math.NaN()},
			},
		},
		"graphNegative": {
			Funcs: []xirho.SysFunc{
				{Func: new(prepf), Weight: 1, Graph: []float64{math.Nextafter(0, -1)}},
			},
		},
		"graphInf": {
			Funcs: []xirho.SysFunc{
				{Func: new(prepf), Weight: 1, Graph: []float64{math.Inf(0)}},
			},
		},
		"graphNan": {
			Funcs: []xirho.SysFunc{
				{Func: new(prepf), Weight: 1, Graph: []float64{math.NaN()}},
			},
		},
	}
	for name, s := range cases {
		t.Run(name, func(t *testing.T) {
			if err := s.Check(); err == nil {
				t.Error("system", s, "did not give check error")
			}
		})
	}
}

type nanf struct {
	p int
	n int64
}

func (v *nanf) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	v.n++
	switch v.p {
	case 1:
		return xirho.P{
			X: math.NaN(),
			Y: in.Y,
			Z: in.Z,
			C: in.C,
		}
	case 2:
		return xirho.P{
			X: in.X,
			Y: in.Y,
			Z: in.Z,
			C: 2,
		}
	}
	return xirho.P{}
}

func (v *nanf) Prep() {}

func TestSystemIter(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	rng := xmath.NewRNG()
	f := nanf{}
	s := xirho.System{
		Funcs: []xirho.SysFunc{
			{Func: &f, Weight: 1, Opacity: 1},
		},
	}
	r := xirho.R{
		Camera:  xirho.Eye(),
		Hist:    xirho.NewHist(1, 1),
		Palette: []color.NRGBA64{{R: 0xffff, A: 0xffff}, {R: 0xffff, A: 0xffff}},
		Procs:   1,
	}
	r.Hist.Reset(1, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	s.Iter(ctx, &r, rng)
	cancel()
	red, _, _, alpha := r.Hist.At(0, 0).RGBA()
	if red == 0 || alpha == 0 {
		t.Error("expected red pixel, got red", red, "alpha", alpha, "with hist", r.Hist, "after", r.Iters(), "iters")
	}
	t.Run("check", func(t *testing.T) {
		s := xirho.System{}
		r := xirho.R{
			Camera:  xirho.Eye(),
			Hist:    xirho.NewHist(1, 1),
			Palette: []color.NRGBA64{{R: 0xffff, A: 0xffff}, {R: 0xffff, A: 0xffff}},
			Procs:   1,
		}
		r.Hist.Reset(1, 1)
		if err := s.Check(); err == nil {
			t.Error("empty system did not give check error")
		}
		defer func() {
			err := recover()
			if err == nil {
				t.Error("iter did not panic")
			}
		}()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		s.Iter(ctx, &r, rng)
	})
	t.Run("fuseSpace", func(t *testing.T) {
		f := nanf{p: 1}
		s := xirho.System{
			Funcs: []xirho.SysFunc{
				{Func: &f, Weight: 1},
			},
		}
		r := xirho.R{
			Camera:  xirho.Eye(),
			Hist:    xirho.NewHist(1, 1),
			Palette: []color.NRGBA64{{R: 0xffff, A: 0xffff}, {R: 0xffff, A: 0xffff}},
			Procs:   1,
		}
		r.Hist.Reset(1, 1)
		ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
		s.Iter(ctx, &r, rng)
		cancel()
		if f.n == 0 {
			t.Log("vacuous condition: no points were calculated")
		}
		if r.Iters() != 0 {
			t.Error("always-invalid function was plotted", r.Iters(), "times of", f.n, "calcs")
		}
	})
	t.Run("fuseColor", func(t *testing.T) {
		f := nanf{p: 2}
		r := xirho.R{
			Camera:  xirho.Eye(),
			Hist:    xirho.NewHist(1, 1),
			Palette: []color.NRGBA64{{R: 0xffff, A: 0xffff}, {R: 0xffff, A: 0xffff}},
			Procs:   1,
		}
		s := xirho.System{
			Funcs: []xirho.SysFunc{
				{Func: &f, Weight: 1},
			},
		}
		r.Hist.Reset(1, 1)
		ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
		s.Iter(ctx, &r, rng)
		cancel()
		if f.n == 0 {
			t.Log("vacuous condition: no points were calculated")
		}
		if r.Iters() != 0 {
			t.Error("always-invalid function was plotted", r.Iters(), "times of", f.n, "calcs")
		}
	})
}
