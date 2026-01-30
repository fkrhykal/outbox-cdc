package entity

import (
	"github.com/fkrhykal/outbox-cdc/internal/inventory/event"
	"github.com/fkrhykal/outbox-cdc/internal/messaging"

	"github.com/google/uuid"
)

type Product struct {
	ID    uuid.UUID
	Name  string
	Stock int
	Price int
}

func (p *Product) Reserve(
	reservationKey uuid.UUID,
	estimatedPrice int,
	quantity int,
) (*Reservation, messaging.FailureEvent) {
	if p.Price != estimatedPrice {
		return nil, &event.MismatchedPrice{
			ID:             uuid.New(),
			ReservationKey: reservationKey,
			EstimedPrice:   estimatedPrice,
			ProductID:      p.ID,
			ActualPrice:    p.Price,
		}
	}
	if p.Stock < quantity {
		return nil, &event.InsuficientStock{
			ID:                uuid.New(),
			ReservationKey:    reservationKey,
			ProductID:         p.ID,
			AvailableStock:    p.Stock,
			RequestedQuantity: quantity,
		}
	}
	p.Stock -= quantity
	return &Reservation{
		ID:        reservationKey,
		ProductID: p.ID,
		Quantity:  quantity,
	}, nil
}
