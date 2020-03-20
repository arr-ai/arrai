package syntax

import "github.com/arr-ai/arrai/rel"

// it is a simple cache component used by import behavior, and it can make cache code simple
type importCache struct {
	cache map[string]rel.Value
}

func newCache() importCache {
	return importCache{cache: make(map[string]rel.Value)}
}

func (service *importCache) exists(key string) (bool, rel.Value) {
	val := service.get(key)
	if val != nil {
		return true, val
	}
	return false, nil
}

func (service *importCache) get(key string) rel.Value {
	return service.cache[key]
}

func (service *importCache) add(key string, val rel.Value) {
	service.cache[key] = val
}
