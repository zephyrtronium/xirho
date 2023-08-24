package fapi_test

import (
	"errors"
	"math"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/fapi"
	"github.com/zephyrtronium/xirho/xi"
)

func TestSetFlag(t *testing.T) {
	for _, c := range typeCases {
		if c.param == reflect.TypeOf(fapi.Flag{}) {
			if len(c.set) == 0 {
				t.Log("no set cases in", c)
				continue
			}
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
		if c.param == reflect.TypeOf(fapi.List{}) {
			if len(c.set) == 0 {
				t.Log("no set cases in", c)
				continue
			}
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
		if c.param == reflect.TypeOf(fapi.Int{}) {
			if len(c.set) == 0 {
				t.Log("no set cases in", c)
				continue
			}
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
		if c.param == reflect.TypeOf(fapi.Angle{}) {
			if len(c.set) == 0 {
				t.Log("no set cases in", c)
				continue
			}
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
		if c.param == reflect.TypeOf(fapi.Real{}) {
			if len(c.set) == 0 {
				t.Log("no set cases in", c)
				continue
			}
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
		if c.param == reflect.TypeOf(fapi.Complex{}) {
			if len(c.set) == 0 {
				t.Log("no set cases in", c)
				continue
			}
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
		if c.param == reflect.TypeOf(fapi.Vec3{}) {
			if len(c.set) == 0 {
				t.Log("no set cases in", c)
				continue
			}
			t.Run(c.name, func(t *testing.T) {
				api := fapi.For(c.v)
				if len(api) != 1 {
					t.Fatalf("wrong number of fields on %#v: expected 1, have %d", c.v, len(api))
				}
				p := api[0].(fapi.Vec3)
				for i, s := range c.set {
					err := p.Set(s.set.([3]float64))
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
		if c.param == reflect.TypeOf(fapi.Affine{}) {
			if len(c.set) == 0 {
				t.Log("no set cases in", c)
				continue
			}
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
		if c.param == reflect.TypeOf(fapi.Func{}) {
			if len(c.set) == 0 {
				t.Log("no set cases in", c)
				continue
			}
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
						err = p.Set(s.set.(xirho.Func))
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
		if c.param == reflect.TypeOf(fapi.FuncList{}) {
			if len(c.set) == 0 {
				t.Log("no set cases in", c)
				continue
			}
			t.Run(c.name, func(t *testing.T) {
				api := fapi.For(c.v)
				if len(api) != 1 {
					t.Fatalf("wrong number of fields on %#v: expected 1, have %d", c.v, len(api))
				}
				p := api[0].(fapi.FuncList)
				for i, s := range c.set {
					err := p.Set(s.set.([]xirho.Func))
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

func TestListOpts(t *testing.T) {
	cases := map[string]struct {
		v    xirho.Func
		opts []string
	}{
		"two":     {v: new(testList2), opts: []string{"1", "2"}},
		"ten":     {v: new(testList10), opts: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}},
		"unnamed": {v: new(testListUnnamed), opts: []string{"1", "2"}},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			api := fapi.For(c.v)
			if len(api) != 1 {
				t.Fatalf("wrong number of fields on %#v: expected 1, have %d", c.v, len(api))
			}
			p := api[0].(fapi.List)
			opts := p.Opts()
			if diff := cmp.Diff(c.opts, opts); diff != "" {
				t.Fatalf("wrong list options for %#v: diff (-expected +got):\n%s", c.v, diff)
			}
			for i, o := range opts {
				p.Set(i)
				if p.String() != o {
					t.Errorf("wrong option %d: expected %q, have %q", i, o, p.String())
				}
			}
		})
	}
}

func TestIntBounds(t *testing.T) {
	cases := map[string]struct {
		v      xirho.Func
		bdd    bool
		lo, hi int64
	}{
		"unbounded":         {v: new(testInt), bdd: false, lo: math.MinInt64, hi: math.MaxInt64},
		"unbounded_unnamed": {v: new(testIntUnnamed), bdd: false, lo: math.MinInt64, hi: math.MaxInt64},
		"bounded":           {v: new(testIntBounded), bdd: true, lo: -1, hi: 1},
		"bounded_unnamed":   {v: new(testIntBoundedUnnamed), bdd: true, lo: -1, hi: 1},
		"singleton":         {v: new(testIntBoundedSingleton), bdd: true, lo: 0, hi: 0},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			api := fapi.For(c.v)
			if len(api) != 1 {
				t.Fatalf("wrong number of fields on %#v: expected 1, have %d", c.v, len(api))
			}
			p := api[0].(fapi.Int)
			if p.Bounded() != c.bdd {
				t.Errorf("%#v has bounded=%v but expected bounded=%v", c.v, p.Bounded(), c.bdd)
			}
			if lo, hi := p.Bounds(); lo != c.lo || hi != c.hi {
				t.Errorf("%#v has wrong bounds: expected %d..%d but got %d..%d", c.v, c.lo, c.hi, lo, hi)
			}
		})
	}
}

func TestRealBounds(t *testing.T) {
	cases := map[string]struct {
		v      xirho.Func
		bdd    bool
		lo, hi float64
	}{
		"unbounded":         {v: new(testReal), bdd: false, lo: math.Inf(-1), hi: math.Inf(0)},
		"unbounded_unnamed": {v: new(testRealUnnamed), bdd: false, lo: math.Inf(-1), hi: math.Inf(0)},
		"bounded":           {v: new(testRealBounded), bdd: true, lo: -1, hi: 1},
		"bounded_unnamed":   {v: new(testRealBoundedUnnamed), bdd: true, lo: -1, hi: 1},
		"singleton":         {v: new(testRealBoundedSingleton), bdd: true, lo: 0, hi: 0},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			api := fapi.For(c.v)
			if len(api) != 1 {
				t.Fatalf("wrong number of fields on %#v: expected 1, have %d", c.v, len(api))
			}
			p := api[0].(fapi.Real)
			if p.Bounded() != c.bdd {
				t.Errorf("%#v has bounded=%v but expected bounded=%v", c.v, p.Bounded(), c.bdd)
			}
			if lo, hi := p.Bounds(); lo != c.lo || hi != c.hi {
				t.Errorf("%#v has wrong bounds: expected %g..%g but got %g..%g", c.v, c.lo, c.hi, lo, hi)
			}
		})
	}
}

func TestFuncOptional(t *testing.T) {
	cases := map[string]struct {
		v xirho.Func
		o bool
	}{
		"required":         {v: new(testFunc), o: false},
		"unnamed":          {v: new(testFuncUnnamed), o: false},
		"optional":         {v: new(testFuncOptional), o: true},
		"optional_unnamed": {v: new(testFuncOptionalUnnamed), o: true},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			api := fapi.For(c.v)
			if len(api) != 1 {
				t.Fatalf("wrong number of fields on %#v: expected 1, have %d", c.v, len(api))
			}
			p := api[0].(fapi.Func)
			if p.IsOptional() != c.o {
				t.Errorf("%#v has optional=%v but expected optional=%v", c.v, p.IsOptional(), c.o)
			}
		})
	}
}

func TestFuncListAppend(t *testing.T) {
	cases := map[string]xirho.Func{
		"named":   new(testFuncList),
		"unnamed": new(testFuncListUnnamed),
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			api := fapi.For(c)
			if len(api) != 1 {
				t.Fatalf("wrong number of fields on %#v: expected 1, have %d", c, len(api))
			}
			p := api[0].(fapi.FuncList)
			s := p.Get()
			p.Append(xi.Spherical{})
			u := p.Get()
			v := append(s, xi.Spherical{})
			if diff := cmp.Diff(v, u); diff != "" {
				t.Errorf("wrong result after appending one item to %#v: diff (-expected +got):\n%s", c, diff)
			}
			s = p.Get()
			p.Append(xi.Spherical{}, xi.Spherical{})
			u = p.Get()
			v = append(s, xi.Spherical{}, xi.Spherical{})
			if diff := cmp.Diff(v, u); diff != "" {
				t.Errorf("wrong result after appending two items to %#v: diff (-expected +got):\n%s", c, diff)
			}
		})
	}
}
