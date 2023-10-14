package redis_streams

func StreamName(eventName string) string {
	return "events:" + eventName
}
