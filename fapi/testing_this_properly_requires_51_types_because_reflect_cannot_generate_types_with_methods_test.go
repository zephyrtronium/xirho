package fapi_test

import (
	"math"
	"reflect"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/fapi"
	"github.com/zephyrtronium/xirho/xi"
)

type (
	testFlag struct {
		V xirho.Flag `xirho:"test"` // ok
	}
	testFlagUnnamed struct {
		V xirho.Flag // ok, named V
	}
	testFlagExtra struct {
		V xirho.Flag `xirho:"test,ignore"` // ok
	}

	testListEmpty struct {
		V xirho.List `xirho:"test"` // error
	}
	testList1 struct {
		V xirho.List `xirho:"test,1"` // error
	}
	testList2 struct {
		V xirho.List `xirho:"test,1,2"` // ok
	}
	testList10 struct {
		V xirho.List `xirho:"test,1,2,3,4,5,6,7,8,9,10"` // ok
	}
	testListUnnamed struct {
		V xirho.List `xirho:",1,2"` // ok, named V
	}

	testInt struct {
		V xirho.Int `xirho:"test"` // ok
	}
	testIntUnnamed struct {
		V xirho.Int // ok, named V
	}
	testIntBounded struct {
		V xirho.Int `xirho:"test,-1,1"` // ok
	}
	testIntBoundedUnnamed struct {
		V xirho.Int `xirho:",-1,1"` // ok, named V
	}
	testIntBoundedBadLower struct {
		V xirho.Int `xirho:"test,-1.5,1"` // error
	}
	testIntBoundedBadUpper struct {
		V xirho.Int `xirho:"test,-1,1.5"` // error
	}
	testIntBoundedLowerOnly struct {
		V xirho.Int `xirho:"test,-1"` // error
	}
	testIntBoundedUpperOnly struct {
		V xirho.Int `xirho:"test,,1"` // error
	}
	testIntBoundedEmpty struct {
		V xirho.Int `xirho:"test,1,-1"` // error
	}
	testIntBoundedSingleton struct {
		V xirho.Int `xirho:"test,0,0"` // ok
	}
	testInt3 struct {
		V xirho.Int `xirho:"test,-1,1,ignore"` // ok
	}

	testAngle struct {
		V xirho.Angle `xirho:"test"` // ok
	}
	testAngleUnnamed struct {
		V xirho.Angle // ok, named V
	}
	testAngleExtra struct {
		V xirho.Angle `xirho:"test,ignore"` // ok
	}

	testReal struct {
		V xirho.Real `xirho:"test"` // ok
	}
	testRealUnnamed struct {
		V xirho.Real // ok, named V
	}
	testRealBounded struct {
		V xirho.Real `xirho:"test,-1,1"` // ok
	}
	testRealBoundedUnnamed struct {
		V xirho.Real `xirho:",-1,1"` // ok, named V
	}
	testRealBoundedBadLower struct {
		V xirho.Real `xirho:"test,-pi,0"` // error
	}
	testRealBoundedBadUpper struct {
		V xirho.Real `xirho:"test,0,pi"` // error
	}
	testRealBoundedLowerOnly struct {
		V xirho.Real `xirho:"test,-1"` // error
	}
	testRealBoundedUpperOnly struct {
		V xirho.Real `xirho:"test,,1"` // error
	}
	testRealBoundedEmpty struct {
		V xirho.Real `xirho:"test,1,-1"` // error
	}
	testRealBoundedSingleton struct {
		V xirho.Real `xirho:"test,0,0"` // ok
	}
	testReal3 struct {
		V xirho.Real `xirho:"test,-1,1,ignore"` // ok
	}

	testComplex struct {
		V xirho.Complex `xirho:"test"` // ok
	}
	testComplexUnnamed struct {
		V xirho.Complex // ok, named V
	}
	testComplexExtra struct {
		V xirho.Complex `xirho:"test,ignore"` // ok
	}

	testVec3 struct {
		V xirho.Vec3 `xirho:"test"` // ok
	}
	testVec3Unnamed struct {
		V xirho.Vec3 // ok, named V
	}
	testVec3Extra struct {
		V xirho.Vec3 `xirho:"test,ignore"` // ok
	}

	testAffine struct {
		V xirho.Affine `xirho:"test"` // ok
	}
	testAffineUnnamed struct {
		V xirho.Affine // ok, named V
	}
	testAffineExtra struct {
		V xirho.Affine `xirho:"test,ignore"` // ok
	}

	testFunc struct {
		V xirho.Func `xirho:"test"` // ok
	}
	testFuncUnnamed struct {
		V xirho.Func // ok, named V
	}
	testFuncOptional struct {
		V xirho.Func `xirho:"test,optional"` // ok
	}
	testFuncOptionalUnnamed struct {
		V xirho.Func `xirho:",optional"` // ok, named V
	}
	testFuncBad struct {
		V xirho.Func `xirho:"test,bad"` // error
	}
	testFuncExtra struct {
		V xirho.Func `xirho:"test,optional,ignore"` // ok
	}

	testFuncList struct {
		V xirho.FuncList `xirho:"test"` // ok
	}
	testFuncListUnnamed struct {
		V xirho.FuncList // ok, named V
	}
	testFuncListExtra struct {
		V xirho.FuncList `xirho:"test,ignore"` // ok
	}
)

type setCase struct {
	set, get interface{}
	err      error
}

var typeCases = []struct {
	// name is the name of the testing type.
	name string
	// v is the testing type.
	v xirho.Func
	// param is the expected param type from the first (and only) field. If
	// nil, expect a panic instead.
	param reflect.Type
	// field is the expected field of the param.
	field string
	// set is a list of values to try setting and associated gets and errors.
	// They must be tried in listed order, so that the gets associated with
	// failed sets have known values. The set value can be asserted to the
	// appropriate type for the parameter's setter, with the exception of Func
	// parameters with nil values. The get value must be checked with
	// reflect.DeepEqual or similar, as it may have a slice type.
	set []setCase
}{
	{
		name:  "flag",
		v:     new(testFlag),
		param: reflect.TypeOf(fapi.Flag{}),
		field: "test",
		set: []setCase{
			{set: true, get: true, err: nil},
			{set: false, get: false, err: nil},
		},
	},
	{
		name:  "flagUnnamed",
		v:     new(testFlagUnnamed),
		param: reflect.TypeOf(fapi.Flag{}),
		field: "V",
		set: []setCase{
			{set: true, get: true, err: nil},
			{set: false, get: false, err: nil},
		},
	},
	{
		name:  "flagExtra",
		v:     new(testFlagExtra),
		param: reflect.TypeOf(fapi.Flag{}),
		field: "test",
		set: []setCase{
			{set: true, get: true, err: nil},
			{set: false, get: false, err: nil},
		},
	},
	{
		name:  "listEmpty",
		v:     new(testListEmpty),
		param: nil,
	},
	{
		name:  "list1",
		v:     new(testList1),
		param: nil,
	},
	{
		name:  "list2",
		v:     new(testList2),
		param: reflect.TypeOf(fapi.List{}),
		field: "test",
		set: []setCase{
			{set: 0, get: 0, err: nil},
			{set: 1, get: 1, err: nil},
			{set: 2, get: 1, err: fapi.OutOfBoundsInt{}},
			{set: -1, get: 1, err: fapi.OutOfBoundsInt{}},
		},
	},
	{
		name:  "list10",
		v:     new(testList10),
		param: reflect.TypeOf(fapi.List{}),
		field: "test",
		set: []setCase{
			{set: 0, get: 0, err: nil},
			{set: 1, get: 1, err: nil},
			{set: 2, get: 2, err: nil},
			{set: 10, get: 2, err: fapi.OutOfBoundsInt{}},
			{set: -1, get: 2, err: fapi.OutOfBoundsInt{}},
		},
	},
	{
		name:  "listUnnamed",
		v:     new(testListUnnamed),
		param: reflect.TypeOf(fapi.List{}),
		field: "V",
		set: []setCase{
			{set: 0, get: 0, err: nil},
			{set: 1, get: 1, err: nil},
			{set: 2, get: 1, err: fapi.OutOfBoundsInt{}},
			{set: -1, get: 1, err: fapi.OutOfBoundsInt{}},
		},
	},
	{
		name:  "int",
		v:     new(testInt),
		param: reflect.TypeOf(fapi.Int{}),
		field: "test",
		set: []setCase{
			{set: int64(-1), get: int64(-1), err: nil},
			{set: int64(0), get: int64(0), err: nil},
			{set: int64(1), get: int64(1), err: nil},
			{set: int64(^uint(0) >> 1), get: int64(^uint(0) >> 1), err: nil},
			{set: int64(-1 << 63), get: int64(-1 << 63), err: nil},
		},
	},
	{
		name:  "intUnnamed",
		v:     new(testIntUnnamed),
		param: reflect.TypeOf(fapi.Int{}),
		field: "V",
		set: []setCase{
			{set: int64(-1), get: int64(-1), err: nil},
			{set: int64(0), get: int64(0), err: nil},
			{set: int64(1), get: int64(1), err: nil},
			{set: int64(^uint(0) >> 1), get: int64(^uint(0) >> 1), err: nil},
			{set: int64(-1 << 63), get: int64(-1 << 63), err: nil},
		},
	},
	{
		name:  "intBounded",
		v:     new(testIntBounded),
		param: reflect.TypeOf(fapi.Int{}),
		field: "test",
		set: []setCase{
			{set: int64(-1), get: int64(-1), err: nil},
			{set: int64(0), get: int64(0), err: nil},
			{set: int64(1), get: int64(1), err: nil},
			{set: int64(2), get: int64(1), err: fapi.OutOfBoundsInt{}},
			{set: int64(-2), get: int64(1), err: fapi.OutOfBoundsInt{}},
		},
	},
	{
		name:  "intBoundedUnnamed",
		v:     new(testIntBoundedUnnamed),
		param: reflect.TypeOf(fapi.Int{}),
		field: "V",
		set: []setCase{
			{set: int64(-1), get: int64(-1), err: nil},
			{set: int64(0), get: int64(0), err: nil},
			{set: int64(1), get: int64(1), err: nil},
			{set: int64(2), get: int64(1), err: fapi.OutOfBoundsInt{}},
			{set: int64(-2), get: int64(1), err: fapi.OutOfBoundsInt{}},
		},
	},
	{
		name:  "intBoundedBadLower",
		v:     new(testIntBoundedBadLower),
		param: nil,
	},
	{
		name:  "intBoundedBadUpper",
		v:     new(testIntBoundedBadUpper),
		param: nil,
	},
	{
		name:  "intBoundedLowerOnly",
		v:     new(testIntBoundedLowerOnly),
		param: nil,
	},
	{
		name:  "intBoundedUpperOnly",
		v:     new(testIntBoundedUpperOnly),
		param: nil,
	},
	{
		name:  "intBoundedEmpty",
		v:     new(testIntBoundedEmpty),
		param: nil,
	},
	{
		name:  "intBoundedSingleton",
		v:     new(testIntBoundedSingleton),
		param: reflect.TypeOf(fapi.Int{}),
		field: "test",
		set: []setCase{
			{set: int64(0), get: int64(0), err: nil},
			{set: int64(1), get: int64(0), err: fapi.OutOfBoundsInt{}},
			{set: int64(-1), get: int64(0), err: fapi.OutOfBoundsInt{}},
		},
	},
	{
		name:  "int3",
		v:     new(testInt3),
		param: reflect.TypeOf(fapi.Int{}),
		field: "test",
		set: []setCase{
			{set: int64(-1), get: int64(-1), err: nil},
			{set: int64(0), get: int64(0), err: nil},
			{set: int64(1), get: int64(1), err: nil},
			{set: int64(2), get: int64(1), err: fapi.OutOfBoundsInt{}},
			{set: int64(-2), get: int64(1), err: fapi.OutOfBoundsInt{}},
		},
	},
	{
		name:  "angle",
		v:     new(testAngle),
		param: reflect.TypeOf(fapi.Angle{}),
		field: "test",
		set: []setCase{
			{set: 0.0, get: 0.0, err: nil},
			{set: math.Pi, get: math.Pi, err: nil},
			{set: math.Nextafter(-math.Pi, 0), get: math.Nextafter(-math.Pi, 0), err: nil},
			{set: -math.Pi, get: math.Pi, err: nil},
			{set: 2 * math.Pi, get: 0.0, err: nil},
			{set: math.Inf(0), get: 0.0, err: fapi.NotFinite{}},
			{set: math.NaN(), get: 0.0, err: fapi.NotFinite{}},
		},
	},
	{
		name:  "angleUnnamed",
		v:     new(testAngleUnnamed),
		param: reflect.TypeOf(fapi.Angle{}),
		field: "V",
		set: []setCase{
			{set: 0.0, get: 0.0, err: nil},
			{set: math.Pi, get: math.Pi, err: nil},
			{set: math.Nextafter(-math.Pi, 0), get: math.Nextafter(-math.Pi, 0), err: nil},
			{set: -math.Pi, get: math.Pi, err: nil},
			{set: 2 * math.Pi, get: 0.0, err: nil},
			{set: math.Inf(0), get: 0.0, err: fapi.NotFinite{}},
			{set: math.NaN(), get: 0.0, err: fapi.NotFinite{}},
		},
	},
	{
		name:  "angleExtra",
		v:     new(testAngleExtra),
		param: reflect.TypeOf(fapi.Angle{}),
		field: "test",
		set: []setCase{
			{set: 0.0, get: 0.0, err: nil},
			{set: math.Pi, get: math.Pi, err: nil},
			{set: math.Nextafter(-math.Pi, 0), get: math.Nextafter(-math.Pi, 0), err: nil},
			{set: -math.Pi, get: math.Pi, err: nil},
			{set: 2 * math.Pi, get: 0.0, err: nil},
			{set: math.Inf(0), get: 0.0, err: fapi.NotFinite{}},
			{set: math.NaN(), get: 0.0, err: fapi.NotFinite{}},
		},
	},
	{
		name:  "real",
		v:     new(testReal),
		param: reflect.TypeOf(fapi.Real{}),
		field: "test",
		set: []setCase{
			{set: 1.0, get: 1.0, err: nil},
			{set: -1.0, get: -1.0, err: nil},
			{set: math.Nextafter(math.Inf(0), 0), get: math.Nextafter(math.Inf(0), 0), err: nil},
			{set: 0.0, get: 0.0, err: nil},
			{set: math.Inf(0), get: 0.0, err: fapi.NotFinite{}},
			{set: math.NaN(), get: 0.0, err: fapi.NotFinite{}},
		},
	},
	{
		name:  "realUnnamed",
		v:     new(testRealUnnamed),
		param: reflect.TypeOf(fapi.Real{}),
		field: "V",
		set: []setCase{
			{set: 1.0, get: 1.0, err: nil},
			{set: -1.0, get: -1.0, err: nil},
			{set: math.Nextafter(math.Inf(0), 0), get: math.Nextafter(math.Inf(0), 0), err: nil},
			{set: 0.0, get: 0.0, err: nil},
			{set: math.Inf(0), get: 0.0, err: fapi.NotFinite{}},
			{set: math.NaN(), get: 0.0, err: fapi.NotFinite{}},
		},
	},
	{
		name:  "realBounded",
		v:     new(testRealBounded),
		param: reflect.TypeOf(fapi.Real{}),
		field: "test",
		set: []setCase{
			{set: -1.0, get: -1.0, err: nil},
			{set: 0.0, get: 0.0, err: nil},
			{set: 1.0, get: 1.0, err: nil},
			{set: math.Nextafter(1, 2), get: 1.0, err: fapi.OutOfBoundsReal{}},
			{set: math.Nextafter(-1, -2), get: 1.0, err: fapi.OutOfBoundsReal{}},
			{set: math.Inf(0), get: 1.0, err: fapi.NotFinite{}},
			{set: math.NaN(), get: 1.0, err: fapi.NotFinite{}},
		},
	},
	{
		name:  "realBoundedUnnamed",
		v:     new(testRealBoundedUnnamed),
		param: reflect.TypeOf(fapi.Real{}),
		field: "V",
		set: []setCase{
			{set: -1.0, get: -1.0, err: nil},
			{set: 0.0, get: 0.0, err: nil},
			{set: 1.0, get: 1.0, err: nil},
			{set: math.Nextafter(1, 2), get: 1.0, err: fapi.OutOfBoundsReal{}},
			{set: math.Nextafter(-1, -2), get: 1.0, err: fapi.OutOfBoundsReal{}},
			{set: math.Inf(0), get: 1.0, err: fapi.NotFinite{}},
			{set: math.NaN(), get: 1.0, err: fapi.NotFinite{}},
		},
	},
	{
		name:  "realBoundedBadLower",
		v:     new(testRealBoundedBadLower),
		param: nil,
	},
	{
		name:  "realBoundedBadUpper",
		v:     new(testRealBoundedBadUpper),
		param: nil,
	},
	{
		name:  "realBoundedLowerOnly",
		v:     new(testRealBoundedLowerOnly),
		param: nil,
	},
	{
		name:  "realBoundedUpperOnly",
		v:     new(testRealBoundedUpperOnly),
		param: nil,
	},
	{
		name:  "realBoundedEmpty",
		v:     new(testRealBoundedEmpty),
		param: nil,
	},
	{
		name:  "realBoundedSingleton",
		v:     new(testRealBoundedSingleton),
		param: reflect.TypeOf(fapi.Real{}),
		field: "test",
		set: []setCase{
			{set: 0.0, get: 0.0, err: nil},
			{set: math.Nextafter(0, -1), get: 0.0, err: fapi.OutOfBoundsReal{}},
			{set: math.Nextafter(0, 1), get: 0.0, err: fapi.OutOfBoundsReal{}},
			{set: math.Inf(0), get: 0.0, err: fapi.NotFinite{}},
			{set: math.NaN(), get: 0.0, err: fapi.NotFinite{}},
		},
	},
	{
		name:  "real3",
		v:     new(testReal3),
		param: reflect.TypeOf(fapi.Real{}),
		field: "test",
		set: []setCase{
			{set: -1.0, get: -1.0, err: nil},
			{set: 0.0, get: 0.0, err: nil},
			{set: 1.0, get: 1.0, err: nil},
			{set: math.Nextafter(1, 2), get: 1.0, err: fapi.OutOfBoundsReal{}},
			{set: math.Nextafter(-1, -2), get: 1.0, err: fapi.OutOfBoundsReal{}},
			{set: math.Inf(0), get: 1.0, err: fapi.NotFinite{}},
			{set: math.NaN(), get: 1.0, err: fapi.NotFinite{}},
		},
	},
	{
		name:  "complex",
		v:     new(testComplex),
		param: reflect.TypeOf(fapi.Complex{}),
		field: "test",
		set: []setCase{
			{set: complex(1, 1), get: complex(1, 1), err: nil},
			{set: complex(0, math.Inf(0)), get: complex(1, 1), err: fapi.NotFinite{}},
			{set: complex(math.NaN(), 0), get: complex(1, 1), err: fapi.NotFinite{}},
		},
	},
	{
		name:  "complexUnnamed",
		v:     new(testComplexUnnamed),
		param: reflect.TypeOf(fapi.Complex{}),
		field: "V",
		set: []setCase{
			{set: complex(1, 1), get: complex(1, 1), err: nil},
			{set: complex(0, math.Inf(0)), get: complex(1, 1), err: fapi.NotFinite{}},
			{set: complex(math.NaN(), 0), get: complex(1, 1), err: fapi.NotFinite{}},
		},
	},
	{
		name:  "complexExtra",
		v:     new(testComplexExtra),
		param: reflect.TypeOf(fapi.Complex{}),
		field: "test",
		set: []setCase{
			{set: complex(1, 1), get: complex(1, 1), err: nil},
			{set: complex(0, math.Inf(0)), get: complex(1, 1), err: fapi.NotFinite{}},
			{set: complex(math.NaN(), 0), get: complex(1, 1), err: fapi.NotFinite{}},
		},
	},
	{
		name:  "vec3",
		v:     new(testVec3),
		param: reflect.TypeOf(fapi.Vec3{}),
		field: "test",
		set: []setCase{
			{set: xirho.Vec3{1, 1, 1}, get: xirho.Vec3{1, 1, 1}, err: nil},
			{set: xirho.Vec3{math.Inf(0), 0, 0}, get: xirho.Vec3{1, 1, 1}, err: fapi.NotFinite{}},
			{set: xirho.Vec3{0, math.Inf(-1), 0}, get: xirho.Vec3{1, 1, 1}, err: fapi.NotFinite{}},
			{set: xirho.Vec3{0, 0, math.NaN()}, get: xirho.Vec3{1, 1, 1}, err: fapi.NotFinite{}},
		},
	},
	{
		name:  "vec3Unnamed",
		v:     new(testVec3Unnamed),
		param: reflect.TypeOf(fapi.Vec3{}),
		field: "V",
		set: []setCase{
			{set: xirho.Vec3{1, 1, 1}, get: xirho.Vec3{1, 1, 1}, err: nil},
			{set: xirho.Vec3{math.Inf(0), 0, 0}, get: xirho.Vec3{1, 1, 1}, err: fapi.NotFinite{}},
			{set: xirho.Vec3{0, math.Inf(-1), 0}, get: xirho.Vec3{1, 1, 1}, err: fapi.NotFinite{}},
			{set: xirho.Vec3{0, 0, math.NaN()}, get: xirho.Vec3{1, 1, 1}, err: fapi.NotFinite{}},
		},
	},
	{
		name:  "vec3Extra",
		v:     new(testVec3Extra),
		param: reflect.TypeOf(fapi.Vec3{}),
		field: "test",
		set: []setCase{
			{set: xirho.Vec3{1, 1, 1}, get: xirho.Vec3{1, 1, 1}, err: nil},
			{set: xirho.Vec3{math.Inf(0), 0, 0}, get: xirho.Vec3{1, 1, 1}, err: fapi.NotFinite{}},
			{set: xirho.Vec3{0, math.Inf(-1), 0}, get: xirho.Vec3{1, 1, 1}, err: fapi.NotFinite{}},
			{set: xirho.Vec3{0, 0, math.NaN()}, get: xirho.Vec3{1, 1, 1}, err: fapi.NotFinite{}},
		},
	},
	{
		name:  "affine",
		v:     new(testAffine),
		param: reflect.TypeOf(fapi.Affine{}),
		field: "test",
		set: []setCase{
			{set: xirho.Eye(), get: xirho.Eye(), err: nil},
			{set: xirho.Affine{0: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{1: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{2: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{3: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{4: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{5: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{6: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{7: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{8: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{9: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{10: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{11: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
		},
	},
	{
		name:  "affineUnnamed",
		v:     new(testAffineUnnamed),
		param: reflect.TypeOf(fapi.Affine{}),
		field: "V",
		set: []setCase{
			{set: xirho.Eye(), get: xirho.Eye(), err: nil},
			{set: xirho.Affine{0: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{1: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{2: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{3: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{4: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{5: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{6: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{7: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{8: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{9: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{10: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{11: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
		},
	},
	{
		name:  "affineExtra",
		v:     new(testAffineExtra),
		param: reflect.TypeOf(fapi.Affine{}),
		field: "test",
		set: []setCase{
			{set: xirho.Eye(), get: xirho.Eye(), err: nil},
			{set: xirho.Affine{0: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{1: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{2: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{3: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{4: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{5: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{6: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{7: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{8: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{9: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{10: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
			{set: xirho.Affine{11: math.Inf(0)}, get: xirho.Eye(), err: fapi.NotFinite{}},
		},
	},
	{
		name:  "func",
		v:     new(testFunc),
		param: reflect.TypeOf(fapi.Func{}),
		field: "test",
		set: []setCase{
			{set: xi.Spherical{}, get: xi.Spherical{}, err: nil},
			{set: nil, get: xi.Spherical{}, err: fapi.NotOptional{}},
		},
	},
	{
		name:  "funcUnnamed",
		v:     new(testFuncUnnamed),
		param: reflect.TypeOf(fapi.Func{}),
		field: "V",
		set: []setCase{
			{set: xi.Spherical{}, get: xi.Spherical{}, err: nil},
			{set: xirho.Func(nil), get: xi.Spherical{}, err: fapi.NotOptional{}},
		},
	},
	{
		name:  "funcOptional",
		v:     new(testFuncOptional),
		param: reflect.TypeOf(fapi.Func{}),
		field: "test",
		set: []setCase{
			{set: xi.Spherical{}, get: xi.Spherical{}, err: nil},
			{set: xirho.Func(nil), get: xirho.Func(nil), err: nil},
		},
	},
	{
		name:  "funcOptionalUnnamed",
		v:     new(testFuncOptionalUnnamed),
		param: reflect.TypeOf(fapi.Func{}),
		field: "V",
		set: []setCase{
			{set: xi.Spherical{}, get: xi.Spherical{}, err: nil},
			{set: xirho.Func(nil), get: xirho.Func(nil), err: nil},
		},
	},
	{
		name:  "funcBad",
		v:     new(testFuncBad),
		param: nil,
	},
	{
		name:  "funcExtra",
		v:     new(testFuncExtra),
		param: reflect.TypeOf(fapi.Func{}),
		field: "test",
		set: []setCase{
			{set: xi.Spherical{}, get: xi.Spherical{}, err: nil},
			{set: xirho.Func(nil), get: xirho.Func(nil), err: nil},
		},
	},
	{
		name:  "funcList",
		v:     new(testFuncList),
		param: reflect.TypeOf(fapi.FuncList{}),
		field: "test",
		set: []setCase{
			{set: xirho.FuncList{xi.Spherical{}}, get: xirho.FuncList{xi.Spherical{}}, err: nil},
			{set: xirho.FuncList(nil), get: xirho.FuncList(nil), err: nil},
		},
	},
	{
		name:  "funcListUnnamed",
		v:     new(testFuncListUnnamed),
		param: reflect.TypeOf(fapi.FuncList{}),
		field: "V",
		set: []setCase{
			{set: xirho.FuncList{xi.Spherical{}}, get: xirho.FuncList{xi.Spherical{}}, err: nil},
			{set: xirho.FuncList(nil), get: xirho.FuncList(nil), err: nil},
		},
	},
	{
		name:  "funcListExtra",
		v:     new(testFuncListExtra),
		param: reflect.TypeOf(fapi.FuncList{}),
		field: "test",
		set: []setCase{
			{set: xirho.FuncList{xi.Spherical{}}, get: xirho.FuncList{xi.Spherical{}}, err: nil},
			{set: xirho.FuncList(nil), get: xirho.FuncList(nil), err: nil},
		},
	},
}

func (*testFlag) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt                 { return xirho.Pt{} }
func (*testFlagUnnamed) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt          { return xirho.Pt{} }
func (*testFlagExtra) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt            { return xirho.Pt{} }
func (*testListEmpty) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt            { return xirho.Pt{} }
func (*testList1) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt                { return xirho.Pt{} }
func (*testList2) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt                { return xirho.Pt{} }
func (*testList10) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt               { return xirho.Pt{} }
func (*testListUnnamed) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt          { return xirho.Pt{} }
func (*testInt) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt                  { return xirho.Pt{} }
func (*testIntUnnamed) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt           { return xirho.Pt{} }
func (*testIntBounded) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt           { return xirho.Pt{} }
func (*testIntBoundedUnnamed) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt    { return xirho.Pt{} }
func (*testIntBoundedBadLower) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt   { return xirho.Pt{} }
func (*testIntBoundedBadUpper) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt   { return xirho.Pt{} }
func (*testIntBoundedLowerOnly) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt  { return xirho.Pt{} }
func (*testIntBoundedUpperOnly) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt  { return xirho.Pt{} }
func (*testIntBoundedEmpty) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt      { return xirho.Pt{} }
func (*testIntBoundedSingleton) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt  { return xirho.Pt{} }
func (*testInt3) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt                 { return xirho.Pt{} }
func (*testAngle) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt                { return xirho.Pt{} }
func (*testAngleUnnamed) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt         { return xirho.Pt{} }
func (*testAngleExtra) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt           { return xirho.Pt{} }
func (*testReal) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt                 { return xirho.Pt{} }
func (*testRealUnnamed) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt          { return xirho.Pt{} }
func (*testRealBounded) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt          { return xirho.Pt{} }
func (*testRealBoundedUnnamed) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt   { return xirho.Pt{} }
func (*testRealBoundedBadLower) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt  { return xirho.Pt{} }
func (*testRealBoundedBadUpper) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt  { return xirho.Pt{} }
func (*testRealBoundedLowerOnly) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt { return xirho.Pt{} }
func (*testRealBoundedUpperOnly) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt { return xirho.Pt{} }
func (*testRealBoundedEmpty) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt     { return xirho.Pt{} }
func (*testRealBoundedSingleton) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt { return xirho.Pt{} }
func (*testReal3) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt                { return xirho.Pt{} }
func (*testComplex) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt              { return xirho.Pt{} }
func (*testComplexUnnamed) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt       { return xirho.Pt{} }
func (*testComplexExtra) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt         { return xirho.Pt{} }
func (*testVec3) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt                 { return xirho.Pt{} }
func (*testVec3Unnamed) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt          { return xirho.Pt{} }
func (*testVec3Extra) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt            { return xirho.Pt{} }
func (*testAffine) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt               { return xirho.Pt{} }
func (*testAffineUnnamed) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt        { return xirho.Pt{} }
func (*testAffineExtra) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt          { return xirho.Pt{} }
func (*testFunc) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt                 { return xirho.Pt{} }
func (*testFuncUnnamed) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt          { return xirho.Pt{} }
func (*testFuncOptional) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt         { return xirho.Pt{} }
func (*testFuncOptionalUnnamed) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt  { return xirho.Pt{} }
func (*testFuncBad) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt              { return xirho.Pt{} }
func (*testFuncExtra) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt            { return xirho.Pt{} }
func (*testFuncList) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt             { return xirho.Pt{} }
func (*testFuncListUnnamed) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt      { return xirho.Pt{} }
func (*testFuncListExtra) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt        { return xirho.Pt{} }

func (*testFlag) Prep()                 {}
func (*testFlagUnnamed) Prep()          {}
func (*testFlagExtra) Prep()            {}
func (*testListEmpty) Prep()            {}
func (*testList1) Prep()                {}
func (*testList2) Prep()                {}
func (*testList10) Prep()               {}
func (*testListUnnamed) Prep()          {}
func (*testInt) Prep()                  {}
func (*testIntUnnamed) Prep()           {}
func (*testIntBounded) Prep()           {}
func (*testIntBoundedUnnamed) Prep()    {}
func (*testIntBoundedBadLower) Prep()   {}
func (*testIntBoundedBadUpper) Prep()   {}
func (*testIntBoundedLowerOnly) Prep()  {}
func (*testIntBoundedUpperOnly) Prep()  {}
func (*testIntBoundedEmpty) Prep()      {}
func (*testIntBoundedSingleton) Prep()  {}
func (*testInt3) Prep()                 {}
func (*testAngle) Prep()                {}
func (*testAngleUnnamed) Prep()         {}
func (*testAngleExtra) Prep()           {}
func (*testReal) Prep()                 {}
func (*testRealUnnamed) Prep()          {}
func (*testRealBounded) Prep()          {}
func (*testRealBoundedUnnamed) Prep()   {}
func (*testRealBoundedBadLower) Prep()  {}
func (*testRealBoundedBadUpper) Prep()  {}
func (*testRealBoundedLowerOnly) Prep() {}
func (*testRealBoundedUpperOnly) Prep() {}
func (*testRealBoundedEmpty) Prep()     {}
func (*testRealBoundedSingleton) Prep() {}
func (*testReal3) Prep()                {}
func (*testComplex) Prep()              {}
func (*testComplexUnnamed) Prep()       {}
func (*testComplexExtra) Prep()         {}
func (*testVec3) Prep()                 {}
func (*testVec3Unnamed) Prep()          {}
func (*testVec3Extra) Prep()            {}
func (*testAffine) Prep()               {}
func (*testAffineUnnamed) Prep()        {}
func (*testAffineExtra) Prep()          {}
func (*testFunc) Prep()                 {}
func (*testFuncUnnamed) Prep()          {}
func (*testFuncOptional) Prep()         {}
func (*testFuncOptionalUnnamed) Prep()  {}
func (*testFuncBad) Prep()              {}
func (*testFuncExtra) Prep()            {}
func (*testFuncList) Prep()             {}
func (*testFuncListUnnamed) Prep()      {}
func (*testFuncListExtra) Prep()        {}
