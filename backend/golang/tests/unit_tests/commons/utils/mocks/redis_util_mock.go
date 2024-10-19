package mockutils

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
)

type RedisUtilMock struct {
	Mock mock.Mock
}

func (util *RedisUtilMock) GetClient() *redis.Client {
	arguments := util.Mock.Called()
	return arguments.Get(0).(*redis.Client)
}

func (util *RedisUtilMock) Close() {
	arguments := util.Mock.Called()
	fmt.Println(arguments)
}
