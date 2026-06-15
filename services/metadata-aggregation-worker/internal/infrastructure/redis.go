package infrastructure

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(redisCfg *RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisCfg.Host, redisCfg.Port),
	})
}
