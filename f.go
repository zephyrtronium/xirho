package xirho

// P is a point in R^3 × [0, 1].
type P struct {
	// X, Y, and Z are spatial coordinates.
	X, Y, Z float64
	// C is the color coordinate in [0, 1].
	C float64
}

// IsValid returns true if its spatial coordinates are finite and its color
// coordinate is in [0, 1].
func (p P) IsValid() bool {
	// x - x is 0 if x is finite and NaN otherwise.
	return p.X-p.X == p.Y-p.Y && p.Z-p.Z == 0 && 0 <= p.C && p.C <= 1
}

// F is a function type.
type F interface {
	// Calc calculates the function at a point.
	Calc(in P, rng *RNG) P

	// Params lists the function parameters. The renderer guarantees that no
	// goroutine is in a call to Calc while any returned Param is in use.
	Params() []Param
}
