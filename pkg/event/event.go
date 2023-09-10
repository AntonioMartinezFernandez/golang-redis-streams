package event

import (
	"encoding"
	"fmt"
	"time"
)

type Type string

const (
	LikeType    Type = "LikeType"
	CommentType Type = "CommentType"
)

type Base struct {
	ID       string
	Type     Type
	DateTime time.Time
	Retry    bool
}

// Event ...
type Event interface {
	GetID() string
	GetType() Type
	GetDateTime() time.Time
	SetID(id string)
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

func (o *Base) GetID() string {
	return o.ID
}

func (o *Base) SetID(id string) {
	o.ID = id
}

func (o *Base) GetType() Type {
	return o.Type
}

func (o *Base) GetDateTime() time.Time {
	return o.DateTime
}

func (o *Base) String() string {

	return fmt.Sprintf("id:%s type:%s", o.ID, o.Type)
}
