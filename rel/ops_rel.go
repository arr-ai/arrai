package rel

import (
	"fmt"

	"github.com/arr-ai/frozen"
	"github.com/go-errors/errors"
)

// RelationAttrs returns the set of names for a relation type, or false if the
// set isn't a regular relation.
func RelationAttrs(a Set) (Names, bool) {
	e := a.Enumerator()
	if !e.MoveNext() {
		return Names{}, true
	}
	names := e.Current().(Tuple).Names()
	for e.MoveNext() {
		if !names.Equal(e.Current().(Tuple).Names()) {
			return Names{}, false
		}
	}
	return names, true
}

// Nest groups the given attributes into nested relations.
func Nest(a Set, attrs Names, attr string) Set {
	if !a.IsTrue() {
		return a
	}
	names, ok := RelationAttrs(a)
	if !ok {
		panic("Tuple names mismatch in nest lhs")
	}
	if !attrs.IsSubsetOf(names) {
		panic(fmt.Errorf("nest attrs (%v) not a subset of relation attrs (%v)", attrs, names))
	}
	key := names.Minus(attrs)
	return Reduce(
		a,
		func(value Value) Value {
			return value.(Tuple).Project(key)
		},
		func(key Value, tuples Set) Set {
			nested := None
			for e := tuples.Enumerator(); e.MoveNext(); {
				nested = nested.With(e.Current().(Tuple).Project(attrs))
			}
			return NewSet(Merge(key.(Tuple), NewTuple(Attr{attr, nested})))
		},
	)
}

// Unnest unpacks the attributes of a nested relation into the outer relation.
func Unnest(a Set, attr string) Set {
	key, ok := RelationAttrs(a)
	if !ok {
		panic("Tuple names mismatch in unnest lhs")
	}
	if !key.Has(attr) {
		panic("Unnest attr not found in relation")
	}
	return Reduce(
		a,
		func(value Value) Value {
			return value.(Tuple).Project(key)
		},
		func(key Value, _ Set) Set {
			unnested := None
			t := key.(Tuple)
			s, _ := t.Get(attr)
			t = t.Without(attr)
			for e := s.(Set).Enumerator(); e.MoveNext(); {
				unnested = unnested.With(Merge(t, e.Current().(Tuple)))
			}
			return unnested
		},
	)
}

// Reduce reduces a set using the given key and reducer functions.
func Reduce(
	a Set,
	getKey func(value Value) Value,
	reduce func(key Value, tuples Set) Set,
) Set {
	var buckets frozen.Map
	for e := a.Enumerator(); e.MoveNext(); {
		value := e.Current()
		key := getKey(value)

		slot, found := buckets.Get(key)
		if !found {
			slot = None
		}

		slot = slot.(Set).With(value)
		buckets = buckets.With(key, slot)
	}

	result := None
	for i := buckets.Range(); i.Next(); {
		result = Union(result, reduce(i.Key().(Value), i.Value().(Set)))
	}
	return result
}

// Joiner returns a function that computes the relational join of a and b.
//
// Defn: Join(a{x…,y…}, b{y…,z…}) = ∀{x…,y…,z…}: {x…,y…} ∈ a ∧ {y…,z…} ∈ b
//         for mutually disjoint x…, y…, z…
//
// The combine function determines how to combine matching tuples from a and b.
func Joiner(combine func(common Names, a, b Tuple) Tuple) func(a, b Set) Set {
	return func(a, b Set) Set {
		aNames, ok := RelationAttrs(a)
		if !ok {
			panic("Tuple names mismatch in join lhs")
		}
		bNames, ok := RelationAttrs(b)
		if !ok {
			panic("Tuple names mismatch in join rhs")
		}
		common := aNames.Intersect(bNames)
		return GenericJoin(
			a, b,
			func(value Value) Value {
				return value.(Tuple).Project(common)
			},
			func(key Value, a, b Set) Set {
				values := []Value{}
				for i := a.Enumerator(); i.MoveNext(); {
					for j := b.Enumerator(); j.MoveNext(); {
						values = append(values, combine(
							common,
							i.Current().(Tuple),
							j.Current().(Tuple),
						))
					}
				}
				return NewSet(values...)
			},
		)
	}
}

var join func(a, b Set) Set = Joiner(func(_ Names, a, b Tuple) Tuple {
	return Merge(a, b)
})

// func Join(a, b Set) Set {
// 	aNames, ok := RelationAttrs(a)
// 	if !ok {
// 		panic("Tuple names mismatch in join lhs")
// 	}
// 	bNames, ok := RelationAttrs(b)
// 	if !ok {
// 		panic("Tuple names mismatch in join rhs")
// 	}
// 	if a.Count() > b.Count() {
// 		a, b = b, a
// 		aNames, bNames = bNames, aNames
// 	}
// 	common := aNames.Intersect(bNames)
// 	var buckets frozen.Map
// 	for e := a.Enumerator(); e.MoveNext(); {
// 		tuple := e.Current().(Tuple)
// 		key := tuple.Project(common)
// 		if bucket, found := buckets.Get(key); found {
// 			buckets, _ = buckets.Set(key, bucket.(Set).With(tuple))
// 		} else {
// 			buckets, _ = buckets.Set(key, NewSet(tuple))
// 		}
// 	}

// 	result := None
// 	for e := b.Enumerator(); e.MoveNext(); {
// 		tuple := e.Current().(Tuple)
// 		key := tuple.Project(common)
// 		if bucket, found := buckets.Get(key); found {
// 			for e := bucket.(Set).Enumerator(); e.MoveNext(); {
// 				if merged := Merge(tuple, e.Current().(Tuple)); merged != nil {
// 					result = result.With(merged)
// 				}
// 			}
// 		}
// 	}
// 	return result
// }

// GenericJoin joins two sets using a key and a joiner
func GenericJoin(
	a, b Set,
	getKey func(value Value) Value,
	join func(key Value, a, b Set) Set,
) Set {
	var mb frozen.MapBuilder
	accumulate := func(s Set, slotKey int) {
		for e := s.Enumerator(); e.MoveNext(); {
			value := e.Current()
			key := getKey(value)

			entry, found := mb.Get(key)
			if !found {
				entry = [2]Set{None, None}
			}
			slots := entry.([2]Set)

			// False denotes lhs accumulator
			slots[slotKey] = slots[slotKey].With(value)
			mb.Put(key, slots)
		}
	}

	const aSlot = 0
	const bSlot = 1

	accumulate(a, aSlot)
	accumulate(b, bSlot)

	result := None
	for i := mb.Finish().Range(); i.Next(); {
		k, v := i.Entry()
		key := k.(Value)
		slots := v.([2]Set)
		aSet := slots[aSlot]
		bSet := slots[bSlot]
		result = Union(result, join(key, aSet, bSet))
	}
	return result
}

// Concatenate is equivalent to a <&> (b => . + {.@ + a count}). Naturally, this
// assumes that every element in b is a tuple with at least an '@' attribute,
// which is numeric.
//
// E.g., [1, 2] + [3] = [1, 2, 3]; "hell" + "o" = "hello"
func Concatenate(a, b Set) (Set, error) {
	offset := a.Count()
	values := make([]Value, 0, a.Count()+b.Count())
	for e := a.Enumerator(); e.MoveNext(); {
		values = append(values, e.Current())
	}
	for e := b.Enumerator(); e.MoveNext(); {
		elt := e.Current()
		if t, ok := elt.(Tuple); ok {
			if pos, found := t.Get("@"); found {
				if n, ok := pos.(Number); ok {
					t = t.With("@", NewNumber(float64(offset)+n.Float64()))
					values = append(values, t)
					continue
				}
			}
		}
		return nil, errors.Errorf("Mismatched elt in set + set: %v", elt)
	}
	return NewSet(values...), nil
}

func NConcatenate(a Set, bs ...Set) (Set, error) {
	for _, b := range bs {
		var err error
		a, err = Concatenate(a, b)
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}
