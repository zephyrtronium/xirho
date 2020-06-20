package xirho

// Ax is an affine transform.
type Ax struct {
	// A is the linear (multiplicative) component of the transform.
	A [3][3]float64
	// B is the translation (additive) component of the transform.
	B [3]float64
}

// Eye returns an identity transform.
func Eye() *Ax {
	return &Ax{
		A: [3][3]float64{
			{1, 0, 0},
			{0, 1, 0},
			{0, 0, 1},
		},
	}
}

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
