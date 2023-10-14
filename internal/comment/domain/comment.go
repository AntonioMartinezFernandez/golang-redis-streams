package comment_domain

import (
	"errors"
	"time"

	pkg_utils "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/utils"
)

type Comment struct {
	id        string
	userId    string
	comment   string
	createdAt time.Time
}

func (c *Comment) Id() string {
	return c.id
}

func (c *Comment) UserId() string {
	return c.userId
}

func (c *Comment) Comment() string {
	return c.comment
}

func (c *Comment) CreatedAt() time.Time {
	return c.createdAt
}

func NewComment(id string, userId string, comment string, createdAt time.Time) (*Comment, error) {
	return &Comment{
		id:        id,
		userId:    userId,
		comment:   comment,
		createdAt: createdAt,
	}, nil
}

func NewCommentFromMap(e map[string]interface{}) (*Comment, error) {
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
	comment, ok := pkg_utils.GetMapValue("comment", attributes).(string)
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

	return NewComment(id, userId, comment, createdAtAsTime)
}
