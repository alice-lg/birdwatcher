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
	val, ok := c.m[key]
	c.a[key] = time.Now().UTC() // Update access
	c.Unlock()

	if !ok { // cache miss
		return NilParse, errors.New("Failed to retrive key '" + key + "' from MemoryCache.")
	}

	// Check if the TTL is still valid
	ttl, ok := val["ttl"].(time.Time)
	if !ok {
		return NilParse, errors.New("Invalid TTL value for key '" + key + "'")
	}

	if ttl.Before(time.Now()) {
		return val, errors.New("TTL expired for key '" + key + "'") // TTL expired
	}

	return val, nil // cache hit
}

// Set a key in the cache.
func (c *MemoryCache) Set(key string, val Parsed, ttl int) error {
	c.Lock()
	defer c.Unlock()

	// Check if the key exists, if not clear the oldest key if
	// the number of entries exceeds maxKeys.
	if _, ok := c.a[key]; !ok {
		if len(c.a) >= c.maxKeys {
			c.expireLRU()
		}
	}

	if ttl == 0 {
		return nil // do not cache
	}
	if ttl < 0 {
		return errors.New("Negative TTL value for key" + key)
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
// 		    should be acquired before calling this function.
func (c *MemoryCache) expireLRU() {
	oldestKey := ""
	oldestTime := time.Now().UTC()
	for key := range c.m {
		if c.a[key].Before(oldestTime) {
			oldestKey = key
			oldestTime = c.a[key]
		}
	}
	if oldestKey == "" {
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

	expiredKeys := []string{}
	for key := range c.m {
		ttl, ok := c.m[key]["ttl"].(time.Time)
		if !ok || ttl.Before(now) {
			expiredKeys = append(expiredKeys, key)
		}
	}

	for _, key := range expiredKeys {
		delete(c.m, key)
		delete(c.a, key)
	}

	return len(expiredKeys)
}
