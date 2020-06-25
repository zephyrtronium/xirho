package xirho

// Param is a function parameter which may vary per function instance. The only
// implementations of Param are Flag, List, Int, Angle, Real, Complex, Vec3,
// Affine, Func, and FuncList.
type Param interface {
	// Name returns the name of the parameter.
	Name() string

	// isParam ensures that no external types may implement Param.
	isParam() sealed
}

// paramName is a shortcut embeddable type for param names.
type paramName string

// Name returns the parameter name.
func (p paramName) Name() string {
	return string(p)
}

// Flag is a boolean function parameter.
type Flag struct {
	// V is a reference to the parameter value.
	V *bool
	paramName
}

// FlagParam creates a Flag function parameter.
func FlagParam(v *bool, name string) Param {
	return Flag{
		V:         v,
		paramName: paramName(name),
	}
}

// List is a function parameter to choose among a list of strings.
type List struct {
	// V is a reference to the selected index.
	V *int
	// Opts is the list of options for display. It should never be modified.
	Opts []string

	paramName
}

// ListParam creates a List function parameter.
func ListParam(idx *int, name string, opts ...string) Param {
	opts = append([]string{}, opts...) // copy
	return List{
		V:         idx,
		paramName: paramName(name),
		Opts:      opts,
	}
}

// Int is an integer function parameter, possibly bounded.
type Int struct {
	// V is a reference to the parameter value.
	V *int64
	// Bounded indicates whether external interfaces should respect Lo and Hi.
	Bounded bool
	// Lo and Hi indicate the minimum and maximum values, inclusive, that an
	// external interface should attempt to assign to V.
	Lo, Hi int64

	paramName
}

// IntParam creates an Int function parameter.
func IntParam(v *int64, name string, bounded bool, lo, hi int64) Param {
	return Int{
		V:         v,
		paramName: paramName(name),
		Bounded:   bounded,
		Lo:        lo,
		Hi:        hi,
	}
}

// Angle is an angle function parameter. External interfaces wrap its value
// into the interval (-pi, pi].
type Angle struct {
	// V is a reference to the parameter value.
	V *float64
	paramName
}

// AngleParam creates an Angle function parameter.
func AngleParam(v *float64, name string) Param {
	return Angle{
		V:         v,
		paramName: paramName(name),
	}
}

// Real is a floating-point function parameter, possibly bounded.
type Real struct {
	// V is a reference to the parameter value.
	V *float64
	// Bounded indicates whether external interfaces should respect Lo and Hi.
	Bounded bool
	// Lo and Hi indicate the minimum and maximum values, inclusive, that an
	// external interface should attempt to assign to V.
	Lo, Hi float64

	paramName
}

// RealParam creates a Real function parameter.
func RealParam(v *float64, name string, bounded bool, lo, hi float64) Param {
	return Real{
		V:         v,
		paramName: paramName(name),
		Bounded:   bounded,
		Lo:        lo,
		Hi:        hi,
	}
}

// Complex is an unconstrained function parameter in R^2.
type Complex struct {
	// V is a reference to the parameter value.
	V *complex128
	paramName
}

// ComplexParam creates a Complex function parameter.
func ComplexParam(v *complex128, name string) Param {
	return Complex{
		V:         v,
		paramName: paramName(name),
	}
}

// Vec3 is an unconstrained function parameter in R^3.
type Vec3 struct {
	// V is a reference to the parameter value.
	V *[3]float64
	paramName
}

// Vec3Param creates a Vec3 function parameter.
func Vec3Param(v *[3]float64, name string) Param {
	return Vec3{
		V:         v,
		paramName: paramName(name),
	}
}

// Affine is an affine transform function parameter.
type Affine struct {
	V *Ax
	paramName
}

// AffineParam creates an Affine function parameter.
func AffineParam(v *Ax, name string) Param {
	return Affine{
		V:         v,
		paramName: paramName(name),
	}
}

// Func is a function parameter that is itself a function.
type Func struct {
	V *F
	paramName
}

// FuncParam creates a Func function parameter.
func FuncParam(v *F, name string) Param {
	return Func{
		V:         v,
		paramName: paramName(name),
	}
}

// FuncList is a function parameter holding a list of functions.
type FuncList struct {
	V *[]F
	paramName
}

// FuncListParam creates a FuncList function parameter.
func FuncListParam(v *[]F, name string) Param {
	return FuncList{
		V:         v,
		paramName: paramName(name),
	}
}

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
