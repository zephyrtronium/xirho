package xi_test

import (
	"reflect"
	"testing"

	"github.com/zephyrtronium/xirho/fapi"
	"github.com/zephyrtronium/xirho/xi"
)

// ExpectAPI tests that for each name, the registered function has the API
// given in expect, such that the names match exactly and the parameter type
// matches the element type.
func ExpectAPI(t *testing.T, expect map[string]fapi.Param, names ...string) {
	t.Helper()
	for _, name := range names {
		name := name
		t.Run(name, func(t *testing.T) {
			expectOne(t, name, expect)
		})
	}
}

// expectOne asserts that the fapi parameters available on the xi function
// type named typ match the given parameters. Each element of expect contains
// an instance of the parameter type the corresponding parameter should be.
func expectOne(t *testing.T, typ string, expect map[string]fapi.Param) {
	t.Helper()
	fn := xi.New(typ)
	if fn == nil {
		t.Fatalf("no registered function named %q", typ)
	}
	api := fapi.For(fn)

	seen := make(map[string]bool, len(expect))
	for _, a := range api {
		p, ok := expect[a.Name()]
		if !ok {
			t.Errorf("unexpected param %s of type %T", a.Name(), a)
			continue
		}
		if reflect.TypeOf(a) != reflect.TypeOf(p) {
			t.Errorf("param %s has wrong type: want %T, got %T", a.Name(), p, a)
		}
		seen[a.Name()] = true
	}

	for name, p := range expect {
		if !seen[name] {
			t.Errorf("missing param %s of type %T", name, p)
		}
	}
}
