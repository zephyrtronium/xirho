package fapi_test

import (
	"reflect"
	"testing"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/fapi"
	"github.com/zephyrtronium/xirho/xmath"
)

type pf struct {
	Not     int32        `xirho:"0"`
	Flag    bool         `xirho:"1"`
	List    int          `xirho:"2,madoka,homura,anime"`
	Int     int64        `xirho:"3"`
	BInt    int64        `xirho:"4,-1,1"`
	Angle   float64      `xirho:"5,angle"`
	Real    float64      `xirho:"6"`
	BReal   float64      `xirho:"7,-1,1"`
	Complex complex128   `xirho:"8"`
	Vec3    [3]float64   `xirho:"9"`
	Affine  xmath.Affine `xirho:"10"`
	Func    xirho.Func   `xirho:"11"`
	NFunc   xirho.Func   `xirho:"12,optional"`
	Funcs   []xirho.Func `xirho:"13"`
}

func (*pf) Calc(in xirho.Pt, rng *xmath.RNG) xirho.Pt {
	return in
}

func (*pf) Prep() {}

func newPf() xirho.Func {
	r := pf{}
	r.Func = &r
	return &r
}

type ef struct{}

func (ef) Calc(in xirho.Pt, rng *xmath.RNG) xirho.Pt {
	return in
}

func (ef) Prep() {}

type uf struct {
	//lint:ignore U1000 field is used to test that we skip unexported fields in reflection
	unexported bool `xirho:"unexported"`
}

func (*uf) Calc(in xirho.Pt, rng *xmath.RNG) xirho.Pt {
	return in
}

func (*uf) Prep() {}

type ff bool

func (*ff) Calc(in xirho.Pt, rng *xmath.RNG) xirho.Pt {
	return in
}

func (*ff) Prep() {}

func TestForCount(t *testing.T) {
	cases := map[string]struct {
		v xirho.Func
		f int
		n int
	}{
		"ef": {v: ef{}, f: 0, n: 0},
		"pf": {v: newPf(), f: 14, n: 13},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			r := reflect.TypeOf(c.v)
			for r.Kind() == reflect.Ptr {
				r = r.Elem()
			}
			if r.NumField() != c.f {
				t.Errorf("wrong number of struct fields: expected %d, have %d", c.f, r.NumField())
			}
			api := fapi.For(c.v)
			if len(api) != c.n {
				t.Errorf("wrong number of api fields: expected %d, have %d", c.n, len(api))
			}
			for i, p := range api {
				switch p.(type) {
				case fapi.Flag, fapi.List, fapi.Int, fapi.Angle, fapi.Real,
					fapi.Complex, fapi.Vec3, fapi.Affine, fapi.Func, fapi.FuncList: // do nothing
				default:
					t.Errorf("unknown parameter type %T for parameter %d named %q", p, i, p.Name())
				}
			}
		})
	}
}

func TestForErrors(t *testing.T) {
	for _, c := range typeCases {
		t.Run(c.name, func(t *testing.T) {
			if c.param != nil {
				_ = fapi.For(c.v)
			} else {
				defer func() {
					if recover() == nil {
						t.Error("expected error, got nil")
					}
				}()
				_ = fapi.For(c.v)
			}
		})
	}
}

func TestForName(t *testing.T) {
	for _, c := range typeCases {
		if c.param != nil {
			t.Run(c.name, func(t *testing.T) {
				api := fapi.For(c.v)
				if len(api) != 1 {
					t.Fatalf("wrong number of fields on %#v: expected 1, have %d", c.v, len(api))
				}
				if api[0].Name() != c.field {
					t.Errorf("wrong field name on %#v: expected %q, have %q", c.v, c.field, api[0].Name())
				}
			})
		}
	}
}

func TestForTypes(t *testing.T) {
	for _, c := range typeCases {
		if c.param != nil {
			t.Run(c.name, func(t *testing.T) {
				api := fapi.For(c.v)
				if len(api) != 1 {
					t.Fatalf("wrong number of fields on %#v: expected 1, have %d", c.v, len(api))
				}
				if reflect.TypeOf(api[0]) != c.param {
					t.Errorf("wrong param type on %#v: expected %v, got %v", c.v, c.param, reflect.TypeOf(api[0]))
				}
			})
		}
	}
}

func TestForUnexported(t *testing.T) {
	v := new(uf)
	api := fapi.For(v)
	if len(api) != 0 {
		t.Errorf("got non-empty parameter list for %#v: %#v", v, api)
	}
}

func TestForNonStruct(t *testing.T) {
	v := new(ff)
	api := fapi.For(v)
	if len(api) != 0 {
		t.Errorf("got non-empty parameter list for %#v: %#v", v, api)
	}
}

// putting this here rather than adding a whole new file

func TestErrorsProduceANonemptyErrorMessage(t *testing.T) {
	v := new(pf)
	p := fapi.For(v)[0]
	cases := map[string]error{
		"OutOfBoundsInt":  fapi.OutOfBoundsInt{Param: p, Value: 1},
		"OutOfBoundsReal": fapi.OutOfBoundsReal{Param: p, Value: 1},
		"NotFinite":       fapi.NotFinite{Param: p},
		"NotOptional":     fapi.NotOptional{Param: p},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if m := c.Error(); m == "" {
				t.Errorf("error of type %T produced empty error message", c)
			}
		})
	}
}
