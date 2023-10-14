package like_application

import (
	"encoding/json"

	like_domain "github.com/AntonioMartinezFernandez/golang-redis-streams/internal/like/domain"

	pkg_domain "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/domain"
)

var _ pkg_domain.DomainMessage = (*LikeCreatedMessage)(nil)

const LikeCreatedMessageType string = "LikeCreated"

type LikeCreatedMessageAttributes struct {
	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	CommentId string `json:"comment_id"`
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
		},
		Metadata: LikeCreatedMessageMetadata{
			CreatedAt: like.CreatedAt().UnixMilli(),
		},
	}
}

func NewLikeCreatedMessageFromMap(messageAsMap map[string]interface{}) (*LikeCreatedMessage, error) {
	newLike, err := like_domain.NewLikeFromMap(messageAsMap)
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
