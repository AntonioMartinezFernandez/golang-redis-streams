package like_application

import (
	"encoding/json"
	"fmt"

	like_domain "github.com/AntonioMartinezFernandez/golang-redis-streams/internal/like/domain"

	redis_streams "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/redis-streams"
)

const LikeCreatedStreamType string = "LikeCreated"

type LikeCreatedStream struct {
	*redis_streams.BaseStream
	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	CommentId string `json:"comment_id"`
	CreatedAt string `json:"created_at"`
}

func NewLikeCreatedStream(comment *like_domain.Like) *LikeCreatedStream {
	baseStream := redis_streams.NewBaseStream(LikeCreatedStreamType)

	return &LikeCreatedStream{
		BaseStream: baseStream,
		Id:         comment.GetId(),
		UserId:     comment.GetUserId(),
		CommentId:  comment.GetCommentId(),
		CreatedAt:  fmt.Sprint(comment.GetCreatedAt().UnixMilli()),
	}
}

func NewLikeCreatedStreamFromMap(stream map[string]interface{}) (*LikeCreatedStream, error) {
	newLike, err := like_domain.NewFromMap(stream)
	if err != nil {
		return nil, err
	}

	return NewLikeCreatedStream(newLike), nil
}

func (cce LikeCreatedStream) AsMap() map[string]interface{} {
	binaryElem, _ := json.Marshal(cce)
	var mapElem map[string]interface{}
	json.Unmarshal(binaryElem, &mapElem)
	return mapElem
}
