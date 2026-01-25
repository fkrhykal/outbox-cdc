package event

import (
	"time"

	"github.com/fkrhykal/outbox-cdc/internal/messaging"
	"github.com/fkrhykal/outbox-cdc/internal/order/entity"
	"github.com/google/uuid"
)

var _ messaging.Event = (*OrderPlaced)(nil)

type OrderPlaced struct {
	ID             uuid.UUID `json:"event_id"`
	OrderID        uuid.UUID `json:"order_id"`
	ItemID         int       `json:"item_id"`
	Quantity       int       `json:"quantity"`
	EstimatedPrice int       `json:"estimated_price"`
	PlacedAt       time.Time `json:"placed_at"`
}

func NewOrderPlaced(o *entity.Order) *OrderPlaced {
	return &OrderPlaced{
		ID:             uuid.New(),
		OrderID:        o.ID,
		ItemID:         o.ItemID,
		Quantity:       o.Quantity,
		EstimatedPrice: o.EstimatedPrice,
		PlacedAt:       o.PlacedAt,
	}
}

// AggregateID implements Event.
func (o OrderPlaced) AggregateID() string {
	return o.OrderID.String()
}

// AggregateType implements Event.
func (o OrderPlaced) AggregateType() string {
	return "order"
}

// EventID implements Event.
func (o OrderPlaced) EventID() string {
	return o.ID.String()
}

// EventType implements Event.
func (o OrderPlaced) EventType() string {
	return "order.placed"
}
