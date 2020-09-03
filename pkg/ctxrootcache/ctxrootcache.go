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
