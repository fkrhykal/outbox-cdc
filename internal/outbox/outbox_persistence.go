package outbox

import "context"

type OutboxPersistence interface {
	Save(ctx context.Context, outbox *Outbox) error
}
