package bird

import (
	"testing"
)

func Test_MemoryCacheAccess(t *testing.T) {

	cache, err := NewMemoryCache()

	parsed := Parsed{
		"foo": 23,
		"bar": 42,
		"baz": true,
	}

	t.Log("Setting memory cache...")
	err = cache.Set("testkey", parsed, 5)
	if err != nil {
		t.Error(err)
	}

	t.Log("Fetching from memory cache...")
	parsed, err = cache.Get("testkey")
	if err != nil {
		t.Error(err)
	}

	cache.Expire()
	t.Log(parsed)
}

func Test_MemoryCacheAccessKeyMissing(t *testing.T) {

	cache, err := NewMemoryCache()

	parsed, err := cache.Get("test_missing_key")
	if !IsSpecial(parsed) {
		t.Error(err)
	}
	t.Log("Cache error:", err)
	t.Log(parsed)
}

func Test_MemoryCacheRoutes(t *testing.T) {
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

	cache, err := NewMemoryCache()

	err = cache.Set("routes_protocol_test", parsed, 5)
	if err != nil {
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
