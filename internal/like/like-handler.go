package like

import (
	"fmt"

	event "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/event"
)

type likeHandler struct {
}

// NewLikeHandler ...
func NewLikeHandler() event.Handler {
	return &likeHandler{}
}

func (h *likeHandler) Handle(e event.Event, retry bool) error {
	event, ok := e.(*LikeEvent)

	if !ok {
		return fmt.Errorf("incorrect event type")
	}

	//! Force handler to fail to be able to check the pending events consumer
	// if event.UserID == 5 && !retry {
	// 	return errors.New("like handler fails ---development purposes---")
	// }

	fmt.Printf("completed like %+v UserID: %v\n", event, event.UserID)

	return nil
}
