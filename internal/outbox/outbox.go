package outbox

type Outbox struct {
	EventID       string
	AggregateID   string
	AggregateType string
	EventType     string
	Payload       []byte
}
