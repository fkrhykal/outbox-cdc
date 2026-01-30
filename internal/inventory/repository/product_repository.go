package repository

import (
	"context"

	"github.com/fkrhykal/outbox-cdc/internal/inventory/entity"
	"github.com/google/uuid"
)

type ProductRepository interface {
	FindByIDLockForUpdate(ctx context.Context, ID uuid.UUID) (*entity.Product, error)
	UpdateStock(ctx context.Context, product *entity.Product) error
}
