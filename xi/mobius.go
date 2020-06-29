package xi

import "github.com/zephyrtronium/xirho"

// Mobius implements Mobius transformations over quaternions.
type Mobius struct {
	Ar   xirho.Real `xirho:"A.scalar"`
	Avec xirho.Vec3 `xirho:"A.vector"`
	Br   xirho.Real `xirho:"B.scalar"`
	Bvec xirho.Vec3 `xirho:"B.vector"`
	Cr   xirho.Real `xirho:"C.scalar"`
	Cvec xirho.Vec3 `xirho:"C.vector"`
	Dr   xirho.Real `xirho:"D.scalar"`
	Dvec xirho.Vec3 `xirho:"D.vector"`

	InZero xirho.List `xirho:"input blank,r,i,j,k"`
}

func newMobius() xirho.F {
	return &Mobius{
		Ar:     1,
		Dr:     1,
		InZero: 3,
	}
}

func (v *Mobius) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	// May the compiler bless me with optimized quaternion operations.
	var nr, ni, nj, nk, dr, di, dj, dk float64
	switch v.InZero {
	case 0:
		// input is 0 + in.X*i + in.Y*j + in.Z*k
		nr = float64(v.Br) - v.Avec[0]*in.X - v.Avec[1]*in.Y - v.Avec[2]*in.Z
		ni = v.Bvec[0] + float64(v.Ar)*in.X + v.Avec[1]*in.Z - v.Avec[2]*in.Y
		nj = v.Bvec[1] + float64(v.Ar)*in.Y - v.Avec[0]*in.Z + v.Avec[2]*in.X
		nk = v.Bvec[2] + float64(v.Ar)*in.Z + v.Avec[0]*in.Y - v.Avec[1]*in.X
		dr = float64(v.Dr) - v.Cvec[0]*in.X - v.Cvec[1]*in.Y - v.Cvec[2]*in.Z
		di = v.Dvec[0] + float64(v.Cr)*in.X + v.Cvec[1]*in.Z - v.Cvec[2]*in.Y
		dj = v.Dvec[1] + float64(v.Cr)*in.Y - v.Cvec[0]*in.Z + v.Cvec[2]*in.X
		dk = v.Dvec[2] + float64(v.Cr)*in.Z + v.Cvec[0]*in.Y - v.Cvec[1]*in.X
	case 1:
		// input is in.X + 0*i + in.Y*j + in.Z*k
		nr = float64(v.Br) + float64(v.Ar)*in.X - v.Avec[1]*in.Y - v.Avec[2]*in.Z
		ni = v.Bvec[0] + v.Avec[0]*in.X + v.Avec[1]*in.Z - v.Avec[2]*in.Y
		nj = v.Bvec[1] + float64(v.Ar)*in.Y - v.Avec[0]*in.Z + v.Avec[1]*in.X
		nk = v.Bvec[2] + float64(v.Ar)*in.Z + v.Avec[0]*in.Y + v.Avec[2]*in.X
		dr = float64(v.Dr) + float64(v.Cr)*in.X - v.Cvec[1]*in.Y - v.Cvec[2]*in.Z
		di = v.Dvec[0] + v.Cvec[0]*in.X + v.Cvec[1]*in.Z - v.Cvec[2]*in.Y
		dj = v.Dvec[1] + float64(v.Cr)*in.Y - v.Cvec[0]*in.Z + v.Cvec[1]*in.X
		dk = v.Dvec[2] + float64(v.Cr)*in.Z + v.Cvec[0]*in.Y + v.Cvec[2]*in.X
	case 2:
		// input is in.X + in.Y*i + 0*j + in.Z*k
		nr = float64(v.Br) + float64(v.Ar)*in.X - v.Avec[0]*in.Y - v.Avec[2]*in.Z
		ni = v.Bvec[0] + float64(v.Ar)*in.Y + v.Avec[0]*in.X + v.Avec[1]*in.Z
		nj = v.Bvec[1] - v.Avec[0]*in.Z + v.Avec[1]*in.X + v.Avec[2]*in.Y
		nk = v.Bvec[2] + float64(v.Ar)*in.Z - v.Avec[1]*in.Y + v.Avec[2]*in.X
		dr = float64(v.Dr) + float64(v.Cr)*in.X - v.Cvec[0]*in.Y - v.Cvec[2]*in.Z
		di = v.Dvec[0] + float64(v.Cr)*in.Y + v.Cvec[0]*in.X + v.Cvec[1]*in.Z
		dj = v.Dvec[1] - v.Cvec[0]*in.Z + v.Cvec[1]*in.X + v.Cvec[2]*in.Y
		dk = v.Dvec[2] + float64(v.Cr)*in.Z - v.Cvec[1]*in.Y + v.Cvec[2]*in.X
	case 3:
		// input is in.X + in.Y*i + in.Z*j + 0*k
		nr = float64(v.Br) + float64(v.Ar)*in.X - v.Avec[0]*in.Y - v.Avec[1]*in.Z
		ni = v.Bvec[0] + float64(v.Ar)*in.Y + v.Avec[0]*in.X - v.Avec[2]*in.Z
		nj = v.Bvec[1] + float64(v.Ar)*in.Z + v.Avec[1]*in.X + v.Avec[2]*in.Y
		nk = v.Bvec[2] + v.Avec[0]*in.Z - v.Avec[1]*in.Y + v.Avec[2]*in.X
		dr = float64(v.Dr) + float64(v.Cr)*in.X - v.Cvec[0]*in.Y - v.Cvec[1]*in.Z
		di = v.Dvec[0] + float64(v.Cr)*in.Y + v.Cvec[0]*in.X - v.Cvec[2]*in.Z
		dj = v.Dvec[1] + float64(v.Cr)*in.Z + v.Cvec[1]*in.X + v.Cvec[2]*in.Y
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
	Register("mobius", newMobius)
	Register("mobiq", newMobius)
}
