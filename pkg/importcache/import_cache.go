package importcache

import (
	"context"
	"sync"

	"github.com/arr-ai/arrai/rel"
)

type importCacheKeyType int

const importCacheKey importCacheKeyType = iota

// it is a simple cache component used by import behavior, and it can make cache code simple
type importCache struct {
	mutex sync.Mutex
	cond  *sync.Cond
	cache map[string]rel.Expr
}

// WithNewImportCache adds an empty cache to the context.
func WithNewImportCache(ctx context.Context) context.Context {
	return context.WithValue(ctx, importCacheKey, newImportCache())
}

// HasImportCacheFrom returns true if there is an import cache in the context.
func HasImportCacheFrom(ctx context.Context) bool {
	return fromCache(ctx) != nil
}

// GetOrAddFromCache gets the cached expression based on the given filepath from the cache in context.
// If the file is not in the cache it will simply store the filepath and the expression in the cache.
// The function panics when cache is not in context.
func GetOrAddFromCache(ctx context.Context, key string, add func() (rel.Expr, error)) (rel.Expr, error) {
	if service := fromCache(ctx); service != nil {
		return service.getOrAdd(key, add)
	}
	panic("GetOrAddFromCache: cache not in context")
}

func newImportCache() *importCache {
	c := &importCache{cache: map[string]rel.Expr{}}
	c.cond = sync.NewCond(&c.mutex)
	return c
}

func fromCache(ctx context.Context) *importCache {
	if c := ctx.Value(importCacheKey); c != nil {
		if cache, is := c.(*importCache); is {
			return cache
		}
	}
	return nil
}

// gerOrAdd tries to get the value for key. If not present, and another
// goroutine is currently computing a value for key, this goroutine will wait
// till it's ready.
func (service *importCache) getOrAdd(key string, add func() (rel.Expr, error)) (rel.Expr, error) {
	adding := false

	service.mutex.Lock()
	defer func() {
		if adding {
			// If panicked trying to add, remove the key from the cache so
			// someone else can have a go.
			service.mutex.Lock()
			delete(service.cache, key)
		}
		service.mutex.Unlock()
	}()

	for {
		if val, has := service.cache[key]; has {
			if val != nil {
				return val, nil
			}
			// Another goroutine is adding an entry.
			service.cond.Wait()
		} else {
			break
		}
	}

	// Indicate that we'll add it.
	service.cache[key] = nil

	// Free the lock while we work.
	service.mutex.Unlock()
	adding = true
	val, err := add()
	if err != nil {
		return nil, err
	}
	adding = false
	service.mutex.Lock()

	if val != nil {
		service.cache[key] = val
	} else {
		delete(service.cache, key)
	}
	service.cond.Broadcast()
	return val, nil
}
