package kafka

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/fkrhykal/outbox-cdc/internal/messaging"
	"github.com/segmentio/kafka-go"
)

var _ messaging.EventListener[messaging.Event] = (*KafkaDebeziumListener[messaging.Event])(nil)

type Payload struct {
	Payload   string `json:"payload"`
	EventType string `json:"event_type"`
}

type Envelop struct {
	Payload Payload `json:"payload"`
}

type KafkaDebeziumListener[E messaging.Event] struct {
	r         *kafka.Reader
	eventType string
}

func NewKafkaOutboxListener[E messaging.Event](
	reader *kafka.Reader,
) *KafkaDebeziumListener[E] {
	var event E
	return &KafkaDebeziumListener[E]{
		r:         reader,
		eventType: event.EventType(),
	}
}

// Listen implements [messaging.EventListener].
func (k *KafkaDebeziumListener[E]) Listen(ctx context.Context, handler messaging.EventHandler[E]) error {
	for {
		message, err := k.r.FetchMessage(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "failed to fetch message", "error", err)
			continue
		}

		var envelop Envelop
		if err := json.Unmarshal(message.Value, &envelop); err != nil {
			slog.ErrorContext(ctx, "failed to parse message value", "error", err)
			continue
		}

		if envelop.Payload.EventType != k.eventType {
			slog.WarnContext(ctx, "event type mismatch: expected %s found %s", k.eventType, envelop.Payload.EventType)
			continue
		}

		var event E
		if err := json.Unmarshal([]byte(envelop.Payload.Payload), &event); err != nil {
			slog.ErrorContext(ctx, "failed to parse event", "error", err)
			continue
		}

		if err := handler(ctx, event); err != nil {
			slog.ErrorContext(ctx, "failed to handle message", "error", err)
			continue
		}

		if err := k.r.CommitMessages(ctx, message); err != nil {
			slog.ErrorContext(ctx, "failed to commit message", "error", err)
		}
	}
}
