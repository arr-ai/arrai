package syntax

import (
	"sync"

	"github.com/arr-ai/arrai/rel"
)

// it is a simple cache component used by import behavior, and it can make cache code simple
type importCache struct {
	mutex sync.Mutex
	cond  *sync.Cond
	cache map[string]rel.Expr
}

func newCache() *importCache {
	c := &importCache{cache: map[string]rel.Expr{}}
	c.cond = sync.NewCond(&c.mutex)
	return c
}

// getOrAdd tries to get the value for key. If not present, and another
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
