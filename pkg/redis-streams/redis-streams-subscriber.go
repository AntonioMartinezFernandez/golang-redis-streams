package redis_streams

type RedisStreamsSubscriber interface {
	MessageTypeName() string
	NewFromMap(map[string]interface{}) interface{}
	Handle(e interface{}) error
}
