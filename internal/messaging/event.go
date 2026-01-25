package messaging

type Event interface {
	EventID() string
	AggregateID() string
	AggregateType() string
	EventType() string
}
