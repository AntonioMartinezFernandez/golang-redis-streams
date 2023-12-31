package redis_streams

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	pkg_domain "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/domain"
	pkg_utils "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/utils"

	"github.com/redis/go-redis/v9"
)

const START string = ">"
const MAX_PROCESS_RETRIES int = 10

var _ Consumer = (*RedisStreamsConsumer)(nil)

type Consumer interface {
	Start()
	Stop()
	RegisterHandler(handler pkg_domain.DomainMessageHandler)
}

type RedisStreamsConsumer struct {
	ctx           context.Context
	wg            *sync.WaitGroup
	logger        *slog.Logger
	client        *redis.Client
	consumerName  string
	consumerGroup string
	streamName    string
	isStarted     bool
	handlers      []pkg_domain.DomainMessageHandler
}

func NewRedisStreamsConsumer(
	ctx context.Context,
	wg *sync.WaitGroup,
	logger *slog.Logger,
	client *redis.Client,
	consumerGroup string,
	messageType string,
) *RedisStreamsConsumer {
	return &RedisStreamsConsumer{
		ctx:           ctx,
		wg:            wg,
		logger:        logger,
		client:        client,
		consumerName:  pkg_utils.NewUlid(),
		consumerGroup: consumerGroup,
		streamName:    StreamName(messageType),
		isStarted:     false,
	}
}

func (rsc *RedisStreamsConsumer) CreateConsumerGroup() {
	rsc.logger.Debug("creating redis streams consumer group", "consumer group", rsc.consumerGroup, "stream name", rsc.streamName)
	if _, err := rsc.client.XGroupCreateMkStream(rsc.ctx, rsc.streamName, rsc.consumerGroup, "0").Result(); err != nil {
		if !strings.Contains(fmt.Sprint(err), "BUSYGROUP") {
			rsc.logger.Error("error creating redis streams consumer group", "consumer group", rsc.consumerGroup, "error", err)
			panic(err)
		}
	}
}

func (rsc *RedisStreamsConsumer) RegisterHandler(handler pkg_domain.DomainMessageHandler) {
	rsc.handlers = append(rsc.handlers, handler)
}

func (rsc *RedisStreamsConsumer) Stop() {
	if !rsc.isStarted {
		rsc.logger.Warn("redis streams consumer is already stopped")
		return
	}

	rsc.logger.Info("stopping redis streams consumer...", "consumer group", rsc.consumerGroup, "stream name", rsc.streamName, "consumer name", rsc.consumerName)

	rsc.isStarted = false
}

func (rsc *RedisStreamsConsumer) Start() {
	if rsc.isStarted {
		rsc.logger.Warn("redis streams consumer is already started")
		return
	}

	rsc.logger.Info("starting redis streams consumer...", "consumer group", rsc.consumerGroup, "stream name", rsc.streamName, "consumer name", rsc.consumerName)

	rsc.isStarted = true
	for rsc.isStarted {
		func() {
			rsc.logger.Debug("new redis streams consumer round", "time", time.Now().Format(time.RFC3339))

			streams, err := rsc.client.XReadGroup(rsc.ctx, &redis.XReadGroupArgs{
				Streams:  []string{rsc.streamName, START},
				Group:    rsc.consumerGroup,
				Consumer: rsc.consumerName,
				Count:    10,
				Block:    0,
			}).Result()

			if err != nil {
				rsc.logger.Warn("error reading stream from redis streams consumer group", "stream", rsc.streamName, "consumer group", rsc.consumerGroup, "consumer name", rsc.consumerName)
				return
			}

			for _, message := range streams[0].Messages {
				rsc.wg.Add(1)
				go rsc.processMessage(message)
			}
			rsc.wg.Wait()
		}()
	}
}

func (rsc *RedisStreamsConsumer) StartPendingMessagesConsumer(timePeriodInSeconds int) {
	rsc.logger.Info("starting redis pending streams messages consumer...", "consumer group", rsc.consumerGroup, "stream name", rsc.streamName, "consumer name", rsc.consumerName, "period in seconds", timePeriodInSeconds)
	ticker := time.Tick(time.Second * time.Duration(timePeriodInSeconds))

	for range ticker {
		rsc.consumePendingMessages()
	}
}

func (rsc *RedisStreamsConsumer) consumePendingMessages() {
	rsc.logger.Debug("processing pending redis streams messages...", "time", time.Now().Format(time.RFC3339))

	var pendingMessagesIds []string
	pendingMessagesToProcess, err := rsc.client.XPendingExt(rsc.ctx, &redis.XPendingExtArgs{
		Stream: rsc.streamName,
		Group:  rsc.consumerGroup,
		Start:  "0",
		End:    "+",
		Count:  10,
	}).Result()

	if err != nil {
		rsc.logger.Error("error retrieving pending redis streams messages", "error", err)
	}

	for _, message := range pendingMessagesToProcess {
		pendingMessagesIds = append(pendingMessagesIds, message.ID)
	}

	if len(pendingMessagesIds) > 0 {
		messages, err := rsc.client.XClaim(rsc.ctx, &redis.XClaimArgs{
			Stream:   rsc.streamName,
			Group:    rsc.consumerGroup,
			Consumer: rsc.consumerName,
			Messages: pendingMessagesIds,
			MinIdle:  30 * time.Second,
		}).Result()

		if err != nil {
			rsc.logger.Error("error processing redis streams pending messages", "error", err)
			return
		}

		for _, message := range messages {
			rsc.wg.Add(1)
			go rsc.processMessage(message)
		}
		rsc.wg.Wait()
	}
}

func (rsc *RedisStreamsConsumer) processMessage(message redis.XMessage) {
	defer rsc.wg.Done()

	err := rsc.handleMessage(message)

	// Backoff policy
	retries := 0
	for err != nil && retries < MAX_PROCESS_RETRIES {
		rsc.logger.Warn("retrying failed process redis streams message", "stream", rsc.streamName, "consumer group", rsc.consumerGroup, "consumer name", rsc.consumerName, "message", message, "error", err, "retry", retries)
		retries++

		<-time.After(time.Duration(500*retries) * time.Millisecond)
		err = rsc.handleMessage(message)
	}

	// Log error
	if err != nil {
		rsc.logger.Error("error processing redis streams message", "stream", rsc.streamName, "consumer group", rsc.consumerGroup, "consumer name", rsc.consumerName, "message", message, "error", err)
	}

	// Mark message as acknowledged
	rsc.client.XAck(rsc.ctx, rsc.streamName, rsc.consumerGroup, message.ID)

	// Delete message from Redis Streams
	rsc.client.XDel(rsc.ctx, rsc.streamName, message.ID)
}

func (rsc *RedisStreamsConsumer) handleMessage(msg redis.XMessage) error {
	// Retrieve "type" value from message
	msgType, ok := msg.Values["type"].(string)
	if !ok {
		return errors.New("error retrieving type from redis streams message")
	}

	// Retrieve "data" value from message
	msgDataAsString, ok := msg.Values["data"].(string)
	if !ok {
		return errors.New("error retrieving data from redis streams message")
	}

	msgDataAsMap := pkg_utils.JsonStrToMap(msgDataAsString)

	rsc.logger.Debug("handling message from redis stream", "type", msgType, "message", msgDataAsMap)

	for _, handler := range rsc.handlers {
		if handler.MessageType() == msgType {
			domainMessage, err := handler.NewDomainMessageFromMap(msgDataAsMap)
			if err != nil {
				return err
			}

			handleErr := handler.Handle(domainMessage)
			if handleErr != nil {
				return handleErr
			}
		}
	}

	return nil
}
