package encoding

import (
	"encoding/json"
	"fmt"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/fapi"
	"github.com/zephyrtronium/xirho/xi"
)

// funcm encodes information about a function and its parameters in an
// intermediary structure to facilitate encoding and decoding.
type funcm struct {
	Name   string                 `json:"name"`
	Params map[string]interface{} `json:"params,omitempty"`
}

// newFuncm creates an encoding wrapper around f. Returns an error if the
// underlying type of f has not been registered with package xi.
func newFuncm(f xirho.F) (*funcm, error) {
	name, ok := xi.NameOf(f)
	if !ok {
		return nil, fmt.Errorf("unregistered function %#v", f)
	}
	api := fapi.For(f)
	r := funcm{Name: name}
	if len(api) == 0 {
		return &r, nil
	}
	r.Params = make(map[string]interface{})
	for _, parm := range api {
		switch p := parm.(type) {
		case fapi.Flag:
			r.Params[p.Name()] = p.Get()
		case fapi.List:
			r.Params[p.Name()] = p.Get()
		case fapi.Int:
			r.Params[p.Name()] = p.Get()
		case fapi.Angle:
			r.Params[p.Name()] = p.Get()
		case fapi.Real:
			r.Params[p.Name()] = p.Get()
		case fapi.Complex:
			v := p.Get()
			r.Params[p.Name()] = [2]float64{real(v), imag(v)}
		case fapi.Vec3:
			r.Params[p.Name()] = p.Get()
		case fapi.Affine:
			r.Params[p.Name()] = p.Get()
		case fapi.Func:
			if p.Get() == nil {
				r.Params[p.Name()] = nil
				continue
			}
			v, err := newFuncm(p.Get())
			if err != nil {
				return nil, err
			}
			r.Params[p.Name()] = v
		case fapi.FuncList:
			pv := p.Get()
			v := make([]*funcm, len(pv))
			for i, x := range pv {
				nf, err := newFuncm(x)
				if err != nil {
					return nil, err
				}
				v[i] = nf
			}
			r.Params[p.Name()] = v
		default:
			panic(fmt.Errorf("xirho: unhandled fapi.Param %#v", p))
		}
	}
	return &r, nil
}

// unf decodes a funcm into the corresponding xirho.F, setting the appropriate
// parameters. The returned error is non-nil if there is no registered function
// with the funcm's name, or if the corresponding function does not have a
// parameter with the name of a funcm parameter, or if any function parameter
// cannot be set to the value in the parameters (e.g. due to bounds).
func unf(f *funcm) (v xirho.F, err error) {
	v = xi.New(f.Name)
	if v == nil {
		return nil, fmt.Errorf("no registered function named %s", f.Name)
	}
	api := fapi.For(v)
	for _, parm := range api {
		x, ok := f.Params[parm.Name()]
		if !ok {
			continue
		}
		switch p := parm.(type) {
		case fapi.Flag:
			t, ok := x.(bool)
			if !ok {
				return nil, fmt.Errorf("expected bool for %s but got %#v", parm.Name(), x)
			}
			if err := p.Set(xirho.Flag(t)); err != nil {
				return nil, err
			}
		case fapi.List:
			t, err := getint(p.Name(), x)
			if err != nil {
				return nil, err
			}
			if err := p.Set(xirho.List(t)); err != nil {
				return nil, err
			}
		case fapi.Int:
			t, err := getint(p.Name(), x)
			if err != nil {
				return nil, err
			}
			if err := p.Set(xirho.Int(t)); err != nil {
				return nil, err
			}
		case fapi.Angle:
			t, err := getfloat(p.Name(), x)
			if err != nil {
				return nil, err
			}
			if err := p.Set(xirho.Angle(t)); err != nil {
				return nil, err
			}
		case fapi.Real:
			t, err := getfloat(p.Name(), x)
			if err != nil {
				return nil, err
			}
			if err := p.Set(xirho.Real(t)); err != nil {
				return nil, err
			}
		case fapi.Complex:
			t, err := getfloatlist(p.Name(), x)
			if err != nil {
				return nil, err
			}
			if len(t) != 2 {
				return nil, fmt.Errorf("expected complex for %s but got %#v", p.Name(), x)
			}
			if err := p.Set(xirho.Complex(complex(t[0], t[1]))); err != nil {
				return nil, err
			}
		case fapi.Vec3:
			t, err := getfloatlist(p.Name(), x)
			if err != nil {
				return nil, err
			}
			if len(t) != 3 {
				return nil, fmt.Errorf("expected vec3 for %s but got %#v", p.Name(), x)
			}
			if err := p.Set(xirho.Vec3{t[0], t[1], t[2]}); err != nil {
				return nil, err
			}
		case fapi.Affine:
			t, err := getfloatlist(p.Name(), x)
			if err != nil {
				return nil, err
			}
			var b xirho.Affine
			if copy(b[:], t) != len(b) {
				return nil, fmt.Errorf("expected affine for %s but got %#v", p.Name(), x)
			}
			if err := p.Set(b); err != nil {
				return nil, err
			}
		case fapi.Func:
			if x == nil {
				// Optional funcs are allowed to be nil. The setter will tell
				// us if it isn't optional.
				if err := p.Set(nil); err != nil {
					return nil, err
				}
				break
			}
			t, err := getfunc(p.Name(), x)
			if err != nil {
				return nil, err
			}
			nf, err := unf(t)
			if err != nil {
				return nil, err
			}
			if err := p.Set(nf); err != nil {
				return nil, err
			}
		case fapi.FuncList:
			fl, ok := x.([]interface{})
			if !ok {
				return nil, fmt.Errorf("expected func list for %s but got %#v", p.Name(), x)
			}
			r := make(xirho.FuncList, len(fl))
			for i, fa := range fl {
				t, err := getfunc(fmt.Sprintf("%s[%d] (of %d)", p.Name(), i, len(fl)), fa)
				if err != nil {
					return nil, err
				}
				r[i], err = unf(t)
				if err != nil {
					return nil, err
				}
			}
			if err := p.Set(r); err != nil {
				return nil, err
			}
		default:
			panic(fmt.Errorf("xirho: unhandled fapi.Param %#v", p))
		}
		delete(f.Params, parm.Name())
	}
	if len(f.Params) != 0 {
		nn, _ := xi.NameOf(v) // err must be nil since xi.New succeeded
		err = fmt.Errorf("unknown params for %s: %v", nn, f.Params)
	}
	return
}

// getint gets an int64 from a decoded JSON numeric value.
func getint(name string, x interface{}) (int64, error) {
	switch t := x.(type) {
	case json.Number:
		r, err := t.Int64()
		if err != nil {
			return 0, err
		}
		return r, nil
	case float64:
		return int64(t), nil
	default:
		return 0, fmt.Errorf("expected int for %s but got %#v", name, x)
	}
}

// getfloat gets a float64 from a decoded JSON numeric value.
func getfloat(name string, x interface{}) (float64, error) {
	switch t := x.(type) {
	case json.Number:
		r, err := t.Float64()
		if err != nil {
			return 0, err
		}
		return r, nil
	case float64:
		return t, nil
	default:
		return 0, fmt.Errorf("expected float for %s but got %#v", name, x)
	}
}

// getfloatlist gets a []float64 from a decoded JSON numeric value.
func getfloatlist(name string, x interface{}) (r []float64, err error) {
	switch t := x.(type) {
	case []interface{}:
		r = make([]float64, len(t))
		for i, v := range t {
			x, err := getfloat(fmt.Sprintf("%s[%d]", name, i), v)
			if err != nil {
				return nil, err
			}
			r[i] = x
		}
		return r, nil
	case []json.Number:
		r = make([]float64, len(t))
		for i, v := range t {
			if r[i], err = v.Float64(); err != nil {
				return nil, err
			}
		}
		return r, nil
	case []float64:
		return t, nil
	default:
		return nil, fmt.Errorf("expected float list for %s but got %#v", name, x)
	}
}

// getfunc gets a funcm from a decoded JSON object.
func getfunc(name string, x interface{}) (r *funcm, err error) {
	v, _ := x.(map[string]interface{})
	// If x isn't a JSON object, then v is nil, so v["name"] is nil, so
	// v["name"].(string) gives "", false. It really does work like that.
	n, ok := v["name"].(string)
	if !ok {
		return nil, fmt.Errorf("expected func for %s but got %#v", name, x)
	}
	switch len(v) {
	case 1: // no params
		return &funcm{Name: n}, nil
	case 2: // yes params
		p, ok := v["params"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expected func for %s but got %#v", name, x)
		}
		return &funcm{Name: n, Params: p}, nil
	default:
		return nil, fmt.Errorf("expected func for %s but got %#v", name, x)
	}
}
