package redis_streams

type StreamToPublish interface {
	AsMap() map[string]interface{}
}
