package xirho

import (
	"context"

	"github.com/zephyrtronium/crazy"
)

// Note that Xoshiro might cause a loss in quality if there are more than
// sixteen functions in a System, due to the exact properties that make it
// suitable for use by multiple goroutines. If this ends up being problematic,
// a better choice might be crazy.MT64.

// RNG is the randomness source type.
type RNG = crazy.Xoshiro

// Note to maintainers: iterator.next is tied to the value of MaxFuncs.

// MaxFuncs is the maximum number of unique functions that a system may hold.
const MaxFuncs = 65536

// System is a generalized iterated function system.
type System struct {
	// Funcs is the system's function list.
	Funcs []F
	// TODO: weights, xaos, &c
}

// iterator manages the iterations of a System by a single goroutine.
type iterator struct {
	System
	// rng is the iterator's source of randomness.
	rng RNG
	// n is the number of iterations this iterator has performed.
	n uint64
	// bk is a buffer to turn 64-bit RNG outputs into 16-bit ones to improve
	// dimensional equidistribution.
	bk uint64
}

// Iter iterates the function system and sends output points over results. It
// continues iterating until the context's Done channel is closed. rng should
// be seeded to a distinct state for each call to this method.
func (s System) Iter(ctx context.Context, results chan<- P, rng RNG) {
	it := iterator{System: s, rng: rng}
	p, k := it.fuse() // p may not be valid!
	done := ctx.Done()
	for {
		select {
		case <-done:
			return
		case results <- p:
			if !p.IsValid() {
				p, k = it.fuse()
				continue
			}
			p = it.Funcs[k].Calc(p, &it.rng)
			k = it.next(k)
		}
	}
}

// fuseLen is the number of iterations to perform before beginning to plot.
const fuseLen = 30

// fuse obtains initial conditions to plot points from the system.
func (it *iterator) fuse() (P, int) {
	d := crazy.Uniform{Source: &it.rng, Low: -1, High: 1}
	p := P{
		X: d.Next(),
		Y: d.Next(),
		Z: d.Next(),
		C: crazy.Uniform0_1{Source: &it.rng}.Next(),
	}
	k := it.next(crazy.RNG{Source: &it.rng}.Intn(len(it.Funcs)))
	for i := 0; i < fuseLen+1; i++ {
		p = it.Funcs[k].Calc(p, &it.rng)
		if !p.IsValid() {
			break
		}
		k = it.next(k)
	}
	return p, k
}

// next obtains the next function to use from the current one.
func (it *iterator) next(k int) int {
	// TODO: weights, xaos
	good := 0xffff - 0xffff%len(it.Funcs)
	for {
		if it.n%16 == 0 {
			it.bk = it.rng.Uint64()
		} else {
			it.bk >>= 16
		}
		k = int(it.bk & 0xffff)
		if k <= good {
			break
		}
	}
	return k % len(it.Funcs)
}

// newRNG creates a new seeded RNG instance.
func newRNG() RNG {
	rng := RNG{}
	crazy.CryptoSeeded(&rng, 8)
	return rng
}
