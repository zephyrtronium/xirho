package xirho

import "math"

// Ax is an affine transform.
type Ax struct {
	// A is the linear (multiplicative) component of the transform.
	A [3][3]float64
	// B is the translation (additive) component of the transform.
	B [3]float64
}

// Eye sets the transform to the identity transform and returns it.
func (ax *Ax) Eye() *Ax {
	*ax = Ax{
		A: [3][3]float64{
			{1, 0, 0},
			{0, 1, 0},
			{0, 0, 1},
		},
	}
	return ax
}

// Translate translates the transform and returns it.
func (ax *Ax) Translate(dx, dy, dz float64) *Ax {
	ax.B = [3]float64{
		ax.B[0] + dx,
		ax.B[1] + dy,
		ax.B[2] + dz,
	}
	return ax
}

// Scale scales the transform along each axis and returns it.
func (ax *Ax) Scale(sx, sy, sz float64) *Ax {
	ax.A = [3][3]float64{
		{ax.A[0][0] * sx, ax.A[0][1] * sy, ax.A[0][2] * sz},
		{ax.A[1][0] * sx, ax.A[1][1] * sy, ax.A[1][2] * sz},
		{ax.A[2][0] * sx, ax.A[2][1] * sy, ax.A[2][2] * sz},
	}
	return ax
}

// RotX rotates the transform about the x axis and returns it. Specifically, tx
// is the rotation in the y/z plane in radians.
func (ax *Ax) RotX(tx float64) *Ax {
	sx, cx := math.Sincos(tx)
	ax.A = [3][3]float64{
		{ax.A[0][0], ax.A[0][1]*cx - ax.A[0][2]*sx, ax.A[0][1]*sx + ax.A[0][2]*cx},
		{ax.A[1][0], ax.A[1][1]*cx - ax.A[1][2]*sx, ax.A[1][1]*sx + ax.A[1][2]*cx},
		{ax.A[2][0], ax.A[2][1]*cx - ax.A[2][2]*sx, ax.A[2][1]*sx + ax.A[2][2]*cx},
	}
	return ax
}

// RotY rotates the transform about the y axis and returns it. Specifically, ty
// is the rotation in the x/z plane in radians.
func (ax *Ax) RotY(ty float64) *Ax {
	sy, cy := math.Sincos(ty)
	ax.A = [3][3]float64{
		{ax.A[0][0]*cy + ax.A[0][2]*sy, ax.A[0][1], ax.A[0][2]*cy - ax.A[0][0]*sy},
		{ax.A[1][0]*cy + ax.A[1][2]*sy, ax.A[1][1], ax.A[1][2]*cy - ax.A[1][0]*sy},
		{ax.A[2][0]*cy + ax.A[2][2]*sy, ax.A[2][1], ax.A[2][2]*cy - ax.A[2][0]*sy},
	}
	return ax
}

// RotZ rotates the transform about the z axis and returns it. Specifically, tz
// is the rotation in the x/y plane in radians.
func (ax *Ax) RotZ(tz float64) *Ax {
	sz, cz := math.Sincos(tz)
	ax.A = [3][3]float64{
		{ax.A[0][0]*cz - ax.A[0][1]*sz, ax.A[0][0]*sz + ax.A[0][1]*cz, ax.A[0][2]},
		{ax.A[1][0]*cz - ax.A[1][1]*sz, ax.A[1][0]*sz + ax.A[1][1]*cz, ax.A[1][2]},
		{ax.A[2][0]*cz - ax.A[2][1]*sz, ax.A[2][0]*sz + ax.A[2][1]*cz, ax.A[2][2]},
	}
	return ax
}

// TODO: shear

// Tx transforms a coordinate with an affine transform.
func Tx(ax *Ax, x, y, z float64) (tx, ty, tz float64) {
	tx = ax.A[0][0]*x + ax.A[0][1]*y + ax.A[0][2]*z + ax.B[0]
	ty = ax.A[1][0]*x + ax.A[1][1]*y + ax.A[1][2]*z + ax.B[1]
	tz = ax.A[2][0]*x + ax.A[2][1]*y + ax.A[2][2]*z + ax.B[2]
	return
}

// TxVec transforms a 3-vector with an affine transform.
func TxVec(ax *Ax, v [3]float64) [3]float64 {
	x, y, z := Tx(ax, v[0], v[1], v[2])
	return [3]float64{x, y, z}
}
