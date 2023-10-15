package redis_streams

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	pkg_domain "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/domain"
)

var _ Publisher = (*RedisStreamsPublisher)(nil)

type Publisher interface {
	Start() error
	Stop() error
	Publish(domainMessage pkg_domain.DomainMessage) error
	PublishDelayed(domainMessage pkg_domain.DomainMessage, delaySeconds int) error
}

type RedisStreamsPublisher struct {
	ctx       context.Context
	wg        *sync.WaitGroup
	client    *redis.Client
	logger    *slog.Logger
	isStarted bool
}

func NewRedisStreamsPublisher(
	ctx context.Context,
	wg *sync.WaitGroup,
	client *redis.Client,
	logger *slog.Logger,
) *RedisStreamsPublisher {
	return &RedisStreamsPublisher{
		ctx:       ctx,
		wg:        wg,
		client:    client,
		logger:    logger,
		isStarted: false,
	}
}

func (rsp *RedisStreamsPublisher) Start() error {
	rsp.logger.Info("starting redis streams publisher...")
	if rsp.isStarted {
		return errors.New("redis streams publisher is already started")
	}

	go rsp.initDelayedMessagesProcessor()

	rsp.isStarted = true
	return nil
}

func (rsp *RedisStreamsPublisher) Stop() error {
	rsp.logger.Info("stopping redis streams publisher...")
	if !rsp.isStarted {
		return errors.New("redis streams publisher is already stopped")
	}
	rsp.isStarted = false
	return nil
}

func (rsp *RedisStreamsPublisher) Publish(domainMessage pkg_domain.DomainMessage) error {
	if !rsp.isStarted {
		return errors.New("redis streams publisher is not started")
	}

	messageAsMap := domainMessage.AsMap()
	messageAsBytes, err := json.Marshal(messageAsMap)
	if err != nil {
		return err
	}

	messageAsString := string(messageAsBytes)
	messageType := domainMessage.GetType()

	rsp.publishMessageToStream(messageAsString, messageType)

	if err != nil {
		rsp.logger.Error("error pusblishing redis streams message", "error", err)
	}

	return err
}

func (rsp *RedisStreamsPublisher) PublishDelayed(domainMessage pkg_domain.DomainMessage, delaySeconds int) error {
	if !rsp.isStarted {
		return errors.New("redis streams publisher is not started")
	}

	messageAsMap := domainMessage.AsMap()
	messageAsBytes, err := json.Marshal(messageAsMap)
	if err != nil {
		return err
	}
	messageAsString := string(messageAsBytes)
	messageType := domainMessage.GetType()
	now := time.Now()
	timeToBeProcessed := now.Add(time.Second * time.Duration(delaySeconds)).UnixMilli()

	delayedMessage := NewDelayedMessage(messageAsString, messageType)
	delayedMessageAsString := delayedMessage.ToString()

	members := []redis.Z{{Score: float64(timeToBeProcessed), Member: delayedMessageAsString}}

	return rsp.client.ZAdd(rsp.ctx, DELAYED_MESSAGES_SORTED_SET_KEY, members...).Err()
}

func (rsp *RedisStreamsPublisher) initDelayedMessagesProcessor() {
	ticker := time.NewTicker(time.Second * time.Duration(1))

	for range ticker.C {
		if !rsp.isStarted {
			return
		}

		// Set processing delayed message flag as running
		ok := rsp.setProcessingDelayedMessageAsRunning()
		if !ok {
			rsp.logger.Info("another instance is processing delayed messages")
			return
		}

		// Retrieve elements with scores equal or less than the current timestamp
		currentTimestamp := time.Now().UnixMilli()

		scoresRange := redis.ZRangeBy{
			Min: "0",
			Max: fmt.Sprint(currentTimestamp),
		}

		keys, err := rsp.client.ZRangeByScore(rsp.ctx, DELAYED_MESSAGES_SORTED_SET_KEY, &scoresRange).Result()
		if err != nil {
			rsp.logger.Error("error retrieving delayed messages from redis", "error", err)
		}

		// Process 1 by 1 every element retrieved
		for _, delayedMessageAsString := range keys {
			err := rsp.processDelayedMessage(delayedMessageAsString)
			if err != nil {
				rsp.logger.Error("error processing delayed message from redis", "message", delayedMessageAsString, "error", err)
			}
		}

		// Remove processed elements
		rsp.client.ZRemRangeByScore(rsp.ctx, DELAYED_MESSAGES_SORTED_SET_KEY, "0", fmt.Sprint(currentTimestamp))

		// Set processing delayed message flag as stopped
		rsp.setProcessingDelayedMessageAsStopped()
	}
}

func (rsp *RedisStreamsPublisher) processDelayedMessage(delayedMessageAsString string) error {
	delayedMessage, err := NewDelayedMessageFromString(delayedMessageAsString)
	if err != nil {
		return err
	}

	return rsp.publishMessageToStream(delayedMessage.Message, delayedMessage.MessageType)
}

func (rsp *RedisStreamsPublisher) publishMessageToStream(messageAsString string, messageType string) error {
	streamName := StreamName(messageType)
	streamValues := map[string]interface{}{
		"type": messageType,
		"data": messageAsString,
	}

	streamMessageId, err := rsp.client.XAdd(rsp.ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: streamValues,
	}).Result()
	if err != nil {
		return err
	}

	rsp.logger.Debug("redis streams message published succesfully", "type", messageType, "message", messageAsString, "stream_message_id", streamMessageId)
	return nil
}

func (rsp *RedisStreamsPublisher) setProcessingDelayedMessageAsRunning() (ok bool) {
	runningFlagValue := "running"
	ttlDuration := time.Second * 120

	flagValue := rsp.client.Get(rsp.ctx, DELAYED_MESSAGES_CONSUMER_RUNNING_FLAG_KEY).Val()

	if flagValue == runningFlagValue {
		return false
	}

	rsp.client.Set(rsp.ctx, DELAYED_MESSAGES_CONSUMER_RUNNING_FLAG_KEY, runningFlagValue, ttlDuration)
	return true
}

func (rsp *RedisStreamsPublisher) setProcessingDelayedMessageAsStopped() {
	stoppedFlagValue := "stopped"
	ttlDuration := time.Second * 120

	rsp.client.Set(rsp.ctx, DELAYED_MESSAGES_CONSUMER_RUNNING_FLAG_KEY, stoppedFlagValue, ttlDuration)
}
