package rel

import (
	"context"
	"testing"

	"github.com/arr-ai/wbnf/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecideLetFanout(t *testing.T) {
	t.Parallel()
	assert.Equal(t, streamEdge, decideLetFanout(1))
	assert.Equal(t, materializeEdge, decideLetFanout(2))
	assert.Equal(t, "materialize", materializeEdge.String())
	assert.Equal(t, "stream", streamEdge.String())
}

func TestZoneMapNumericColumn(t *testing.T) {
	t.Parallel()
	r := newPositionalRelation(2, row(1, 9), row(5, 3), row(2, 7))
	z := r.zone(0)
	require.True(t, z.ok)
	assert.Equal(t, NewNumber(1), z.min)
	assert.Equal(t, NewNumber(5), z.max)
	assert.False(t, r.zone(3).ok)
}

func TestPlanCacheGuard(t *testing.T) {
	t.Parallel()
	c := &planCache{}
	live := true
	guard := func() bool { return live }
	c.put("k", guard, NewNumber(1))
	v, ok := c.get("k", guard)
	require.True(t, ok)
	assert.True(t, v.Equal(NewNumber(1)))
	live = false
	_, ok = c.get("k", func() bool { return true })
	assert.False(t, ok, "stale guard must not reuse the plan")
}

func TestPruneStackedProjects(t *testing.T) {
	t.Parallel()
	if !fastPaths {
		t.Skip("slowpath does not prune")
	}
	sc := *parser.NewScanner("")
	dot := func(name string) Expr { return NewDotExpr(sc, NewIdentExpr(sc, "."), name) }
	attr := func(name string) AttrExpr {
		a, err := NewAttrExpr(sc, name, dot(name))
		require.NoError(t, err)
		return a
	}
	innerFn := NewFunction(sc, IdentPattern("."), NewTupleExpr(sc, attr("a"), attr("b"), attr("c")))
	outerFn := NewFunction(sc, IdentPattern("."), NewTupleExpr(sc, attr("a")))
	base := mustRel(t,
		NewTuple(NewAttr("a", NewNumber(1)), NewAttr("b", NewNumber(2)), NewAttr("c", NewNumber(3))),
	)
	inner := NewDArrowExpr(sc, base, innerFn)
	outer := NewDArrowExpr(sc, inner, outerFn)
	d, ok := outer.(*DArrowExpr)
	require.True(t, ok)
	in, ok := d.lhs.(*DArrowExpr)
	require.True(t, ok)
	te, ok := in.fn.body.(*TupleExpr)
	require.True(t, ok)
	require.Equal(t, 1, len(te.attrs), "inner project should keep only a")
	assert.Equal(t, "a", te.attrs[0].name)
}

func eqDot(sc parser.Scanner, attr string, key Expr) CompareExpr {
	eq := func(a, b Value) (bool, error) { return a.Equal(b), nil }
	return NewCompareExpr(sc,
		[]Expr{NewDotExpr(sc, NewIdentExpr(sc, "."), attr), key},
		[]CompareFunc{eq},
		[]string{"="},
	)
}

func TestMatchEqAttrPredicatesAndTree(t *testing.T) {
	t.Parallel()
	sc := *parser.NewScanner("")
	one := ExprAsFunction(eqDot(sc, "a", NewNumber(1)))
	require.Len(t, matchEqAttrPredicates(one), 1)

	two := ExprAsFunction(NewAndExpr(sc, eqDot(sc, "a", NewNumber(1)), eqDot(sc, "b", NewNumber(2))))
	ps := matchEqAttrPredicates(two)
	require.Len(t, ps, 2)
	assert.Equal(t, "a", ps[0].attr)
	assert.Equal(t, "b", ps[1].attr)

	three := ExprAsFunction(NewAndExpr(sc,
		NewAndExpr(sc, eqDot(sc, "a", NewNumber(1)), eqDot(sc, "b", NewNumber(2))),
		eqDot(sc, "c", NewNumber(3)),
	))
	require.Len(t, matchEqAttrPredicates(three), 3)

	neq := func(a, b Value) (bool, error) { return !a.Equal(b), nil }
	mixed := NewCompareExpr(sc,
		[]Expr{NewDotExpr(sc, NewIdentExpr(sc, "."), "b"), NewNumber(0)},
		[]CompareFunc{neq},
		[]string{"!="},
	)
	assert.Nil(t, matchEqAttrPredicates(ExprAsFunction(NewAndExpr(sc, eqDot(sc, "a", NewNumber(1)), mixed))))
}

func TestWhereConjEqAttrDoesNotInflate(t *testing.T) {
	if !fastPaths {
		t.Skip("slowpath scans")
	}
	var n int
	tupleInflateHook = func() { n++ }
	t.Cleanup(func() { tupleInflateHook = nil })

	sc := *parser.NewScanner("")
	r := mustRel(t,
		NewTuple(NewAttr("a", NewNumber(1)), NewAttr("b", NewString([]rune("x")))),
		NewTuple(NewAttr("a", NewNumber(2)), NewAttr("b", NewString([]rune("y")))),
		NewTuple(NewAttr("a", NewNumber(2)), NewAttr("b", NewString([]rune("z")))),
	).(Relation)
	pred := NewAndExpr(sc, eqDot(sc, "a", NewNumber(2)), eqDot(sc, "b", NewString([]rune("y"))))
	v, err := NewWhereExpr(sc, r, pred).Eval(context.Background(), EmptyScope)
	require.NoError(t, err)
	assert.Equal(t, 1, v.(Set).Count())
	assert.Equal(t, 0, n, "conjunctive eq-attr where must not inflate rows")
}

func TestWhereConjEqAttrPlanCache(t *testing.T) {
	t.Parallel()
	if !fastPaths {
		t.Skip("slowpath scans")
	}
	sc := *parser.NewScanner("")
	r := mustRel(t,
		NewTuple(NewAttr("a", NewNumber(1)), NewAttr("b", NewNumber(10))),
		NewTuple(NewAttr("a", NewNumber(2)), NewAttr("b", NewNumber(20))),
		NewTuple(NewAttr("a", NewNumber(2)), NewAttr("b", NewNumber(21))),
	).(Relation)
	pred := NewAndExpr(sc, eqDot(sc, "a", NewNumber(2)), eqDot(sc, "b", NewNumber(20)))
	w := NewWhereExpr(sc, r, pred)
	v1, err := w.Eval(context.Background(), EmptyScope)
	require.NoError(t, err)
	assert.Equal(t, 1, v1.(Set).Count())
	hits := r.rows.planHitCount()
	v2, err := w.Eval(context.Background(), EmptyScope)
	require.NoError(t, err)
	assert.Equal(t, 1, v2.(Set).Count())
	assert.Equal(t, hits+1, r.rows.planHitCount(), "second conj where must reuse the fact-keyed plan")

	swapped := NewWhereExpr(sc, r, NewAndExpr(sc, eqDot(sc, "b", NewNumber(20)), eqDot(sc, "a", NewNumber(2))))
	_, err = swapped.Eval(context.Background(), EmptyScope)
	require.NoError(t, err)
	assert.Equal(t, hits+2, r.rows.planHitCount(), "conjunct order must share the fact key")
}

func TestEvalAtRowDoesNotInflate(t *testing.T) {
	if !fastPaths {
		t.Skip("slowpath inflates")
	}
	var n int
	tupleInflateHook = func() { n++ }
	t.Cleanup(func() { tupleInflateHook = nil })

	sc := *parser.NewScanner("")
	r := mustRel(t,
		NewTuple(NewAttr("a", NewNumber(1)), NewAttr("b", NewNumber(10))),
		NewTuple(NewAttr("a", NewNumber(2)), NewAttr("b", NewNumber(20))),
	).(Relation)
	fn := NewFunction(sc, IdentPattern("."), NewAddExpr(sc,
		NewDotExpr(sc, NewIdentExpr(sc, "."), "a"),
		NewDotExpr(sc, NewIdentExpr(sc, "."), "b"),
	))
	v, err := NewDArrowExpr(sc, r, fn).Eval(context.Background(), EmptyScope)
	require.NoError(t, err)
	assert.Equal(t, 2, v.(Set).Count())
	assert.Equal(t, 0, n, "column-only => must not inflate")
}

func TestNestOnStoreDoesNotInflate(t *testing.T) {
	if !fastPaths {
		t.Skip("slowpath reduces tuples")
	}
	var n int
	tupleInflateHook = func() { n++ }
	t.Cleanup(func() { tupleInflateHook = nil })

	r := mustRel(t,
		NewTuple(NewAttr("a", NewNumber(1)), NewAttr("b", NewNumber(10))),
		NewTuple(NewAttr("a", NewNumber(1)), NewAttr("b", NewNumber(11))),
		NewTuple(NewAttr("a", NewNumber(2)), NewAttr("b", NewNumber(20))),
	).(Relation)
	got := Nest(r, r.attrSet, NewNames("b"), "rows")
	assert.Equal(t, 2, got.(Set).Count())
	assert.Equal(t, 0, n, "nest must not inflate source rows")
}

func TestOrderByColumnKeysFromStore(t *testing.T) {
	if !fastPaths {
		t.Skip("slowpath enumerates")
	}
	var n int
	tupleInflateHook = func() { n++ }
	t.Cleanup(func() { tupleInflateHook = nil })

	r := mustRel(t,
		NewTuple(NewAttr("a", NewNumber(3))),
		NewTuple(NewAttr("a", NewNumber(1))),
		NewTuple(NewAttr("a", NewNumber(2))),
	).(Relation)
	values, ok := r.orderByColumn("a")
	require.True(t, ok)
	require.Len(t, values, 3)
	assert.Equal(t, NewNumber(1), values[0].(Tuple).MustGet("a"))
	assert.Equal(t, NewNumber(2), values[1].(Tuple).MustGet("a"))
	assert.Equal(t, NewNumber(3), values[2].(Tuple).MustGet("a"))
	assert.Equal(t, 3, n, "orderby boxes the result array only, not the keys")
}

func TestWhereColumnPredDoesNotInflate(t *testing.T) {
	if !fastPaths {
		t.Skip("slowpath scans")
	}
	var n int
	tupleInflateHook = func() { n++ }
	t.Cleanup(func() { tupleInflateHook = nil })

	sc := *parser.NewScanner("")
	r := mustRel(t,
		NewTuple(NewAttr("t", NewString([]rune("x"))), NewAttr("rest", NewNumber(1))),
		NewTuple(NewAttr("t", NewString([]rune("y"))), NewAttr("rest", NewNumber(0))),
	).(Relation)
	pred := NewDotExpr(sc, NewIdentExpr(sc, "."), "rest")
	v, err := NewWhereExpr(sc, r, pred).Eval(context.Background(), EmptyScope)
	require.NoError(t, err)
	assert.Equal(t, 1, v.(Set).Count())
	assert.Equal(t, 0, n, "column truthy where must not inflate rows")
}

func TestTuplePatternProjectDoesNotInflate(t *testing.T) {
	if !fastPaths {
		t.Skip("slowpath binds")
	}
	var n int
	tupleInflateHook = func() { n++ }
	t.Cleanup(func() { tupleInflateHook = nil })

	sc := *parser.NewScanner("")
	r := mustRel(t,
		NewTuple(NewAttr("a", NewNumber(1)), NewAttr("b", NewNumber(2))),
		NewTuple(NewAttr("a", NewNumber(3)), NewAttr("b", NewNumber(4))),
	).(Relation)
	tp, err := NewTuplePattern(
		NewTuplePatternAttr("a", NewFallbackPattern(IdentPattern("a"), nil)),
		NewTuplePatternAttr("b", NewFallbackPattern(IdentPattern("b"), nil)),
	)
	require.NoError(t, err)
	fn := NewFunction(sc, tp, NewTupleExpr(sc,
		mustAttr(t, sc, "x", NewIdentExpr(sc, "a")),
		mustAttr(t, sc, "y", NewIdentExpr(sc, "b")),
	))
	v, err := NewDArrowExpr(sc, r, fn).Eval(context.Background(), EmptyScope)
	require.NoError(t, err)
	assert.Equal(t, 2, v.(Set).Count())
	assert.Equal(t, 0, n, "tuple-pattern project must not inflate rows")
}

func TestProjectCanonicalDoesNotInflate(t *testing.T) {
	if !fastPaths {
		t.Skip("slowpath enumerates")
	}
	var n int
	tupleInflateHook = func() { n++ }
	t.Cleanup(func() { tupleInflateHook = nil })

	sc := *parser.NewScanner("")
	r := mustRel(t,
		NewTuple(NewAttr("k", NewString([]rune("a"))), NewAttr("v", NewNumber(1))),
		NewTuple(NewAttr("k", NewString([]rune("b"))), NewAttr("v", NewNumber(2))),
	).(Relation)
	fn := NewFunction(sc, IdentPattern("."), NewTupleExpr(sc,
		mustAttr(t, sc, "@", NewDotExpr(sc, NewIdentExpr(sc, "."), "k")),
		mustAttr(t, sc, "@value", NewDotExpr(sc, NewIdentExpr(sc, "."), "v")),
	))
	v, err := NewDArrowExpr(sc, r, fn).Eval(context.Background(), EmptyScope)
	require.NoError(t, err)
	assert.Equal(t, 2, v.(Set).Count())
	assert.Equal(t, 0, n, "dict-dots must not inflate source rows")
}

func mustAttr(t *testing.T, sc parser.Scanner, name string, e Expr) AttrExpr {
	t.Helper()
	a, err := NewAttrExpr(sc, name, e)
	require.NoError(t, err)
	return a
}

func TestProjectColumnDoesNotInflate(t *testing.T) {
	if !fastPaths {
		t.Skip("slowpath enumerates")
	}
	var n int
	tupleInflateHook = func() { n++ }
	t.Cleanup(func() { tupleInflateHook = nil })

	sc := *parser.NewScanner("")
	r := mustRel(t,
		NewTuple(NewAttr("a", NewNumber(1)), NewAttr("b", NewString([]rune("x")))),
		NewTuple(NewAttr("a", NewNumber(2)), NewAttr("b", NewString([]rune("y")))),
		NewTuple(NewAttr("a", NewNumber(2)), NewAttr("b", NewString([]rune("z")))),
	).(Relation)
	fn := NewFunction(sc, IdentPattern("."), NewDotExpr(sc, NewIdentExpr(sc, "."), "a"))
	v, err := NewDArrowExpr(sc, r, fn).Eval(context.Background(), EmptyScope)
	require.NoError(t, err)
	assert.Equal(t, 2, v.(Set).Count())
	assert.Equal(t, 0, n, "column extract must not inflate rows")
}

func TestWhereIndexPlanCacheThroughEval(t *testing.T) {
	t.Parallel()
	if !fastPaths {
		t.Skip("slowpath scans")
	}
	sc := *parser.NewScanner("")
	r := mustRel(t,
		NewTuple(NewAttr("k", NewNumber(1)), NewAttr("a", NewNumber(10))),
		NewTuple(NewAttr("k", NewNumber(2)), NewAttr("a", NewNumber(20))),
		NewTuple(NewAttr("k", NewNumber(1)), NewAttr("a", NewNumber(11))),
	).(Relation)
	eq := func(a, b Value) (bool, error) { return a.Equal(b), nil }
	pred := ExprAsFunction(NewCompareExpr(sc,
		[]Expr{NewDotExpr(sc, NewIdentExpr(sc, "."), "k"), NewNumber(1)},
		[]CompareFunc{eq},
		[]string{"="},
	))
	w := NewWhereExpr(sc, r, pred)
	v1, err := w.Eval(context.Background(), EmptyScope)
	require.NoError(t, err)
	assert.Equal(t, 2, v1.(Set).Count())
	hits := r.rows.planHitCount()
	v2, err := w.Eval(context.Background(), EmptyScope)
	require.NoError(t, err)
	assert.Equal(t, 2, v2.(Set).Count())
	assert.Equal(t, hits+1, r.rows.planHitCount(), "second where must reuse the fact-keyed plan")
}

func TestMaterializeFreezesDemandedIndex(t *testing.T) {
	t.Parallel()
	r := mustRel(t,
		NewTuple(NewAttr("k", NewNumber(1))),
		NewTuple(NewAttr("k", NewNumber(2))),
	).(Relation)
	_, err := materializeValue(r, []string{"k"})
	require.NoError(t, err)
	assert.NotNil(t, r.rows.cachedGroup(valueProjector{r.getAttrIndex("k")}),
		"breaker must freeze the demanded group index")
}

func TestCountIdentUsesAndFanout(t *testing.T) {
	t.Parallel()
	sc := *parser.NewScanner("")
	x := NewIdentExpr(sc, "x")
	body := NewAddExpr(sc, NewCountExpr(sc, x), NewCountExpr(sc, x))
	assert.Equal(t, 2, countIdentUses(body, "x"))
	assert.Equal(t, 0, countIdentUses(body, "y"))
	fn := NewFunction(sc, IdentPattern("x"), body).(*Function)
	assert.Equal(t, materializeEdge, fn.recordedFanout())
	once := NewFunction(sc, IdentPattern("x"), NewCountExpr(sc, x)).(*Function)
	assert.Equal(t, streamEdge, once.recordedFanout())
}
