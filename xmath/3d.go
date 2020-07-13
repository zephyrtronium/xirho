package xmath

import "math"

// R3 calculates the Euclidean distance from the origin to the given point.
func R3(x, y, z float64) float64 {
	// yx, zx := y/x, z/x
	// return math.Abs(x) * math.Sqrt(1 + yx*yx + zx*zx)
	return math.Sqrt(x*x + y*y + z*z)
}

// Spherical converts rectangular (x, y, z) coordinates to spherical (r, θ, φ).
func Spherical(x, y, z float64) (r, theta, phi float64) {
	r = R3(x, y, z)
	theta = math.Atan2(y, x)
	phi = math.Acos(z / r)
	return
}

// FromSpherical converts spherical (r, θ, φ) coordinates to rectangular
// (x, y, z).
func FromSpherical(r, theta, phi float64) (x, y, z float64) {
	st, ct := math.Sincos(theta)
	sp, cp := math.Sincos(phi)
	x = r * ct * sp
	y = r * st * sp
	z = r * cp
	return
}
