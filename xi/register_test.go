package xi_test

import (
	"sort"
	"testing"

	"github.com/zephyrtronium/xirho"
	"github.com/zephyrtronium/xirho/xi"
)

type T struct{}
type U struct{}
type N struct{}

func (T) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt { return in }
func (U) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt { return in }
func (N) Calc(in xirho.Pt, rng *xirho.RNG) xirho.Pt { return in }

func (T) Prep() {}
func (U) Prep() {}
func (N) Prep() {}

func newt() xirho.Func { return T{} }

const name1 = "T-for-testing"
const name2 = name1 + "-init-alt"

func init() {
	// Register in an init to ensure this type is always available regardless
	// of the order in which tests run.
	if err := xi.Register(name1, newt); err != nil {
		panic(err)
	}
	if err := xi.Register(name2, newt); err != nil {
		panic(err)
	}
}

func TestRegister(t *testing.T) {
	if err := xi.Register(name1+"-test-alt", newt); err != nil {
		t.Error("error registering alternate name:", err)
	}
	if err := xi.Register(name1, newt); err == nil {
		t.Error("expected error registering used name but got nil")
	}
}

func TestNew(t *testing.T) {
	if f := xi.New(name1); f == nil {
		t.Errorf("expected function for name %q but got nil", name1)
	}
	if f := xi.New(name2); f == nil {
		t.Errorf("expected function for name %q but got nil", name2)
	}
	if f := xi.New(name1 + "-does-not-exist"); f != nil {
		t.Errorf("expected no function but got %#v", f)
	}
}

func TestNameOf(t *testing.T) {
	name, ok := xi.NameOf(T{})
	if !ok {
		t.Error("no name for registered function")
	}
	if name != name1 {
		t.Errorf("wrong name %q, expected %q", name, name1)
	}
	name, ok = xi.NameOf(U{})
	if ok {
		t.Errorf("got name %q for unregistered function", name)
	}
}

func TestNames(t *testing.T) {
	names := xi.Names(false)
	if !sort.StringsAreSorted(names) {
		t.Errorf("names not sorted: %q", names)
	}
	if k := findname(names, name1); k < 0 {
		t.Errorf("no name %q in %q", name1, names)
	}
	if k := findname(names, name2); k < 0 {
		t.Errorf("no name %q in %q", name2, names)
	}
	const nname = "T-for-testing-names"
	if k := findname(names, nname); k >= 0 {
		t.Errorf("testing name %q already at position %d in %q", nname, k, names)
	}
	if err := xi.Register(nname, newt); err != nil {
		t.Errorf("couldn't register %q: %v", nname, err)
	}
	names = xi.Names(false)
	if k := findname(names, nname); k < 0 {
		t.Errorf("didn't get new name %q in %q", nname, names)
	}
}

func TestUniqueNames(t *testing.T) {
	names := xi.Names(true)
	if !sort.StringsAreSorted(names) {
		t.Error("names not sorted:", names)
	}
	if k := findname(names, name1); k < 0 {
		t.Errorf("no name %q in %q", name1, names)
	}
	if k := findname(names, name2); k >= 0 {
		t.Errorf("unexpected name %q in %q", name2, names)
	}
	const nname = "T-for-testing-unique-names"
	if k := findname(names, nname); k >= 0 {
		t.Errorf("testing name %q already at position %d in %q", nname, k, names)
	}
	if err := xi.Register(nname, func() xirho.Func { return N{} }); err != nil {
		t.Errorf("couldn't register %q: %v", nname, err)
	}
	names = xi.Names(true)
	if k := findname(names, nname); k < 0 {
		t.Errorf("didn't get new name %q in %q", nname, names)
	}
}

func findname(names []string, name string) int {
	for k, v := range names {
		if v == name {
			return k
		}
	}
	return -1
}
