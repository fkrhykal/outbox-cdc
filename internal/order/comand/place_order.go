package command

import (
	"context"

	"github.com/google/uuid"
)

type PlaceOrder struct {
	ProductID      uuid.UUID `json:"product_id"`
	EstimatedPrice int       `json:"estimated_price"`
	Quantity       int       `json:"quantity"`
}

type PlacedOrder struct {
	ID uuid.UUID `json:"id"`
}

type PlaceOrderHandler interface {
	PlaceOrder(ctx context.Context, cmd *PlaceOrder) (*PlacedOrder, error)
}
