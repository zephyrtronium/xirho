package hist

import (
	"fmt"
	"image/color"
	"sync"
	"testing"
	"unsafe"
)

func TestNewZero(t *testing.T) {
	cases := []Size{
		{W: 0, H: 0, OSA: 1},
		{W: 1, H: 0, OSA: 1},
		{W: 0, H: 1, OSA: 1},
		{W: 1, H: 1, OSA: 0},
	}
	for _, c := range cases {
		h := New(c)
		if h.arr != nil {
			t.Errorf("non-nil array pointer %p after allocating %dx%d:%d hist", h.arr, c.W, c.H, c.OSA)
		}
		if !h.IsEmpty() {
			t.Errorf("non-empty after allocating %dx%d:%d hist", c.W, c.H, c.OSA)
		}
	}
}

func TestNew(t *testing.T) {
	z := []int{1, 10}
	for _, i := range z {
		for _, j := range z {
			w, h := 1<<i, 1<<j
			t.Run(fmt.Sprintf("%dx%d", w, h), func(t *testing.T) {
				h := New(Size{W: w, H: h, OSA: 1})
				if h.arr != unsafe.Pointer(&h.counts[0]) {
					t.Errorf("wrong array pointer after allocating %dx%d hist: want %p, got %p", w, h, h.arr, &h.counts[0])
				}
				if h.IsEmpty() {
					t.Errorf("empty after allocating %dx%d hist", w, h)
				}
			})
		}
	}
}

func TestMemFor(t *testing.T) {
	z := []int{0, 1, 10}
	for _, i := range z {
		for _, j := range z {
			w, h := 1<<i, 1<<j
			t.Run(fmt.Sprintf("%dx%d", w, h), func(t *testing.T) {
				est := uintptr(MemFor(w, h))
				hist := New(Size{W: w, H: h, OSA: 1})
				act := unsafe.Sizeof(*hist) + uintptr(len(hist.counts))*unsafe.Sizeof(hist.counts[0])
				if est != act {
					t.Error("wrong histogram size; estimated", est, "but actual size is", act)
				}
			})
		}
	}
}

func zerobin(bin *bin) bool {
	return bin.r.Load() == 0 && bin.g.Load() == 0 && bin.b.Load() == 0 && bin.n.Load() == 0
}

func TestHistReset(t *testing.T) {
	h := New(Size{W: 1, H: 1, OSA: 1})
	h.counts[0].r.Store(1)
	h.counts[0].g.Store(2)
	h.counts[0].b.Store(3)
	h.counts[0].n.Store(4)
	ar := &h.counts[0]
	h.Reset(Size{W: 1, H: 1, OSA: 1})
	if &h.counts[0] != ar {
		t.Error("same-size reset reallocated memory")
	}
	if !zerobin(&h.counts[0]) {
		bin := &h.counts[0]
		t.Error("reset failed to zero bin: have", bin.r.Load(), bin.g.Load(), bin.b.Load(), bin.n.Load())
	}
	h.Reset(Size{W: 1, H: 2, OSA: 1})
	if &h.counts[0] == ar {
		t.Error("different-size reset failed to reallocate memory")
	}
	ar = &h.counts[0]
	h.Reset(Size{W: 2, H: 1, OSA: 1})
	if &h.counts[0] != ar {
		t.Error("same-size different-dim reset reallocated memory")
	}

	h.Reset(Size{W: 2, H: 2, OSA: 4})
	if len(h.counts) != 2*2*4*4 {
		t.Errorf("wrong number of bins after setting osa: want %d, got %d", 2*2*4*4, len(h.counts))
	}
}

func TestHistAdd(t *testing.T) {
	h := New(Size{W: 1, H: 1, OSA: 1})
	c := color.RGBA64{R: 1, G: 10, B: 100, A: 1000}
	h.Add(0, 0, c)
	bin := &h.counts[0]
	if bin.r.Load() != uint64(c.R) {
		t.Error("wrong red: want", uint64(c.R), "have", bin.r.Load())
	}
	if bin.g.Load() != uint64(c.G) {
		t.Error("wrong green: want", uint64(c.G), "have", bin.g.Load())
	}
	if bin.b.Load() != uint64(c.B) {
		t.Error("wrong blue: want", uint64(c.B), "have", bin.b.Load())
	}
	if bin.n.Load() != uint64(c.A) {
		t.Error("wrong alpha: want", uint64(c.A), "have", bin.n.Load())
	}
	for i := 1; i < 10; i++ {
		h.Add(0, 0, c)
	}
	bin = &h.counts[0]
	if bin.r.Load() != 10*uint64(c.R) {
		t.Error("wrong red: want", 10*uint64(c.R), "have", bin.r.Load())
	}
	if bin.g.Load() != 10*uint64(c.G) {
		t.Error("wrong green: want", 10*uint64(c.G), "have", bin.g.Load())
	}
	if bin.b.Load() != 10*uint64(c.B) {
		t.Error("wrong blue: want", 10*uint64(c.B), "have", bin.b.Load())
	}
	if bin.n.Load() != 10*uint64(c.A) {
		t.Error("wrong alpha: want", 10*uint64(c.A), "have", bin.n.Load())
	}

	t.Run("concurrent", func(t *testing.T) {
		h.Reset(Size{W: 1, H: 1, OSA: 1})
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
		bin := &h.counts[0]
		if bin.r.Load() != procs*iters*uint64(c.R) {
			t.Error("wrong red: want", procs*iters*uint64(c.R), "have", bin.r.Load())
		}
		if bin.g.Load() != procs*iters*uint64(c.G) {
			t.Error("wrong green: want", procs*iters*uint64(c.G), "have", bin.g.Load())
		}
		if bin.b.Load() != procs*iters*uint64(c.B) {
			t.Error("wrong blue: want", procs*iters*uint64(c.B), "have", bin.b.Load())
		}
		if bin.n.Load() != procs*iters*uint64(c.A) {
			t.Error("wrong alpha: want", procs*iters*uint64(c.A), "have", bin.n.Load())
		}
	})
}
