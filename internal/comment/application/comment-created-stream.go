package comment_application

import (
	"encoding/json"
	"fmt"

	comment_domain "github.com/AntonioMartinezFernandez/golang-redis-streams/internal/comment/domain"

	redis_streams "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/redis-streams"
)

const CommentCreatedStreamType string = "CommentCreated"

type CommentCreatedStream struct {
	*redis_streams.BaseStream
	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"created_at"`
}

func NewCommentCreatedStream(comment *comment_domain.Comment) *CommentCreatedStream {
	baseStream := redis_streams.NewBaseStream(CommentCreatedStreamType)

	return &CommentCreatedStream{
		BaseStream: baseStream,
		Id:         comment.GetId(),
		UserId:     comment.GetUserId(),
		Comment:    comment.GetComment(),
		CreatedAt:  fmt.Sprint(comment.GetCreatedAt().UnixMilli()),
	}
}

func NewCommentCreatedStreamFromMap(stream map[string]interface{}) (*CommentCreatedStream, error) {
	newComment, err := comment_domain.NewFromMap(stream)
	if err != nil {
		return nil, err
	}

	return NewCommentCreatedStream(newComment), nil
}

func (cce CommentCreatedStream) AsMap() map[string]interface{} {
	binaryElem, _ := json.Marshal(cce)
	var mapElem map[string]interface{}
	json.Unmarshal(binaryElem, &mapElem)
	return mapElem
}
