package command

import (
	"context"

	"github.com/google/uuid"
)

type PlaceReservation struct {
	ReservationKey uuid.UUID
	ProductID      uuid.UUID
	Quantity       int
	EstimatedPrice int
}

type PlaceReservationHandler interface {
	PlaceReservation(ctx context.Context, command *PlaceReservation) error
}
