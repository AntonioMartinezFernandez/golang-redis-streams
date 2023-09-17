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

func (scocc *SaveLikeOnLikeCreated) MessageTypeName() string {
	return LikeCreatedStreamType
}

func (scocc *SaveLikeOnLikeCreated) NewFromMap(eventAsMap map[string]interface{}) interface{} {
	likeCreated, err := NewLikeCreatedStreamFromMap(eventAsMap)
	if err != nil {
		return nil
	}
	return likeCreated
}

func (scocc *SaveLikeOnLikeCreated) Handle(event interface{}) error {
	likeCreated, ok := event.(*LikeCreatedStream)
	if !ok {
		return errors.New("event cannot be casted as LikeCreatedStream")
	}

	fmt.Println("Running SaveLikeOnLikeCreated handler for comment with id", likeCreated.Id)

	return nil
}
