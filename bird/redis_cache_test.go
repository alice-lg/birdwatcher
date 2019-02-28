package bird

import (
	"testing"
)

func Test_RedisCacheAccess(t *testing.T) {

	cache, err := NewRedisCache(CacheConfig{
		RedisServer: "localhost:6379",
	})

	if err != nil {
		t.Log("Redis server not available:", err)
		t.Log("Skipping redis tests.")
		return
	}

	parsed := Parsed{
		"foo": 23,
		"bar": 42,
		"baz": true,
	}

	t.Log("Setting redis cache...")
	err = cache.Set("testkey", parsed, 5)
	if err != nil {
		t.Error(err)
	}

	t.Log("Fetching from redis...")
	parsed, err = cache.Get("testkey")
	if err != nil {
		t.Error(err)
	}

	t.Log(parsed)
}

func Test_RedisCacheAccessKeyMissing(t *testing.T) {

	cache, err := NewRedisCache(CacheConfig{
		RedisServer: "localhost:6379",
	})

	if err != nil {
		t.Log("Redis server not available:", err)
		t.Log("Skipping redis tests.")
		return
	}

	parsed, err := cache.Get("test_missing_key")
	if err == nil {
		t.Error(err)
	}
	t.Log("Cache error:", err)
	t.Log(parsed)
}

func Test_RedisCacheRoutes(t *testing.T) {
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

	cache, err := NewRedisCache(CacheConfig{
		RedisServer: "localhost:6379",
	})

	if err != nil {
		t.Log("Redis server not available:", err)
		t.Log("Skipping redis tests.")
		return
	}

	err = cache.Set("routes_protocol_test", parsed, 5)
	if err != nil {
		t.Error(err)
	}

	parsed, err = cache.Get("routes_protocol_test")
	if err != nil {
		t.Error(err)
		return
	}

	routes, ok := parsed["routes"].([]interface{})
	if !ok {
		t.Error("Error getting routes")
	}

	t.Log("Retrieved routes:", len(routes))
}
