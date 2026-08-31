package syntax

import "testing"

// Pattern-stress benchmarks for 🎯T15: cond arms trying and failing
// structural matches, tuple destructuring, and array patterns.

// Tuple destructuring per element, one arm matching.
const benchPatTupleCode = `(//seq.repeat(4000, [0]) => (a: .@, b: .@ % 7, c: $"x${.@ % 3}")) => cond . {
	(a: 0, ...): 0,
	(:a, b: 0, ...): a,
	(:a, :b, c: "x0"): a + b,
	(:a, :b, :c): a - b,
}`

// Arms that mostly fail before one matches: the discarded-match path.
const benchPatMissCode = `(//seq.repeat(4000, [0]) => (@: .@, @item: .@ % 11)) >> cond . {
	9993: 1,
	9994: 2,
	9995: 3,
	9996: 4,
	9997: 5,
	_: 0,
}`

// Array destructuring with binds and tails.
const benchPatArrayCode = `(//seq.repeat(3000, [0]) => (@: .@, @item: [.@, .@ % 5, .@ % 3])) >> cond . {
	[0, ...]: 0,
	[x, 0, ...t]: x + (t count),
	[x, y, z]: x + y + z,
	_: -1,
}`

func BenchmarkPatTuple(b *testing.B) { benchEvalLoop(b, benchPatTupleCode) }
func BenchmarkPatMiss(b *testing.B)  { benchEvalLoop(b, benchPatMissCode) }
func BenchmarkPatArray(b *testing.B) { benchEvalLoop(b, benchPatArrayCode) }
