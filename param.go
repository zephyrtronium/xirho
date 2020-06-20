package xirho

// Param is a function parameter which may vary per function instance. The only
// implementations of Param are Int, Angle, Real, Complex, Vec3, and Affine.
type Param interface {
	// Name returns the name
	Name() string

	// isParam ensures that no external types may implement Param.
	isParam() sealed
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

	name string
}

// IntParam creates an Int function parameter.
func IntParam(v *int64, name string, bounded bool, lo, hi int64) Param {
	return Int{
		V:       v,
		name:    name,
		Bounded: bounded,
		Lo:      lo,
		Hi:      hi,
	}
}

func (p Int) Name() string {
	return p.name
}

// Angle is an angle function parameter. External interfaces wrap its value
// into the interval (-pi, pi].
type Angle struct {
	// V is a reference to the parameter value.
	V *float64

	name string
}

// AngleParam creates an Angle function parameter.
func AngleParam(v *float64, name string) Param {
	return Angle{
		V:    v,
		name: name,
	}
}

func (p Angle) Name() string {
	return p.name
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

	name string
}

// RealParam creates a Real function parameter.
func RealParam(v *float64, name string, bounded bool, lo, hi float64) Param {
	return Real{
		V:       v,
		name:    name,
		Bounded: bounded,
		Lo:      lo,
		Hi:      hi,
	}
}

func (p Real) Name() string {
	return p.name
}

// Complex is an unconstrained function parameter in R^2.
type Complex struct {
	// V is a reference to the parameter value.
	V *complex128

	name string
}

// ComplexParam creates a Complex function parameter.
func ComplexParam(v *complex128, name string) Param {
	return Complex{
		V:    v,
		name: name,
	}
}

func (p Complex) Name() string {
	return p.name
}

// Vec3 is an unconstrained function parameter in R^3.
type Vec3 struct {
	// V is a reference to the parameter value.
	V *[3]float64

	name string
}

// Vec3Param creates a Vec3 function parameter.
func Vec3Param(v *[3]float64, name string) Param {
	return Vec3{
		V:    v,
		name: name,
	}
}

func (p Vec3) Name() string {
	return p.name
}

// Affine is an affine transform function parameter.
type Affine struct {
	V *Ax

	name string
}

// AffineParam creates an Affine function parameter.
func AffineParam(v *Ax, name string) Param {
	return Affine{
		V:    v,
		name: name,
	}
}

func (p Affine) Name() string {
	return p.name
}

// sealed prevents external types from implementing Param.
type sealed struct{}

func (Int) isParam() sealed     { panic(nil) }
func (Angle) isParam() sealed   { panic(nil) }
func (Real) isParam() sealed    { panic(nil) }
func (Complex) isParam() sealed { panic(nil) }
func (Vec3) isParam() sealed    { panic(nil) }
func (Affine) isParam() sealed  { panic(nil) }
