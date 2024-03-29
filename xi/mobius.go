package xi

import (
	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xmath"
)

// Mobius implements Mobius transformations over quaternions.
type Mobius struct {
	Ar   float64    `xirho:"A.scalar"`
	Avec [3]float64 `xirho:"A.vector"`
	Br   float64    `xirho:"B.scalar"`
	Bvec [3]float64 `xirho:"B.vector"`
	Cr   float64    `xirho:"C.scalar"`
	Cvec [3]float64 `xirho:"C.vector"`
	Dr   float64    `xirho:"D.scalar"`
	Dvec [3]float64 `xirho:"D.vector"`

	InZero int `xirho:"input blank,r,i,j,k"`
}

func newMobius() xirho.Func {
	return &Mobius{
		Ar:     1,
		Dr:     1,
		InZero: 3,
	}
}

func (v *Mobius) Calc(in xirho.Pt, rng *xmath.RNG) xirho.Pt {
	// May the compiler bless me with optimized quaternion operations.
	var nr, ni, nj, nk, dr, di, dj, dk float64
	switch v.InZero {
	case 0:
		// input is 0 + in.X*i + in.Y*j + in.Z*k
		nr = v.Br - v.Avec[0]*in.X - v.Avec[1]*in.Y - v.Avec[2]*in.Z
		ni = v.Bvec[0] + v.Ar*in.X + v.Avec[1]*in.Z - v.Avec[2]*in.Y
		nj = v.Bvec[1] + v.Ar*in.Y - v.Avec[0]*in.Z + v.Avec[2]*in.X
		nk = v.Bvec[2] + v.Ar*in.Z + v.Avec[0]*in.Y - v.Avec[1]*in.X
		dr = v.Dr - v.Cvec[0]*in.X - v.Cvec[1]*in.Y - v.Cvec[2]*in.Z
		di = v.Dvec[0] + v.Cr*in.X + v.Cvec[1]*in.Z - v.Cvec[2]*in.Y
		dj = v.Dvec[1] + v.Cr*in.Y - v.Cvec[0]*in.Z + v.Cvec[2]*in.X
		dk = v.Dvec[2] + v.Cr*in.Z + v.Cvec[0]*in.Y - v.Cvec[1]*in.X
	case 1:
		// input is in.X + 0*i + in.Y*j + in.Z*k
		nr = v.Br + v.Ar*in.X - v.Avec[1]*in.Y - v.Avec[2]*in.Z
		ni = v.Bvec[0] + v.Avec[0]*in.X + v.Avec[1]*in.Z - v.Avec[2]*in.Y
		nj = v.Bvec[1] + v.Ar*in.Y - v.Avec[0]*in.Z + v.Avec[1]*in.X
		nk = v.Bvec[2] + v.Ar*in.Z + v.Avec[0]*in.Y + v.Avec[2]*in.X
		dr = v.Dr + v.Cr*in.X - v.Cvec[1]*in.Y - v.Cvec[2]*in.Z
		di = v.Dvec[0] + v.Cvec[0]*in.X + v.Cvec[1]*in.Z - v.Cvec[2]*in.Y
		dj = v.Dvec[1] + v.Cr*in.Y - v.Cvec[0]*in.Z + v.Cvec[1]*in.X
		dk = v.Dvec[2] + v.Cr*in.Z + v.Cvec[0]*in.Y + v.Cvec[2]*in.X
	case 2:
		// input is in.X + in.Y*i + 0*j + in.Z*k
		nr = v.Br + v.Ar*in.X - v.Avec[0]*in.Y - v.Avec[2]*in.Z
		ni = v.Bvec[0] + v.Ar*in.Y + v.Avec[0]*in.X + v.Avec[1]*in.Z
		nj = v.Bvec[1] - v.Avec[0]*in.Z + v.Avec[1]*in.X + v.Avec[2]*in.Y
		nk = v.Bvec[2] + v.Ar*in.Z - v.Avec[1]*in.Y + v.Avec[2]*in.X
		dr = v.Dr + v.Cr*in.X - v.Cvec[0]*in.Y - v.Cvec[2]*in.Z
		di = v.Dvec[0] + v.Cr*in.Y + v.Cvec[0]*in.X + v.Cvec[1]*in.Z
		dj = v.Dvec[1] - v.Cvec[0]*in.Z + v.Cvec[1]*in.X + v.Cvec[2]*in.Y
		dk = v.Dvec[2] + v.Cr*in.Z - v.Cvec[1]*in.Y + v.Cvec[2]*in.X
	case 3:
		// input is in.X + in.Y*i + in.Z*j + 0*k
		nr = v.Br + v.Ar*in.X - v.Avec[0]*in.Y - v.Avec[1]*in.Z
		ni = v.Bvec[0] + v.Ar*in.Y + v.Avec[0]*in.X - v.Avec[2]*in.Z
		nj = v.Bvec[1] + v.Ar*in.Z + v.Avec[1]*in.X + v.Avec[2]*in.Y
		nk = v.Bvec[2] + v.Avec[0]*in.Z - v.Avec[1]*in.Y + v.Avec[2]*in.X
		dr = v.Dr + v.Cr*in.X - v.Cvec[0]*in.Y - v.Cvec[1]*in.Z
		di = v.Dvec[0] + v.Cr*in.Y + v.Cvec[0]*in.X - v.Cvec[2]*in.Z
		dj = v.Dvec[1] + v.Cr*in.Z + v.Cvec[1]*in.X + v.Cvec[2]*in.Y
		dk = v.Dvec[2] + v.Cvec[0]*in.Z - v.Cvec[1]*in.Y + v.Cvec[2]*in.X
	}
	rr := dr*dr + di*di + dj*dj + dk*dk
	dr /= rr
	di /= -rr
	dj /= -rr
	dk /= -rr
	in.X = nr*dr - ni*di - nj*dj - nk*dk
	in.Y = nr*di + ni*dr + nj*dk - nk*dj
	in.Z = nr*dj - ni*dk + nj*dr + nk*di
	// outk := nr*dk + ni*dj - nj*di + nk*dr
	return in
}

func (v *Mobius) Prep() {}

func init() {
	must("mobius", newMobius)
	must("mobiq", newMobius)
}
