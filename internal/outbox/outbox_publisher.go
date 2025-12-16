package outbox

import (
	"context"
	"encoding/json"

	"github.com/fkrhykal/outbox-cdc/internal/event"
	"github.com/fkrhykal/outbox-cdc/internal/messaging"
)

var _ messaging.EventPublisher[event.Event] = (*OutboxEventPublisher)(nil)

type OutboxEventPublisher struct {
	outboxPersistence OutboxPersistence
}

func NewOutboxEventPublisher(outboxPersistence OutboxPersistence) *OutboxEventPublisher {
	return &OutboxEventPublisher{outboxPersistence: outboxPersistence}
}

// Publish implements messaging.EventPublisher.
func (o *OutboxEventPublisher) Publish(ctx context.Context, event event.Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	outbox := &Outbox{
		EventID:       event.EventID(),
		AggregateID:   event.AggregateID(),
		AggregateType: event.AggregateType(),
		EventType:     event.EventType(),
		Payload:       payload,
	}
	return o.outboxPersistence.Save(ctx, outbox)
}
