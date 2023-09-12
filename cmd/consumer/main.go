package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	config "github.com/AntonioMartinezFernandez/golang-redis-streams/config"
	slogger "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/logger"
	redis_streams "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/redis-streams"
	utils "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/utils"

	comment_application "github.com/AntonioMartinezFernandez/golang-redis-streams/internal/comment/application"
	like_application "github.com/AntonioMartinezFernandez/golang-redis-streams/internal/like/application"

	"github.com/redis/go-redis/v9"
)

var (
	ctx           context.Context
	wg            sync.WaitGroup
	env_vars      config.Config
	client        *redis.Client
	consumerGroup string
	logger        *slog.Logger
)

func init() {
	ctx = context.Background()
	env_vars = config.LoadEnvConfig()
	consumerGroup = env_vars.ConsumerGroup

	// Init logger
	logger = slogger.NewLogger(env_vars.LogLevel)

	// Init Redis client
	var err error
	client, err = utils.NewRedisClient(env_vars.RedisHost, env_vars.RedisPort)
	if err != nil {
		panic(err)
	}
}

func main() {
	// Init Redis Stream consumer for Comments
	commentCreatedStreamName := comment_application.CommentCreatedStreamType
	comments_rsc := redis_streams.NewRedisStreamsConsumer(ctx, &wg, logger, client, consumerGroup, commentCreatedStreamName)
	comments_rsc.RegisterSubscriber(comment_application.NewSaveCommentOnCommentCreated())
	comments_rsc.CreateConsumerGroup()
	go comments_rsc.Start()
	go comments_rsc.StartPendingMessagesConsumer(60)

	// Init Redis Stream consumer for Likes
	likeCreatedStreamName := like_application.LikeCreatedStreamType
	likes_rsc := redis_streams.NewRedisStreamsConsumer(ctx, &wg, logger, client, consumerGroup, likeCreatedStreamName)
	likes_rsc.RegisterSubscriber(like_application.NewSaveLikeOnLikeCreated())
	likes_rsc.CreateConsumerGroup()
	go likes_rsc.Start()
	go likes_rsc.StartPendingMessagesConsumer(60)

	//Gracefully shutdown
	chanOS := make(chan os.Signal, 1)
	signal.Notify(chanOS, syscall.SIGINT, syscall.SIGTERM)

	<-chanOS
	comments_rsc.Stop()
	// likes_rsc.Stop()

	wg.Wait()
	client.Close()
}
