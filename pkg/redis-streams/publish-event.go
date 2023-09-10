package redis_streams

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisStreamEvent struct {
	ctx        context.Context
	client     *redis.Client
	streamName string
	event      map[string]interface{}
}

func NewRedisStreamEvent(
	ctx context.Context,
	client *redis.Client,
	streamName string,
	event map[string]interface{},
) *RedisStreamEvent {
	return &RedisStreamEvent{
		ctx:        ctx,
		client:     client,
		streamName: streamName,
		event:      event,
	}
}

func (rse *RedisStreamEvent) Publish() (string, error) {
	return rse.client.XAdd(rse.ctx, &redis.XAddArgs{
		Stream: rse.streamName,
		Values: rse.event,
	}).Result()
}
