package like

import (
	"time"
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
