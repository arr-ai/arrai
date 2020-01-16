package rel

import (
	"sort"

	"github.com/go-errors/errors"
)

// Intersect returns every Value from a that is also in b.
func Intersect(a, b Set) Set {
	return a.Where(func(v Value) bool { return b.Has(v) })
}

// Union returns every Values that is in either input Set or both.
func Union(a, b Set) Set {
	for e := b.Enumerator(); e.MoveNext(); {
		a = a.With(e.Current())
	}
	return a
}

// Difference returns every Value from the first Set that is not in the second.
func Difference(a, b Set) Set {
	return a.Where(func(v Value) bool { return !b.Has(v) })
}

// SymmetricDifference returns Values in either Set, but not in both.
func SymmetricDifference(a, b Set) Set {
	return Union(Difference(a, b), Difference(b, a))
}

// Concatenate is equivalent to a <&> (b => . + {.@ + a count}). Naturally, this
// assumes that every element in b is a tuple with at least an '@' attribute,
// which is numeric.
//
// E.g., [1, 2] + [3] = [1, 2, 3]; "hell" + "o" = "hello"
func Concatenate(a, b Set) (Set, error) {
	offset := a.Count()
	for e := b.Enumerator(); e.MoveNext(); {
		elt := e.Current()
		if t, ok := elt.(Tuple); ok {
			if pos, found := t.Get("@"); found {
				if n, ok := pos.(Number); ok {
					t = t.With("@", NewNumber(float64(offset)+n.Float64()))
					a = a.With(t)
					continue
				}
			}
		}
		return nil, errors.Errorf("Mismatched elt in set + set: %v", elt)
	}
	return a, nil
}

// Order returns a slice with the sets Values sorted by the given key.
func Order(s Set, key func(v Value) (Value, error)) ([]Value, error) {
	o := newOrderer(s.Count())
	i := 0
	for e := s.Enumerator(); e.MoveNext(); {
		value := e.Current()
		o.values[i] = value
		var err error
		o.keys[i], err = key(value)
		if err != nil {
			return nil, err
		}
		i++
	}
	sort.Sort(o)
	return o.values, nil
}

type orderer struct {
	values []Value
	keys   []Value
}

func newOrderer(n int) *orderer {
	buf := make([]Value, 2*n)
	return &orderer{buf[:n], buf[n:]}
}

// Len is the number of elements in the collection.
func (o *orderer) Len() int {
	return len(o.values)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (o *orderer) Less(i, j int) bool {
	return o.keys[i].Less(o.keys[j])
}

// Swap swaps the elements with indexes i and j.
func (o *orderer) Swap(i, j int) {
	o.values[i], o.values[j] = o.values[j], o.values[i]
	o.keys[i], o.keys[j] = o.keys[j], o.keys[i]
}

// PowerSet computes the power set of a set.
func PowerSet(s Set) Set {
	result := NewSet(None)
	for e := s.Enumerator(); e.MoveNext(); {
		c := e.Current()
		newSets := NewSet()
		for s := result.Enumerator(); s.MoveNext(); {
			newSets = newSets.With(s.Current().(Set).With(c))
		}
		result = Union(result, newSets)
	}
	return result
}
