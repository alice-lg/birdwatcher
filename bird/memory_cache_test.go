package bird

import (
	"testing"
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
}
