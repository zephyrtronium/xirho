// Package fapi creates a generic public API for xirho functions.
//
// The For function uses reflection to gather a list of modifiable parameters.
// This can be used to implement serialization formats or to provide user
// interfaces for functions. Parameter types are based on semantics rather than
// on representation to allow for more natural user interfaces.
//
package fapi

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/zephyrtronium/xirho"
)

// For collects a xirho.Func's exported parameters. Each returned parameter has
// its name set according to the first comma-separated section of the field's
// "xirho" struct tag, defaulting to the field name. E.g., the JuliaN variation
// is defined as such:
//
//		type JuliaN struct {
//			Power xirho.Int  `xirho:"power"`
//			Dist  xirho.Real `xirho:"dist"`
//		}
//
// Certain parameters provide additional options; see the documentation for
// each for details.
func For(f xirho.Func) []Param {
	var r []Param
	val := reflect.ValueOf(f)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()
	if typ.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < typ.NumField(); i++ {
		if p := getParam(typ.Field(i), val.Field(i)); p != nil {
			r = append(r, p)
		}
	}
	return r
}

// getParam gets a Param for a single struct field, or nil if it does not have
// a xirho.Param type.
func getParam(f reflect.StructField, v reflect.Value) Param {
	if !v.CanInterface() {
		// If we can't Interface(), then the field is unexported.
		return nil
	}
	val := v.Addr().Interface()
	tag := strings.Split(f.Tag.Get("xirho"), ",")
	name := pname(tag, f.Name)
	switch f.Type {
	case rFlag:
		return flagFor(name, val.(*xirho.Flag))
	case rList:
		if len(tag) < 3 {
			panic(fmt.Errorf("xirho: list must have at least 2 options (after name); have %q", tag))
		}
		return listFor(name, val.(*xirho.List), tag[1:]...)
	case rInt:
		switch len(tag) {
		case 0, 1:
			return intFor(name, val.(*xirho.Int), false, 0, 0)
		case 2:
			panic(fmt.Errorf("xirho: 2 tag fields in %q is probably a mistake; need 0, 1, or 3", f.Tag))
		default:
			var lo, hi int64
			var err error
			if lo, err = strconv.ParseInt(tag[1], 0, 64); err != nil {
				panic(fmt.Errorf("xirho: error parsing Int lo: %w", err))
			}
			if hi, err = strconv.ParseInt(tag[2], 0, 64); err != nil {
				panic(fmt.Errorf("xirho: error parsing Int hi: %w", err))
			}
			if lo > hi {
				panic(fmt.Errorf("xirho: Int lo > hi"))
			}
			return intFor(name, val.(*xirho.Int), true, xirho.Int(lo), xirho.Int(hi))
		}
	case rAngle:
		return angleFor(name, val.(*xirho.Angle))
	case rReal:
		switch len(tag) {
		case 0, 1:
			return realFor(name, val.(*xirho.Real), false, 0, 0)
		case 2:
			panic(fmt.Errorf("xirho: 2 tag fields in %q is probably a mistake; need 0, 1, or 3", f.Tag))
		default:
			var lo, hi float64
			var err error
			if lo, err = strconv.ParseFloat(tag[1], 64); err != nil {
				panic(fmt.Errorf("xirho: error parsing Real lo: %w", err))
			}
			if hi, err = strconv.ParseFloat(tag[2], 64); err != nil {
				panic(fmt.Errorf("xirho: error parsing Real hi: %w", err))
			}
			if lo > hi {
				panic(fmt.Errorf("xirho: Real lo > hi"))
			}
			return realFor(name, val.(*xirho.Real), true, xirho.Real(lo), xirho.Real(hi))
		}
	case rComplex:
		return complexFor(name, val.(*xirho.Complex))
	case rVec3:
		return vec3For(name, val.(*xirho.Vec3))
	case rAffine:
		return affineFor(name, val.(*xirho.Affine))
	case rFunc:
		opt := false
		if len(tag) >= 2 {
			if tag[1] != "optional" {
				panic(fmt.Errorf(`xirho: bad value %q for func tag; did you mean "optional"?`, tag[1]))
			}
			opt = true
		}
		return funcFor(name, opt, val.(*xirho.Func))
	case rFuncList:
		return funcListFor(name, val.(*xirho.FuncList))
	default:
		return nil
	}
}

// pname gets the name of a parameter.
func pname(tag []string, name string) string {
	if len(tag) > 0 && tag[0] != "" {
		return tag[0]
	}
	return name
}

// Reflected Param types.
var (
	rFlag     = reflect.TypeOf(xirho.Flag(false))
	rList     = reflect.TypeOf(xirho.List(0))
	rInt      = reflect.TypeOf(xirho.Int(0))
	rAngle    = reflect.TypeOf(xirho.Angle(0))
	rReal     = reflect.TypeOf(xirho.Real(0))
	rComplex  = reflect.TypeOf(xirho.Complex(0))
	rVec3     = reflect.TypeOf(xirho.Vec3{})
	rAffine   = reflect.TypeOf(xirho.Affine{})
	rFunc     = reflect.TypeOf((*xirho.Func)(nil)).Elem()
	rFuncList = reflect.TypeOf(xirho.FuncList(nil))
)
