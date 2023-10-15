package redis_streams

import "encoding/json"

type DelayedMessage struct {
	Message     string `json:"message"`
	MessageType string `json:"type"`
}

func NewDelayedMessage(messageAsString string, messageType string) *DelayedMessage {
	return &DelayedMessage{
		Message:     messageAsString,
		MessageType: messageType,
	}
}

func NewDelayedMessageFromString(delayedMessageAsString string) (*DelayedMessage, error) {
	var delayedMessage DelayedMessage
	if err := json.Unmarshal([]byte(delayedMessageAsString), &delayedMessage); err != nil {
		return nil, err
	}
	return &delayedMessage, nil
}

func (dm *DelayedMessage) ToString() string {
	delayedMessageString, _ := json.Marshal(dm)
	return string(delayedMessageString)
}
