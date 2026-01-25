package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fkrhykal/outbox-cdc/internal/order/entity"
	"github.com/fkrhykal/outbox-cdc/internal/order/repository"
)

var _ repository.OrderRepository = (*PgOrderRepository)(nil)

type PgOrderRepository struct {
	pgRepository
}

func NewPgOrderRepository(db *sql.DB) *PgOrderRepository {
	return &PgOrderRepository{
		pgRepository: pgRepository{
			db: db,
		},
	}
}

// Save implements repository.OrderRepository.
func (p *PgOrderRepository) Save(ctx context.Context, order *entity.Order) error {
	query := `
		INSERT INTO orders(id, item_id, quantity, estimated_price, placed_at)
		VALUES($1, $2, $3, $4, $5)
	`
	_, err := p.Executor(ctx).
		ExecContext(ctx, query, order.ID, order.ItemID, order.Quantity, order.EstimatedPrice, order.PlacedAt)
	if err != nil {
		return fmt.Errorf("failed to insert order record: %w", err)
	}
	return nil
}
