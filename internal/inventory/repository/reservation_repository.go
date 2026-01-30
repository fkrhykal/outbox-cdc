package repository

import (
	"context"

	"github.com/fkrhykal/outbox-cdc/internal/inventory/entity"
)

type ReservationRepository interface {
	Save(ctx context.Context, itemReservation *entity.Reservation) error
}
