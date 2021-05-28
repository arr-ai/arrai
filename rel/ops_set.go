package rel

import (
	"sort"

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
		m := frozen.StringMapBuilder{}
		for i := keys.Range(); i.Next(); {
			key := i.Value().(string)
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
				func(_, left, right interface{}) interface{} {
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
		m := frozen.StringMapBuilder{}
		m.Put(a.unionSetSubsetBucket(), a)
		m.Put(b.unionSetSubsetBucket(), b)
		return newSetFromBuckets(m.Finish())
	default:
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
		m := frozen.StringMapBuilder{}
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
		return newSetFromBuckets(au.m.Without(frozen.NewSet(key)))
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
	o := newOrderer(s.Count(), less)
	for i, e := 0, s.Enumerator(); e.MoveNext(); i++ {
		value := e.Current()
		o.values[i] = value
		var err error
		o.keys[i], err = key(value)
		if err != nil {
			return nil, err
		}
	}
	sort.Sort(o)
	return o.values, nil
}

func OrderedValueEnumerator(e ValueEnumerator, less Less) ValueEnumerator {
	if less == nil {
		return e
	}
	var values []Value
	for e.MoveNext() {
		values = append(values, e.Current())
	}
	return &valueSliceEnumerator{values: values, i: -1}
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

type orderer struct {
	values []Value
	keys   []Value
	less   func(a, b Value) bool
}

func newOrderer(n int, less func(a, b Value) bool) *orderer {
	buf := make([]Value, 2*n)
	return &orderer{values: buf[:n], keys: buf[n:], less: less}
}

// Len is the number of elements in the collection.
func (o *orderer) Len() int {
	return len(o.values)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (o *orderer) Less(i, j int) bool {
	return o.less(o.keys[i], o.keys[j])
}

// Swap swaps the elements with indexes i and j.
func (o *orderer) Swap(i, j int) {
	o.values[i], o.values[j] = o.values[j], o.values[i]
	o.keys[i], o.keys[j] = o.keys[j], o.keys[i]
}

// PowerSet computes the power set of a set.
func PowerSet(s Set) (Set, error) {
	if _, is := s.(EmptySet); is {
		return NewSet(None)
	}
	if gs, ok := s.(GenericSet); ok {
		var sb frozen.SetBuilder
		for i := gs.set.Powerset().Range(); i.Next(); {
			sb.Add(newSetFromFrozenSet(i.Value().(frozen.Set)))
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
