// Package xirho implements an iterated function system fractal renderer.
//
// An iterated function system is a collection of functions from points to
// points. Starting with a randomly selected point, we choose a function at
// random, apply that function to the point, and plot its new location, then
// repeat ad infinitum. With some additional steps, the result images can be
// stunning.
//
// Xirho does not include a designer to produce systems to render. Existing
// parameters can be loaded through the encoding and encoding/flame
// subpackages, or programmed by hand.
//
// To use xirho to render a system, create an R containing the System and a
// Hist to plot points, then call R.Render with a context. Alternatively, the
// System.Iter method can be used for more fine-grained control of the
// rendering process.
package xirho
