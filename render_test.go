package xirho_test

import (
	"context"
	"image"
	"image/color"
	"testing"
	"time"

	"github.com/zephyrtronium/xirho"
	"golang.org/x/image/draw"
)

func TestRender(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	r := xirho.Render{
		Hist:    xirho.NewHist(1, 1),
		Palette: color.Palette{color.RGBA64{R: 0xffff, A: 0xffff}, color.RGBA64{R: 0xffff, A: 0xffff}},
	}
	f := nanf{}
	s := xirho.System{
		Nodes: []xirho.Node{
			{Func: &f, Opacity: 1, Weight: 1},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	r.Render(ctx, s, 0)
	cancel()

	if !f.f {
		t.Error("system was not prepped")
	}
	if f.n == 0 {
		t.Error("no points were calculated")
	}
	if r.Iters() == 0 {
		// This could happen if the context closes during the fuse. It would be
		// nice not to call this a failure in that case, but there isn't a
		// consistent way to check.
		t.Error("calculated", f.n, "points but plotted none")
	}
	if r.Iters() != r.Hits() {
		t.Error("iters and hits should be equal, but got", r.Iters(), "iters and", r.Hits(), "hits")
	}
	tm := xirho.ToneMap{1, 1, 0}
	red, _, _, alpha := r.Hist.Image(tm, 1, r.Iters(), 1).At(0, 0).RGBA()
	if red == 0 || alpha == 0 {
		t.Error("expected solid red pixel, got red", red, "alpha", alpha)
	}
}

func TestRenderAsync(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	r := xirho.Render{
		Hist:    xirho.NewHist(1, 1),
		Palette: color.Palette{color.RGBA64{R: 0xffff, A: 0xffff}, color.RGBA64{R: 0xffff, A: 0xffff}},
	}
	f := nanf{}
	s := xirho.System{
		Nodes: []xirho.Node{
			{Func: &f, Opacity: 1, Weight: 1},
		},
	}
	change := make(chan xirho.ChangeRender, 1)
	plot := make(chan xirho.PlotOnto)
	imgs := make(chan draw.Image, 1)
	change <- xirho.ChangeRender{
		System: s,
		Procs:  4,
	}
	ctx, cancel := context.WithCancel(context.Background())
	go r.RenderAsync(ctx, change, plot, imgs)
	time.Sleep(150 * time.Millisecond)

	img := image.NewNRGBA64(image.Rect(0, 0, 1, 1))
	img.SetNRGBA64(0, 0, color.NRGBA64{A: 0xffff})
	plot <- xirho.PlotOnto{
		Image:   img,
		Scale:   draw.NearestNeighbor,
		ToneMap: xirho.ToneMap{Brightness: 1, Gamma: 1},
	}
	p, ok := <-imgs
	if !ok {
		cancel()
		t.Fatal("renderer closed imgs early")
	}
	iters := r.Iters()
	if r.Hits() == 0 {
		t.Log("note: renderer plotted no points")
	}
	red, green, blue, alpha := p.At(0, 0).RGBA()
	if red == 0 || green != 0 || blue != 0 || alpha == 0 {
		t.Errorf("expected solid red pixel, got rgba64=#%04x%04x%04x%04x", red, green, blue, alpha)
	}
	time.Sleep(150 * time.Millisecond)
	if iters == r.Iters() {
		t.Error("renderer did not continue after", iters, "iters")
	}
	// TODO: many other things to test: changing render, pause/resume, coalescing ops, ...
	cancel()
	// Make sure the renderer closes channels after the context cancels.
	for range imgs {
		// do nothing
	}
}
