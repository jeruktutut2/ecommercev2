package mockhelpers

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
)

type RedisHelperMock struct {
	Mock mock.Mock
}

func (helper *RedisHelperMock) Set(client *redis.Client, ctx context.Context, key string, value interface{}, expiration time.Duration) (result string, err error) {
	arguments := helper.Mock.Called(client, ctx, key, value, expiration)
	return arguments.Get(0).(string), arguments.Error(1)
}

func (helper *RedisHelperMock) Del(client *redis.Client, ctx context.Context, key string) (result int64, err error) {
	arguments := helper.Mock.Called(client, ctx, key)
	return arguments.Get(0).(int64), arguments.Error(1)
}
