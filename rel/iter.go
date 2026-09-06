package rel

import "iter"

// All returns a single-pass iterator over a set's values.
//
// Measured against the MoveNext/Current enumerators (iter_bench_test.go,
// under testing.B.Loop so inlining artifacts don't flatter it): ranging over
// All(s) is ~6% faster — one indirect call per element instead of two
// interface calls, and no enumerator box — with element boxing unchanged.
// Adopt it for clarity and the modest win, with two rules:
//
//   - Never in a loop whose body writes captured locals: range-over-func
//     turns each written capture into a heap cell per loop entry, measured
//     as a net regression on small-set-heavy workloads (Union, DArrow).
//   - It is not the fast path. A concrete indexed loop over the arena is
//     5× faster than either pattern because the per-element tuple box is
//     the real cost and only fully-inlined concrete code avoids it. Hot
//     evaluator loops should stay concrete.
//
// Kinds without a concrete sequence fall back to their enumerator, which is
// never worse than calling it directly.
func All(s Set) iter.Seq[Value] {
	switch s := s.(type) {
	case Relation:
		return func(yield func(Value) bool) {
			for i := 0; i < s.rows.n; i++ {
				if !yield(s.tuple(s.rows.rowAt(i))) {
					return
				}
			}
		}
	case String:
		return func(yield func(Value) bool) {
			for i, n := 0, s.size(); i < n; i++ {
				if r := s.runeAt(i); r >= 0 && !yield(NewStringCharTuple(s.offset+i, r)) {
					return
				}
			}
		}
	case Array:
		return func(yield func(Value) bool) {
			for i, v := range s.values {
				if v != nil && !yield(NewArrayItemTuple(s.offset+i, v)) {
					return
				}
			}
		}
	case GenericSet:
		return func(yield func(Value) bool) {
			for i := s.set.Range(); i.Next(); {
				if !yield(i.Value()) {
					return
				}
			}
		}
	default:
		return func(yield func(Value) bool) {
			for e := s.Enumerator(); e.MoveNext(); {
				if !yield(e.Current()) {
					return
				}
			}
		}
	}
}
