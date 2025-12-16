package messaging

import (
	"context"

	"github.com/fkrhykal/outbox-cdc/internal/event"
)

type EventPublisher[E event.Event] interface {
	Publish(ctx context.Context, event E) error
}
