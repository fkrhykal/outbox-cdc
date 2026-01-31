package service

import (
	"context"
	"fmt"

	"github.com/fkrhykal/outbox-cdc/internal/inventory/entity"
	"github.com/fkrhykal/outbox-cdc/internal/inventory/query"
	"github.com/fkrhykal/outbox-cdc/internal/inventory/repository"
	"github.com/fkrhykal/outbox-cdc/internal/utils"
)

var _ query.GetProductQueryHandler = (*ProductService)(nil)

type ProductService struct {
	productRepositry repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepositry: productRepository,
	}
}

// GetProducts implements [query.GetProductQueryHandler].
func (p *ProductService) GetProducts(ctx context.Context) ([]query.ProductQueryResult, error) {
	products, err := p.productRepositry.FindWithLimit(ctx, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	return utils.Map(products, func(p entity.Product) query.ProductQueryResult {
		return query.ProductQueryResult{
			ID:    p.ID,
			Name:  p.Name,
			Price: p.Price,
			Stock: p.Stock,
		}
	}), nil
}
