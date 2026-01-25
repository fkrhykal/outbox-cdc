package command

import (
	"context"

	"github.com/google/uuid"
)

type PlaceItemReservation struct {
	ReservationKey uuid.UUID
	ProductID      uuid.UUID
	Quantity       int
	EstimatedPrice int
}

type PlaceItemReservationHandler interface {
	PlaceItemReservation(ctx context.Context, command *PlaceItemReservation) error
}
