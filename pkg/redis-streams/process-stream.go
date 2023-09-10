package redis_streams

import (
	"context"
	"fmt"
	"sync"

	"github.com/AntonioMartinezFernandez/golang-redis-streams/internal/comment"
	"github.com/AntonioMartinezFernandez/golang-redis-streams/internal/like"
	"github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/event"
	"github.com/redis/go-redis/v9"
)

func processStream(
	ctx context.Context,
	wg *sync.WaitGroup,
	client *redis.Client,
	stream redis.XMessage,
	streamName string,
	consumerGroup string,
	retry bool,
	handlerFactory func(t event.Type) event.Handler,
) {
	defer wg.Done()

	typeEvent := stream.Values["type"].(string)

	newEvent, _ := newEvent(event.Type(typeEvent))

	err := newEvent.UnmarshalBinary([]byte(stream.Values["data"].(string)))
	if err != nil {
		fmt.Printf("error on unmarshal stream:%v\n", stream.ID)
		return
	}

	newEvent.SetID(stream.ID)

	h := handlerFactory(newEvent.GetType())
	err = h.Handle(newEvent, retry)
	if err != nil {
		fmt.Printf("error on process event:%v\n", newEvent)
		fmt.Println(err)
		return
	}

	// Mark event as acknowledged
	client.XAck(ctx, streamName, consumerGroup, stream.ID)

	// Delete event from Redis Streams
	client.XDel(ctx, streamName, stream.ID)
}

func newEvent(t event.Type) (event.Event, error) {
	b := &event.Base{
		Type: t,
	}

	switch t {

	case like.LikeType:
		return &like.LikeEvent{
			Base: b,
		}, nil

	case comment.CommentType:
		return &comment.CommentEvent{
			Base: b,
		}, nil

	}

	return nil, fmt.Errorf("type %v not supported", t)
}
