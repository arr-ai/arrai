package rel

// Targeted micro-benchmarks to isolate specific performance factors in the
// frozen generics migration. These benchmarks focus on the hotspots identified
// by CPU and memory profiling of the join benchmarks.

import (
	"testing"

	"github.com/arr-ai/frozen"
)

// --- Tuple creation benchmarks ---

func BenchmarkNewTuple2Attrs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewTuple(
			NewAttr("a", NewNumber(float64(i))),
			NewAttr("b", NewNumber(float64(i+1))),
		)
	}
}

func BenchmarkNewTuple5Attrs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewTuple(
			NewAttr("a", NewNumber(float64(i))),
			NewAttr("b", NewNumber(float64(i+1))),
			NewAttr("c", NewNumber(float64(i+2))),
			NewAttr("d", NewNumber(float64(i+3))),
			NewAttr("e", NewNumber(float64(i+4))),
		)
	}
}

// --- Tuple access benchmarks ---

func BenchmarkTupleGet(b *testing.B) {
	t := NewTuple(
		NewAttr("a", NewNumber(1)),
		NewAttr("b", NewNumber(2)),
		NewAttr("c", NewNumber(3)),
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t.Get("b")
	}
}

func BenchmarkTupleEnumerator(b *testing.B) {
	t := NewTuple(
		NewAttr("a", NewNumber(1)),
		NewAttr("b", NewNumber(2)),
		NewAttr("c", NewNumber(3)),
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for e := t.Enumerator(); e.MoveNext(); {
			e.Current()
		}
	}
}

// --- Tuple equality benchmarks ---

func BenchmarkTupleEqual2Attrs(b *testing.B) {
	t1 := NewTuple(
		NewAttr("a", NewNumber(1)),
		NewAttr("b", NewNumber(2)),
	)
	t2 := NewTuple(
		NewAttr("a", NewNumber(1)),
		NewAttr("b", NewNumber(2)),
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t1.Equal(t2)
	}
}

func BenchmarkTupleEqual5Attrs(b *testing.B) {
	t1 := NewTuple(
		NewAttr("a", NewNumber(1)),
		NewAttr("b", NewNumber(2)),
		NewAttr("c", NewNumber(3)),
		NewAttr("d", NewNumber(4)),
		NewAttr("e", NewNumber(5)),
	)
	t2 := NewTuple(
		NewAttr("a", NewNumber(1)),
		NewAttr("b", NewNumber(2)),
		NewAttr("c", NewNumber(3)),
		NewAttr("d", NewNumber(4)),
		NewAttr("e", NewNumber(5)),
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t1.Equal(t2)
	}
}

// --- Hash benchmarks ---

func BenchmarkTupleHash2Attrs(b *testing.B) {
	t := NewTuple(
		NewAttr("a", NewNumber(1)),
		NewAttr("b", NewNumber(2)),
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t.Hash(0)
	}
}

func BenchmarkNumberHash(b *testing.B) {
	n := NewNumber(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n.Hash(0)
	}
}

// --- Names benchmarks ---

func BenchmarkNamesCreate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewNames("a", "b", "c")
	}
}

func BenchmarkNamesHas(b *testing.B) {
	n := NewNames("a", "b", "c", "d", "e")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n.Has("c")
	}
}

func BenchmarkNamesEqual(b *testing.B) {
	n1 := NewNames("a", "b", "c")
	n2 := NewNames("a", "b", "c")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n1.Equal(n2)
	}
}

func BenchmarkNamesIntersect(b *testing.B) {
	n1 := NewNames("a", "b", "c")
	n2 := NewNames("b", "c", "d")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n1.Intersect(n2)
	}
}

func BenchmarkNamesMinus(b *testing.B) {
	n1 := NewNames("a", "b", "c")
	n2 := NewNames("b")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n1.Minus(n2)
	}
}

func BenchmarkNamesToSlice(b *testing.B) {
	n := NewNames("a", "b", "c", "d", "e")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n.Names()
	}
}

// --- frozen.Set[Value] benchmarks (isolating HAMT overhead) ---

func BenchmarkFrozenSetWithValue(b *testing.B) {
	// Benchmark adding single elements to a frozen set
	s := frozen.Set[Value]{}
	vals := make([]Value, 1000)
	for i := range vals {
		vals[i] = NewNumber(float64(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = s.With(vals[i%1000])
	}
}

func BenchmarkFrozenSetBuilderValue(b *testing.B) {
	vals := make([]Value, 1000)
	for i := range vals {
		vals[i] = NewNumber(float64(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sb := frozen.NewSetBuilder[Value](1000)
		for _, v := range vals {
			sb.Add(v)
		}
		sb.Finish()
	}
}

func BenchmarkFrozenSetHasValue(b *testing.B) {
	sb := frozen.NewSetBuilder[Value](1000)
	for i := 0; i < 1000; i++ {
		sb.Add(NewNumber(float64(i)))
	}
	s := sb.Finish()
	target := NewNumber(500)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Has(target)
	}
}

func BenchmarkFrozenSetRangeValue(b *testing.B) {
	sb := frozen.NewSetBuilder[Value](1000)
	for i := 0; i < 1000; i++ {
		sb.Add(NewNumber(float64(i)))
	}
	s := sb.Finish()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for iter := s.Range(); iter.Next(); {
			_ = iter.Value()
		}
	}
}

// --- frozen.Map[string, Value] benchmarks (isolating tuple map overhead) ---

func BenchmarkFrozenMapGetStringValue(b *testing.B) {
	mb := frozen.MapBuilder[string, Value]{}
	mb.Put("a", NewNumber(1))
	mb.Put("b", NewNumber(2))
	mb.Put("c", NewNumber(3))
	m := mb.Finish()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get("b")
	}
}

func BenchmarkFrozenMapRangeStringValue(b *testing.B) {
	mb := frozen.MapBuilder[string, Value]{}
	mb.Put("a", NewNumber(1))
	mb.Put("b", NewNumber(2))
	mb.Put("c", NewNumber(3))
	m := mb.Finish()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for iter := m.Range(); iter.Next(); {
			_ = iter.Key()
			_ = iter.Value()
		}
	}
}

func BenchmarkFrozenMapWithStringValue(b *testing.B) {
	mb := frozen.MapBuilder[string, Value]{}
	mb.Put("a", NewNumber(1))
	mb.Put("b", NewNumber(2))
	m := mb.Finish()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.With("c", NewNumber(float64(i)))
	}
}

// --- frozen.Set[string] benchmarks (isolating Names/string set overhead) ---

func BenchmarkFrozenSetHasString(b *testing.B) {
	sb := frozen.SetBuilder[string]{}
	sb.Add("a")
	sb.Add("b")
	sb.Add("c")
	sb.Add("d")
	sb.Add("e")
	s := sb.Finish()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Has("c")
	}
}

func BenchmarkFrozenSetRangeString(b *testing.B) {
	sb := frozen.SetBuilder[string]{}
	sb.Add("a")
	sb.Add("b")
	sb.Add("c")
	s := sb.Finish()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for iter := s.Range(); iter.Next(); {
			_ = iter.Value()
		}
	}
}

// --- SetBuilder benchmarks (isolating the arrai set builder overhead) ---

func BenchmarkSetBuilderAdd100(b *testing.B) {
	tuples := make([]Value, 100)
	for i := range tuples {
		tuples[i] = NewTuple(
			NewAttr("a", NewNumber(float64(i))),
			NewAttr("b", NewNumber(float64(i+1))),
		)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sb := NewSetBuilder()
		for _, t := range tuples {
			sb.Add(t)
		}
		sb.Finish()
	}
}

func BenchmarkSetBuilderAdd1000(b *testing.B) {
	tuples := make([]Value, 1000)
	for i := range tuples {
		tuples[i] = NewTuple(
			NewAttr("a", NewNumber(float64(i))),
			NewAttr("b", NewNumber(float64(i+1))),
		)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sb := NewSetBuilder()
		for _, t := range tuples {
			sb.Add(t)
		}
		sb.Finish()
	}
}

// --- Relation builder benchmarks ---

func BenchmarkRelationBuilderAdd100(b *testing.B) {
	tuples := make([]Tuple, 100)
	for i := range tuples {
		tuples[i] = NewTuple(
			NewAttr("a", NewNumber(float64(i))),
			NewAttr("b", NewNumber(float64(i+1))),
		)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb := newRelationBuilder([]string{"a", "b"}, 100)
		for _, t := range tuples {
			rb.Add(t)
		}
		rb.Finish()
	}
}

func BenchmarkRelationBuilderAdd1000(b *testing.B) {
	tuples := make([]Tuple, 1000)
	for i := range tuples {
		tuples[i] = NewTuple(
			NewAttr("a", NewNumber(float64(i))),
			NewAttr("b", NewNumber(float64(i+1))),
		)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb := newRelationBuilder([]string{"a", "b"}, 1000)
		for _, t := range tuples {
			rb.Add(t)
		}
		rb.Finish()
	}
}

// --- Relation enumeration benchmarks ---

func BenchmarkRelationEnumerate100(b *testing.B) {
	r := makeRelation(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for e := r.Enumerator(); e.MoveNext(); {
			e.Current()
		}
	}
}

func BenchmarkRelationEnumerate1000(b *testing.B) {
	r := makeRelation(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for e := r.Enumerator(); e.MoveNext(); {
			e.Current()
		}
	}
}

func BenchmarkGenericSetEnumerate100(b *testing.B) {
	g := makeGenericSet(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for e := g.Enumerator(); e.MoveNext(); {
			e.Current()
		}
	}
}

func BenchmarkGenericSetEnumerate1000(b *testing.B) {
	g := makeGenericSet(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for e := g.Enumerator(); e.MoveNext(); {
			e.Current()
		}
	}
}

// --- Set.With benchmarks (measures HAMT path copying) ---

func BenchmarkRelationWith(b *testing.B) {
	r := makeRelation(100)
	newTuples := make([]Tuple, 100)
	for i := range newTuples {
		newTuples[i] = NewTuple(
			NewAttr("a", NewNumber(float64(i+1000))),
			NewAttr("b", NewNumber(float64(i+1001))),
		)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.With(newTuples[i%100])
	}
}

func BenchmarkGenericSetWith(b *testing.B) {
	g := makeGenericSet(100)
	newTuples := make([]Value, 100)
	for i := range newTuples {
		newTuples[i] = NewTuple(
			NewAttr("a", NewNumber(float64(i+1000))),
			NewAttr("b", NewNumber(float64(i+1001))),
		)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.With(newTuples[i%100])
	}
}

// --- Merge/Combine benchmarks (hot in GenericJoin) ---

func BenchmarkMergeTuples(b *testing.B) {
	t1 := NewTuple(
		NewAttr("a", NewNumber(1)),
		NewAttr("b", NewNumber(2)),
	)
	t2 := NewTuple(
		NewAttr("b", NewNumber(2)),
		NewAttr("c", NewNumber(3)),
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MergeTuples(t1, t2)
	}
}

// --- Value equality dispatch benchmarks ---

func BenchmarkNumberEqual(b *testing.B) {
	n1 := NewNumber(42)
	n2 := NewNumber(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n1.Equal(n2)
	}
}

func BenchmarkValueEqualViaInterface(b *testing.B) {
	var v1 Value = NewNumber(42)
	var v2 Value = NewNumber(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v1.Equal(v2)
	}
}

// --- frozen.Set[Value] equality (the #1 allocation hotspot) ---

func BenchmarkFrozenSetEqualValue100(b *testing.B) {
	sb1 := frozen.NewSetBuilder[Value](100)
	sb2 := frozen.NewSetBuilder[Value](100)
	for i := 0; i < 100; i++ {
		sb1.Add(NewNumber(float64(i)))
		sb2.Add(NewNumber(float64(i)))
	}
	s1 := sb1.Finish()
	s2 := sb2.Finish()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s1.Equal(s2)
	}
}

func BenchmarkFrozenSetEqualValue1000(b *testing.B) {
	sb1 := frozen.NewSetBuilder[Value](1000)
	sb2 := frozen.NewSetBuilder[Value](1000)
	for i := 0; i < 1000; i++ {
		sb1.Add(NewNumber(float64(i)))
		sb2.Add(NewNumber(float64(i)))
	}
	s1 := sb1.Finish()
	s2 := sb2.Finish()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s1.Equal(s2)
	}
}

// --- frozen.Map[string, Value] equality (used by tuple equality) ---

func BenchmarkFrozenMapEqualStringValue(b *testing.B) {
	build := func() frozen.Map[string, Value] {
		mb := frozen.MapBuilder[string, Value]{}
		mb.Put("a", NewNumber(1))
		mb.Put("b", NewNumber(2))
		mb.Put("c", NewNumber(3))
		return mb.Finish()
	}
	m1 := build()
	m2 := build()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m1.Equal(m2)
	}
}

// --- Helper functions ---

func makeRelation(n int) Relation {
	rb := newRelationBuilder([]string{"a", "b"}, n)
	for i := 0; i < n; i++ {
		rb.Add(NewTuple(
			NewAttr("a", NewNumber(float64(i))),
			NewAttr("b", NewNumber(float64(i+1))),
		))
	}
	s, err := rb.Finish()
	if err != nil {
		panic(err)
	}
	return s.(Relation)
}

func makeGenericSet(n int) GenericSet {
	sb := frozen.NewSetBuilder[Value](n)
	for i := 0; i < n; i++ {
		sb.Add(NewTuple(
			NewAttr("a", NewNumber(float64(i))),
			NewAttr("b", NewNumber(float64(i+1))),
		))
	}
	return GenericSet{sb.Finish()}
}
