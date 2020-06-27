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

// F is a function ("variation") type.
//
// Functions may be parametrized in a number of ways with the Param types Flag,
// List, Int, Angle, Real, Complex, Vec3, Affine, Func, and FuncList. If fields
// of these types are exported, package fapi can collect them to enable a user
// interface for setting and displaying such parameters.
type F interface {
	// Calc calculates the function at a point.
	Calc(in P, rng *RNG) P
	// Prep is called once prior to iteration so that a function can cache
	// expensive calculations.
	Prep()
}

// Param is a function parameter which may vary per function instance. The only
// implementations of Param are Flag, List, Int, Angle, Real, Complex, Vec3,
// Affine, Func, and FuncList.
type Param interface {
	// isParam ensures that no external types may implement Param.
	isParam() sealed
}

// Flag is a boolean function parameter.
type Flag bool

// List is a function parameter to choose among a fixed set of options.
type List int

// Int is an integer function parameter, possibly bounded.
type Int int64

// Angle is an angle function parameter. External interfaces wrap its value
// into the interval (-pi, pi].
type Angle float64

// Real is a floating-point function parameter, possibly bounded.
type Real float64

// Complex is an unconstrained function parameter in R^2.
type Complex complex128

// Vec3 is an unconstrained function parameter in R^3.
type Vec3 [3]float64

// Affine is an affine transform function parameter.
type Affine = Ax

// Func must be a struct wrapping its value because F is an interface type,
// which means it cannot have methods – i.e. isParam() – defined on it.

// Func is a function parameter that is itself a function. Note that unlike
// other parameter types, Func is a struct wrapping its value.
type Func struct {
	F
}

// FuncList is a function parameter holding a list of functions.
type FuncList []F

// sealed prevents external types from implementing Param.
type sealed struct{}

func (Flag) isParam() sealed     { panic(nil) }
func (List) isParam() sealed     { panic(nil) }
func (Int) isParam() sealed      { panic(nil) }
func (Angle) isParam() sealed    { panic(nil) }
func (Real) isParam() sealed     { panic(nil) }
func (Complex) isParam() sealed  { panic(nil) }
func (Vec3) isParam() sealed     { panic(nil) }
func (Affine) isParam() sealed   { panic(nil) }
func (Func) isParam() sealed     { panic(nil) }
func (FuncList) isParam() sealed { panic(nil) }
