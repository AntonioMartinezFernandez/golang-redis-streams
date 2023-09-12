package redis_streams

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

var _ Publisher = (*RedisStreamsPublisher)(nil)

type Publisher interface {
	Publish(streamName string, event map[string]interface{}) error
}

type RedisStreamsPublisher struct {
	ctx    context.Context
	client *redis.Client
	logger *slog.Logger
}

func NewRedisStreamsPublisher(
	ctx context.Context,
	client *redis.Client,
	logger *slog.Logger,
) *RedisStreamsPublisher {
	return &RedisStreamsPublisher{
		ctx:    ctx,
		client: client,
		logger: logger,
	}
}

func (rse *RedisStreamsPublisher) Publish(streamName string, event map[string]interface{}) error {
	timestampId, err := rse.client.XAdd(rse.ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: event,
	}).Result()

	if err != nil {
		rse.logger.Error("error pusblishing redis streams event", "error", err)
	} else {
		rse.logger.Debug("redis streams event published succesfully", "event", event, "redis_timestamp_id", timestampId)
	}

	return err
}
