package bird

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client    *redis.Client
	keyPrefix string
}

func NewRedisCache(config CacheConfig) (*RedisCache, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisServer,
		Password: config.RedisPassword,
		DB:       config.RedisDb,
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	cache := &RedisCache{
		client: client,
	}

	return cache, nil
}

// Get retrievs a birdwatcher `Parsed` result from
// the redis cache.
func (self *RedisCache) Get(key string) (Parsed, error) {
	ctx := context.Background()
	key = self.keyPrefix + key //"B" + IPVersion + "_" + key
	data, err := self.client.Get(ctx, key).Result()
	if err != nil {
		return NilParse, err
	}

	parsed := Parsed{}
	err = json.Unmarshal([]byte(data), &parsed)

	ttl, err := parseCacheTTL(parsed["ttl"])
	if err != nil {
		return NilParse, fmt.Errorf("invalid TTL value for key: %s", key)
	}

	if ttl.Before(time.Now()) {
		return NilParse, err // TTL expired
	}

	return parsed, err // cache hit
}

// Set adds a birdwatcher `Parsed` result
// to the redis cache.
func (self *RedisCache) Set(key string, parsed Parsed, ttl int) error {
	switch {
	case ttl == 0:
		return nil // do not cache

	case ttl > 0:
		key = self.keyPrefix + key //TODO "B" + IPVersion + "_" + key
		payload, err := json.Marshal(parsed)
		if err != nil {
			return err
		}

		ctx := context.Background()
		_, err = self.client.Set(
			ctx, key, payload, time.Duration(ttl)*time.Minute).Result()
		return err

	default: // ttl negative - invalid
		return fmt.Errorf("negative TTL value for key: %s", key)
	}
}

func (self *RedisCache) Expire() int {
	log.Printf("Cannot expire entries in RedisCache backend, redis does this automatically")
	return 0
}

// Helperfunction to decode the cache ttl stored
// in the cache - which will most likely just be
// RFC3339 timestamp.
func parseCacheTTL(cacheTTL interface{}) (time.Time, error) {
	if cacheTTL == nil {
		// We preseve the nil value as a zero value
		return time.Time{}, nil
	}

	switch cacheTTL.(type) {
	case string:
		ttl, err := time.Parse(time.RFC3339, cacheTTL.(string))
		if err != nil {
			return time.Time{}, err
		}
		return ttl, nil
	case time.Time:
		return cacheTTL.(time.Time), nil
	}
	return time.Time{}, nil
}
