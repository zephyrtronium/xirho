package fapi

import "github.com/zephyrtronium/xirho"

// Param is a user-controlled function parameter. Each implementing type wraps
// a corresponding xirho.Param type.
type Param interface {
	// Name returns the parameter name.
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
	V *xirho.Flag
	paramName
}

// flagFor creates a Flag function parameter.
func flagFor(name string, v *xirho.Flag) Param {
	return Flag{
		V:         v,
		paramName: paramName(name),
	}
}

// List is a function parameter to choose among a list of strings. After the
// parameter name, a List field may include any number of additional
// comma-separated tags to define the display names of each option. For
// example, to define a parameter allowing the user to choose between "fast",
// "accurate", or "balanced", do:
//
//		type Example struct {
//			Method xirho.List `xirho:"method,fast,accurate,balanced"`
//		}
type List struct {
	V *xirho.List
	// opts is the list of options for display.
	opts []string

	paramName
}

// listFor creates a List function parameter.
func listFor(name string, idx *xirho.List, opts ...string) Param {
	opts = append([]string{}, opts...) // copy
	return List{
		V:         idx,
		paramName: paramName(name),
		opts:      opts,
	}
}

// Int is an integer function parameter. After the parameter name, an Int field
// may include two additional comma-separated tags to define the lowest and
// highest permitted values. For example, to define a parameter allowing the
// user to choose any integer in [-3, 12], do:
//
//		type Example struct {
//			P xirho.Int `xirho:"p,-3,12"`
//		}
type Int struct {
	V *xirho.Int
	// Bounded indicates whether external interfaces should respect Lo and Hi.
	Bounded bool
	// Lo and Hi indicate the minimum and maximum values, inclusive, that an
	// external interface should attempt to assign to V.
	Lo, Hi int64

	paramName
}

// intFor creates an Int function parameter.
func intFor(name string, v *xirho.Int, bounded bool, lo, hi int64) Param {
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
	V *xirho.Angle
	paramName
}

// angleFor creates an Angle function parameter.
func angleFor(name string, v *xirho.Angle) Param {
	return Angle{
		V:         v,
		paramName: paramName(name),
	}
}

// Real is a floating-point function parameter. After the parameter name, a
// Real field may include two additional comma-separated tags to define the
// lowest and highest permitted values. For example, to define a parameter
// allowing the user to choose any real in [-2π, 2π], do:
//
//		type Example struct {
//			P xirho.Real `xirho:"p,-6.283185307179586,6.283185307179586"`
//		}
type Real struct {
	V *xirho.Real
	// Bounded indicates whether external interfaces should respect Lo and Hi.
	Bounded bool
	// Lo and Hi indicate the minimum and maximum values, inclusive, that an
	// external interface should attempt to assign to V.
	Lo, Hi float64

	paramName
}

// realFor creates a Real function parameter.
func realFor(name string, v *xirho.Real, bounded bool, lo, hi float64) Param {
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
	V *xirho.Complex
	paramName
}

// complexFor creates a Complex function parameter.
func complexFor(name string, v *xirho.Complex) Param {
	return Complex{
		V:         v,
		paramName: paramName(name),
	}
}

// Vec3 is an unconstrained function parameter in R^3.
type Vec3 struct {
	V *xirho.Vec3
	paramName
}

// vec3For creates a Vec3 function parameter.
func vec3For(name string, v *xirho.Vec3) Param {
	return Vec3{
		V:         v,
		paramName: paramName(name),
	}
}

// Affine is an affine transform function parameter.
type Affine struct {
	V *xirho.Affine
	paramName
}

// affineFor creates an Affine function parameter.
func affineFor(name string, v *xirho.Affine) Param {
	return Affine{
		V:         v,
		paramName: paramName(name),
	}
}

// Func is a function parameter that is itself a function.
type Func struct {
	V *xirho.Func
	paramName
}

// funcFor creates a Func function parameter.
func funcFor(name string, v *xirho.Func) Param {
	return Func{
		V:         v,
		paramName: paramName(name),
	}
}

// FuncList is a function parameter holding a list of functions.
type FuncList struct {
	V *xirho.FuncList
	paramName
}

// funcListFor creates a FuncList function parameter.
func funcListFor(name string, v *xirho.FuncList) Param {
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
