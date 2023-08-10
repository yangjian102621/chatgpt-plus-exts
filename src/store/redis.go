package store

import (
	"chatgpt-plus-exts/core"
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(config *core.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.RedisConfig.Address,
		Password: config.RedisConfig.Password,
		DB:       config.RedisConfig.Db,
	})
}
