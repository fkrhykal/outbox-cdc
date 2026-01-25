package inbox

import (
	"context"
	"fmt"

	"github.com/fkrhykal/outbox-cdc/data"
	"github.com/fkrhykal/outbox-cdc/internal/messaging"
)

type InboxHandlerConfig[T any] struct {
	ProcessKey       string
	InboxPersistence InboxPersistence
	TxManager        data.TxManager[T]
}

func WithTransactionInboxHandler[E messaging.Event, T any](
	handler messaging.EventHandler[E],
	config *InboxHandlerConfig[T],
) messaging.EventHandler[E] {
	return func(ctx context.Context, event E) error {
		txCtx, err := config.TxManager.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to initiate transaction: %w", err)
		}
		defer txCtx.Rollback()
		exist, err := config.InboxPersistence.ExistByID(txCtx, event.EventID())
		if err != nil {
			return fmt.Errorf("failed to find inbox message: %w", err)
		}
		if exist {
			return nil
		}
		if err := handler(txCtx, event); err != nil {
			return fmt.Errorf("internal inbox transaction handler error: %w", err)
		}
		if err := txCtx.Commit(); err != nil {
			return fmt.Errorf("failed to commit database transaction: %w", err)
		}
		return nil
	}
}
