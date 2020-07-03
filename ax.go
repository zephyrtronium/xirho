package xirho

import (
	"math"

	"golang.org/x/image/math/f64"
)

// Ax is an affine transform.
type Ax f64.Aff4

// Eye returns an identity transform.
func Eye() Ax {
	return Ax{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
	}
}

// Eye sets the transform to the identity transform and returns it.
func (ax *Ax) Eye() *Ax {
	*ax = Ax{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
	}
	return ax
}

// Translate translates the transform and returns it.
func (ax *Ax) Translate(dx, dy, dz float64) *Ax {
	ax[3] += dx
	ax[7] += dy
	ax[11] += dz
	return ax
}

// Scale scales the transform along each axis and returns it.
func (ax *Ax) Scale(sx, sy, sz float64) *Ax {
	ax[0] *= sx
	ax[1] *= sy
	ax[2] *= sz
	ax[4] *= sx
	ax[5] *= sy
	ax[6] *= sz
	ax[8] *= sx
	ax[9] *= sy
	ax[10] *= sz
	return ax
}

// RotX rotates the transform about the x axis and returns it. Specifically, tx
// is the counter-clockwise rotation in the y/z plane in radians. The
// translation vector is also rotated appropriately.
func (ax *Ax) RotX(tx float64) *Ax {
	// Rotation matrix:
	//	[1	0	0	0]
	//	[0	cx	sx	0]
	//	[0	-sx	cx	0]
	//	[0	0	0	1]
	sx, cx := math.Sincos(tx)
	ax[1], ax[2] = ax[1]*cx-ax[2]*sx, ax[1]*sx+ax[2]*cx
	ax[5], ax[6] = ax[5]*cx-ax[6]*sx, ax[5]*sx+ax[6]*cx
	ax[9], ax[10] = ax[9]*cx-ax[10]*sx, ax[9]*sx+ax[10]*cx
	ax[7], ax[11] = ax[7]*cx+ax[11]*sx, -ax[7]*sx+ax[11]*cx
	return ax
}

// RotY rotates the transform about the y axis and returns it. Specifically, ty
// is the rotation in the x/z plane in radians.
func (ax *Ax) RotY(ty float64) *Ax {
	// Rotation matrix:
	//	[cy	0	-sy	0]
	//	[0	1	0	0]
	//	[sy	0	cy	0]
	//	[0	0	0	1]
	sy, cy := math.Sincos(ty)
	ax[0], ax[2] = ax[0]*cy+ax[2]*sy, ax[2]*cy-ax[0]*sy
	ax[4], ax[6] = ax[4]*cy+ax[6]*sy, ax[6]*cy-ax[4]*sy
	ax[8], ax[10] = ax[8]*cy+ax[10]*sy, ax[10]*cy-ax[8]*sy
	ax[3], ax[11] = ax[3]*cy-ax[11]*sy, ax[3]*sy+ax[11]*cy
	return ax
}

// RotZ rotates the transform about the z axis and returns it. Specifically, tz
// is the rotation in the x/y plane in radians.
func (ax *Ax) RotZ(tz float64) *Ax {
	// Rotation matrix:
	//	[	cz	sz	0	0]
	//	[	-sz	cz	0	0]
	//	[	0	0	1	0]
	//	[	0	0	0	1]
	sz, cz := math.Sincos(tz)
	ax[0], ax[1] = ax[0]*cz-ax[1]*sz, ax[0]*sz+ax[1]*cz
	ax[4], ax[5] = ax[4]*cz-ax[5]*sz, ax[4]*sz+ax[5]*cz
	ax[8], ax[9] = ax[8]*cz-ax[9]*sz, ax[8]*sz+ax[9]*cz
	ax[3], ax[7] = ax[3]*cz+ax[7]*sz, -ax[3]*sz+ax[7]*cz
	return ax
}

// Pitch rotates the transform about the x axis centered on the transform's
// translation point and returns it.
func (ax *Ax) Pitch(tx float64) *Ax {
	sx, cx := math.Sincos(tx)
	ax[1], ax[2] = ax[1]*cx-ax[2]*sx, ax[1]*sx+ax[2]*cx
	ax[5], ax[6] = ax[5]*cx-ax[6]*sx, ax[5]*sx+ax[6]*cx
	ax[9], ax[10] = ax[9]*cx-ax[10]*sx, ax[9]*sx+ax[10]*cx
	return ax
}

// Roll rotates the transform about the y axis centered on the transform's
// translation point and returns it.
func (ax *Ax) Roll(ty float64) *Ax {
	sy, cy := math.Sincos(ty)
	ax[0], ax[2] = ax[0]*cy+ax[2]*sy, ax[2]*cy-ax[0]*sy
	ax[4], ax[6] = ax[4]*cy+ax[6]*sy, ax[6]*cy-ax[4]*sy
	ax[8], ax[10] = ax[8]*cy+ax[10]*sy, ax[10]*cy-ax[8]*sy
	return ax
}

// Yaw rotates the transform about the z axis centered on the transform's
// translation point and returns it.
func (ax *Ax) Yaw(tz float64) *Ax {
	sz, cz := math.Sincos(tz)
	ax[0], ax[1] = ax[0]*cz-ax[1]*sz, ax[0]*sz+ax[1]*cz
	ax[4], ax[5] = ax[4]*cz-ax[5]*sz, ax[4]*sz+ax[5]*cz
	ax[8], ax[9] = ax[8]*cz-ax[9]*sz, ax[8]*sz+ax[9]*cz
	return ax
}

// TODO: shear

// Tx transforms a coordinate with an affine transform.
func Tx(ax *Ax, x, y, z float64) (tx, ty, tz float64) {
	tx = ax[0]*x + ax[1]*y + ax[2]*z + ax[3]
	ty = ax[4]*x + ax[5]*y + ax[6]*z + ax[7]
	tz = ax[8]*x + ax[9]*y + ax[10]*z + ax[11]
	return
}

// TxVec transforms a 3-vector with an affine transform.
func TxVec(ax *Ax, v [3]float64) [3]float64 {
	x, y, z := Tx(ax, v[0], v[1], v[2])
	return [3]float64{x, y, z}
}
