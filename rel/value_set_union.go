package rel

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/arr-ai/frozen"
	"github.com/arr-ai/wbnf/parser"
)

type UnionSet struct {
	m frozen.StringMap
}

func newSetFromBuckets(m frozen.StringMap) Set {
	switch m.Count() {
	case 0:
		return None
	case 1:
		if i := m.Range(); i.Next() {
			return i.Value().(Set)
		}
	}
	return UnionSet{m}
}

func toUnionSetWithItem(s Set, v Value) Set {
	valueBucket := v.getBucket().String()
	subsetBucket := s.unionSetSubsetBucket()
	if valueBucket == subsetBucket {
		panic(fmt.Errorf("toUnionSetWithItem expects that the set bucket and value bucket are different"))
	}
	mb := frozen.StringMapBuilder{}
	mb.Put(subsetBucket, s)
	mb.Put(valueBucket, MustNewSet(v))
	return newSetFromBuckets(mb.Finish())
}

func (u UnionSet) Count() int {
	count := 0
	for i := u.bucketRange(); i.next(); {
		count += i.subset().Count()
	}
	return count
}

func (u UnionSet) Has(v Value) bool {
	b, has := u.m.Get(v.getBucket().String())
	return has && b.(Set).Has(v)
}

type unionSetEnumerator struct {
	set     *unionSetBucketRange
	current ValueEnumerator
}

func (e *unionSetEnumerator) MoveNext() bool {
	if e.current != nil && e.current.MoveNext() {
		return true
	}
	if !e.set.next() {
		return false
	}
	e.current = e.set.subset().Enumerator()
	return e.current.MoveNext()
}

func (e *unionSetEnumerator) Current() Value {
	return e.current.Current()
}

func (u UnionSet) Enumerator() ValueEnumerator {
	return &unionSetEnumerator{
		set:     u.bucketRange(),
		current: nil,
	}
}

type unionSetOrderedEnumerator struct {
	set     frozen.Iterator
	current ValueEnumerator
}

func (e *unionSetOrderedEnumerator) MoveNext() bool {
	if e.current != nil && e.current.MoveNext() {
		return true
	}
	if !e.set.Next() {
		return false
	}
	e.current = e.set.Value().(Set).ArrayEnumerator()
	return e.current.MoveNext()
}

func (e *unionSetOrderedEnumerator) Current() Value {
	return e.current.Current()
}

func (u UnionSet) ArrayEnumerator() ValueEnumerator {
	return &unionSetOrderedEnumerator{
		// ordered by rel.Set because the bucket keys are strings
		// which wouldn't provide the correct sorting based on type.
		set: u.m.Values().OrderedRange(
			func(a, b interface{}) bool { return a.(Set).Less(b.(Set)) },
		),
		current: nil,
	}
}

func (u UnionSet) With(v Value) Set {
	bucket := v.getBucket().String()
	return newSetFromBuckets(u.m.With(bucket, u.getSubset(bucket).With(v)))
}

func (u UnionSet) unionWithSubset(subset Set) Set {
	bucket := subset.unionSetSubsetBucket()
	return newSetFromBuckets(u.m.With(bucket, Union(u.getSubset(bucket), subset)))
}

func (u UnionSet) getSubset(t string) Set {
	if subset, has := u.m.Get(t); has {
		return subset.(Set)
	}
	return None
}

func (u UnionSet) Without(v Value) Set {
	if !u.Has(v) {
		return u
	}
	bucket := v.getBucket().String()
	newSet := u.getSubset(bucket).Without(v)
	if !newSet.IsTrue() {
		return newSetFromBuckets(u.m.Without(frozen.NewSetFromStrings(bucket)))
	}
	return newSetFromBuckets(u.m.With(bucket, newSet))
}

func (u UnionSet) Map(f func(Value) (Value, error)) (Set, error) {
	sb := NewSetBuilder()
	for i := u.bucketRange(); i.next(); {
		mappedSubset, err := i.subset().Map(f)
		if err != nil {
			return nil, err
		}
		for j := mappedSubset.Enumerator(); j.MoveNext(); {
			sb.Add(j.Current())
		}
	}
	return sb.Finish()
}

func (u UnionSet) Where(f func(Value) (bool, error)) (Set, error) {
	var mb frozen.StringMapBuilder
	for i := u.bucketRange(); i.next(); {
		v, err := i.subset().Where(f)
		if err != nil {
			return nil, err
		}
		if v.IsTrue() {
			mb.Put(i.bucketKey(), v)
		}
	}
	return newSetFromBuckets(mb.Finish()), nil
}

func (u UnionSet) CallAll(ctx context.Context, arg Value, sb SetBuilder) error {
	for i := u.m.Range(); i.Next(); {
		if err := i.Value().(Set).CallAll(ctx, arg, sb); err != nil {
			return err
		}
	}
	return nil
}

func (UnionSet) unionSetSubsetBucket() string {
	panic("UnionSet.unionSetSubsetBucket should not be called")
}

var unionSetKind = registerKind(210, reflect.TypeOf(UnionSet{}))

func (u UnionSet) Kind() int {
	return unionSetKind
}

func (u UnionSet) IsTrue() bool {
	return !u.m.IsEmpty()
}

func (u UnionSet) Less(v Value) bool {
	if u.Kind() != v.Kind() {
		return u.Kind() < v.Kind()
	}
	x := v.(UnionSet)
	less := func(a, b interface{}) bool {
		return a.(Set).Less(b.(Set))
	}
	a := u.m.Values().OrderedRange(less)
	b := x.m.Values().OrderedRange(less)
	for {
		aHasMore, bHasMore := a.Next(), b.Next()
		switch {
		case !aHasMore:
			return bHasMore
		case !bHasMore:
			return false
		}
		aSubset, bSubset := a.Value().(Set), b.Value().(Set)
		if aSubset.Less(bSubset) {
			return true
		}
		if bSubset.Less(aSubset) {
			return false
		}
	}
}

func (u UnionSet) Negate() Value {
	if !u.IsTrue() {
		return u
	}
	return NewTuple(NewAttr(negateTag, u))
}

func (u UnionSet) Export(ctx context.Context) interface{} {
	if u.m.IsEmpty() {
		return []interface{}{}
	}
	result := make([]interface{}, 0, u.Count())
	for e := u.Enumerator(); e.MoveNext(); {
		result = append(result, e.Current().Export(ctx))
	}
	return result
}

func (UnionSet) getSetBuilder() setBuilder {
	panic("UnionSet.getSetBuilder should not be called")
}

func (UnionSet) getBucket() fmt.Stringer {
	panic("UnionSet.getBucket should not be called")
}

func (u UnionSet) Eval(ctx context.Context, local Scope) (Value, error) {
	return u, nil
}

func (u UnionSet) Source() parser.Scanner {
	return *parser.NewScanner("")
}

func (u UnionSet) String() string {
	var sb strings.Builder
	sb.WriteString("{")
	for i, v := range u.OrderedValues() {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(v.String())
	}
	sb.WriteString("}")
	return sb.String()
}

func (u UnionSet) Equal(s interface{}) bool {
	if t, ok := s.(UnionSet); ok {
		return u.m.Equal(t.m)
	}
	return false
}

func (u UnionSet) Hash(seed uintptr) uintptr {
	h := seed
	for e := u.Enumerator(); e.MoveNext(); {
		h ^= e.Current().Hash(0)
	}
	return h
}

func (u UnionSet) OrderedValues() []Value {
	sets := make([]Set, 0, u.m.Count())
	count := 0
	for i := u.bucketRange(); i.next(); {
		s := i.subset()
		count += s.Count()
		sets = append(sets, s)
	}
	values := make([]Value, 0, count)
	for _, set := range sets {
		for i := set.Enumerator(); i.MoveNext(); {
			values = append(values, i.Current())
		}
	}
	sort.Sort(ValueList(values))
	return values
}

type unionSetBucketRange struct {
	i *frozen.StringMapIterator
}

func (u UnionSet) bucketRange() *unionSetBucketRange {
	return &unionSetBucketRange{u.m.Range()}
}

func (ur *unionSetBucketRange) next() bool {
	return ur.i.Next()
}

func (ur *unionSetBucketRange) subset() Set {
	return ur.i.Value().(Set)
}

func (ur *unionSetBucketRange) bucketKey() string {
	return ur.i.Key()
}
