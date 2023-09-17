package redis_streams

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

var _ Publisher = (*RedisStreamsPublisher)(nil)

type Publisher interface {
	Publish(streamName string, event StreamToPublish) error
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

func (rse *RedisStreamsPublisher) Publish(streamName string, event StreamToPublish) error {
	eventAsMap := event.AsMap()

	timestampId, err := rse.client.XAdd(rse.ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: eventAsMap,
	}).Result()

	if err != nil {
		rse.logger.Error("error pusblishing redis streams event", "error", err)
	} else {
		rse.logger.Debug("redis streams event published succesfully", "event", event, "redis_timestamp_id", timestampId)
	}

	return err
}
