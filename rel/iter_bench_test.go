package rel

import (
	"iter"
	"testing"
)

// A fair three-way comparison for the iterator question (#741): the
// MoveNext/Current enumerator (two interface calls per element plus an
// allocated box), a range-over-func iter.Seq (one indirect call per
// element, no box), and a raw indexed loop over the arena as the floor.
// Same work per element: inflate the row to a tuple and read one attribute.

// allRows is the experimental iter.Seq form of Relation iteration.
func (r Relation) allRows() iter.Seq[Tuple] {
	return func(yield func(Tuple) bool) {
		for i := 0; i < r.rows.n; i++ {
			if !yield(r.tuple(r.rows.rowAt(i))) {
				return
			}
		}
	}
}

func benchConsume(t Tuple) int {
	if v, has := t.Get("qty"); has && v != nil {
		return 1
	}
	return 0
}

func BenchmarkIterEnumerator(b *testing.B) {
	r := benchRelation(b)
	b.ReportAllocs()
	for b.Loop() {
		n := 0
		for e := r.Enumerator(); e.MoveNext(); {
			n += benchConsume(e.Current().(Tuple))
		}
		if n != benchRelRows {
			b.Fatal(n)
		}
	}
}

func BenchmarkIterSeq(b *testing.B) {
	r := benchRelation(b)
	b.ReportAllocs()
	for b.Loop() {
		n := 0
		for t := range r.allRows() {
			n += benchConsume(t)
		}
		if n != benchRelRows {
			b.Fatal(n)
		}
	}
}

func BenchmarkIterIndexed(b *testing.B) {
	r := benchRelation(b)
	b.ReportAllocs()
	for b.Loop() {
		n := 0
		for j := 0; j < r.rows.n; j++ {
			n += benchConsume(r.tuple(r.rows.rowAt(j)))
		}
		if n != benchRelRows {
			b.Fatal(n)
		}
	}
}

// The same comparison through the Set interface, as a consumer outside the
// package would see it: the enumerator behind two interface dispatches
// versus an iter.Seq obtained once and then called directly.
func BenchmarkIterEnumeratorViaSet(b *testing.B) {
	var s Set = benchRelation(b)
	b.ReportAllocs()
	for b.Loop() {
		n := 0
		for e := s.Enumerator(); e.MoveNext(); {
			n += benchConsume(e.Current().(Tuple))
		}
		if n != benchRelRows {
			b.Fatal(n)
		}
	}
}

// tupleSeqer models Set gaining an All() method: the seq is obtained
// through an interface, so the compiler cannot see the yield's caller.
type tupleSeqer interface{ allRows() iter.Seq[Tuple] }

func BenchmarkIterSeqViaInterface(b *testing.B) {
	var s tupleSeqer = benchRelation(b)
	b.ReportAllocs()
	for b.Loop() {
		n := 0
		for t := range s.allRows() {
			n += benchConsume(t)
		}
		if n != benchRelRows {
			b.Fatal(n)
		}
	}
}
