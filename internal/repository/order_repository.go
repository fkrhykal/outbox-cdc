package repository

import (
	"context"

	"github.com/fkrhykal/outbox-cdc/internal/entity"
)

type OrderRepository interface {
	Save(ctx context.Context, order *entity.Order) error
}
