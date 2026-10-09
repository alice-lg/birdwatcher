package bird

import (
	"errors"
	"sync"
	"time"
)

// MemoryCache is a simple in-memory cache for parsed BIRD output.
// Limiting the number of cached results is using a simple LRU algorithm.
type MemoryCache struct {
	sync.Mutex
	m map[string]Parsed    // Cached data
	a map[string]time.Time // Access times

	maxKeys int // Maximum number of keys to cache
}

// NewMemoryCache creates a new MemoryCache with a maximum number of keys.
func NewMemoryCache(maxKeys int) *MemoryCache {
	var cache *MemoryCache
	cache = &MemoryCache{
		m: make(map[string]Parsed),
		a: make(map[string]time.Time),

		maxKeys: maxKeys,
	}
	return cache
}

// Get a key from the cache.
func (c *MemoryCache) Get(key string) (Parsed, error) {
	c.Lock()
	defer c.Unlock()

	val, ok := c.m[key]

	if !ok { // cache miss
		return NilParse, errors.New("Failed to retrive key '" + key + "' from MemoryCache.")
	}

	// Check if the TTL is still valid
	ttl, ok := val["ttl"].(time.Time)
	if !ok {
		delete(c.m, key)
		delete(c.a, key)
		return NilParse, errors.New("Invalid TTL value for key '" + key + "'")
	}

	now := time.Now().UTC()
	if ttl.Before(now) {
		delete(c.m, key)
		delete(c.a, key)
		return val, errors.New("TTL expired for key '" + key + "'") // TTL expired
	}

	c.a[key] = now  // Update access only for a valid cache hit
	return val, nil // cache hit
}

// Set a key in the cache.
func (c *MemoryCache) Set(key string, val Parsed, ttl int) error {
	if ttl == 0 {
		return nil // do not cache
	}
	if ttl < 0 {
		return errors.New("Negative TTL value for key" + key)
	}

	c.Lock()
	defer c.Unlock()

	// Check if the key exists, if not clear the oldest key if
	// the number of entries exceeds maxKeys.
	if _, ok := c.m[key]; !ok {
		if len(c.m) >= c.maxKeys {
			c.expireLRU()
		}
	}

	cachedAt := time.Now().UTC()
	cacheTTL := cachedAt.Add(time.Duration(ttl) * time.Minute)

	// This is not a really ... clean way of doing this.
	val["ttl"] = cacheTTL
	val["cached_at"] = cachedAt

	c.m[key] = val
	c.a[key] = cachedAt

	return nil
}

// Expire oldest key in cache.
// WARNING: this is not thread safe and a mutex
//
//	should be acquired before calling this function.
func (c *MemoryCache) expireLRU() {
	oldestKey := ""
	var oldestTime time.Time
	found := false
	for key := range c.m {
		if !found || c.a[key].Before(oldestTime) {
			oldestKey = key
			oldestTime = c.a[key]
			found = true
		}
	}
	if !found {
		return // Nothing to do here.
	}
	delete(c.m, oldestKey)
	delete(c.a, oldestKey)
}

// Expire all keys in cache that are older than the
// TTL value.
func (c *MemoryCache) Expire() int {
	c.Lock()
	defer c.Unlock()

	now := time.Now().UTC()

	count := 0
	for key := range c.m {
		ttl, ok := c.m[key]["ttl"].(time.Time)
		if !ok || ttl.Before(now) {
			delete(c.m, key)
			delete(c.a, key)
			count++
		}
	}

	return count
}
