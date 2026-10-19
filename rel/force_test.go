package rel

import (
	"context"
	"testing"

	"github.com/arr-ai/wbnf/parser"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestForceMaterialisesPipelines(t *testing.T) {
	t.Parallel()
	inc := func(_, v Value) (Value, error) { return NewNumber(v.(Number).Float64() + 1), nil }
	base := NewArray(NewNumber(1), NewNumber(2)).(Array)
	p := newSeqPipeline(base, inc)

	f, err := Force(p)
	require.NoError(t, err)
	a, ok := f.(Array)
	require.True(t, ok, "got %T", f)
	assert.True(t, a.Equal(NewArray(NewNumber(2), NewNumber(3))))

	tup := NewTuple(Attr{"xs", p}, Attr{"n", NewNumber(1)})
	f, err = Force(tup)
	require.NoError(t, err)
	xs, _ := f.(Tuple).Get("xs")
	_, ok = xs.(Array)
	assert.True(t, ok, "nested: got %T", xs)
	assert.True(t, f.Equal(tup))

	outer := NewArray(p, NewNumber(9)).(Array)
	f, err = Force(outer)
	require.NoError(t, err)
	_, ok = f.(Array).Values()[0].(Array)
	assert.True(t, ok)
	assert.True(t, f.Equal(outer))
}

func TestForceLeavesConcreteValuesAlone(t *testing.T) {
	t.Parallel()
	tup := NewTuple(Attr{"a", NewArray(NewNumber(1))})
	f, err := Force(tup)
	require.NoError(t, err)
	assert.Same(t, tup.(*GenericTuple), f.(*GenericTuple))
}

func TestForceResult(t *testing.T) {
	t.Parallel()
	inc := func(_, v Value) (Value, error) { return NewNumber(v.(Number).Float64() + 1), nil }
	p := newSeqPipeline(NewArray(NewNumber(1)).(Array), inc)

	e, err := ForceResult(NewLiteralExpr(parser.Scanner{}, p))
	require.NoError(t, err)
	lit, ok := e.(LiteralExpr)
	require.True(t, ok, "got %T", e)
	_, ok = lit.Literal().(Array)
	assert.True(t, ok, "literal: got %T", lit.Literal())

	id := NewIdentExpr(parser.Scanner{}, "x")
	e, err = ForceResult(id)
	require.NoError(t, err)
	again, err := ForceResult(e)
	require.NoError(t, err)
	assert.Equal(t, e, again, "idempotent")
	assert.Equal(t, id, Unforced(e))
	v, err := e.Eval(context.Background(), EmptyScope.With("x", p))
	require.NoError(t, err)
	_, ok = v.(Array)
	assert.True(t, ok, "eval: got %T", v)
}
