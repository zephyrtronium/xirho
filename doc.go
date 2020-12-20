// Package xirho implements an iterated function system fractal art renderer.
//
// An iterated function system is a collection of functions from points to
// points. Starting with a randomly selected point, we choose a function at
// random, apply that function to the point, and plot its new location, then
// repeat ad infinitum. With some additional steps, the result images can be
// stunning.
//
// The mathematical terminology used in xirho's documentation and API is as
// follows. A point is an element of R³ × [0, 1], i.e. a 3D point plus a color
// coordinate. A function, sometimes function type, is a procedure which maps
// points to points, possibly using additional fixed parameters to control the
// exact mapping. (Other IFS implementations typically refer to functions in
// this sense as variations.) A node is a particular instance of a function and
// its fixed parameters. An iterated function system, or just system, is a
// non-empty list of nodes, a Markov chain giving the probability of the
// algorithm transitioning from each node in the list to each other node in the
// list, an additional node applied to each output point to serve as a possibly
// nonlinear camera, and a mapping of color coordinates to colors. The Markov
// chain of a system may also be called the weights graph, or just the graph.
//
// Xirho does not include a designer to produce systems to render. Existing
// parameters can be loaded through the encoding and encoding/flame
// subpackages, or programmed by hand.
//
// To use xirho to render a system, create a Render containing the System and a
// Hist to plot points, then call its Render method with a non-trivial context.
// (The context closing is the only way that Render returns.) Alternatively,
// the RenderAsync method provides an API to manage rendering concurrently,
// e.g. to support a UI. For fine-grained control of the rendering process,
// the System.Iter method can be used directly.
package xirho
