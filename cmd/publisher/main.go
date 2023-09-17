package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	config "github.com/AntonioMartinezFernandez/golang-redis-streams/config"
	utils "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/utils"

	slogger "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/logger"
	redis_streams "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/redis-streams"

	comment_application "github.com/AntonioMartinezFernandez/golang-redis-streams/internal/comment/application"
	comment_domain "github.com/AntonioMartinezFernandez/golang-redis-streams/internal/comment/domain"
	like_application "github.com/AntonioMartinezFernandez/golang-redis-streams/internal/like/application"
	like_domain "github.com/AntonioMartinezFernandez/golang-redis-streams/internal/like/domain"

	"github.com/redis/go-redis/v9"
)

var (
	env_vars              config.Config
	ctx                   context.Context
	client                *redis.Client
	redisStreamsPublisher *redis_streams.RedisStreamsPublisher
	logger                *slog.Logger
)

func init() {
	InitDependencies()
}

func main() {
	startTime := time.Now()
	commentIds := []string{}
	streamsToPublishByType := 500000

	// Publish CommentCreated streams
	for i := 1; i <= streamsToPublishByType; i++ {
		comment, _ := comment_domain.NewComment(utils.NewUuid(), fmt.Sprint(i), utils.RandomHello(), time.Now())
		commentCreatedStream := comment_application.NewCommentCreatedStream(comment)
		_ = redisStreamsPublisher.Publish(commentCreatedStream.GetBaseStreamType(), commentCreatedStream)

		commentIds = append(commentIds, comment.GetId())
	}

	// Publish LikeCreated streams
	for i := 1; i <= streamsToPublishByType; i++ {
		like, _ := like_domain.NewLike(utils.NewUuid(), fmt.Sprint(i), commentIds[i-1], time.Now())
		likeCreatedStream := like_application.NewLikeCreatedStream(like)
		_ = redisStreamsPublisher.Publish(likeCreatedStream.GetBaseStreamType(), likeCreatedStream)
	}

	fmt.Printf("Published %d streams in %s\n", streamsToPublishByType*2, time.Since(startTime))
}

func InitDependencies() {
	ctx = context.Background()
	env_vars = config.LoadEnvConfig()

	// Init logger
	logger = slogger.NewLogger(env_vars.LogLevel)

	// Init Redis client
	var err error
	client, err = utils.NewRedisClient(env_vars.RedisHost, env_vars.RedisPort)
	if err != nil {
		panic(err)
	}

	// Init Redis Streams publisher
	redisStreamsPublisher = redis_streams.NewRedisStreamsPublisher(ctx, client, logger)
}
