package repository

import (
	"context"

	"github.com/fkrhykal/outbox-cdc/internal/order/entity"
)

type OrderRepository interface {
	Save(ctx context.Context, order *entity.Order) error
}
