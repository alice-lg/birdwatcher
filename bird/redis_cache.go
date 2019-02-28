package bird

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(config CacheConfig) (*RedisCache, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisServer,
		Password: config.RedisPassword,
		DB:       config.RedisDb,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	cache := &RedisCache{
		client: client,
	}

	return cache, nil
}

func (self *RedisCache) Get(key string) (Parsed, error) {
	data, err := self.client.Get(key).Result()
	if err != nil {
		return NilParse, err
	}

	parsed := Parsed{}
	err = json.Unmarshal([]byte(data), &parsed)

	return parsed, err
}

func (self *RedisCache) Set(key string, parsed Parsed) error {
	payload, err := json.Marshal(parsed)
	if err != nil {
		return err
	}

	_, err = self.client.Set(key, payload, time.Minute*5).Result()
	return err
}
