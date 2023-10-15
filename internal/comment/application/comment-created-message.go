package comment_application

import (
	"encoding/json"
	"errors"
	"time"

	comment_domain "github.com/AntonioMartinezFernandez/golang-redis-streams/internal/comment/domain"

	pkg_domain "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/domain"
	pkg_utils "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/utils"
)

var _ pkg_domain.DomainMessage = (*CommentCreatedMessage)(nil)

const CommentCreatedMessageType string = "comment-created"

type CommentCreatedMessageAttributes struct {
	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	Comment   string `json:"comment"`
	CreatedAt int64  `json:"created_at"`
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
			Id:        comment.Id(),
			UserId:    comment.UserId(),
			Comment:   comment.Comment(),
			CreatedAt: comment.CreatedAt().UnixMilli(),
		},
		Metadata: CommentCreatedMessageMetadata{
			CreatedAt: time.Now().UnixMilli(),
		},
	}
}

func NewCommentCreatedMessageFromMap(messageAsMap map[string]interface{}) (*CommentCreatedMessage, error) {
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
	comment, ok := pkg_utils.GetMapValue("comment", attributes).(string)
	if !ok {
		return nil, errors.New("invalid comment")
	}
	createdAtAsFloat, ok := pkg_utils.GetMapValue("created_at", attributes).(float64)
	if !ok {
		return nil, errors.New("invalid created at")
	}
	createdAtAsTime := time.UnixMilli(int64(createdAtAsFloat))

	newComment, err := comment_domain.NewComment(id, userId, comment, createdAtAsTime)
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
