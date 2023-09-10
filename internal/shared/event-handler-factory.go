package shared

import (
	event "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/event"

	comment "github.com/AntonioMartinezFernandez/golang-redis-streams/internal/comment"
	like "github.com/AntonioMartinezFernandez/golang-redis-streams/internal/like"
)

func HandlerFactory() func(t event.Type) event.Handler {

	return func(t event.Type) event.Handler {
		switch t {
		case like.LikeType:
			return like.NewLikeHandler()
		case comment.CommentType:
			return comment.NewCommentHandler()
		default:
			return event.NewDefaultHandler()
		}
	}
}
