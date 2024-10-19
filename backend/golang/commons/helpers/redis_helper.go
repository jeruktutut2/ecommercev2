package helpers

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisHelper interface {
	Set(client *redis.Client, ctx context.Context, key string, value interface{}, expiration time.Duration) (result string, err error)
	Del(client *redis.Client, ctx context.Context, key string) (result int64, err error)
}

type RedisHelperImplementation struct {
}

func NewRedisHelper() RedisHelper {
	return &RedisHelperImplementation{}
}

func (helper *RedisHelperImplementation) Set(client *redis.Client, ctx context.Context, key string, value interface{}, expiration time.Duration) (result string, err error) {
	return client.Set(ctx, key, value, expiration).Result()
}

func (helper *RedisHelperImplementation) Del(client *redis.Client, ctx context.Context, key string) (result int64, err error) {
	return client.Del(ctx, key).Result()
}
