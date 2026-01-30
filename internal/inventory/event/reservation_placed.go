package event

import (
	"github.com/fkrhykal/outbox-cdc/internal/messaging"
	"github.com/google/uuid"
)

var _ messaging.Event = (*ReservationPlaced)(nil)

type ReservationPlaced struct {
	ID            uuid.UUID `json:"event_id"`
	ProductID     uuid.UUID `json:"product_id"`
	ReservationID uuid.UUID `json:"reservation_id"`
	PriceLevel    int       `json:"price_level"`
	Quantity      int       `json:"quantity"`
}

// AggregateID implements [messaging.Event].
func (i ReservationPlaced) AggregateID() string {
	return i.ReservationID.String()
}

// AggregateType implements [messaging.Event].
func (i ReservationPlaced) AggregateType() string {
	return "reservation"
}

// EventID implements [messaging.Event].
func (i ReservationPlaced) EventID() string {
	return i.ID.String()
}

// EventType implements [messaging.Event].
func (i ReservationPlaced) EventType() string {
	return "reservation.placed"
}
