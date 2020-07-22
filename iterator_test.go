package xirho

import (
	"testing"
	"time"

	"github.com/zephyrtronium/xirho/xmath"
)

type givef struct{}

func (givef) Calc(in P, rng *RNG) P {
	return P{}
}

func (givef) Prep() {}

func TestIteratorPrep(t *testing.T) {
	cases := map[string]System{
		// iterator.prep has code to handle zero functions, but System.Check
		// errs in that case, so we don't test it.
		"one": {
			Funcs: []SysFunc{
				{Func: givef{}, Opacity: 1, Weight: 1},
			},
		},
		"four": {
			Funcs: []SysFunc{
				{Func: givef{}, Opacity: 1, Weight: 1},
				{Func: givef{}, Opacity: 1, Weight: 1},
				{Func: givef{}, Opacity: 1, Weight: 1},
				{Func: givef{}, Opacity: 1, Weight: 1},
			},
		},
		"zero": {
			Funcs: []SysFunc{
				{Func: givef{}, Opacity: 1, Weight: 0},
				{Func: givef{}, Opacity: 1, Weight: 0},
				{Func: givef{}, Opacity: 1, Weight: 0},
				{Func: givef{}, Opacity: 1, Weight: 0},
			},
		},
	}
	for name, s := range cases {
		t.Run(name, func(t *testing.T) {
			it := iterator{rng: xmath.NewRNG()}
			it.prep(s)
			for i, v := range it.op {
				if v > 1<<53 {
					t.Error("opacity", i, "too large:", v, ">", 1<<53)
				}
			}
			ng := 0
			for _, v := range it.w {
				if v >= 1<<53 {
					ng++
				}
			}
			if ng == 0 {
				t.Errorf("no weights guaranteeing selection: graph is %x", it.w)
			}
		})
	}
}

func TestIteratorNext(t *testing.T) {
	// This is a Las Vegas algorithm with theoretically infinite running time
	// if our RNG is good enough (which it isn't). Skip if short.
	if testing.Short() {
		t.SkipNow()
	}
	s := System{
		Funcs: []SysFunc{
			{Func: givef{}, Weight: 1e4},
			{Func: givef{}, Weight: 1},
			{Func: givef{}, Weight: 1e-4},
			{Func: givef{}, Weight: 1e-4},
			{Func: givef{}, Weight: 1e-4},
			{Func: givef{}, Weight: 1e-4},
		},
	}
	it := iterator{rng: xmath.NewRNG()}
	it.prep(s)
	for i := 1; i < len(it.w); i++ {
		if it.w[i] == it.w[i-1] {
			t.Error("weight", i, "equals its predecessor")
		}
	}
	if t.Failed() {
		t.Fatalf("weight graph was %x", it.w)
	}
	m := make([][]bool, 0, len(s.Funcs))
	for range s.Funcs {
		m = append(m, make([]bool, len(s.Funcs)))
	}
	n := 0
	k := 0
	// The failure case here is that this loop is infinite, so the test will
	// time out. To have an explicit condition, we will instead run for 30
	// seconds; on my PC, it generally takes between 0.2 and 0.5 seconds.
	start := time.Now()
	for n < len(m) {
		j := it.next(k)
		if !m[k][j] {
			n++
			m[k][j] = true
		}
		k = j
		if time.Since(start) > 30*time.Second {
			t.Fatal("took too long; selections were", m)
		}
	}
}

func TestIteratorFinal(t *testing.T) {
	s := System{
		Funcs: []SysFunc{
			{Func: givef{}, Weight: 1},
		},
	}
	it := iterator{rng: xmath.NewRNG()}
	it.prep(s)
	p := P{1, 1, 1, 1}
	fp := it.doFinal(p)
	if fp != p {
		t.Error("point modified by missing final: want", p, "have", fp)
	}
	s.Final = givef{}
	it.prep(s)
	fp = it.doFinal(p)
	if fp != (P{}) {
		t.Error("point not modified by missing final: want", P{}, "have", fp, "with input", p)
	}
}
