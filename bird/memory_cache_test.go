package bird

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestMemoryCacheAccess(t *testing.T) {

	cache := NewMemoryCache(100)

	parsed := Parsed{
		"foo": 23,
		"bar": 42,
		"baz": true,
	}

	t.Log("Setting memory cache...")
	if err := cache.Set("testkey", parsed, 5); err != nil {
		t.Error(err)
	}

	t.Log("Fetching from memory cache...")
	parsed, err := cache.Get("testkey")
	if err != nil {
		t.Error(err)
	}

	cache.Expire()
	t.Log(parsed)
}

func TestMemoryCacheAccessKeyMissing(t *testing.T) {
	cache := NewMemoryCache(100)
	parsed, err := cache.Get("test_missing_key")
	if !IsSpecial(parsed) {
		t.Error(err)
	}
	t.Log("Cache error:", err)
	t.Log(parsed)
}

func TestMemoryCacheRoutes(t *testing.T) {
	f, err := openFile("routes_bird1_ipv4.sample")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	parsed := parseRoutes(f)
	_, ok := parsed["routes"].([]Parsed)
	if !ok {
		t.Fatal("Error getting routes")
	}

	cache := NewMemoryCache(100)

	if err := cache.Set("routes_protocol_test", parsed, 5); err != nil {
		t.Error(err)
	}

	parsed, err = cache.Get("routes_protocol_test")
	if err != nil {
		t.Error(err)
		return
	}
	routes, ok := parsed["routes"].([]Parsed)
	if !ok {
		t.Error("Error getting routes")
	}
	t.Log("Retrieved routes:", len(routes))
}

func TestMemoryCacheMaxEntries(t *testing.T) {
	cache := NewMemoryCache(2)

	parsed := Parsed{
		"foo": 23,
		"bar": 42,
	}

	// Set 3 entries
	if err := cache.Set("testkey1", parsed, 5); err != nil {
		t.Error(err)
	}
	if err := cache.Set("testkey2", parsed, 5); err != nil {
		t.Error(err)
	}
	if err := cache.Set("testkey3", parsed, 5); err != nil {
		t.Error(err)
	}

	// Check that the first entry is gone
	_, err := cache.Get("testkey1")
	if err == nil {
		t.Error("Expected error, got nil")
	}
	// Check that the second entry is still there
	value, err := cache.Get("testkey2")
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	if value["foo"] != 23 {
		t.Error("Expected 23, got", value["foo"])
	}
}

func TestMemoryCacheNoCache(t *testing.T) {
	cache := NewMemoryCache(2)

	parsed := Parsed{
		"foo": 23,
		"bar": 42,
	}

	// Set 2 entries
	if err := cache.Set("testkey1", parsed, 0); err != nil {
		t.Error(err)
	}
	if err := cache.Set("testkey2", parsed, 0); err != nil {
		t.Error(err)
	}

	// Check that the first entry is not in Cache
	_, err := cache.Get("testkey1")
	if err == nil {
		t.Error("Expected error, got nil")
	}
	// Check that the second entry is not in Cache
	_, err = cache.Get("testkey2")
	if err == nil {
		t.Error("Expected error, got nil")
	}
	// Check that the third entry is not in Cache (also never set)
	_, err = cache.Get("testkey3")
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if len(cache.m) != 0 || len(cache.a) != 0 {
		t.Errorf("Expected no cache entries or access times, got %d entries and %d access times", len(cache.m), len(cache.a))
	}
}

func TestMemoryCacheMissesDoNotTrackAccess(t *testing.T) {
	cache := NewMemoryCache(2)
	for i := 0; i < 100; i++ {
		if _, err := cache.Get(fmt.Sprintf("missing-%d", i)); err == nil {
			t.Fatal("Expected a cache miss")
		}
	}

	if count := cache.Expire(); count != 0 {
		t.Errorf("Expected no expired entries, got %d", count)
	}
	if len(cache.m) != 0 || len(cache.a) != 0 {
		t.Errorf("Cache misses left %d entries and %d access times", len(cache.m), len(cache.a))
	}
}

func TestMemoryCacheMaxEntriesAfterMiss(t *testing.T) {
	cache := NewMemoryCache(2)
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key-%d", i)
		if _, err := cache.Get(key); err == nil {
			t.Fatal("Expected a cache miss")
		}
		if err := cache.Set(key, Parsed{"value": i}, 5); err != nil {
			t.Fatal(err)
		}
		if len(cache.m) > 2 || len(cache.a) != len(cache.m) {
			t.Fatalf("Cache exceeded its capacity or retained access times: %d entries and %d access times", len(cache.m), len(cache.a))
		}
	}
}

func TestMemoryCacheRemovesInvalidEntriesOnAccess(t *testing.T) {
	tests := []struct {
		name string
		ttl  interface{}
	}{
		{"expired", time.Now().Add(-time.Minute)},
		{"invalid", "not a timestamp"},
		{"missing", nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := NewMemoryCache(2)
			parsed := Parsed{"value": 1}
			if err := cache.Set("key", parsed, 5); err != nil {
				t.Fatal(err)
			}
			parsed["ttl"] = test.ttl

			if _, err := cache.Get("key"); err == nil {
				t.Fatal("Expected an invalid or expired TTL error")
			}
			if len(cache.m) != 0 || len(cache.a) != 0 {
				t.Errorf("Invalid entry left %d entries and %d access times", len(cache.m), len(cache.a))
			}
		})
	}
}

func TestMemoryCacheNonCacheableWritesDoNotEvict(t *testing.T) {
	for _, ttl := range []int{0, -1} {
		t.Run(fmt.Sprintf("ttl-%d", ttl), func(t *testing.T) {
			cache := NewMemoryCache(1)
			if err := cache.Set("existing", Parsed{"value": 1}, 5); err != nil {
				t.Fatal(err)
			}

			err := cache.Set("other", Parsed{"value": 2}, ttl)
			if (err != nil) != (ttl < 0) {
				t.Fatalf("Unexpected error for TTL %d: %v", ttl, err)
			}
			if _, err := cache.Get("existing"); err != nil {
				t.Fatalf("Non-cacheable write evicted an existing entry: %v", err)
			}
			if len(cache.m) != 1 || len(cache.a) != 1 {
				t.Errorf("Expected one entry and access time, got %d entries and %d access times", len(cache.m), len(cache.a))
			}
		})
	}
}

func TestMemoryCacheUpdatesExistingEntryAtCapacity(t *testing.T) {
	cache := NewMemoryCache(2)
	if err := cache.Set("first", Parsed{"value": 1}, 5); err != nil {
		t.Fatal(err)
	}
	if err := cache.Set("second", Parsed{"value": 2}, 5); err != nil {
		t.Fatal(err)
	}
	if err := cache.Set("first", Parsed{"value": 3}, 5); err != nil {
		t.Fatal(err)
	}

	for key, expected := range map[string]int{"first": 3, "second": 2} {
		parsed, err := cache.Get(key)
		if err != nil || parsed["value"] != expected {
			t.Errorf("Unexpected value for %s: %v, error: %v", key, parsed, err)
		}
	}
	if len(cache.m) != 2 || len(cache.a) != 2 {
		t.Errorf("Expected two entries and access times, got %d entries and %d access times", len(cache.m), len(cache.a))
	}
}

func TestMemoryCacheLRUUsesSuccessfulAccess(t *testing.T) {
	cache := NewMemoryCache(2)
	if err := cache.Set("first", Parsed{"value": 1}, 5); err != nil {
		t.Fatal(err)
	}
	if err := cache.Set("second", Parsed{"value": 2}, 5); err != nil {
		t.Fatal(err)
	}
	cache.a["first"] = time.Now().Add(-2 * time.Minute)
	cache.a["second"] = time.Now().Add(-time.Minute)

	if _, err := cache.Get("first"); err != nil {
		t.Fatal(err)
	}
	if err := cache.Set("third", Parsed{"value": 3}, 5); err != nil {
		t.Fatal(err)
	}

	if _, err := cache.Get("second"); err == nil {
		t.Error("Expected the least recently used entry to be evicted")
	}
	if _, ok := cache.a["second"]; ok {
		t.Error("Evicted entry retained an access time")
	}
	for _, key := range []string{"first", "third"} {
		if _, err := cache.Get(key); err != nil {
			t.Errorf("Expected %s to remain cached: %v", key, err)
		}
	}
}

func TestMemoryCacheLRUEvictsEmptyKey(t *testing.T) {
	cache := NewMemoryCache(1)
	if err := cache.Set("", Parsed{"value": 1}, 5); err != nil {
		t.Fatal(err)
	}
	if err := cache.Set("other", Parsed{"value": 2}, 5); err != nil {
		t.Fatal(err)
	}
	if len(cache.m) != 1 || len(cache.a) != 1 {
		t.Errorf("Expected one entry and access time, got %d entries and %d access times", len(cache.m), len(cache.a))
	}
	if _, err := cache.Get(""); err == nil {
		t.Error("Expected the empty key to be evicted")
	}
}

func TestMemoryCacheExpireRemovesAccessTimes(t *testing.T) {
	cache := NewMemoryCache(3)
	for _, key := range []string{"expired", "invalid", "fresh"} {
		if err := cache.Set(key, Parsed{"value": key}, 5); err != nil {
			t.Fatal(err)
		}
	}
	cache.m["expired"]["ttl"] = time.Now().Add(-time.Minute)
	delete(cache.m["invalid"], "ttl")

	if count := cache.Expire(); count != 2 {
		t.Errorf("Expected two expired entries, got %d", count)
	}
	if len(cache.m) != 1 || len(cache.a) != 1 {
		t.Errorf("Expected one entry and access time, got %d entries and %d access times", len(cache.m), len(cache.a))
	}
	if _, err := cache.Get("fresh"); err != nil {
		t.Errorf("Expected the fresh entry to remain cached: %v", err)
	}
}

func TestMemoryCacheConcurrentAccounting(t *testing.T) {
	cache := NewMemoryCache(8)
	var wg sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				key := fmt.Sprintf("worker-%d-key-%d", worker, i)
				cache.Get(key)
				if err := cache.Set(key, Parsed{"value": i}, 5); err != nil {
					t.Error(err)
				}
				cache.Get(key)
				cache.Expire()
			}
		}(worker)
	}
	wg.Wait()

	if len(cache.m) > 8 || len(cache.a) != len(cache.m) {
		t.Errorf("Concurrent access left %d entries and %d access times", len(cache.m), len(cache.a))
	}
	for key := range cache.m {
		if _, ok := cache.a[key]; !ok {
			t.Errorf("Cached entry %s has no access time", key)
		}
	}
}
