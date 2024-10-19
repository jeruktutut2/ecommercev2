package utils

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisUtil interface {
	GetClient() *redis.Client
	Close()
}

type RedisUtilImplementation struct {
	Client *redis.Client
}

func NewRedisConnection() RedisUtil {
	println(time.Now().String()+" redis: connecting to ", os.Getenv("ECOMMERCEV2_REDIS_HOST")+":"+os.Getenv("ECOMMERCEV2_REDIS_PORT"))
	db, err := strconv.Atoi(os.Getenv("ECOMMERCEV2_REDIS_DATABASE"))
	if err != nil {
		log.Fatalln("error when converting db redis: " + err.Error())
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("ECOMMERCEV2_REDIS_HOST") + ":" + os.Getenv("ECOMMERCEV2_REDIS_PORT"),
		Password: "", // no password set
		DB:       db, // use default DB
	})
	ctx := context.Background()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalln("redis connection error:", err)
	}
	println(time.Now().String()+" redis: connected to", os.Getenv("ECOMMERCEV2_REDIS_HOST")+":"+os.Getenv("ECOMMERCEV2_REDIS_PORT"))
	return &RedisUtilImplementation{
		Client: rdb,
	}
}

func (util *RedisUtilImplementation) GetClient() *redis.Client {
	return util.Client
}

func (util *RedisUtilImplementation) Close() {
	err := util.Client.Close()
	if err != nil {
		panic("redis close connection error: " + err.Error())
	}
	println(time.Now().String(), "redis: closed properly")
}
