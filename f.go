package xirho

// Pt is a point in R^3 × [0, 1].
type Pt struct {
	// X, Y, and Z are spatial coordinates.
	X, Y, Z float64
	// C is the color coordinate in [0, 1].
	C float64
}

// IsValid returns true if its spatial coordinates are finite and its color
// coordinate is in [0, 1].
func (p Pt) IsValid() bool {
	// x - x is 0 if x is finite and NaN otherwise.
	return p.X-p.X == p.Y-p.Y && p.Z-p.Z == 0 && 0 <= p.C && p.C <= 1
}

// Func is a function ("variation") type.
//
// Functions may be parametrized in a number of ways with the types Flag, List,
// Int, Angle, Real, Complex, Vec3, Affine, Func, and FuncList. If fields of
// these types are exported, package fapi can collect them to enable a user
// interface for setting and displaying such parameters.
type Func interface {
	// Calc calculates the function at a point.
	Calc(in Pt, rng *RNG) Pt
	// Prep is called once prior to iteration so that a function can cache
	// expensive calculations.
	Prep()
}
