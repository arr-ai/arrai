package syntax

import "testing"

func TestDArrowExprEmpty(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{}`, `{} => .`)
}

func TestDArrowExprIdent(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{1, 2, 3}`, `{1,2,3} => .`)
}

func TestDArrowExprDouble(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{2, 4, 6}`, `{1,2,3} => \i i * 2`)
}

func TestDArrowExprIdentHoles(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t, `{1, , , 2}`, `{1,,,2} => .`)
}

// `set => (@: .i, @item: .v)` — the idiom for turning a set into an array —
// must produce a real Array, not a two-column Relation.
func TestDArrowExprProjectToArrayItemShape(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`['a']`,
		`{(index: 0, val: 'a')} => (@: .index, @item: .val)`,
	)
	AssertCodesEvalToSameValue(t,
		`['a', 'b']`,
		`{(index: 0, val: 'a'), (index: 1, val: 'b')} => (@: .index, @item: .val)`,
	)
}

// identDots into a reserved @-shape stays on the store (🎯T29.3).
func TestDArrowExprDictDots(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`{'a': 1, 'b': 2}`,
		`{(k: 'a', v: 1), (k: 'b', v: 2)} => (@: .k, @value: .v)`,
	)
	AssertCodesEvalToSameValue(t,
		`{'a': 1, 'b': 2}`,
		`{(k: 'a', v: 1), (k: 'b', v: 2)} => (@value: .v, @: .k)`,
	)
	AssertCodesEvalToSameValue(t,
		`{'a': {(x: 1), (x: 2)}, 'b': {(x: 3)}}`,
		`let r = {(appName: 'a', x: 1), (appName: 'a', x: 2), (appName: 'b', x: 3)}; `+
			`(r nest ~|appName|rows) => (@: .appName, @value: .rows)`,
	)
}

// A destructuring => pattern that doesn't match an element's shape re-runs
// the bind with explain=true (explainBind) to produce a detailed message,
// rather than the generic "pattern did not match" used internally to try
// candidate patterns cheaply.
// `rel => .attr` extracts the column as a set of values (🎯T29.2).
// Tuple-pattern => that is project/rename is identDots (🎯T29.4).
func TestDArrowExprTuplePatternProject(t *testing.T) {
	t.Parallel()

	AssertCodesEvalToSameValue(t,
		`{(x: 1, y: 2)}`,
		`{(a: 1, b: 2)} => \(:a, :b) (x: a, y: b)`,
	)
	AssertCodesEvalToSameValue(t,
		`{(x: 1, y: 2)}`,
		`{(a: 1, b: 2)} => \(:b, :a) (x: a, y: b)`,
	)
	AssertCodesEvalToSameValue(t,
		`{(name: 'n', value: 1)}`,
		`{(appAnnoName: 'n', appAnnoValue: 1, extra: 0)} => `+
			`\(:appAnnoName, :appAnnoValue, ...) (name: appAnnoName, value: appAnnoValue)`,
	)
	AssertCodesEvalToSameValue(t,
		`{(appName: 'a', appSrc: 's')}`,
		`{(appName: 'a', appSrc: 's')} => \(:appName, :appSrc) (:appName, :appSrc)`,
	)
	AssertCodeErrors(t, ``,
		`{(a: 1, b: 2, c: 3)} => \(:a, :b) (x: a, y: b)`)
}

// Column-only bodies eval against the row cursor (🎯T29.9).
func TestDArrowExprRowCursor(t *testing.T) {
	t.Parallel()
	const rel = `let r = {(a: 1, b: 10), (a: 2, b: 20)};`
	AssertCodesEvalToSameValue(t, `{11, 22}`, rel+`r => .a + .b`)
	AssertCodesEvalToSameValue(t, `{(s: 2), (s: 3)}`, rel+`r => (s: .a + 1)`)
	AssertCodesEvalToSameValue(t, `{(a: 1, b: 10, c: 2), (a: 2, b: 20, c: 3)}`,
		rel+`r => . +> (c: .a + 1)`)
	AssertCodesEvalToSameValue(t, `{(a: true, b: 10), (a: true, b: 20)}`,
		rel+`r => . +> (a: .a > 0)`)
	AssertCodesEvalToSameValue(t, `{(a: 2, b: 20)}`, rel+`r where .a + .b > 20`)
	AssertCodesEvalToSameValue(t, `{(a: 1, b: 10), (a: 2, b: 20)}`, rel+`r => .`)
}

func TestDArrowExprColumnExtract(t *testing.T) {
	t.Parallel()

	const rel = `let r = {(a: 1, b: 'x'), (a: 2, b: 'y'), (a: 2, b: 'z')};`
	AssertCodesEvalToSameValue(t, `{1, 2}`, rel+`r => .a`)
	AssertCodesEvalToSameValue(t, `{'x', 'y', 'z'}`, rel+`r => .b`)
	AssertCodesEvalToSameValue(t, `{1, 2}`, rel+`r => \row row.a`)
	AssertCodesEvalToSameValue(t, `{}`, `{} => .a`)
	AssertCodesEvalToSameValue(t, `{2}`,
		rel+`let k = 2; (r where .a = k) => .a`)

	// Reconstruct-shaped tag extract.
	AssertCodesEvalToSameValue(t, `{'t'}`,
		`{(stmtKey: 1, stmtTag: 't'), (stmtKey: 2, stmtTag: 't')} => .stmtTag`)

	AssertCodeErrors(t, ``, rel+`r => .c`)
}

func TestDArrowExprPatternMismatchExplains(t *testing.T) {
	t.Parallel()

	AssertCodeErrors(t, "couldn't find x in tuple (y: 1)", `{(y: 1)} => \(x: a) a`)
}
