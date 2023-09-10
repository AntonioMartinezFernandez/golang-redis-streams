package utils

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Create a new instace of Redis client
func NewRedisClient(host string, port string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "",
		DB:       0, // use default DB
	})

	_, err := client.Ping(context.Background()).Result()
	return client, err

}
