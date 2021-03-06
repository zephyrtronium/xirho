package xi

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/zephyrtronium/xirho"
)

// name2f maps names of registered functions to default factories.
var name2f = make(map[string]func() xirho.Func)

// f2name maps registered function types to names.
var f2name = make(map[reflect.Type]string)

// Register registers a function type by a factory which produces it. An error
// is returned if there is already a function with the given name. If the same
// underlying type is registered multiple times, the name passed in the first
// associated call to Register is used by NameOf.
func Register(name string, factory func() xirho.Func) error {
	if f, ok := name2f[name]; ok {
		t1 := reflect.TypeOf(factory())
		for t1.Kind() == reflect.Ptr {
			t1 = t1.Elem()
		}
		t2 := reflect.TypeOf(f())
		for t2.Kind() == reflect.Ptr {
			t2 = t2.Elem()
		}
		return fmt.Errorf("xirho: attempted to register %q for %s/%s but already registered for %s/%s", name, t1.PkgPath(), t1.Name(), t2.PkgPath(), t2.Name())
	}
	name2f[name] = factory
	t := reflect.TypeOf(factory())
	if _, ok := f2name[t]; !ok {
		f2name[t] = name
	}
	return nil
}

// New creates a new function using a registered name. The result is nil if
// there is no function registered with the given name.
func New(name string) xirho.Func {
	f := name2f[name]
	if f == nil {
		return nil
	}
	return f()
}

// NameOf returns the name of a registered function. ok is false if there is no
// such function.
func NameOf(f xirho.Func) (name string, ok bool) {
	name, ok = f2name[reflect.TypeOf(f)]
	return
}

// Names returns a list of all registered function names. If unique is true,
// then only the first name registered for each underlying type is used. The
// names are returned in lexicographically increasing order.
func Names(unique bool) []string {
	var r []string
	if unique {
		r = make([]string, 0, len(f2name))
		for _, name := range f2name {
			r = append(r, name)
		}
	} else {
		r = make([]string, 0, len(name2f))
		for name := range name2f {
			r = append(r, name)
		}
	}
	sort.Strings(r)
	return r
}

// must registers a function and panics if Register returns a non-nil error.
func must(name string, factory func() xirho.Func) {
	if err := Register(name, factory); err != nil {
		panic(err)
	}
}
