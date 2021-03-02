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

// Fit returns the width and height of an integer rectangle having an aspect
// ratio closest to aspect and with width or height equal to w or h,
// respectively.
func Fit(w, h int, aspect float64) (int, int) {
	if aspect >= 1 {
		h = int(float64(w)/aspect + 0.5)
	} else {
		w = int(float64(h)*aspect + 0.5)
	}
	return w, h
}
