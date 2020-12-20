// Package xi provides function registration and implementations for xirho.
//
// External packages adding new function types only need to Register a factory
// here to automatically work with xirho serialization and user interfaces.
// Conversely, user interfaces will be primarily interested in the New, NameOf,
// and Names functions, as most of the rest of the work is done in package
// fapi.
//
// Lastly, some users may want to use the function implementation types
// directly, either to decode systems serialized in another format (e.g.
// package encoding/flame) or to build a system by hand. As a general rule,
// if a function type has any parameters, it implements xirho.Func as a
// pointer, and if it has none, then it implements it as a (size zero) value.
//
package xi
