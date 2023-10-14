package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	config "github.com/AntonioMartinezFernandez/golang-redis-streams/config"

	pkg_logger "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/logger"
	pkg_redis_streams "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/redis-streams"
	pkg_utils "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/utils"

	comment_application "github.com/AntonioMartinezFernandez/golang-redis-streams/internal/comment/application"
	comment_domain "github.com/AntonioMartinezFernandez/golang-redis-streams/internal/comment/domain"
	like_application "github.com/AntonioMartinezFernandez/golang-redis-streams/internal/like/application"
	like_domain "github.com/AntonioMartinezFernandez/golang-redis-streams/internal/like/domain"

	"github.com/redis/go-redis/v9"
)

var (
	env_vars              config.Config
	ctx                   context.Context
	redisClient           *redis.Client
	redisStreamsPublisher *pkg_redis_streams.RedisStreamsPublisher
	logger                *slog.Logger
)

func init() {
	InitDependencies()
}

func main() {
	var msg_counter uint32 = 0
	startTime := time.Now()
	commentIds := []string{}
	messagesToPublishByType := 5000

	// Publish CommentCreated messages
	for i := 1; i <= messagesToPublishByType; i++ {
		commentId := pkg_utils.NewUuid()
		userId := fmt.Sprint(i)
		commentMsg := pkg_utils.RandomHello()
		now := time.Now()

		commentIds = append(commentIds, commentId)

		comment, _ := comment_domain.NewComment(commentId, userId, commentMsg, now)

		CommentCreatedMessage := comment_application.NewCommentCreatedMessage(comment)
		_ = redisStreamsPublisher.Publish(CommentCreatedMessage)

		msg_counter++
	}

	// Publish LikeCreated messages
	for i := 1; i <= messagesToPublishByType; i++ {
		likeId := pkg_utils.NewUuid()
		userId := fmt.Sprint(i)
		commentId := commentIds[i-1]
		now := time.Now()

		like, _ := like_domain.NewLike(likeId, userId, commentId, now)

		LikeCreatedMessage := like_application.NewLikeCreatedMessage(like)
		_ = redisStreamsPublisher.Publish(LikeCreatedMessage)

		msg_counter++
	}

	fmt.Printf("Published %d messages in %s\n", msg_counter, time.Since(startTime))
}

func InitDependencies() {
	ctx = context.Background()
	env_vars = config.LoadEnvConfig()

	// Init logger
	logger = pkg_logger.NewLogger(env_vars.LogLevel)

	// Init Redis client
	var err error
	redisClient, err = pkg_utils.NewRedisClient(env_vars.RedisHost, env_vars.RedisPort, "")
	if err != nil {
		panic(err)
	}

	// Init Redis Streams publisher
	redisStreamsPublisher = pkg_redis_streams.NewRedisStreamsPublisher(ctx, redisClient, logger)
}
