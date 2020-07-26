package fapi_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/fapi"
)

type pf struct {
	Not     bool           `xirho:"0"`
	Flag    xirho.Flag     `xirho:"1"`
	List    xirho.List     `xirho:"2,madoka,homura,anime"`
	Int     xirho.Int      `xirho:"3"`
	BInt    xirho.Int      `xirho:"4,-1,1"`
	Angle   xirho.Angle    `xirho:"5"`
	Real    xirho.Real     `xirho:"6"`
	BReal   xirho.Real     `xirho:"7,-1,1"`
	Complex xirho.Complex  `xirho:"8"`
	Vec3    xirho.Vec3     `xirho:"9"`
	Affine  xirho.Affine   `xirho:"10"`
	Func    xirho.Func     `xirho:"11"`
	NFunc   xirho.Func     `xirho:"12,optional"`
	Funcs   xirho.FuncList `xirho:"13"`
	AlsoNot xirho.F        `xirho:"14"`
}

func (*pf) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	return in
}

func (*pf) Prep() {}

func newPf() xirho.F {
	r := pf{}
	r.Func.F = &r
	r.AlsoNot = &r
	return &r
}

type ef struct{}

func (ef) Calc(in xirho.P, rng *xirho.RNG) xirho.P {
	return in
}

func (ef) Prep() {}

func TestForCount(t *testing.T) {
	cases := map[string]struct {
		v xirho.F
		f int
		n int
	}{
		"ef": {v: ef{}, f: 0, n: 0},
		"pf": {v: newPf(), f: 15, n: 13},
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

func TestSetFlag(t *testing.T) {
	for _, c := range typeCases {
		if len(c.set) != 0 && c.param == reflect.TypeOf(fapi.Flag{}) {
			t.Run(c.name, func(t *testing.T) {
				api := fapi.For(c.v)
				if len(api) != 1 {
					t.Fatalf("wrong number of fields on %#v: expected 1, have %d", c.v, len(api))
				}
				p := api[0].(fapi.Flag)
				for i, s := range c.set {
					err := p.Set(s.set.(bool))
					if (err != nil && s.err != nil && !errors.As(err, &s.err)) || (err == nil && s.err != nil) || (err != nil && s.err == nil) {
						t.Errorf("wrong error for set case %d: expected %T, got %T", i, s.err, err)
					}
					if !cmp.Equal(s.get, p.Get()) {
						t.Errorf("wrong get after set %v (with expected error %T): expected %v, got %v", s.set, s.err, s.get, p.Get())
					}
				}
			})
		}
	}
}

func TestSetList(t *testing.T) {
	for _, c := range typeCases {
		if len(c.set) != 0 && c.param == reflect.TypeOf(fapi.List{}) {
			t.Run(c.name, func(t *testing.T) {
				api := fapi.For(c.v)
				if len(api) != 1 {
					t.Fatalf("wrong number of fields on %#v: expected 1, have %d", c.v, len(api))
				}
				p := api[0].(fapi.List)
				for i, s := range c.set {
					err := p.Set(s.set.(int))
					if (err != nil && s.err != nil && !errors.As(err, &s.err)) || (err == nil && s.err != nil) || (err != nil && s.err == nil) {
						t.Errorf("wrong error for set case %d: expected %T, got %T", i, s.err, err)
					}
					if !cmp.Equal(s.get, p.Get()) {
						t.Errorf("wrong get after set %v (with expected error %T): expected %v, got %v", s.set, s.err, s.get, p.Get())
					}
				}
			})
		}
	}
}

func TestSetInt(t *testing.T) {
	for _, c := range typeCases {
		if len(c.set) != 0 && c.param == reflect.TypeOf(fapi.Int{}) {
			t.Run(c.name, func(t *testing.T) {
				api := fapi.For(c.v)
				if len(api) != 1 {
					t.Fatalf("wrong number of fields on %#v: expected 1, have %d", c.v, len(api))
				}
				p := api[0].(fapi.Int)
				for i, s := range c.set {
					err := p.Set(s.set.(int64))
					if (err != nil && s.err != nil && !errors.As(err, &s.err)) || (err == nil && s.err != nil) || (err != nil && s.err == nil) {
						t.Errorf("wrong error for set case %d: expected %T, got %T", i, s.err, err)
					}
					if !cmp.Equal(s.get, p.Get()) {
						t.Errorf("wrong get after set %v (with expected error %T): expected %v, got %v", s.set, s.err, s.get, p.Get())
					}
				}
			})
		}
	}
}

func TestSetAngle(t *testing.T) {
	for _, c := range typeCases {
		if len(c.set) != 0 && c.param == reflect.TypeOf(fapi.Angle{}) {
			t.Run(c.name, func(t *testing.T) {
				api := fapi.For(c.v)
				if len(api) != 1 {
					t.Fatalf("wrong number of fields on %#v: expected 1, have %d", c.v, len(api))
				}
				p := api[0].(fapi.Angle)
				for i, s := range c.set {
					err := p.Set(s.set.(float64))
					if (err != nil && s.err != nil && !errors.As(err, &s.err)) || (err == nil && s.err != nil) || (err != nil && s.err == nil) {
						t.Errorf("wrong error for set case %d: expected %T, got %T", i, s.err, err)
					}
					if !cmp.Equal(s.get, p.Get(), cmpopts.EquateNaNs()) {
						t.Errorf("wrong get after set %v (with expected error %T): expected %v, got %v", s.set, s.err, s.get, p.Get())
					}
				}
			})
		}
	}
}

func TestSetReal(t *testing.T) {
	for _, c := range typeCases {
		if len(c.set) != 0 && c.param == reflect.TypeOf(fapi.Real{}) {
			t.Run(c.name, func(t *testing.T) {
				api := fapi.For(c.v)
				if len(api) != 1 {
					t.Fatalf("wrong number of fields on %#v: expected 1, have %d", c.v, len(api))
				}
				p := api[0].(fapi.Real)
				for i, s := range c.set {
					err := p.Set(s.set.(float64))
					if (err != nil && s.err != nil && !errors.As(err, &s.err)) || (err == nil && s.err != nil) || (err != nil && s.err == nil) {
						t.Errorf("wrong error for set case %d: expected %T, got %T", i, s.err, err)
					}
					if !cmp.Equal(s.get, p.Get(), cmpopts.EquateNaNs()) {
						t.Errorf("wrong get after set %v (with expected error %T): expected %v, got %v", s.set, s.err, s.get, p.Get())
					}
				}
			})
		}
	}
}

func TestSetComplex(t *testing.T) {
	for _, c := range typeCases {
		if len(c.set) != 0 && c.param == reflect.TypeOf(fapi.Complex{}) {
			t.Run(c.name, func(t *testing.T) {
				api := fapi.For(c.v)
				if len(api) != 1 {
					t.Fatalf("wrong number of fields on %#v: expected 1, have %d", c.v, len(api))
				}
				p := api[0].(fapi.Complex)
				for i, s := range c.set {
					err := p.Set(s.set.(complex128))
					if (err != nil && s.err != nil && !errors.As(err, &s.err)) || (err == nil && s.err != nil) || (err != nil && s.err == nil) {
						t.Errorf("wrong error for set case %d: expected %T, got %T", i, s.err, err)
					}
					if !cmp.Equal(s.get, p.Get(), cmpopts.EquateNaNs()) {
						t.Errorf("wrong get after set %v (with expected error %T): expected %v, got %v", s.set, s.err, s.get, p.Get())
					}
				}
			})
		}
	}
}

func TestSetVec3(t *testing.T) {
	for _, c := range typeCases {
		if len(c.set) != 0 && c.param == reflect.TypeOf(fapi.Vec3{}) {
			t.Run(c.name, func(t *testing.T) {
				api := fapi.For(c.v)
				if len(api) != 1 {
					t.Fatalf("wrong number of fields on %#v: expected 1, have %d", c.v, len(api))
				}
				p := api[0].(fapi.Vec3)
				for i, s := range c.set {
					err := p.Set(s.set.(xirho.Vec3))
					if (err != nil && s.err != nil && !errors.As(err, &s.err)) || (err == nil && s.err != nil) || (err != nil && s.err == nil) {
						t.Errorf("wrong error for set case %d: expected %T, got %T", i, s.err, err)
					}
					if diff := cmp.Diff(s.get, p.Get(), cmpopts.EquateNaNs()); diff != "" {
						t.Errorf("wrong get after set %v (with expected error %T): diff (-expected +got):\n%s", s.set, s.err, diff)
					}
				}
			})
		}
	}
}

func TestSetAffine(t *testing.T) {
	for _, c := range typeCases {
		if len(c.set) != 0 && c.param == reflect.TypeOf(fapi.Affine{}) {
			t.Run(c.name, func(t *testing.T) {
				api := fapi.For(c.v)
				if len(api) != 1 {
					t.Fatalf("wrong number of fields on %#v: expected 1, have %d", c.v, len(api))
				}
				p := api[0].(fapi.Affine)
				for i, s := range c.set {
					err := p.Set(s.set.(xirho.Affine))
					if (err != nil && s.err != nil && !errors.As(err, &s.err)) || (err == nil && s.err != nil) || (err != nil && s.err == nil) {
						t.Errorf("wrong error for set case %d: expected %T, got %T", i, s.err, err)
					}
					if diff := cmp.Diff(s.get, p.Get(), cmpopts.EquateNaNs()); diff != "" {
						t.Errorf("wrong get after set %v (with expected error %T): diff (-expected +got):\n%s", s.set, s.err, diff)
					}
				}
			})
		}
	}
}

func TestSetFunc(t *testing.T) {
	for _, c := range typeCases {
		if len(c.set) != 0 && c.param == reflect.TypeOf(fapi.Func{}) {
			t.Run(c.name, func(t *testing.T) {
				api := fapi.For(c.v)
				if len(api) != 1 {
					t.Fatalf("wrong number of fields on %#v: expected 1, have %d", c.v, len(api))
				}
				p := api[0].(fapi.Func)
				for i, s := range c.set {
					var err error
					if s.set != nil {
						// special handling because nil does not assert
						err = p.Set(s.set.(xirho.F))
					} else {
						err = p.Set(nil)
					}
					if (err != nil && s.err != nil && !errors.As(err, &s.err)) || (err == nil && s.err != nil) || (err != nil && s.err == nil) {
						t.Errorf("wrong error for set case %d: expected %T, got %T", i, s.err, err)
					}
					if diff := cmp.Diff(s.get, p.Get()); diff != "" {
						t.Errorf("wrong get after set %v (with expected error %T): diff (-expected +got):\n%s", s.set, s.err, diff)
					}
				}
			})
		}
	}
}

func TestSetFuncList(t *testing.T) {
	for _, c := range typeCases {
		if len(c.set) != 0 && c.param == reflect.TypeOf(fapi.FuncList{}) {
			t.Run(c.name, func(t *testing.T) {
				api := fapi.For(c.v)
				if len(api) != 1 {
					t.Fatalf("wrong number of fields on %#v: expected 1, have %d", c.v, len(api))
				}
				p := api[0].(fapi.FuncList)
				for i, s := range c.set {
					err := p.Set(s.set.(xirho.FuncList))
					if (err != nil && s.err != nil && !errors.As(err, &s.err)) || (err == nil && s.err != nil) || (err != nil && s.err == nil) {
						t.Errorf("wrong error for set case %d: expected %T, got %T", i, s.err, err)
					}
					if diff := cmp.Diff(s.get, p.Get()); diff != "" {
						t.Errorf("wrong get after set %v (with expected error %T): diff (-expected +got):\n%s", s.set, s.err, diff)
					}
				}
			})
		}
	}
}
