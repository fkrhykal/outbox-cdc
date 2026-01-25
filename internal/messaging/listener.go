package messaging

import (
	"context"
)

type EventHandler[E Event] func(ctx context.Context, event E) error

type EventListener[E Event] interface {
	Listen(ctx context.Context, handler EventHandler[E]) error
}
