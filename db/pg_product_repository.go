package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fkrhykal/outbox-cdc/internal/inventory/entity"
	"github.com/fkrhykal/outbox-cdc/internal/inventory/repository"
	"github.com/google/uuid"
)

var _ repository.ProductRepository = (*PgProductRepository)(nil)

type PgProductRepository struct {
	pgRepository
}

func NewPgProductRepository(db *sql.DB) *PgProductRepository {
	return &PgProductRepository{
		pgRepository: pgRepository{
			db: db,
		},
	}
}

// FindByIDLockForUpdate implements [repository.ProductRepository].
func (p *PgProductRepository) FindByIDLockForUpdate(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	query := `SELECT id, name, price, stock FROM products WHERE id = $1 FOR UPDATE`
	product := new(entity.Product)
	err := p.Executor(ctx).
		QueryRowContext(ctx, query, id).
		Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Stock,
		)
	if err == nil {
		return product, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return nil, err
}

// FindWithLimit implements [repository.ProductRepository].
func (p *PgProductRepository) FindWithLimit(ctx context.Context, limit int) ([]entity.Product, error) {
	query := `SELECT id, name, price, stock FROM products LIMIT $1`
	products := make([]entity.Product, 0, limit)
	rows, err := p.Executor(ctx).QueryContext(ctx, query, limit)
	if errors.Is(err, sql.ErrNoRows) {
		return products, nil
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		product := entity.Product{}
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

// UpdateStock implements [repository.ProductRepository].
func (p *PgProductRepository) UpdateStock(ctx context.Context, product *entity.Product) error {
	query := `UPDATE products SET stock = $1 WHERE id = $1`
	resut, err := p.Executor(ctx).ExecContext(ctx, query, product.Stock, product.ID)
	if err != nil {
		return err
	}
	affectedRow, err := resut.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRow != 1 {
		return fmt.Errorf("affected row mismatch: expected 1 found %d", affectedRow)
	}
	return nil
}
