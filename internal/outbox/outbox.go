package outbox

import (
	"encoding/json"
	"fmt"

	"github.com/fkrhykal/outbox-cdc/internal/messaging"
)

type Outbox struct {
	EventID       string
	AggregateID   string
	AggregateType string
	EventType     string
	Payload       []byte
}

func Event[E messaging.Event](event E) (*Outbox, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event %s: %w", event.EventType(), err)
	}
	return &Outbox{
		EventID:       event.EventID(),
		AggregateID:   event.AggregateID(),
		AggregateType: event.AggregateType(),
		EventType:     event.EventType(),
		Payload:       payload,
	}, nil
}
