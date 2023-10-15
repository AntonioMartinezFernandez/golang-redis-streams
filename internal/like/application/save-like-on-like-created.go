package like_application

import (
	"errors"
	"fmt"

	pkg_domain "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/domain"
)

var _ pkg_domain.DomainMessageHandler = (*SaveLikeOnLikeCreated)(nil)

type SaveLikeOnLikeCreated struct {
}

func NewSaveLikeOnLikeCreated() *SaveLikeOnLikeCreated {
	return &SaveLikeOnLikeCreated{}
}

func (slolc *SaveLikeOnLikeCreated) MessageType() string {
	return LikeCreatedMessageType
}

func (slolc *SaveLikeOnLikeCreated) NewDomainMessageFromMap(messageAsMap map[string]interface{}) (pkg_domain.DomainMessage, error) {
	likeCreated, err := NewLikeCreatedMessageFromMap(messageAsMap)
	if err != nil {
		return nil, err
	}
	return likeCreated, nil
}

func (slolc *SaveLikeOnLikeCreated) Handle(message interface{}) error {
	likeCreated, ok := message.(*LikeCreatedMessage)
	if !ok {
		return errors.New("message cannot be casted as LikeCreatedMessage")
	}

	fmt.Printf("Running SaveLikeOnLikeCreated handler for like id %s and comment id %s\n", likeCreated.Attributes.Id, likeCreated.Attributes.CommentId)

	return nil
}
