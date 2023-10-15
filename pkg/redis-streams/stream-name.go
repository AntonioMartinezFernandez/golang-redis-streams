package redis_streams

const MESSAGES_KEY = "messages:"
const MESSAGES_STREAMS_KEY = MESSAGES_KEY + "streams:"
const DELAYED_MESSAGES_CONSUMER_RUNNING_FLAG_KEY = MESSAGES_KEY + "DELAYED-CONSUMER-RUNNING"
const DELAYED_MESSAGES_SORTED_SET_KEY = MESSAGES_KEY + "WAITING-TO-BE-PROCESSED"

func StreamName(messageType string) string {
	return MESSAGES_STREAMS_KEY + messageType
}
