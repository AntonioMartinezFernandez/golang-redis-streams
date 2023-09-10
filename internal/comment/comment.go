package comment

import (
	event "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/event"
	"github.com/vmihailenco/msgpack/v4"
)

const CommentType event.Type = "CommentType"

type CommentEvent struct {
	*event.Base
	UserID  uint64
	Comment string
}

func (o *CommentEvent) MarshalBinary() (data []byte, err error) {
	return msgpack.Marshal(o)
}

func (o *CommentEvent) UnmarshalBinary(data []byte) error {
	return msgpack.Unmarshal(data, o)
}

func NewCommentEvent() (CommentEvent, error) {
	b := &event.Base{
		Type: event.Type("CommentType"),
	}

	return CommentEvent{
		Base: b,
	}, nil
}
