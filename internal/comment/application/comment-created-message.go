package comment_application

import (
	"encoding/json"

	comment_domain "github.com/AntonioMartinezFernandez/golang-redis-streams/internal/comment/domain"

	pkg_domain "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/domain"
)

var _ pkg_domain.DomainMessage = (*CommentCreatedMessage)(nil)

const CommentCreatedMessageType string = "comment-created"

type CommentCreatedMessageAttributes struct {
	Id      string `json:"id"`
	UserId  string `json:"user_id"`
	Comment string `json:"comment"`
}

type CommentCreatedMessageMetadata struct {
	CreatedAt int64 `json:"created_at"`
}

type CommentCreatedMessage struct {
	*pkg_domain.DomainMessageBase

	Attributes CommentCreatedMessageAttributes `json:"attributes"`
	Metadata   CommentCreatedMessageMetadata   `json:"metadata"`
}

func NewCommentCreatedMessage(comment *comment_domain.Comment) *CommentCreatedMessage {
	DomainMessageBase := pkg_domain.NewDomainMessageBase(CommentCreatedMessageType)

	return &CommentCreatedMessage{
		DomainMessageBase: DomainMessageBase,

		Attributes: CommentCreatedMessageAttributes{
			Id:      comment.Id(),
			UserId:  comment.UserId(),
			Comment: comment.Comment(),
		},
		Metadata: CommentCreatedMessageMetadata{
			CreatedAt: comment.CreatedAt().UnixMilli(),
		},
	}
}

func NewCommentCreatedMessageFromMap(stream map[string]interface{}) (*CommentCreatedMessage, error) {
	newComment, err := comment_domain.NewCommentFromMap(stream)
	if err != nil {
		return nil, err
	}

	return NewCommentCreatedMessage(newComment), nil
}

func (ccm CommentCreatedMessage) AsMap() map[string]interface{} {
	msgAsBytes, _ := json.Marshal(ccm)
	var msgAsMap map[string]interface{}
	json.Unmarshal(msgAsBytes, &msgAsMap)
	return msgAsMap
}
