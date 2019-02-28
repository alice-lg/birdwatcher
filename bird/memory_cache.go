package bird

import (
	"errors"
	"sync"
	"time"
)

// Implementation of the MemoryCache backend.

type MemoryCache struct {
	sync.RWMutex
	m map[string]Parsed
}

func NewMemoryCache() (*MemoryCache, error) {
	var cache *MemoryCache
	cache = &MemoryCache{m: make(map[string]Parsed)}
	return cache, nil
}

func (c *MemoryCache) Get(key string) (Parsed, error) {
	c.RLock()
	val, ok := c.m[key]
	c.RUnlock()
	if !ok {
		return NilParse, errors.New("Could not retrive key" + key + "from MemoryCache.")
	}

	ttl, correct := val["ttl"].(time.Time)
	if !correct {
		return NilParse, errors.New("Invalid TTL value for key" + key)
	}

	if ttl.Before(time.Now()) {
		return NilParse, nil // TTL expired
	} else {
		return val, nil // cache hit
	}
}

func (c *MemoryCache) Set(key string, val Parsed, ttl int) error {
	switch {
	case ttl == 0:
		return nil // do not cache
	case ttl > 0:
		cachedAt := time.Now().UTC()
		cacheTtl := cachedAt.Add(time.Duration(ttl) * time.Minute)

		c.Lock()
		// This is not a really ... clean way of doing this.
		val["ttl"] = cacheTtl
		val["cached_at"] = cachedAt

		c.m[key] = val
		c.Unlock()
		return nil
	default: // ttl negative - invalid
		return errors.New("Negative TTL value for key" + key)
	}
}
