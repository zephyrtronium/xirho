# xirho

xirho is a simple generalized iterated function system plotter. It makes pretty pictures out of math.

---

![Spherical gasket](img/spherical.png)

---

An [iterated function system](https://en.wikipedia.org/wiki/Iterated_function_system), or IFS, is just a list of functions that turn points into other points.

Take a random place in 3D space. Pick a function in the system at random. Apply the function to the point, resulting in a new point. Plot the new point. Repeat.

That's it.

---

![Sierpinski gasket](img/sierpinski.png)

The original treatments of IFS were mostly concerned with affine transformations: simple functions describing uniform scaling, rotation, shearing, and translation. We can get some pretty good images out of just these; the Sierpinski gasket just above is an example.

The "fractal flame" algorithm is a way of generalizing IFS: allowing arbitrary functions in the system and adding color and tone mapping to the output. The spherical gasket at the top of this page is very similar to the Sierpinski gasket, just replacing a couple affine transformations with a simple nonlinear function.

---

![Grand Julian](img/grandjulian.png)

xirho is a pretty basic fractal flame renderer, with only a handful of function types available (for now). It doesn't support designing a system interactively (yet), or even loading systems from some serialization format (yet) – each example is programmed by hand. The renderer is flexible thanks to Go's powerful type system, and it's fast because of its completely lock-free parallel design.

Ultimately, xirho is a pet project that I've wanted to implement for a decade to address the bugs in [Apophysis](https://en.wikipedia.org/wiki/Apophysis_(software)) and the commercial limitations in [Chaotica](https://www.chaoticafractals.com/) (which is still an outstanding piece of software!). xirho isn't intended to be the fastest IFS renderer, nor the most versatile, and it's explicitly avoiding compatibility with those other tools and their mimics. But it's something I enjoy working on, so it will get better than it is.

---

![Disc Julian](img/discjulian.png)
