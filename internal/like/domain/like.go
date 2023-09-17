package like

import (
	"errors"
	"strconv"
	"time"

	"github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/utils"
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

func NewFromMap(e map[string]interface{}) (*Like, error) {
	id, ok := utils.GetMapValue("id", e).(string)
	if !ok {
		return nil, errors.New("invalid id")
	}
	userId, ok := utils.GetMapValue("user_id", e).(string)
	if !ok {
		return nil, errors.New("invalid user id")
	}
	commentId, ok := utils.GetMapValue("comment_id", e).(string)
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

	return NewLike(id, userId, commentId, createdAtAsTime)
}

func (c *Like) GetId() string {
	return c.id
}

func (c *Like) GetUserId() string {
	return c.userId
}

func (c *Like) GetCommentId() string {
	return c.commentId
}

func (c *Like) GetCreatedAt() time.Time {
	return c.createdAt
}
