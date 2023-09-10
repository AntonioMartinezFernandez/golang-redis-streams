package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	config "github.com/AntonioMartinezFernandez/golang-redis-streams/config"
	utils "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/utils"

	redis_streams "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/redis-streams"

	"github.com/redis/go-redis/v9"
)

var (
	env_vars      config.Config
	ctx           context.Context
	wg            sync.WaitGroup
	client        *redis.Client
	streamName    string
	consumerGroup string
	consumerName  string = utils.NewUuid()
)

func init() {
	// Init context
	ctx = context.Background()

	// Init variables
	env_vars = config.LoadEnvConfig()
	streamName = env_vars.Stream
	consumerGroup = env_vars.StreamGroup

	// Init Redis client
	var err error
	client, err = utils.NewRedisClient(env_vars.RedisHost, env_vars.RedisPort)
	if err != nil {
		panic(err)
	}

}

func main() {
	fmt.Printf("Starting Consumer: %v  ConsumerGroup: %v  Stream: %v\n", consumerName, consumerGroup, streamName)

	redis_streams.CreateConsumerGroup(ctx, client, streamName, consumerGroup)
	go redis_streams.StartConsumer(ctx, &wg, client, streamName, consumerGroup, consumerName)
	go redis_streams.StartPendingEventsConsumer(ctx, &wg, client, streamName, consumerGroup, consumerName, 30)

	//Gracefully disconection
	chanOS := make(chan os.Signal, 1)
	signal.Notify(chanOS, syscall.SIGINT, syscall.SIGTERM)

	<-chanOS
	fmt.Printf("Stopping Consumer:%v  ConsumerGroup: %v  Stream: %v\n", consumerName, consumerGroup, streamName)
	wg.Wait()
	client.Close()
}
