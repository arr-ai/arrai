package rel

import (
	"slices"

	"github.com/arr-ai/frozen"
)

// Intersect returns every Value from a that is also in b.
func Intersect(a, b Set) Set {
	if _, is := a.(EmptySet); is {
		return a
	}
	if _, is := b.(EmptySet); is {
		return b
	}
	if ga, ok := a.(GenericSet); ok {
		if gb, ok := b.(GenericSet); ok {
			return newSetFromFrozenSet(ga.set.Intersection(gb.set))
		}
	}

	au, aUnion := a.(UnionSet)
	bu, bUnion := b.(UnionSet)
	switch {
	case aUnion && bUnion:
		keys := au.m.Keys().Intersection(bu.m.Keys())
		m := frozen.MapBuilder[string, any]{}
		for i := keys.Range(); i.Next(); {
			key := i.Value()
			if subset := Intersect(au.getSubset(key), bu.getSubset(key)); subset.IsTrue() {
				m.Put(key, subset)
			}
		}
		return newSetFromBuckets(m.Finish())
	case aUnion || bUnion:
		if bUnion {
			a, b = b, a
		}
		return Intersect(a.(UnionSet).getSubset(b.unionSetSubsetBucket()), b)
	}

	result, err := a.Where(func(v Value) (bool, error) { return b.Has(v), nil })
	if err != nil {
		panic(err)
	}
	return result
}

// NIntersect returns every Value from a that is also in all bs.
func NIntersect(a Set, bs ...Set) Set {
	for _, b := range bs {
		a = Intersect(a, b)
	}
	return a
}

// Union returns every value that is in either input Set (or both).
func Union(a, b Set) Set {
	if _, is := a.(EmptySet); is {
		return b
	}
	if _, is := b.(EmptySet); is {
		return a
	}
	if ga, ok := a.(GenericSet); ok {
		if gb, ok := b.(GenericSet); ok {
			// Adding a small set to a large one goes element by element:
			// frozen's tree merge rebuilds far more than it shares, which
			// made accumulator folds (acc | {x}) quadratic with a huge
			// constant — measured 878ms/691MB for a 3,000-step fold, against
			// ~3ms via With.
			// A union of two canonical generic sets is canonical: no
			// element's bucket changes, so skip re-canonicalising — which
			// otherwise rebuilds the whole accumulator on every step of an
			// acc | {x} fold (measured 878ms/691MB for 3,000 steps).
			if fastPaths && ga.canonical && gb.canonical {
				big, small := ga.set, gb.set
				if big.Count() < small.Count() {
					big, small = small, big
				}
				if small.Count() <= 8 && big.Count() >= 4*small.Count() {
					for i := small.Range(); i.Next(); {
						big = big.With(i.Value())
					}
				} else {
					big = big.Union(small)
				}
				if u := (GenericSet{set: big, canonical: true}); u.set.Count() != 0 {
					return u
				}
				return None
			}
			return CanonicalSet(newSetFromFrozenSet(ga.set.Union(gb.set)))
		}
	}

	au, aUnion := a.(UnionSet)
	bu, bUnion := b.(UnionSet)
	switch {
	case aUnion && bUnion:
		return newSetFromBuckets(
			au.m.Merge(
				bu.m,
				func(_ string, left, right any) any {
					return Union(left.(Set), right.(Set))
				},
			),
		)
	case aUnion != bUnion:
		if bUnion {
			a, b = b, a
		}
		return a.(UnionSet).unionWithSubset(b)
	case a.unionSetSubsetBucket() != b.unionSetSubsetBucket():
		m := frozen.MapBuilder[string, any]{}
		m.Put(a.unionSetSubsetBucket(), a)
		m.Put(b.unionSetSubsetBucket(), b)
		return newSetFromBuckets(m.Finish())
	default:
		// Not range-over-All: the body writes the captured accumulator,
		// which would heap-allocate a cell per Union call.
		for e := b.Enumerator(); e.MoveNext(); {
			a = a.With(e.Current())
		}
		return CanonicalSet(a)
	}
}

func NUnion(sets ...Set) Set {
	result := None
	for _, s := range sets {
		result = Union(result, s)
	}
	return result
}

// Difference returns every Value from the first Set that is not in the second.
func Difference(a, b Set) Set {
	if _, is := a.(EmptySet); is {
		return a
	}
	if _, is := b.(EmptySet); is {
		return a
	}
	if ga, ok := a.(GenericSet); ok {
		if gb, ok := b.(GenericSet); ok {
			return newSetFromFrozenSet(ga.set.Difference(gb.set))
		}
	}
	au, aUnion := a.(UnionSet)
	bu, bUnion := b.(UnionSet)
	switch {
	case aUnion && bUnion:
		m := frozen.MapBuilder[string, any]{}
		for i := au.m.Range(); i.Next(); {
			bucket, subset := i.Entry()
			if d := Difference(subset.(Set), bu.getSubset(bucket)); d.IsTrue() {
				m.Put(bucket, d)
			}
		}
		return newSetFromBuckets(m.Finish())
	case aUnion:
		key := b.unionSetSubsetBucket()
		if diff := Difference(au.getSubset(key), b); diff.IsTrue() {
			return newSetFromBuckets(au.m.With(key, diff))
		}
		return newSetFromBuckets(au.m.Without(key))
	case bUnion:
		return Difference(a, bu.getSubset(a.unionSetSubsetBucket()))
	default:
		result, err := a.Where(func(v Value) (bool, error) { return !b.Has(v), nil })
		if err != nil {
			panic(err)
		}
		return result
	}
}

// SymmetricDifference returns Values in either Set, but not in both.
func SymmetricDifference(a, b Set) Set {
	if _, is := a.(EmptySet); is {
		return b
	}
	if _, is := b.(EmptySet); is {
		return a
	}
	if ga, ok := a.(GenericSet); ok {
		if gb, ok := b.(GenericSet); ok {
			return newSetFromFrozenSet(ga.set.SymmetricDifference(gb.set))
		}
	}
	return Union(Difference(a, b), Difference(b, a))
}

// OrderBy returns a slice with the sets Values sorted by the given key.
func OrderBy(s Set, key func(v Value) (Value, error), less func(a, b Value) bool) ([]Value, error) {
	type kv struct{ key, value Value }
	pairs := make([]kv, s.Count())
	for i, e := 0, s.Enumerator(); e.MoveNext(); i++ {
		value := e.Current()
		k, err := key(value)
		if err != nil {
			return nil, err
		}
		pairs[i] = kv{key: k, value: value}
	}
	slices.SortStableFunc(pairs, func(a, b kv) int {
		if less(a.key, b.key) {
			return -1
		} else if less(b.key, a.key) {
			return 1
		}
		return 0
	})
	values := make([]Value, len(pairs))
	for i, p := range pairs {
		values[i] = p.value
	}
	return values, nil
}

func OrderedValueEnumerator(e ValueEnumerator, less Less) ValueEnumerator {
	if less == nil {
		return e
	}
	var values []Value
	for e.MoveNext() {
		values = append(values, e.Current())
	}
	slices.SortFunc(values, func(a, b Value) int {
		if less(a, b) {
			return -1
		} else if less(b, a) {
			return 1
		}
		return 0
	})
	return &valueSliceEnumerator{values: values, i: -1}
}

func ValueLess(a, b Value) bool {
	return a.Less(b)
}

type valueSliceEnumerator struct {
	values []Value
	i      int
}

func (e *valueSliceEnumerator) MoveNext() bool {
	if e.i >= len(e.values)-1 {
		return false
	}
	e.i++
	return true
}

func (e *valueSliceEnumerator) Current() Value {
	return e.values[e.i]
}

// PowerSet computes the power set of a set.
func PowerSet(s Set) (Set, error) {
	if _, is := s.(EmptySet); is {
		return NewSet(None)
	}
	if gs, ok := s.(GenericSet); ok {
		var sb frozen.SetBuilder[Value]
		for i := frozen.Powerset(gs.set).Range(); i.Next(); {
			sb.Add(newSetFromFrozenSet(i.Value()))
		}
		return newSetFromFrozenSet(sb.Finish()), nil
	}
	result, err := NewSet(None)
	if err != nil {
		return nil, err
	}
	for e := s.Enumerator(); e.MoveNext(); {
		c := e.Current()
		newSets := None
		for s := result.Enumerator(); s.MoveNext(); {
			newSets = newSets.With(s.Current().(Set).With(c))
		}
		result = Union(result, newSets)
	}
	return result, nil
}
