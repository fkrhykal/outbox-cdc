package outbox

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fkrhykal/outbox-cdc/internal/messaging"
)

var _ messaging.EventPublisher[messaging.Event] = (*OutboxEventPublisher)(nil)

type OutboxEventPublisher struct {
	outboxPersistence OutboxRepository
}

func NewOutboxEventPublisher(outboxPersistence OutboxRepository) *OutboxEventPublisher {
	return &OutboxEventPublisher{outboxPersistence: outboxPersistence}
}

// Publish implements messaging.EventPublisher.
func (o *OutboxEventPublisher) Publish(ctx context.Context, event messaging.Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed marshal event %s: %w", event.EventType(), err)
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
