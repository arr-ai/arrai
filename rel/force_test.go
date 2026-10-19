package rel

import (
	"testing"

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
