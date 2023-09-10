package comment

import (
	"fmt"

	event "github.com/AntonioMartinezFernandez/golang-redis-streams/pkg/event"
)

type commentHandler struct {
}

// NewCommentHandler ...
func NewCommentHandler() event.Handler {
	return &commentHandler{}
}

func (h *commentHandler) Handle(e event.Event, retry bool) error {
	event, ok := e.(*CommentEvent)

	if !ok {
		return fmt.Errorf("incorrect event type")
	}

	//! Force handler to fail to be able to check the pending events consumer
	// if event.UserID == 5 && !retry {
	// 	return errors.New("comment handler fails ---development purposes---")
	// }

	fmt.Printf("processed event %+v UserID: %v Comment:%v \n", event, event.UserID, event.Comment)

	return nil
}
