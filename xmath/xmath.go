// Package xmath provides mathematics routines convenient to xirho and its components.
package xmath

import "math"

// Angle wraps a float64 value to the interval (-pi, pi]. The result may be NaN
// if the argument is not finite.
func Angle(x float64) float64 {
	x = math.Mod(x, 2*math.Pi)
	if x > math.Pi {
		x -= 2 * math.Pi
	} else if x <= -math.Pi {
		x += 2 * math.Pi
	}
	return x
}

// IsFinite returns whether the argument is finite.
func IsFinite(x float64) bool {
	// If x is ±inf or nan, then x-x is nan; otherwise, x-x is 0.
	return x-x == 0
}
