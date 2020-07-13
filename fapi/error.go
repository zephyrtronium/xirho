package fapi

import "fmt"

// OutOfBoundsInt is an error returned when attempting to set an Int or List to
// a value which is out of the parameter's defined limits.
type OutOfBoundsInt struct {
	// Param is the parameter which the caller attempted to set.
	Param Param
	// Value is the value which the caller attempted to use.
	Value int64
	// Lo and Hi are the minimum and maximum allowed values, inclusive.
	Lo, Hi int64
}

// Error returns a formatted error message.
func (err OutOfBoundsInt) Error() string {
	return fmt.Sprintf("cannot set %s to %d: value must be between %d and %d, inclusive", err.Param.Name(), err.Value, err.Lo, err.Hi)
}

// OutOfBoundsReal is an error returned when attempting to set a Real to a
// value which is out of the parameter's defined limits.
type OutOfBoundsReal struct {
	// Param is the parameter which the caller attempted to set.
	Param Param
	// Value is the value which the caller attempted to use.
	Value float64
	// Lo and Hi are the minimum and maximum allowed values, inclusive.
	Lo, Hi float64
}

// Error returns a formatted error message.
func (err OutOfBoundsReal) Error() string {
	return fmt.Sprintf("cannot set %s to %g: value must be between %g and %g, inclusive", err.Param.Name(), err.Value, err.Lo, err.Hi)
}

// NotFinite is an error returned when attempting to set an Angle, Real,
// Complex, Vec3, or Affine parameter to a value with a component that is not
// finite.
type NotFinite struct {
	// Param is the parameter which the caller attempted to set.
	Param Param
}

func (err NotFinite) Error() string {
	return fmt.Sprintf("cannot set %s to a value which is not finite", err.Param.Name())
}

// NotOptional is an error returned when attempting to set a Func to nil when
// the Func is not marked as optional.
type NotOptional struct {
	// Param is the parameter which the caller attempted to set.
	Param Param
}

// Error returns a formatted error message.
func (err NotOptional) Error() string {
	return err.Param.Name() + " is not optional"
}
