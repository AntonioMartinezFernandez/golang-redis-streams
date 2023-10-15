package comment_domain

import (
	"time"
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
