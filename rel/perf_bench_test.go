package rel

// Supplementary performance benchmarks targeting hotspots identified by
// profiling the frozen generics migration. Complements micro_bench_test.go
// with size-parameterised variants and join-phase isolation benchmarks.

import (
	"fmt"
	"testing"

	"github.com/arr-ai/frozen"
)

// ---------------------------------------------------------------------------
// Helpers (unique names to avoid conflicts with micro_bench_test.go)
// ---------------------------------------------------------------------------

// perfMakeTuple builds a GenericTuple with n attributes named a0..aN.
func perfMakeTuple(n int) Tuple {
	attrs := make([]Attr, n)
	for i := 0; i < n; i++ {
		attrs[i] = NewAttr(fmt.Sprintf("a%d", i), NewNumber(float64(i)))
	}
	return NewTuple(attrs...)
}

// perfMakeGenericSet builds a GenericSet of count tuples, each with nAttrs attributes.
func perfMakeGenericSet(count, nAttrs int) GenericSet {
	sb := frozen.NewSetBuilder[Value](count)
	for i := 0; i < count; i++ {
		attrs := make([]Attr, nAttrs)
		for j := 0; j < nAttrs; j++ {
			attrs[j] = NewAttr(fmt.Sprintf("a%d", j), NewNumber(float64(i*nAttrs+j)))
		}
		sb.Add(NewTuple(attrs...))
	}
	return GenericSet{sb.Finish()}
}

// perfMakeNumberSet builds a GenericSet of count Numbers.
func perfMakeNumberSet(count int) GenericSet {
	sb := frozen.NewSetBuilder[Value](count)
	for i := 0; i < count; i++ {
		sb.Add(NewNumber(float64(i)))
	}
	return GenericSet{sb.Finish()}
}

// perfMakeNames creates a Names with n names: a0..aN.
func perfMakeNames(n int) Names {
	s := make([]string, n)
	for i := 0; i < n; i++ {
		s[i] = fmt.Sprintf("a%d", i)
	}
	return NewNames(s...)
}

// ---------------------------------------------------------------------------
// 1. Value boxing and equality — parameterised variants
// ---------------------------------------------------------------------------

func BenchmarkPerfNumberEqualMatch(b *testing.B) {
	b.ReportAllocs()
	a := NewNumber(42)
	c := NewNumber(42)
	for i := 0; i < b.N; i++ {
		_ = a.Equal(c)
	}
}

func BenchmarkPerfNumberEqualMismatch(b *testing.B) {
	b.ReportAllocs()
	a := NewNumber(42)
	c := NewNumber(99)
	for i := 0; i < b.N; i++ {
		_ = a.Equal(c)
	}
}

func BenchmarkPerfTupleEqual(b *testing.B) {
	for _, n := range []int{2, 5, 10} {
		b.Run(fmt.Sprintf("attrs=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			t1 := perfMakeTuple(n)
			t2 := perfMakeTuple(n)
			for i := 0; i < b.N; i++ {
				_ = t1.Equal(t2)
			}
		})
	}
}

func BenchmarkPerfStringEqual(b *testing.B) {
	for _, n := range []int{10, 100, 1000} {
		b.Run(fmt.Sprintf("len=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			runes := make([]rune, n)
			for i := range runes {
				runes[i] = 'a'
			}
			s1 := NewString(runes)
			s2 := NewString(runes)
			for i := 0; i < b.N; i++ {
				_ = s1.(Value).Equal(s2.(Value))
			}
		})
	}
}

// ---------------------------------------------------------------------------
// 2. Tuple operations — parameterised
// ---------------------------------------------------------------------------

func BenchmarkPerfTupleCreate(b *testing.B) {
	for _, n := range []int{2, 5, 10} {
		attrs := make([]Attr, n)
		for i := 0; i < n; i++ {
			attrs[i] = NewAttr(fmt.Sprintf("a%d", i), NewNumber(float64(i)))
		}
		b.Run(fmt.Sprintf("attrs=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = NewTuple(attrs...)
			}
		})
	}
}

func BenchmarkPerfTupleGet(b *testing.B) {
	for _, n := range []int{2, 5, 10} {
		b.Run(fmt.Sprintf("attrs=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			t := perfMakeTuple(n)
			key := fmt.Sprintf("a%d", n/2)
			for i := 0; i < b.N; i++ {
				_, _ = t.Get(key)
			}
		})
	}
}

func BenchmarkPerfTupleMustGet(b *testing.B) {
	for _, n := range []int{2, 5, 10} {
		b.Run(fmt.Sprintf("attrs=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			t := perfMakeTuple(n)
			key := fmt.Sprintf("a%d", n/2)
			for i := 0; i < b.N; i++ {
				_ = t.MustGet(key)
			}
		})
	}
}

func BenchmarkPerfTupleProject(b *testing.B) {
	for _, totalAttrs := range []int{3, 5, 10} {
		for _, projSize := range []int{1, totalAttrs / 2, totalAttrs - 1} {
			if projSize < 1 {
				projSize = 1
			}
			b.Run(fmt.Sprintf("total=%d/proj=%d", totalAttrs, projSize), func(b *testing.B) {
				b.ReportAllocs()
				t := perfMakeTuple(totalAttrs)
				names := perfMakeNames(projSize)
				for i := 0; i < b.N; i++ {
					_ = t.Project(names)
				}
			})
		}
	}
}

func BenchmarkPerfTupleNames(b *testing.B) {
	for _, n := range []int{2, 5, 10} {
		b.Run(fmt.Sprintf("attrs=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			t := perfMakeTuple(n)
			for i := 0; i < b.N; i++ {
				_ = t.Names()
			}
		})
	}
}

func BenchmarkPerfTupleEnumerator(b *testing.B) {
	for _, n := range []int{2, 5, 10} {
		b.Run(fmt.Sprintf("attrs=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			t := perfMakeTuple(n)
			for i := 0; i < b.N; i++ {
				e := t.Enumerator()
				for e.MoveNext() {
					_, _ = e.Current()
				}
			}
		})
	}
}

func BenchmarkPerfTupleWith(b *testing.B) {
	for _, n := range []int{2, 5, 10} {
		b.Run(fmt.Sprintf("attrs=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			t := perfMakeTuple(n)
			val := NewNumber(999)
			for i := 0; i < b.N; i++ {
				_ = t.With("newattr", val)
			}
		})
	}
}

// BenchmarkPerfTupleGetBucket measures getBucket() on tuples, which
// internally calls Names().OrderedNames(). Called per element in SetBuilder.Add.
func BenchmarkPerfTupleGetBucket(b *testing.B) {
	for _, n := range []int{2, 5, 10} {
		b.Run(fmt.Sprintf("attrs=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			t := perfMakeTuple(n)
			for i := 0; i < b.N; i++ {
				_ = t.(Value).getBucket()
			}
		})
	}
}

// ---------------------------------------------------------------------------
// 3. Set operations (GenericSet backed by frozen.Set[Value])
// ---------------------------------------------------------------------------

func BenchmarkPerfSetWith(b *testing.B) {
	for _, n := range []int{10, 100, 1000} {
		b.Run(fmt.Sprintf("size=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			s := perfMakeNumberSet(n)
			extra := NewNumber(float64(n + 1))
			for i := 0; i < b.N; i++ {
				_ = s.With(extra)
			}
		})
	}
}

func BenchmarkPerfSetBuilderAddNumbers(b *testing.B) {
	for _, n := range []int{100, 1000} {
		vals := make([]Value, n)
		for i := 0; i < n; i++ {
			vals[i] = NewNumber(float64(i))
		}
		b.Run(fmt.Sprintf("size=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				sb := frozen.NewSetBuilder[Value](n)
				for _, v := range vals {
					sb.Add(v)
				}
				_ = sb.Finish()
			}
		})
	}
}

// BenchmarkPerfSetBuilderAddTuples uses arrai's SetBuilder (buckets by type).
func BenchmarkPerfSetBuilderAddTuples(b *testing.B) {
	for _, n := range []int{100, 1000} {
		tuples := make([]Value, n)
		for i := 0; i < n; i++ {
			tuples[i] = NewTuple(
				NewAttr("a", NewNumber(float64(i))),
				NewAttr("b", NewNumber(float64(i*2))),
				NewAttr("c", NewNumber(float64(i*3))),
			)
		}
		b.Run(fmt.Sprintf("size=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				sb := NewSetBuilder()
				for _, t := range tuples {
					sb.Add(t)
				}
				_, _ = sb.Finish()
			}
		})
	}
}

func BenchmarkPerfSetHas(b *testing.B) {
	for _, n := range []int{100, 1000} {
		b.Run(fmt.Sprintf("size=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			s := perfMakeNumberSet(n)
			target := NewNumber(float64(n / 2))
			for i := 0; i < b.N; i++ {
				_ = s.Has(target)
			}
		})
	}
}

func BenchmarkPerfSetEnumerator(b *testing.B) {
	for _, n := range []int{100, 1000} {
		b.Run(fmt.Sprintf("size=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			s := perfMakeNumberSet(n)
			for i := 0; i < b.N; i++ {
				e := s.Enumerator()
				for e.MoveNext() {
					_ = e.Current()
				}
			}
		})
	}
}

func BenchmarkPerfSetEqual(b *testing.B) {
	for _, n := range []int{100, 1000} {
		b.Run(fmt.Sprintf("size=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			s1 := perfMakeNumberSet(n)
			s2 := perfMakeNumberSet(n)
			for i := 0; i < b.N; i++ {
				_ = s1.Equal(s2)
			}
		})
	}
}

func BenchmarkPerfSetUnion(b *testing.B) {
	for _, n := range []int{100, 1000} {
		b.Run(fmt.Sprintf("size=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			s1 := perfMakeNumberSet(n)
			sb := frozen.NewSetBuilder[Value](n)
			for i := n / 2; i < n+n/2; i++ {
				sb.Add(NewNumber(float64(i)))
			}
			s2 := GenericSet{sb.Finish()}
			for i := 0; i < b.N; i++ {
				_ = Union(s1, s2)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// 4. Names operations — parameterised
// ---------------------------------------------------------------------------

func BenchmarkPerfNamesIntersect(b *testing.B) {
	for _, n := range []int{3, 5, 10} {
		b.Run(fmt.Sprintf("size=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			names1 := perfMakeNames(n)
			s := make([]string, n)
			for i := 0; i < n; i++ {
				s[i] = fmt.Sprintf("a%d", i+n/2)
			}
			names2 := NewNames(s...)
			for i := 0; i < b.N; i++ {
				_ = names1.Intersect(names2)
			}
		})
	}
}

func BenchmarkPerfNamesEqual(b *testing.B) {
	for _, n := range []int{3, 5, 10} {
		b.Run(fmt.Sprintf("size=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			names1 := perfMakeNames(n)
			names2 := perfMakeNames(n)
			for i := 0; i < b.N; i++ {
				_ = names1.Equal(names2)
			}
		})
	}
}

func BenchmarkPerfNamesOrderedNames(b *testing.B) {
	for _, n := range []int{3, 5, 10} {
		b.Run(fmt.Sprintf("size=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			names := perfMakeNames(n)
			for i := 0; i < b.N; i++ {
				_ = names.OrderedNames()
			}
		})
	}
}

func BenchmarkPerfNamesCount(b *testing.B) {
	b.ReportAllocs()
	names := perfMakeNames(5)
	for i := 0; i < b.N; i++ {
		_ = names.Count()
	}
}

func BenchmarkPerfNamesHash(b *testing.B) {
	b.ReportAllocs()
	names := perfMakeNames(5)
	for i := 0; i < b.N; i++ {
		_ = names.Hash(0)
	}
}

func BenchmarkPerfNamesHas(b *testing.B) {
	for _, n := range []int{3, 5, 10} {
		b.Run(fmt.Sprintf("size=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			names := perfMakeNames(n)
			target := fmt.Sprintf("a%d", n/2)
			for i := 0; i < b.N; i++ {
				_ = names.Has(target)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// 5. Join phases isolated
// ---------------------------------------------------------------------------

// BenchmarkPerfRelationAttrsGenericSet measures RelationAttrs() on a
// GenericSet. Must iterate all elements, calling Names() per element.
func BenchmarkPerfRelationAttrsGenericSet(b *testing.B) {
	for _, n := range []int{100, 1000} {
		b.Run(fmt.Sprintf("size=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			s := perfMakeGenericSet(n, 3)
			for i := 0; i < b.N; i++ {
				_, _ = RelationAttrs(s)
			}
		})
	}
}

// BenchmarkPerfRelationAttrsRelation measures RelationAttrs on a Relation.
func BenchmarkPerfRelationAttrsRelation(b *testing.B) {
	for _, n := range []int{100, 1000} {
		b.Run(fmt.Sprintf("size=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			_, r := generateRelations(n / 3)
			for i := 0; i < b.N; i++ {
				_, _ = RelationAttrs(r)
			}
		})
	}
}

// BenchmarkPerfGenericJoinGroupBy isolates the grouping phase of GenericJoin.
func BenchmarkPerfGenericJoinGroupBy(b *testing.B) {
	for _, n := range []int{100, 500} {
		b.Run(fmt.Sprintf("size=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			g1, g2 := genericSetTuples(n)
			names := NewNames("b")
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var mb frozen.MapBuilder[Value, any]
				accumulate := func(s Set, slotKey int) {
					for e := s.Enumerator(); e.MoveNext(); {
						value := e.Current()
						key := value.(Tuple).Project(names)
						entry, found := mb.Get(key)
						if !found {
							entry = [2]Set{None, None}
						}
						slots := entry.([2]Set)
						slots[slotKey] = slots[slotKey].With(value)
						mb.Put(key, slots)
					}
				}
				accumulate(g1, 0)
				accumulate(g2, 1)
				_ = mb.Finish()
			}
		})
	}
}

// BenchmarkPerfRelationJoinDirect measures Relation.Join directly.
func BenchmarkPerfRelationJoinDirect(b *testing.B) {
	for _, n := range []int{100, 500} {
		b.Run(fmt.Sprintf("size=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			r1, r2 := generateRelations(n)
			common := r1.attrs.intersect(r2.attrs)
			left := r1.attrs
			right := r2.attrs.minus(r1.attrs)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = r1.Join(r2, common, left, right)
			}
		})
	}
}

// BenchmarkPerfTupleProjectInLoop simulates repeated Project during GenericJoin.
func BenchmarkPerfTupleProjectInLoop(b *testing.B) {
	for _, n := range []int{100, 1000} {
		b.Run(fmt.Sprintf("size=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			s := perfMakeGenericSet(n, 4)
			projNames := NewNames("a0", "a1")
			tuples := make([]Tuple, 0, n)
			for e := s.Enumerator(); e.MoveNext(); {
				tuples = append(tuples, e.Current().(Tuple))
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for _, t := range tuples {
					_ = t.Project(projNames)
				}
			}
		})
	}
}

// BenchmarkPerfMergeTuples measures Merge with varying attr counts.
func BenchmarkPerfMergeTuples(b *testing.B) {
	for _, n := range []int{3, 5} {
		b.Run(fmt.Sprintf("attrs=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			attrs1 := make([]Attr, n)
			attrs2 := make([]Attr, n)
			attrs1[0] = NewAttr("common", NewNumber(1))
			attrs2[0] = NewAttr("common", NewNumber(1))
			for i := 1; i < n; i++ {
				attrs1[i] = NewAttr(fmt.Sprintf("left%d", i), NewNumber(float64(i)))
				attrs2[i] = NewAttr(fmt.Sprintf("right%d", i), NewNumber(float64(i)))
			}
			t1 := NewTuple(attrs1...)
			t2 := NewTuple(attrs2...)
			for i := 0; i < b.N; i++ {
				_ = Merge(t1, t2)
			}
		})
	}
}

// BenchmarkPerfValueHash measures hashing cost for different value types.
func BenchmarkPerfValueHash(b *testing.B) {
	b.Run("Number", func(b *testing.B) {
		b.ReportAllocs()
		v := NewNumber(42)
		for i := 0; i < b.N; i++ {
			_ = v.Hash(0)
		}
	})
	b.Run("Tuple/3", func(b *testing.B) {
		b.ReportAllocs()
		v := perfMakeTuple(3)
		for i := 0; i < b.N; i++ {
			_ = v.Hash(0)
		}
	})
	b.Run("Tuple/5", func(b *testing.B) {
		b.ReportAllocs()
		v := perfMakeTuple(5)
		for i := 0; i < b.N; i++ {
			_ = v.Hash(0)
		}
	})
}
