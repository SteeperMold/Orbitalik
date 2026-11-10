package infrastructure

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(redisCfg *RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisCfg.Host, redisCfg.Port),
	})

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal(err)
	}

	return rdb
}
