package main

import (
	"context"
	"fmt"

	config "github.com/AntonioMartinezFernandez/golang-redis-streams/config"
	utils "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/utils"

	"github.com/redis/go-redis/v9"
)

var (
	env_vars   config.Config
	ctx        context.Context
	streamName string
	client     *redis.Client
	someId     uint64
)

func init() {
	// Init context
	ctx = context.Background()

	// Init variables
	env_vars = config.LoadEnvConfig()
	streamName = env_vars.Stream
	someId = 0

	// Init Redis client
	var err error
	client, err = utils.NewRedisClient(env_vars.RedisHost, env_vars.RedisPort)
	if err != nil {
		panic(err)
	}
}

func main() {
	for i := 0; i < 15; i++ {
		someId++

		// Generate random event
		event, eventType := utils.GenerateRamdomEvent(ctx, client, streamName, someId)

		// Publish event
		eventOffset, err := event.Publish()

		logResultOrError(err, eventOffset, string(eventType))
	}
}

func logResultOrError(err error, eventOffset string, eventType string) {
	if err != nil {
		fmt.Printf("producer event error: %v\n", err)
	} else {
		fmt.Printf("producer event success  Type:%v  offset:%v\n", string(eventType), eventOffset)
	}
}
