package utils

import (
	"context"
	"math/rand"
	"time"

	comment "github.com/AntonioMartinezFernandez/golang-redis-streams/internal/comment"
	like "github.com/AntonioMartinezFernandez/golang-redis-streams/internal/like"
	event "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/event"
	redis_streams "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/redis-streams"

	"github.com/redis/go-redis/v9"
)

func GenerateRamdomEvent(
	ctx context.Context,
	client *redis.Client,
	streamName string,
	userID uint64,
) (*redis_streams.RedisStreamEvent, event.Type) {
	// Random event type
	eventType := []event.Type{event.LikeType, event.CommentType}[rand.Intn(2)]

	if eventType == event.LikeType {
		// Return like event
		return redis_streams.NewRedisStreamEvent(
			ctx,
			client,
			streamName,
			map[string]interface{}{
				"type": string(eventType),
				"data": &like.LikeEvent{
					Base: &event.Base{
						Type:     eventType,
						DateTime: time.Now(),
					},
					UserID: userID,
				},
			},
		), eventType
	} else {
		// Random greeting
		commentMessage := []string{
			"Kaixo redis streams!",
			"Hola redis streams!",
			"Hello redis streams!",
			"Konichiwa redis streams!",
			"Sveika redis streams!",
		}[rand.Intn(5)]

		// Return comment event
		return redis_streams.NewRedisStreamEvent(
			ctx,
			client,
			streamName,
			map[string]interface{}{
				"type": string(eventType),
				"data": &comment.CommentEvent{
					Base: &event.Base{
						Type:     eventType,
						DateTime: time.Now(),
					},
					UserID:  userID,
					Comment: commentMessage,
				},
			},
		), eventType
	}
}
