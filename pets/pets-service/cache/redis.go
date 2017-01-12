package cache

import (
	"gopkg.in/redis.v5"
	"time"
)

const defaultTTL = 20 * time.Second

type RedisCache struct {
	redisClient *redis.Client
}

func NewRedisCache(uri string) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     uri,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &RedisCache{
		redisClient: client,
	}, nil
}


func (c *RedisCache) Get(key string) (string, error) {
	val, err := c.redisClient.Get(key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (c *RedisCache) StoreWithDefaultTTL(key string, value string) error {
	return c.Store(key, value, defaultTTL)
}

func (c *RedisCache) Store(key string, value string, expiration time.Duration) error {
	return c.redisClient.Set(key, value, expiration).Err()
}
