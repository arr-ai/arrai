package rel

import "fmt"

// CompareOps maps every comparison operator's surface syntax to its implementation. It is the
// single source of truth for both compiling a CompareExpr (syntax.compile.go) and reconstructing
// one from a compiled plan (decodeCompare), so the two can't drift out of sync with each other.
var CompareOps = map[string]CompareFunc{
	"<:": func(a, b Value) (bool, error) {
		set, is := b.(Set)
		if !is {
			return false, fmt.Errorf("<: rhs not a set: %v", b)
		}
		return set.Has(a), nil
	},
	"!<:": func(a, b Value) (bool, error) {
		set, is := b.(Set)
		if !is {
			return false, fmt.Errorf("!<: rhs not a set: %v", b)
		}
		return !set.Has(a), nil
	},
	"=":  func(a, b Value) (bool, error) { return a.Equal(b), nil },
	"!=": func(a, b Value) (bool, error) { return !a.Equal(b), nil },
	"<":  func(a, b Value) (bool, error) { return a.Less(b), nil },
	">":  func(a, b Value) (bool, error) { return b.Less(a), nil },
	"<=": func(a, b Value) (bool, error) { return !b.Less(a), nil },
	">=": func(a, b Value) (bool, error) { return !a.Less(b), nil },

	"(<)":   func(a, b Value) (bool, error) { return subset(a, b), nil },
	"(>)":   func(a, b Value) (bool, error) { return subset(b, a), nil },
	"(<=)":  func(a, b Value) (bool, error) { return subsetOrEqual(a, b), nil },
	"(>=)":  func(a, b Value) (bool, error) { return subsetOrEqual(b, a), nil },
	"(<>)":  func(a, b Value) (bool, error) { return subsetOrSuperset(a, b), nil },
	"(<>=)": func(a, b Value) (bool, error) { return subsetSupersetOrEqual(b, a), nil },

	"!(<)":   func(a, b Value) (bool, error) { return !subset(a, b), nil },
	"!(>)":   func(a, b Value) (bool, error) { return !subset(b, a), nil },
	"!(<=)":  func(a, b Value) (bool, error) { return !subsetOrEqual(a, b), nil },
	"!(>=)":  func(a, b Value) (bool, error) { return !subsetOrEqual(b, a), nil },
	"!(<>)":  func(a, b Value) (bool, error) { return !subsetOrSuperset(a, b), nil },
	"!(<>=)": func(a, b Value) (bool, error) { return !subsetSupersetOrEqual(b, a), nil },
}

func subset(a, b Value) bool {
	s := a.(Set)
	t := b.(Set)
	if t.Count() == 0 {
		return false
	}
	for e := s.Enumerator(); e.MoveNext(); {
		if !t.Has(e.Current()) {
			return false
		}
	}
	return s.Count() < t.Count()
}

func subsetOrEqual(a, b Value) bool {
	s := a.(Set)
	t := b.(Set)
	if t.Count() == 0 {
		return s.Count() == 0
	}
	for e := s.Enumerator(); e.MoveNext(); {
		if !t.Has(e.Current()) {
			return false
		}
	}
	return s.Count() <= t.Count()
}

func subsetOrSuperset(a, b Value) bool {
	return subset(a, b) || subset(b, a) && !a.Equal(b)
}

func subsetSupersetOrEqual(a, b Value) bool {
	return subset(a, b) || subset(b, a) || a.Equal(b)
}
