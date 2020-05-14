package syntax

import (
	"github.com/arr-ai/arrai/rel"
)

func subset(a, b rel.Value) bool {
	s := a.(rel.Set)
	t := b.(rel.Set)
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

func subsetOrEqual(a, b rel.Value) bool {
	s := a.(rel.Set)
	t := b.(rel.Set)
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

func subsetOrSuperset(a, b rel.Value) bool {
	return subset(a, b) || subset(b, a) && !a.Equal(b)
}

func subsetSupersetOrEqual(a, b rel.Value) bool {
	return subset(a, b) || subset(b, a) || a.Equal(b)
}
