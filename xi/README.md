# Ξ

Package xi provides function registration and implementations for xirho.

The following function types (variations) are implemented here:

- Affine (the "triangles" in Apophysis)
- Blur
- Bubble
- CElliptic (similar to elliptic)
- ColorSpeed (like color and color symmetry in Apophysis)
- Curl
- Cylinder
- Disc
- Flatten
- JuliaN
- Log
- Mobius (a 3D version, like the mobiq plugin)
- Perspective (like in the Apophysis render settings)
- Polar
- Rod
- Scale (like linear or linear3D)
- Spherical
- Splits (the 3D version)
- Sum (roughly implements the behavior of multiple variations in Apophysis)
- Then (turns any variation into a pre- or post- variant, and more general besides)

## Adding new variations

Xi is designed so that external packages may add any number of variations during initialization. For example, in a package providing function types named "madoka" and "homura", one could do:

```go
func init() {
    xi.Register("madoka", func() xirho.F { return Madoka{} })
    xi.Register("homura", newHomura)
}
```

Then madoka and homura, including any of their parameters, will be marshaled and unmarshaled automatically by package encoding. User interfaces will automatically have access to madoka and homura and their parameters through the same mechanisms by which they use any other function types.

Doing this also opens the option of creating Go plugins to distribute functions types easily. If the package providing madoka and homura is called anime, then one could add a file called e.g. `plugin.go` along these lines:

```go
// +build ignore

package main

import _ "example.org/anime"
```

Then `go build -buildmode=plugin plugin.go` will create a Go plugin that a renderer program could load dynamically and automatically have madoka and homura. Note that Go does not implement `-buildmode=plugin` on Windows.

If xi is the best place for a new variation, e.g. because it is very common in Apophysis, then make sure to follow the following steps:

- Create the function type with `Calc` and `Prep` methods.
- Register any factories in a `func init()`. There must be at least one factory, to ensure that the type implements `xirho.F`.
- Add it to this README, in the list near the top.
- Create Flame parsers in package `xirho/encoding/flame` and add it to the README there.