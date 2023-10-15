package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	wg                    sync.WaitGroup
	redisClient           *redis.Client
	redisStreamsPublisher *pkg_redis_streams.RedisStreamsPublisher
	logger                *slog.Logger
)

func init() {
	InitDependencies()
}

func main() {
	redisStreamsPublisher.Start()

	var msg_counter uint32 = 0
	var delayed_msg_counter uint32 = 0

	startTime := time.Now()
	commentIds := []string{}
	messagesToPublishByType := 50000

	// Publish INSTANT CommentCreated messages
	for i := 1; i <= messagesToPublishByType; i++ {
		commentId := pkg_utils.NewUuid()
		userId := fmt.Sprint(i)
		commentMsg := pkg_utils.RandomHello()
		now := time.Now()

		commentIds = append(commentIds, commentId)

		comment, _ := comment_domain.NewComment(commentId, userId, commentMsg, now)

		CommentCreatedMessage := comment_application.NewCommentCreatedMessage(comment)
		err := redisStreamsPublisher.Publish(CommentCreatedMessage)
		if err != nil {
			logger.Error("error publishing delayed message through redis streams ", "error", err)
		}

		msg_counter++
	}

	// Publish 5 seconds DELAYED LikeCreated messages
	for i := 1; i <= messagesToPublishByType; i++ {
		likeId := pkg_utils.NewUuid()
		userId := fmt.Sprint(i)
		commentId := commentIds[i-1]
		now := time.Now()

		like, _ := like_domain.NewLike(likeId, userId, commentId, now)

		LikeCreatedMessage := like_application.NewLikeCreatedMessage(like)
		err := redisStreamsPublisher.PublishDelayed(LikeCreatedMessage, 5)
		if err != nil {
			logger.Error("error publishing delayed message through redis streams ", "error", err)
		}

		delayed_msg_counter++
	}

	fmt.Printf("Published %d messages\n", msg_counter)
	fmt.Printf("Published %d delayed messages\n", delayed_msg_counter)
	fmt.Printf("Processed in %s\n", time.Since(startTime))

	//Gracefully shutdown
	chanOS := make(chan os.Signal, 1)
	signal.Notify(chanOS, syscall.SIGINT, syscall.SIGTERM)

	<-chanOS
	redisStreamsPublisher.Stop()

	wg.Wait()
	redisClient.Close()
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
	redisStreamsPublisher = pkg_redis_streams.NewRedisStreamsPublisher(ctx, &wg, redisClient, logger)
}
