package fapi

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/zephyrtronium/xirho"
)

// For collects a xirho.F's exported parameters. Each returned parameter has
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
func For(f xirho.F) []Param {
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
		opts := []string{}
		if len(tag) >= 1 {
			opts = tag[1:]
		}
		return listFor(name, val.(*xirho.List), opts...)
	case rInt:
		bdd := len(tag) == 3
		var lo, hi int64
		if bdd {
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
		}
		return intFor(name, val.(*xirho.Int), bdd, lo, hi)
	case rAngle:
		return angleFor(name, val.(*xirho.Angle))
	case rReal:
		bdd := len(tag) == 3
		var lo, hi float64
		if bdd {
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
		}
		return realFor(name, val.(*xirho.Real), bdd, lo, hi)
	case rComplex:
		return complexFor(name, val.(*xirho.Complex))
	case rVec3:
		return vec3For(name, val.(*xirho.Vec3))
	case rAffine:
		return affineFor(name, val.(*xirho.Affine))
	case rFunc:
		opt := len(tag) == 2 && tag[1] == "optional"
		return funcFor(name, opt, val.(*xirho.Func))
	case rFuncList:
		return funcListFor(name, val.(*xirho.FuncList))
	default:
		return nil
	}
}

// pname gets the name of a parameter.
func pname(tag []string, name string) string {
	if len(tag) > 0 {
		return tag[0]
	}
	return name
}

// Reflected Param types.
var (
	rFlag     = reflect.TypeOf(zFlag)
	rList     = reflect.TypeOf(zList)
	rInt      = reflect.TypeOf(zInt)
	rAngle    = reflect.TypeOf(zAngle)
	rReal     = reflect.TypeOf(zReal)
	rComplex  = reflect.TypeOf(zComplex)
	rVec3     = reflect.TypeOf(zVec3)
	rAffine   = reflect.TypeOf(zAffine)
	rFunc     = reflect.TypeOf(zFunc)
	rFuncList = reflect.TypeOf(zFuncList)
)

// Zero values for Param types.
var (
	zFlag     xirho.Flag
	zList     xirho.List
	zInt      xirho.Int
	zAngle    xirho.Angle
	zReal     xirho.Real
	zComplex  xirho.Complex
	zVec3     xirho.Vec3
	zAffine   xirho.Affine
	zFunc     xirho.Func
	zFuncList xirho.FuncList
)
