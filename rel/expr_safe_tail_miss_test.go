package rel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSafeTailGetMissDoesNotBuildError(t *testing.T) {
	ctx := context.Background()
	tuple := NewTuple(NewAttr("a", NewNumber(1)), NewAttr("c", NewNumber(3)))

	safe := NewSafeTailGet(true, "b")
	v, err := safe.apply(ctx, tuple, EmptyScope)
	require.NoError(t, err)
	assert.Nil(t, v)

	hit, err := NewSafeTailGet(true, "a").apply(ctx, tuple, EmptyScope)
	require.NoError(t, err)
	assert.True(t, NewNumber(1).Equal(hit))

	// An unsafe miss still reports which attribute is missing and what is available.
	unsafe := NewSafeTailGet(false, "b")
	_, err = unsafe.apply(ctx, tuple, EmptyScope)
	require.Error(t, err)
	ctxErr, ok := err.(ContextErr)
	require.True(t, ok)
	_, isMissing := ctxErr.NextErr().(MissingAttrError)
	assert.True(t, isMissing)
	assert.Contains(t, err.Error(), `Missing attr "b"`)

	// A safe step still reports a failure that is not a missing attribute.
	_, err = safe.apply(ctx, NewNumber(1), EmptyScope)
	assert.Error(t, err)

	missAllocs := testing.AllocsPerRun(100, func() {
		if _, err := safe.apply(ctx, tuple, EmptyScope); err != nil {
			t.Fatal(err)
		}
	})
	errAllocs := testing.AllocsPerRun(100, func() {
		if _, err := unsafe.apply(ctx, tuple, EmptyScope); err == nil {
			t.Fatal("expected a missing-attr error")
		}
	})
	assert.Less(t, missAllocs, errAllocs)
}
