package ctxrootcache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrimaryRootComputesOnceAndReuses(t *testing.T) {
	ctx := WithPrimaryRootCache(context.Background())

	calls := 0
	first := PrimaryRoot(ctx, func() string {
		calls++
		return "/outer/project"
	})
	require.Equal(t, "/outer/project", first)
	require.Equal(t, 1, calls)

	// A later call simulating a nested import (whose own compute() would
	// find a different, wrong root, e.g. a dependency's own go.mod) must
	// reuse the first result rather than recomputing.
	second := PrimaryRoot(ctx, func() string {
		calls++
		return "/some/nested/dependency/cache/dir"
	})
	require.Equal(t, "/outer/project", second)
	require.Equal(t, 1, calls, "PrimaryRoot must compute the root at most once per context")
}

func TestPrimaryRootReusesEmptyResult(t *testing.T) {
	ctx := WithPrimaryRootCache(context.Background())

	calls := 0
	first := PrimaryRoot(ctx, func() string {
		calls++
		return ""
	})
	require.Equal(t, "", first)

	second := PrimaryRoot(ctx, func() string {
		calls++
		return "/would/be/found/later"
	})
	require.Equal(t, "", second,
		`an initial "not found" result must stick, not be overwritten by a later lucky compute`)
	require.Equal(t, 1, calls)
}

func TestPrimaryRootWithoutCacheComputesEveryTime(t *testing.T) {
	ctx := context.Background()

	calls := 0
	compute := func() string {
		calls++
		return "root"
	}
	PrimaryRoot(ctx, compute)
	PrimaryRoot(ctx, compute)
	require.Equal(t, 2, calls)
}
