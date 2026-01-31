package query

import (
	"context"

	"github.com/google/uuid"
)

type ProductQueryResult struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Price int       `json:"price"`
	Stock int       `json:"stock"`
}

type GetProductQueryHandler interface {
	GetProducts(ctx context.Context) ([]ProductQueryResult, error)
}
