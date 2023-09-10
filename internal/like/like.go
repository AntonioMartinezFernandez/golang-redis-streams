package like

import (
	event "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/event"
	"github.com/vmihailenco/msgpack/v4"
)

const LikeType event.Type = "LikeType"

type LikeEvent struct {
	*event.Base
	UserID uint64
}

func (o *LikeEvent) MarshalBinary() (data []byte, err error) {
	return msgpack.Marshal(o)
}

func (o *LikeEvent) UnmarshalBinary(data []byte) error {
	return msgpack.Unmarshal(data, o)
}

func NewLikeEvent() (LikeEvent, error) {
	b := &event.Base{
		Type: LikeType,
	}

	return LikeEvent{
		Base: b,
	}, nil
}
