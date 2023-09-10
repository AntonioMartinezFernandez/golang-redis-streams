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

func StartPendingEventsConsumer(
	ctx context.Context,
	wg *sync.WaitGroup,
	client *redis.Client,
	streamName string,
	consumerGroup string,
	consumerName string,
	timePeriodInSeconds int,
) {

	ticker := time.Tick(time.Second * time.Duration(timePeriodInSeconds))

	for range ticker {
		func() {

			var streamsRetry []string
			pendingStreams, err := client.XPendingExt(ctx, &redis.XPendingExtArgs{
				Stream: streamName,
				Group:  consumerGroup,
				Start:  "0",
				End:    "+",
				Count:  10,
				//Consumer string
			}).Result()

			if err != nil {
				panic(err)
			}

			for _, stream := range pendingStreams {
				streamsRetry = append(streamsRetry, stream.ID)
			}

			if len(streamsRetry) > 0 {

				streams, err := client.XClaim(ctx, &redis.XClaimArgs{
					Stream:   streamName,
					Group:    consumerGroup,
					Consumer: consumerName,
					Messages: streamsRetry,
					MinIdle:  30 * time.Second,
				}).Result()

				if err != nil {
					log.Printf("err on process pending: %+v\n", err)
					return
				}

				for _, stream := range streams {
					wg.Add(1)
					go processStream(ctx, wg, client, stream, streamName, consumerGroup, true, shared.HandlerFactory())
				}
				wg.Wait()
			}

			fmt.Println("process pending streams at ", time.Now().Format(time.RFC3339))

		}()
	}
}
