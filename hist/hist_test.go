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

func TestHistReset(t *testing.T) {
	h := New(Size{W: 1, H: 1, OSA: 1})
	h.counts[0] = bin{r: 1, g: 2, b: 3, n: 4}
	ar := &h.counts[0]
	h.Reset(Size{W: 1, H: 1, OSA: 1})
	if &h.counts[0] != ar {
		t.Error("same-size reset reallocated memory")
	}
	if h.counts[0] != (bin{}) {
		t.Error("reset failed to zero bin: have", h.counts[0])
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
