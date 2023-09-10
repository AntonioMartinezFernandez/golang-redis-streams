package redis_streams

import (
	"context"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

func CreateConsumerGroup(
	ctx context.Context,
	client *redis.Client,
	streamName string,
	consumerGroup string,
) {
	if _, err := client.XGroupCreateMkStream(ctx, streamName, consumerGroup, "0").Result(); err != nil {
		if !strings.Contains(fmt.Sprint(err), "BUSYGROUP") {
			fmt.Printf("Error on create Consumer Group: %v ...\n", consumerGroup)
			panic(err)
		}
	}
}
