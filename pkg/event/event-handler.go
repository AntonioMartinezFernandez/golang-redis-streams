package event

import (
	"fmt"
)

type Handler interface {
	Handle(e Event, retry bool) error
}

type defaultHandler struct {
}

// NewViewHandler ...
func NewDefaultHandler() Handler {
	return &defaultHandler{}
}

func (h *defaultHandler) Handle(e Event, retry bool) error {
	fmt.Printf("undefined event %+v\n", e)
	return nil
}
