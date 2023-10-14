package comment_application

import (
	"errors"
	"fmt"

	pkg_domain "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/domain"
)

var _ pkg_domain.DomainMessageHandler = (*SaveCommentOnCommentCreated)(nil)

type SaveCommentOnCommentCreated struct {
}

func NewSaveCommentOnCommentCreated() *SaveCommentOnCommentCreated {
	return &SaveCommentOnCommentCreated{}
}

func (scocc *SaveCommentOnCommentCreated) MessageType() string {
	return CommentCreatedMessageType
}

func (scocc *SaveCommentOnCommentCreated) NewDomainMessageFromMap(messageAsMap map[string]interface{}) (pkg_domain.DomainMessage, error) {
	commentCreated, err := NewCommentCreatedMessageFromMap(messageAsMap)
	if err != nil {
		return nil, err
	}
	return commentCreated, nil
}

func (scocc *SaveCommentOnCommentCreated) Handle(message interface{}) error {
	commentCreated, ok := message.(*CommentCreatedMessage)
	if !ok {
		return errors.New("message cannot be casted as CommentCreatedMessage")
	}

	fmt.Println("Running SaveCommentOnCommentCreated handler for comment with id", commentCreated.Id)

	return nil
}
