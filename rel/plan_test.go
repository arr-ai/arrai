package rel

import (
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
