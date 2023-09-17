package comment_application

import (
	"errors"
	"fmt"

	redis_streams "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/redis-streams"
)

var _ redis_streams.RedisStreamsSubscriber = (*SaveCommentOnCommentCreated)(nil)

type SaveCommentOnCommentCreated struct {
}

func NewSaveCommentOnCommentCreated() *SaveCommentOnCommentCreated {
	return &SaveCommentOnCommentCreated{}
}

func (scocc *SaveCommentOnCommentCreated) MessageTypeName() string {
	return CommentCreatedStreamType
}

func (scocc *SaveCommentOnCommentCreated) NewStreamEventFromMap(eventAsMap map[string]interface{}) redis_streams.StreamToPublish {
	commentCreated, err := NewCommentCreatedStreamFromMap(eventAsMap)
	if err != nil {
		return nil
	}
	return commentCreated
}

func (scocc *SaveCommentOnCommentCreated) Handle(event interface{}) error {
	commentCreated, ok := event.(*CommentCreatedStream)
	if !ok {
		return errors.New("event cannot be casted as CommentCreatedStream")
	}

	fmt.Println("Running SaveCommentOnCommentCreated handler for comment with id", commentCreated.Id)

	return nil
}
