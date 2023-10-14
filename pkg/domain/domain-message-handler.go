package domain

type DomainMessageHandler interface {
	MessageType() string
	NewDomainMessageFromMap(map[string]interface{}) (DomainMessage, error)
	Handle(e interface{}) error
}
