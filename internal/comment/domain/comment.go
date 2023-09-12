package comment_domain

import (
	"errors"
	"strconv"
	"time"

	"github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/utils"
)

type Comment struct {
	id        string
	userId    string
	comment   string
	createdAt time.Time
}

func NewComment(id string, userId string, comment string, createdAt time.Time) (*Comment, error) {
	return &Comment{
		id:        id,
		userId:    userId,
		comment:   comment,
		createdAt: createdAt,
	}, nil
}

func NewFromMap(e map[string]interface{}) (*Comment, error) {
	id, ok := utils.GetMapValue("id", e).(string)
	if !ok {
		return nil, errors.New("invalid id")
	}
	userId, ok := utils.GetMapValue("user_id", e).(string)
	if !ok {
		return nil, errors.New("invalid user id")
	}
	comment, ok := utils.GetMapValue("comment", e).(string)
	if !ok {
		return nil, errors.New("invalid comment")
	}
	createdAt, ok := utils.GetMapValue("created_at", e).(string)
	if !ok {
		return nil, errors.New("invalid created at")
	}
	createdAtAsInt, err := strconv.Atoi(createdAt)
	if err != nil {
		return nil, errors.New("invalid created at")
	}
	createdAtAsTime := time.UnixMilli(int64(createdAtAsInt))

	return NewComment(id, userId, comment, createdAtAsTime)
}

func (c *Comment) GetId() string {
	return c.id
}

func (c *Comment) GetUserId() string {
	return c.userId
}

func (c *Comment) GetComment() string {
	return c.comment
}

func (c *Comment) GetCreatedAt() time.Time {
	return c.createdAt
}
