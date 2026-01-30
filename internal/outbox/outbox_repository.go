package outbox

import (
	"context"

	"github.com/fkrhykal/outbox-cdc/internal/messaging"
)

type OutboxRepository interface {
	Save(ctx context.Context, outbox *Outbox) error
	SaveEvent(ctx context.Context, event messaging.Event) error
}
