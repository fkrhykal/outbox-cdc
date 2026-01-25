package repository

import (
	"context"

	"github.com/fkrhykal/outbox-cdc/internal/inventory/entity"
)

type ItemReservationRepository interface {
	Save(ctx context.Context, itemReservation *entity.ItemReservation) error
}
