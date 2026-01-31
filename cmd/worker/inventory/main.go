package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/fkrhykal/outbox-cdc/internal/order/event"
	messaging "github.com/fkrhykal/outbox-cdc/kafka"
	"github.com/segmentio/kafka-go"
)

func main() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:29092"},
		Topic:   "order",
		GroupID: "log",
	})

	kafkaListener := messaging.NewKafkaOutboxListener[event.OrderPlaced](reader)

	kafkaListener.Listen(context.Background(), func(ctx context.Context, event event.OrderPlaced) error {
		res, err := json.Marshal(event)
		if err != nil {
			return err
		}
		log.Printf("Received %s: %s\n", event.EventType(), string(res))
		return nil
	})
}
