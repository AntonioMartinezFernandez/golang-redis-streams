package like_application

import (
	"errors"
	"fmt"

	redis_streams "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/redis-streams"
)

var _ redis_streams.RedisStreamsSubscriber = (*SaveLikeOnLikeCreated)(nil)

type SaveLikeOnLikeCreated struct {
}

func NewSaveLikeOnLikeCreated() *SaveLikeOnLikeCreated {
	return &SaveLikeOnLikeCreated{}
}

func (slolc *SaveLikeOnLikeCreated) MessageTypeName() string {
	return LikeCreatedStreamType
}

func (slolc *SaveLikeOnLikeCreated) NewStreamEventFromMap(eventAsMap map[string]interface{}) redis_streams.StreamToPublish {
	likeCreated, err := NewLikeCreatedStreamFromMap(eventAsMap)
	if err != nil {
		return nil
	}
	return likeCreated
}

func (slolc *SaveLikeOnLikeCreated) Handle(event interface{}) error {
	likeCreated, ok := event.(*LikeCreatedStream)
	if !ok {
		return errors.New("event cannot be casted as LikeCreatedStream")
	}

	fmt.Println("Running SaveLikeOnLikeCreated handler for like with id", likeCreated.Id)

	return nil
}
