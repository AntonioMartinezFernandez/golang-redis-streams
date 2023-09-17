package redis_streams

/*
! Do not use nested JSON structures in the streams to avoid Redis marshalling errors
https://stackoverflow.com/a/67782704
*/

import (
	"time"

	utils "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/utils"
)

type BaseStreamHandler func(e BaseStream) error

type BaseStream struct {
	BaseStreamId       string    `json:"base_stream_id"`
	BaseStreamType     string    `json:"base_stream_type"`
	BaseStreamDateTime time.Time `json:"base_stream_date_time"`
	BaseStreamRetry    bool      `json:"base_stream_retry"`
}

func NewBaseStream(streamType string) *BaseStream {
	return &BaseStream{
		BaseStreamId:       utils.NewUuid(),
		BaseStreamType:     streamType,
		BaseStreamDateTime: time.Now(),
		BaseStreamRetry:    false,
	}
}

func (bs *BaseStream) GetBaseStreamId() string {
	return bs.BaseStreamId
}

func (bs *BaseStream) GetBaseStreamType() string {
	return bs.BaseStreamType
}

func (bs *BaseStream) GetBaseStreamDateTime() time.Time {
	return bs.BaseStreamDateTime
}

func (bs *BaseStream) GetBaseStreamRetry() bool {
	return bs.BaseStreamRetry
}
