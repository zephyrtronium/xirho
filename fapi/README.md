# xirho/fapi

Package fapi creates a generic public API for xirho function types.

Fapi generates an abstracted set and get layer over function parameters of Flag, List, Int, Angle, Real, Complex, Vec3, Affine, Func, and FuncList types from package xirho. There is a corresponding type in fapi for each, meaning that type switches can enumerate every possibility to use the complete API of any xirho function.

Typical use of package fapi will look something like this:

```go
api := fapi.For(someFunc)
for _, param := range api {
    switch p := param.(type) {
    case fapi.Flag:
        // p.Get(), p.Set(), p.Name()
    case fapi.List:
        // p.Get(), p.Set(), p.Name(), p.String(), p.Opts()
    // ...
    default:
        panic("unknown parameter type")
    }
}
```

For a more complete example, see xirho/encoding.Unmarshal.
