package messaging

import (
	"context"
)

type EventPublisher[E Event] interface {
	Publish(ctx context.Context, event E) error
}
