package redis_streams

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/redis/go-redis/v9"

	pkg_domain "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/domain"
)

var _ Publisher = (*RedisStreamsPublisher)(nil)

type Publisher interface {
	Publish(domainMessage pkg_domain.DomainMessage) error
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

func (rsp *RedisStreamsPublisher) Publish(domainMessage pkg_domain.DomainMessage) error {
	messageAsMap := domainMessage.AsMap()
	messageAsBytes, err := json.Marshal(messageAsMap)
	if err != nil {
		return err
	}
	messageAsString := string(messageAsBytes)

	streamName := domainMessage.GetType()
	streamValues := map[string]interface{}{
		"type": domainMessage.GetType(),
		"data": messageAsString,
	}

	streamMessageId, err := rsp.client.XAdd(rsp.ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: streamValues,
	}).Result()

	if err != nil {
		rsp.logger.Error("error pusblishing redis streams message", "error", err)
	} else {
		rsp.logger.Debug("redis streams message published succesfully", "message", domainMessage, "stream_message_id", streamMessageId)
	}

	return err
}
