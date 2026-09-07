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
