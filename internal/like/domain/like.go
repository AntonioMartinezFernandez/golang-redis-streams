package like

import (
	"errors"
	"time"

	pkg_utils "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/utils"
)

type Like struct {
	id        string
	userId    string
	commentId string
	createdAt time.Time
}

func NewLike(id string, userId string, commentId string, createdAt time.Time) (*Like, error) {
	return &Like{
		id:        id,
		userId:    userId,
		commentId: commentId,
		createdAt: createdAt,
	}, nil
}

func (l *Like) Id() string {
	return l.id
}

func (l *Like) UserId() string {
	return l.userId
}

func (l *Like) CommentId() string {
	return l.commentId
}

func (l *Like) CreatedAt() time.Time {
	return l.createdAt
}

func NewLikeFromMap(e map[string]interface{}) (*Like, error) {
	attributes, ok := pkg_utils.GetMapValue("attributes", e).(map[string]interface{})
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

	metadata, ok := pkg_utils.GetMapValue("metadata", e).(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid metadata")
	}

	createdAtAsFloat, ok := pkg_utils.GetMapValue("created_at", metadata).(float64)
	if !ok {
		return nil, errors.New("invalid created at")
	}
	createdAtAsTime := time.UnixMilli(int64(createdAtAsFloat))

	return NewLike(id, userId, commentId, createdAtAsTime)
}
