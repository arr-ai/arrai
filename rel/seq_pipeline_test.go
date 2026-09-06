package rel

import (
	"context"
	"testing"

	"github.com/arr-ai/wbnf/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSeqPipelineCountDoesNotForce(t *testing.T) {
	t.Parallel()
	base := NewOffsetArray(0, NewNumber(1), NewNumber(2), NewNumber(3)).(Array)
	calls := 0
	p := newSeqPipeline(base, func(_, v Value) (Value, error) {
		calls++
		return v, nil
	})
	assert.Equal(t, 3, p.Count())
	assert.Equal(t, 0, calls, "Count must not run the map")
}

func TestSeqPipelineNStageMaterialisesOnce(t *testing.T) {
	t.Parallel()
	base := NewOffsetArray(0, NewNumber(1), NewNumber(2)).(Array)
	p := newSeqPipeline(base, func(_, v Value) (Value, error) {
		return NewNumber(v.(Number).Float64() + 1), nil
	})
	p = p.then(func(_, v Value) (Value, error) {
		return NewNumber(v.(Number).Float64() * 2), nil
	})
	got, err := p.force()
	require.NoError(t, err)
	want := NewArray(NewNumber(4), NewNumber(6))
	assert.True(t, got.Equal(want), "%s vs %s", got, want)
	got2, err := p.force()
	require.NoError(t, err)
	assert.Equal(t, got.Hash128(), got2.Hash128())
}

func TestSeqArrowExprReturnsPipeline(t *testing.T) {
	if !fastPaths {
		t.Skip("slowpath materialises each stage")
	}
	t.Parallel()
	scanner := *parser.NewScanner("")
	arr := NewArray(NewNumber(1), NewNumber(2), NewNumber(3))
	body := NewAddExpr(scanner, NewIdentExpr(scanner, "."), NewNumber(1))
	e := NewSeqArrowExpr(false)(scanner, arr, body)
	v, err := e.Eval(context.Background(), EmptyScope)
	require.NoError(t, err)
	p, ok := v.(seqPipeline)
	require.True(t, ok, "got %T", v)
	assert.Equal(t, 3, p.Count())
	forced, err := p.force()
	require.NoError(t, err)
	assert.True(t, forced.Equal(NewArray(NewNumber(2), NewNumber(3), NewNumber(4))))
}

func TestDictPipelineCountDoesNotForce(t *testing.T) {
	t.Parallel()
	base := MustNewDict(false,
		NewDictEntryTuple(NewString([]rune("a")), NewNumber(1)),
		NewDictEntryTuple(NewString([]rune("b")), NewNumber(2)),
	).(Dict)
	calls := 0
	p := newDictPipeline(base, func(_, v Value) (Value, error) {
		calls++
		return v, nil
	})
	assert.Equal(t, 2, p.Count())
	assert.Equal(t, 0, calls, "Count must not run the map")
}

func TestDictPipelineNStageMaterialisesOnce(t *testing.T) {
	t.Parallel()
	base := MustNewDict(false,
		NewDictEntryTuple(NewString([]rune("a")), NewNumber(1)),
	).(Dict)
	p := newDictPipeline(base, func(_, v Value) (Value, error) {
		return NewNumber(v.(Number).Float64() + 1), nil
	})
	p = p.then(func(_, v Value) (Value, error) {
		return NewNumber(v.(Number).Float64() * 2), nil
	})
	got, err := p.force()
	require.NoError(t, err)
	want := MustNewDict(false, NewDictEntryTuple(NewString([]rune("a")), NewNumber(4)))
	assert.True(t, got.Equal(want), "%s vs %s", got, want)
}

func TestSeqPipelineForceError(t *testing.T) {
	t.Parallel()
	base := NewOffsetArray(0, NewNumber(1)).(Array)
	p := newSeqPipeline(base, func(_, _ Value) (Value, error) {
		return nil, assert.AnError
	})
	assert.Equal(t, 1, p.Count())
	_, err := p.force()
	require.Error(t, err)
}
