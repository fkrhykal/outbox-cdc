package event

import "github.com/fkrhykal/outbox-cdc/internal/messaging"

var _ messaging.Event = (*ItemReserved)(nil)

type ItemReserved struct {
}

// AggregateID implements [messaging.Event].
func (i ItemReserved) AggregateID() string {
	panic("unimplemented")
}

// AggregateType implements [messaging.Event].
func (i ItemReserved) AggregateType() string {
	panic("unimplemented")
}

// EventID implements [messaging.Event].
func (i ItemReserved) EventID() string {
	panic("unimplemented")
}

// EventType implements [messaging.Event].
func (i ItemReserved) EventType() string {
	panic("unimplemented")
}
