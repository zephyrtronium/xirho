package xirho

import (
	"fmt"
	"image/color"
	"sync"
	"testing"
	"unsafe"
)

func TestHistMem(t *testing.T) {
	z := []int{0, 1, 10}
	for _, i := range z {
		for _, j := range z {
			w, h := 1<<i, 1<<j
			t.Run(fmt.Sprintf("%dx%d", w, h), func(t *testing.T) {
				est := uintptr(HistMem(w, h))
				hist := NewHist(w, h)
				act := unsafe.Sizeof(*hist) + uintptr(len(hist.counts))*unsafe.Sizeof(hist.counts[0])
				if est != act {
					t.Error("wrong histogram size; estimated", est, "but actual size is", act)
				}
			})
		}
	}
}

func TestHistReset(t *testing.T) {
	h := NewHist(1, 1)
	h.counts[0] = histBin{r: 1, g: 2, b: 3, n: 4}
	ar := &h.counts[0]
	h.Reset(1, 1)
	if &h.counts[0] != ar {
		t.Error("same-size reset reallocated memory")
	}
	if h.counts[0] != (histBin{}) {
		t.Error("reset failed to zero bin: have", h.counts[0])
	}
	h.Reset(1, 2)
	if &h.counts[0] == ar {
		t.Error("different-size reset failed to reallocate memory")
	}
	ar = &h.counts[0]
	h.Reset(2, 1)
	if &h.counts[0] == ar {
		t.Error("same-size different-dim reset failed to reallocate memory")
	}
}

func TestHistAdd(t *testing.T) {
	h := NewHist(1, 1)
	c := color.NRGBA64{R: 1, G: 10, B: 100, A: 1000}
	h.Add(0, 0, c)
	bin := h.counts[0]
	if bin.r != uint64(c.R) {
		t.Error("wrong red: want", uint64(c.R), "have", bin.r)
	}
	if bin.g != uint64(c.G) {
		t.Error("wrong green: want", uint64(c.G), "have", bin.g)
	}
	if bin.b != uint64(c.B) {
		t.Error("wrong blue: want", uint64(c.B), "have", bin.b)
	}
	if bin.n != uint64(c.A) {
		t.Error("wrong alpha: want", uint64(c.A), "have", bin.n)
	}
	for i := 1; i < 10; i++ {
		h.Add(0, 0, c)
	}
	bin = h.counts[0]
	if bin.r != 10*uint64(c.R) {
		t.Error("wrong red: want", 10*uint64(c.R), "have", bin.r)
	}
	if bin.g != 10*uint64(c.G) {
		t.Error("wrong green: want", 10*uint64(c.G), "have", bin.g)
	}
	if bin.b != 10*uint64(c.B) {
		t.Error("wrong blue: want", 10*uint64(c.B), "have", bin.b)
	}
	if bin.n != 10*uint64(c.A) {
		t.Error("wrong alpha: want", 10*uint64(c.A), "have", bin.n)
	}

	t.Run("concurrent", func(t *testing.T) {
		h.Reset(1, 1)
		const procs = 10
		const iters = 10000
		var wg sync.WaitGroup
		wg.Add(procs)
		for i := 0; i < procs; i++ {
			go func() {
				for i := 0; i < iters; i++ {
					h.Add(0, 0, c)
				}
				wg.Done()
			}()
		}
		wg.Wait()
		bin := h.counts[0]
		if bin.r != procs*iters*uint64(c.R) {
			t.Error("wrong red: want", procs*iters*uint64(c.R), "have", bin.r)
		}
		if bin.g != procs*iters*uint64(c.G) {
			t.Error("wrong green: want", procs*iters*uint64(c.G), "have", bin.g)
		}
		if bin.b != procs*iters*uint64(c.B) {
			t.Error("wrong blue: want", procs*iters*uint64(c.B), "have", bin.b)
		}
		if bin.n != procs*iters*uint64(c.A) {
			t.Error("wrong alpha: want", procs*iters*uint64(c.A), "have", bin.n)
		}
	})
}
