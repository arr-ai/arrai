package ctxrootcache

import (
	"context"
	"errors"
	"sync"
)

type rootCacheKeyType int

const rootCacheKey rootCacheKeyType = iota

var errNoRootCache = errors.New("root cache not in context")

// WithRootCache adds a module root cache in the context. This is meant to avoid
// global sync.Map which breaks with parallelized tests.
func WithRootCache(ctx context.Context) context.Context {
	return context.WithValue(ctx, rootCacheKey, &sync.Map{})
}

// StoreRoot stores the path that leads to the modulePath to the map.
func StoreRoot(ctx context.Context, path, modulePath string) error {
	cache := rootCacheFrom(ctx)
	if cache == nil {
		return errNoRootCache
	}
	cache.Store(path, modulePath)
	return nil
}

// LoadRoot loads the modulePath that the provided path will lead to.
func LoadRoot(ctx context.Context, path string) (string, bool, error) {
	cache := rootCacheFrom(ctx)
	if cache == nil {
		return "", false, errNoRootCache
	}
	val, exists := cache.Load(path)
	if exists {
		return val.(string), exists, nil
	}
	return "", false, nil
}

func rootCacheFrom(ctx context.Context) *sync.Map {
	m := ctx.Value(rootCacheKey)
	if m == nil {
		return nil
	}
	return m.(*sync.Map)
}

type primaryRootKeyType int

const primaryRootKey primaryRootKeyType = iota

// primaryRoot memoizes the module root used to resolve go.mod-pinned module
// imports for an entire run. Go's module graph has exactly one resolved
// version per module (via minimal version selection) rooted at the main
// module's go.mod; consulting a transitive dependency's own go.mod instead
// (found by walking up from the source file that contains the import,
// which for a dependency is some read-only directory in the module cache)
// would both disagree with that resolution and fail outright, since `go
// list -mod=mod` needs to write go.sum and the module cache is read-only.
type primaryRoot struct {
	mu   sync.Mutex
	root string
	done bool
}

// WithPrimaryRootCache adds an empty, not-yet-resolved primary module root
// to the context.
func WithPrimaryRootCache(ctx context.Context) context.Context {
	return context.WithValue(ctx, primaryRootKey, &primaryRoot{})
}

// PrimaryRoot returns the primary module root for ctx, computing it via
// compute() on the first call and reusing that result (even "not found",
// i.e. computed == "") for every later call, regardless of what compute
// would return this time. If ctx has no primary root cache (WithPrimaryRootCache
// was never called), compute() runs every time.
func PrimaryRoot(ctx context.Context, compute func() string) string {
	p, ok := ctx.Value(primaryRootKey).(*primaryRoot)
	if !ok {
		return compute()
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.done {
		p.root = compute()
		p.done = true
	}
	return p.root
}
