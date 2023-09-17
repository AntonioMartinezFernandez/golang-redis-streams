package redis_streams

type RedisStreamsSubscriber interface {
	MessageTypeName() string
	NewStreamEventFromMap(map[string]interface{}) StreamToPublish
	Handle(e interface{}) error
}
