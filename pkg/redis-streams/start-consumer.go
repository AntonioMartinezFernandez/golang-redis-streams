package redis_streams

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/AntonioMartinezFernandez/golang-redis-streams/internal/shared"
	"github.com/redis/go-redis/v9"
)

const START string = ">"

// start consume events
func StartConsumer(
	ctx context.Context,
	wg *sync.WaitGroup,
	client *redis.Client,
	streamName string,
	consumerGroup string,
	consumerName string,
) {
	for {
		func() {
			fmt.Println("new round ", time.Now().Format(time.RFC3339))

			streams, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
				Streams:  []string{streamName, START},
				Group:    consumerGroup,
				Consumer: consumerName,
				Count:    10,
				Block:    0,
			}).Result()

			if err != nil {
				log.Printf("err on consume events: %+v\n", err)
				return
			}

			for _, stream := range streams[0].Messages {
				wg.Add(1)
				go processStream(ctx, wg, client, stream, streamName, consumerGroup, false, shared.HandlerFactory())
			}
			wg.Wait()
		}()
	}
}
