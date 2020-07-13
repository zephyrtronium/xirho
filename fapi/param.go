package fapi

import (
	"math"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xmath"
)

// Param is a user-controlled function parameter. Each implementing type has
// a setter and a getter for a corresponding xirho.Param type.
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
	v *xirho.Flag
	paramName
}

// flagFor creates a Flag function parameter.
func flagFor(name string, v *xirho.Flag) Param {
	return Flag{
		v:         v,
		paramName: paramName(name),
	}
}

// Set sets the parameter value. The returned error is always nil.
func (p Flag) Set(v xirho.Flag) error {
	*p.v = v
	return nil
}

// Get gets the current parameter value.
func (p Flag) Get() xirho.Flag {
	return *p.v
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
	v *xirho.List
	// opts is the list of options for display.
	opts []string

	paramName
}

// listFor creates a List function parameter.
func listFor(name string, idx *xirho.List, opts ...string) Param {
	if len(opts) == 0 {
		panic("xirho: list needs at least one option")
	}
	opts = append([]string{}, opts...) // copy
	return List{
		v:         idx,
		paramName: paramName(name),
		opts:      opts,
	}
}

// Set sets the list value. If v is less than zero or larger than the number
// of available options, an error of type OutOfBoundsInt is returned instead.
func (p List) Set(v xirho.List) error {
	if v < 0 || int64(v) >= int64(len(p.opts)) {
		return OutOfBoundsInt{
			Param: p,
			Value: int64(v),
			Lo:    0,
			Hi:    int64(len(p.opts) - 1),
		}
	}
	*p.v = v
	return nil
}

// Get gets the list integer value.
func (p List) Get() xirho.List {
	return *p.v
}

// String gets the list's selected string.
func (p List) String() string {
	return p.opts[*p.v]
}

// Opts returns a copy of the list's options.
func (p List) Opts() []string {
	return append([]string(nil), p.opts...)
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
	v *xirho.Int
	// bdd indicates whether external interfaces should respect Lo and Hi.
	bdd bool
	// lo and hi indicate the minimum and maximum values, inclusive, that an
	// external interface should attempt to assign to V.
	lo, hi xirho.Int

	paramName
}

// intFor creates an Int function parameter.
func intFor(name string, v *xirho.Int, bounded bool, lo, hi xirho.Int) Param {
	return Int{
		v:         v,
		paramName: paramName(name),
		bdd:       bounded,
		lo:        lo,
		hi:        hi,
	}
}

// Set sets the int value. If the Int is bounded and v is out of its bounds, an
// error of type OutOfBoundsInt is returned instead.
func (p Int) Set(v xirho.Int) error {
	if p.bdd && (v < p.lo || p.hi < v) {
		return OutOfBoundsInt{
			Param: p,
			Value: int64(v),
			Lo:    int64(p.lo),
			Hi:    int64(p.hi),
		}
	}
	*p.v = v
	return nil
}

// Get gets the int value.
func (p Int) Get() xirho.Int {
	return *p.v
}

// Bounded returns whether the parameter is bounded.
func (p Int) Bounded() bool {
	return p.bdd
}

// Bounds returns the parameter bounds. If the parameter is not bounded, the
// returned bounds are the minimum and maximum values of int64.
func (p Int) Bounds() (lo, hi xirho.Int) {
	if !p.bdd {
		return -1 << 63, 1<<63 - 1
	}
	return p.lo, p.hi
}

// Angle is an angle function parameter. External interfaces wrap its value
// into the interval (-pi, pi].
type Angle struct {
	v *xirho.Angle
	paramName
}

// angleFor creates an Angle function parameter.
func angleFor(name string, v *xirho.Angle) Param {
	return Angle{
		v:         v,
		paramName: paramName(name),
	}
}

// Set sets the angle value wrapped into the interval (-pi, pi].
func (p Angle) Set(v xirho.Angle) error {
	x := xmath.Angle(float64(v))
	// TODO: error if x is not finite
	*p.v = xirho.Angle(x)
	return nil
}

// Get gets the angle value.
func (p Angle) Get() xirho.Angle {
	return *p.v
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
	v *xirho.Real
	// bdd indicates whether external interfaces should respect Lo and Hi.
	bdd bool
	// lo and hi indicate the minimum and maximum values, inclusive, that an
	// external interface should attempt to assign to V.
	lo, hi xirho.Real

	paramName
}

// realFor creates a Real function parameter.
func realFor(name string, v *xirho.Real, bounded bool, lo, hi xirho.Real) Param {
	return Real{
		v:         v,
		paramName: paramName(name),
		bdd:       bounded,
		lo:        lo,
		hi:        hi,
	}
}

// Set sets the real value. If the Real is bounded and v is out of its bounds,
// an error of type OutOfBoundsReal is returned instead.
func (p Real) Set(v xirho.Real) error {
	if p.bdd && (v < p.lo || p.hi < v) {
		return OutOfBoundsReal{
			Param: p,
			Value: float64(v),
			Lo:    float64(p.lo),
			Hi:    float64(p.hi),
		}
	}
	*p.v = v
	return nil
}

// Get gets the real value.
func (p Real) Get() xirho.Real {
	return *p.v
}

// Bounded returns whether the parameter is bounded.
func (p Real) Bounded() bool {
	return p.bdd
}

// Bounds returns the parameter bounds. If the parameter is not bounded, the
// returned bounds are -inf and +inf.
func (p Real) Bounds() (lo, hi xirho.Real) {
	if !p.bdd {
		return xirho.Real(math.Inf(-1)), xirho.Real(math.Inf(1))
	}
	return p.lo, p.hi
}

// Complex is an unconstrained function parameter in R^2.
type Complex struct {
	v *xirho.Complex
	paramName
}

// complexFor creates a Complex function parameter.
func complexFor(name string, v *xirho.Complex) Param {
	return Complex{
		v:         v,
		paramName: paramName(name),
	}
}

// Set sets the complex value.
func (p Complex) Set(v xirho.Complex) error {
	*p.v = v
	return nil
}

// Get gets the complex value.
func (p Complex) Get() xirho.Complex {
	return *p.v
}

// Vec3 is an unconstrained function parameter in R^3.
type Vec3 struct {
	v *xirho.Vec3
	paramName
}

// vec3For creates a Vec3 function parameter.
func vec3For(name string, v *xirho.Vec3) Param {
	return Vec3{
		v:         v,
		paramName: paramName(name),
	}
}

// Set sets the vector value.
func (p Vec3) Set(v xirho.Vec3) error {
	*p.v = v
	return nil
}

// Get gets the vector value.
func (p Vec3) Get() xirho.Vec3 {
	return *p.v
}

// Affine is an affine transform function parameter.
type Affine struct {
	v *xirho.Affine
	paramName
}

// affineFor creates an Affine function parameter.
func affineFor(name string, v *xirho.Affine) Param {
	return Affine{
		v:         v,
		paramName: paramName(name),
	}
}

// Set sets the affine transform value.
func (p Affine) Set(v xirho.Affine) error {
	*p.v = v
	return nil
}

// Get gets the affine transform value.
func (p Affine) Get() xirho.Affine {
	return *p.v
}

// Func is a function parameter that is itself a function. After the parameter
// name, a Func field may include an additional comma-separated tag containing
// the string "optional". Func fields marked optional may be set to nil. For
// example:
//
//		type Example struct {
//			NillableFunc `xirho:"func,optional"`
//		}
type Func struct {
	v *xirho.Func
	// opt indicates whether the parameter is allowed to be nil.
	opt bool

	paramName
}

// funcFor creates a Func function parameter.
func funcFor(name string, opt bool, v *xirho.Func) Param {
	return Func{
		v:         v,
		opt:       opt,
		paramName: paramName(name),
	}
}

// Set sets the function value. If the parameter is not optional and v.F is
// nil, an error of type NotOptional is returned instead.
func (p Func) Set(v xirho.F) error {
	if !p.opt && v == nil {
		return NotOptional{Param: p}
	}
	*p.v = xirho.Func{F: v}
	return nil
}

// Get gets the function value.
func (p Func) Get() xirho.F {
	return p.v.F
}

// IsOptional returns whether the function may be set to nil.
func (p Func) IsOptional() bool {
	return p.opt
}

// FuncList is a function parameter holding a list of functions.
type FuncList struct {
	v *xirho.FuncList
	paramName
}

// funcListFor creates a FuncList function parameter.
func funcListFor(name string, v *xirho.FuncList) Param {
	return FuncList{
		v:         v,
		paramName: paramName(name),
	}
}

// Set sets the function list value.
func (p FuncList) Set(v xirho.FuncList) error {
	*p.v = v
	return nil
}

// Get returns a copy of the function list value.
func (p FuncList) Get() xirho.FuncList {
	return append(xirho.FuncList(nil), *p.v...)
}

// Append appends functions to the list.
func (p FuncList) Append(v ...xirho.F) {
	*p.v = append(*p.v, v...)
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
