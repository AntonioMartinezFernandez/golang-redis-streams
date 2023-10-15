package like_application

import (
	"encoding/json"
	"errors"
	"time"

	like_domain "github.com/AntonioMartinezFernandez/golang-redis-streams/internal/like/domain"

	pkg_domain "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/domain"
	pkg_utils "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/utils"
)

var _ pkg_domain.DomainMessage = (*LikeCreatedMessage)(nil)

const LikeCreatedMessageType string = "like-created"

type LikeCreatedMessageAttributes struct {
	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	CommentId string `json:"comment_id"`
	CreatedAt int64  `json:"created_at"`
}

type LikeCreatedMessageMetadata struct {
	CreatedAt int64 `json:"created_at"`
}

type LikeCreatedMessage struct {
	*pkg_domain.DomainMessageBase

	Attributes LikeCreatedMessageAttributes `json:"attributes"`
	Metadata   LikeCreatedMessageMetadata   `json:"metadata"`
}

func NewLikeCreatedMessage(like *like_domain.Like) *LikeCreatedMessage {
	DomainMessageBase := pkg_domain.NewDomainMessageBase(LikeCreatedMessageType)

	return &LikeCreatedMessage{
		DomainMessageBase: DomainMessageBase,

		Attributes: LikeCreatedMessageAttributes{
			Id:        like.Id(),
			UserId:    like.UserId(),
			CommentId: like.CommentId(),
			CreatedAt: like.CreatedAt().UnixMilli(),
		},
		Metadata: LikeCreatedMessageMetadata{
			CreatedAt: time.Now().UnixMilli(),
		},
	}
}

func NewLikeCreatedMessageFromMap(messageAsMap map[string]interface{}) (*LikeCreatedMessage, error) {
	attributes, ok := pkg_utils.GetMapValue("attributes", messageAsMap).(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid attributes")
	}

	id, ok := pkg_utils.GetMapValue("id", attributes).(string)
	if !ok {
		return nil, errors.New("invalid id")
	}

	userId, ok := pkg_utils.GetMapValue("user_id", attributes).(string)
	if !ok {
		return nil, errors.New("invalid user id")
	}

	commentId, ok := pkg_utils.GetMapValue("comment_id", attributes).(string)
	if !ok {
		return nil, errors.New("invalid comment")
	}

	createdAtAsFloat, ok := pkg_utils.GetMapValue("created_at", attributes).(float64)
	if !ok {
		return nil, errors.New("invalid created at")
	}
	createdAtAsTime := time.UnixMilli(int64(createdAtAsFloat))

	newLike, err := like_domain.NewLike(id, userId, commentId, createdAtAsTime)
	if err != nil {
		return nil, err
	}

	return NewLikeCreatedMessage(newLike), nil
}

func (lcm LikeCreatedMessage) AsMap() map[string]interface{} {
	msgAsBytes, _ := json.Marshal(lcm)
	var msgAsMap map[string]interface{}
	json.Unmarshal(msgAsBytes, &msgAsMap)
	return msgAsMap
}
