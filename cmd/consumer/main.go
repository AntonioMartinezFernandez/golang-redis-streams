package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	config "github.com/AntonioMartinezFernandez/golang-redis-streams/config"

	pkg_logger "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/logger"
	pkg_redis_streams "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/redis-streams"
	pkg_utils "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/utils"

	comment_application "github.com/AntonioMartinezFernandez/golang-redis-streams/internal/comment/application"
	like_application "github.com/AntonioMartinezFernandez/golang-redis-streams/internal/like/application"

	"github.com/redis/go-redis/v9"
)

var (
	ctx           context.Context
	wg            sync.WaitGroup
	env_vars      config.Config
	redisClient   *redis.Client
	consumerGroup string
	logger        *slog.Logger
)

func init() {
	InitDependencies()
}

func main() {
	// Init Redis Stream consumer for CommentCreated messages
	commentCreatedMessageType := comment_application.CommentCreatedMessageType
	saveCommentOnCommentCreatedHandler := comment_application.NewSaveCommentOnCommentCreated()

	commentCreatedMessagesConsumer := pkg_redis_streams.NewRedisStreamsConsumer(ctx, &wg, logger, redisClient, consumerGroup, commentCreatedMessageType)
	commentCreatedMessagesConsumer.RegisterHandler(saveCommentOnCommentCreatedHandler)
	commentCreatedMessagesConsumer.CreateConsumerGroup()

	go commentCreatedMessagesConsumer.Start()
	go commentCreatedMessagesConsumer.StartPendingMessagesConsumer(60)

	// Init Redis Stream consumer for LikeCreated messages
	likeCreatedMessageType := like_application.LikeCreatedMessageType
	saveLikeOnLikeCreatedHandler := like_application.NewSaveLikeOnLikeCreated()

	likeCreatedMessagesConsumer := pkg_redis_streams.NewRedisStreamsConsumer(ctx, &wg, logger, redisClient, consumerGroup, likeCreatedMessageType)
	likeCreatedMessagesConsumer.RegisterHandler(saveLikeOnLikeCreatedHandler)
	likeCreatedMessagesConsumer.CreateConsumerGroup()

	go likeCreatedMessagesConsumer.Start()
	go likeCreatedMessagesConsumer.StartPendingMessagesConsumer(60)

	//Gracefully shutdown
	chanOS := make(chan os.Signal, 1)
	signal.Notify(chanOS, syscall.SIGINT, syscall.SIGTERM)

	<-chanOS
	commentCreatedMessagesConsumer.Stop()
	likeCreatedMessagesConsumer.Stop()

	wg.Wait()
	redisClient.Close()
}

func InitDependencies() {
	ctx = context.Background()
	env_vars = config.LoadEnvConfig()
	consumerGroup = env_vars.ConsumerGroup

	// Init logger
	logger = pkg_logger.NewLogger(env_vars.LogLevel)

	// Init Redis client
	var err error
	redisClient, err = pkg_utils.NewRedisClient(env_vars.RedisHost, env_vars.RedisPort, "")
	if err != nil {
		panic(err)
	}
}
