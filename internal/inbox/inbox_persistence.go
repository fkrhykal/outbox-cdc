package inbox

import (
	"context"
)

type InboxPersistence interface {
	ExistByID(ctx context.Context, ID string) (bool, error)
	FindById(ctx context.Context, ID string) (*InboxMessage, error)
	Save(ctx context.Context, inbox *InboxMessage) error
}
