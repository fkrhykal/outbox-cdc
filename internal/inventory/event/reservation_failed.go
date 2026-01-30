package event

import (
	"fmt"

	"github.com/fkrhykal/outbox-cdc/internal/messaging"
	"github.com/google/uuid"
)

var _ messaging.FailureEvent = (*ProductNotFound)(nil)

type ProductNotFound struct {
	ID             uuid.UUID `json:"event_id"`
	ReservationKey uuid.UUID `json:"reservation_key"`
	ProductID      uuid.UUID `json:"product_id"`
}

// Error implements [error].
func (p *ProductNotFound) Error() string {
	return fmt.Sprintf("product with id %s not found", p.ProductID)
}

// AggregateID implements [messaging.FailureEvent].
func (p *ProductNotFound) AggregateID() string {
	return p.ProductID.String()
}

// AggregateType implements [messaging.FailureEvent].
func (p *ProductNotFound) AggregateType() string {
	return "product"
}

// EventID implements [messaging.FailureEvent].
func (p *ProductNotFound) EventID() string {
	return p.ID.String()
}

// EventType implements [messaging.FailureEvent].
func (p *ProductNotFound) EventType() string {
	return "product_not_found"
}

var _ messaging.FailureEvent = (*MismatchedPrice)(nil)

type MismatchedPrice struct {
	ID             uuid.UUID `json:"event_id"`
	ReservationKey uuid.UUID `json:"reservation_key"`
	ProductID      uuid.UUID `json:"prodcut_id"`
	ActualPrice    int       `json:"actual_price"`
	EstimedPrice   int       `json:"estimated_price"`
}

// AggregateID implements [messaging.FailureEvent].
func (m *MismatchedPrice) AggregateID() string {
	return m.ProductID.String()
}

// AggregateType implements [messaging.FailureEvent].
func (m *MismatchedPrice) AggregateType() string {
	return "product"
}

// EventID implements [messaging.FailureEvent].
func (m *MismatchedPrice) EventID() string {
	return m.ID.String()
}

// EventType implements [messaging.FailureEvent].
func (m *MismatchedPrice) EventType() string {
	return "mismatched_price"
}

// Error implements [error].
func (m *MismatchedPrice) Error() string {
	return fmt.Sprintf("mismatched price: expected %d found %d", m.EstimedPrice, m.ActualPrice)
}

var _ messaging.FailureEvent = (*InsuficientStock)(nil)

type InsuficientStock struct {
	ID                uuid.UUID `json:"event_id"`
	ReservationKey    uuid.UUID `json:"reservation_key"`
	ProductID         uuid.UUID `json:"product_id"`
	AvailableStock    int       `json:"available_stock"`
	RequestedQuantity int       `json:"requested_quantity"`
}

// AggregateID implements [messaging.FailureEvent].
func (o *InsuficientStock) AggregateID() string {
	return o.ProductID.String()
}

// AggregateType implements [messaging.FailureEvent].
func (o *InsuficientStock) AggregateType() string {
	return "product"
}

// EventID implements [messaging.FailureEvent].
func (o *InsuficientStock) EventID() string {
	return o.ID.String()
}

// EventType implements [messaging.FailureEvent].
func (o *InsuficientStock) EventType() string {
	return "insuficient_stock"
}

// Error implements [error].
func (o *InsuficientStock) Error() string {
	return fmt.Sprintf("insuficient stock: requested %d but only %d available", o.RequestedQuantity, o.AvailableStock)
}
