package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// RedisCache manages interactions with Redis
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache initializes a new Redis cache
func NewRedisCache(addr string, password string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisCache{client: client}, nil
}

// GetFromCache attempts to retrieve data from Redis
func (r *RedisCache) GetFromCache(key string) ([]byte, error) {
	return r.client.Get(ctx, key).Bytes()
}

// SetToCache stores data in Redis with an expiration time
func (r *RedisCache) SetToCache(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, expiration).Err()
}
