package syntax

import "testing"

// Fold benchmarks for 🎯T10: accumulators built one step at a time through
// persistent structures. Relations and strings extend in place since the
// arena/appendBuf work; these measure what still copies per step.

// Dict accumulation: the groupByKey pattern model pipelines use.
const benchFoldDictCode = `let rec f = \n \acc cond {n = 0: acc, _: f(n - 1, acc +> {$"k${n % 500}": n})};
(f(3000, {}) count)`

// Generic-set accumulation.
const benchFoldSetCode = `let rec f = \n \acc cond {n = 0: acc, _: f(n - 1, acc | {$"v${n}"})};
(f(3000, {}) count)`

// Relation accumulation (the arena's in-place case, as the control).
const benchFoldRelCode = `let rec f = \n \acc cond {n = 0: acc, _: f(n - 1, acc | {(a: n, b: n % 7)})};
(f(3000, {}) count)`

// Array append via ++ (appendBuf case, control).
const benchFoldArrayCode = `let rec f = \n \acc cond {n = 0: acc, _: f(n - 1, acc ++ [n])};
(f(3000, []) count)`

func BenchmarkFoldDict(b *testing.B)  { benchEvalLoop(b, benchFoldDictCode) }
func BenchmarkFoldSet(b *testing.B)   { benchEvalLoop(b, benchFoldSetCode) }
func BenchmarkFoldRel(b *testing.B)   { benchEvalLoop(b, benchFoldRelCode) }
func BenchmarkFoldArray(b *testing.B) { benchEvalLoop(b, benchFoldArrayCode) }
