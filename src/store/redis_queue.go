package store

import (
	"chatgpt-plus-exts/utils"
	"context"
	"github.com/go-redis/redis/v8"
)

type RedisQueue struct {
	name   string
	client *redis.Client
	ctx    context.Context
}

type RedisMQs struct {
	WeChat     *RedisQueue
	MidJourney *RedisQueue
}

func NewRedisMQs(client *redis.Client) *RedisMQs {
	return &RedisMQs{
		WeChat:     &RedisQueue{name: "wechat_message_queue", client: client, ctx: context.Background()},
		MidJourney: &RedisQueue{name: "midjourney_message_queue", client: client, ctx: context.Background()},
	}
}

func (q *RedisQueue) Push(value interface{}) {
	q.client.RPush(q.ctx, q.name, utils.JsonEncode(value))
}

func (q *RedisQueue) Take(value interface{}) error {
	result, err := q.client.BLPop(q.ctx, 0, q.name).Result()
	if err != nil {
		return err
	}
	return utils.JsonDecode(result[1], value)
}
