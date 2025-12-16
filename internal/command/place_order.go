package command

import (
	"context"

	"github.com/google/uuid"
)

type PlaceOrder struct {
	ItemID         int
	EstimatedPrice int
	Quantity       int
}

type PlacedOrder struct {
	ID uuid.UUID
}

type PlaceOrderHandler interface {
	PlaceOrder(ctx context.Context, cmd *PlaceOrder) (*PlacedOrder, error)
}
