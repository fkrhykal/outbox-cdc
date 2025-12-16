package service

import (
	"context"
	"time"

	"github.com/fkrhykal/outbox-cdc/data"
	"github.com/fkrhykal/outbox-cdc/internal/command"
	"github.com/fkrhykal/outbox-cdc/internal/entity"
	"github.com/fkrhykal/outbox-cdc/internal/event"
	"github.com/fkrhykal/outbox-cdc/internal/messaging"
	"github.com/fkrhykal/outbox-cdc/internal/repository"
	"github.com/google/uuid"
)

var _ command.PlaceOrderHandler = (*OrderService[any])(nil)

type OrderService[T any] struct {
	txManager      data.TxManager[T]
	orderRepostory repository.OrderRepository
	publisher      messaging.EventPublisher[event.Event]
}

func NewOrderService[T any](
	txManager data.TxManager[T],
	orderRepository repository.OrderRepository,
	publisher messaging.EventPublisher[event.Event],
) *OrderService[T] {
	return &OrderService[T]{
		txManager:      txManager,
		orderRepostory: orderRepository,
		publisher:      publisher,
	}
}

// PlaceOrder implements command.PlaceOrderHandler.
func (o *OrderService[T]) PlaceOrder(ctx context.Context, cmd *command.PlaceOrder) (*command.PlacedOrder, error) {
	txCtx, err := o.txManager.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer txCtx.Rollback()

	order := &entity.Order{
		ID:             uuid.New(),
		ItemID:         cmd.ItemID,
		Quantity:       cmd.Quantity,
		EstimatedPrice: cmd.EstimatedPrice,
		PlacedAt:       time.Now(),
	}
	if err := o.orderRepostory.Save(txCtx, order); err != nil {
		return nil, err
	}

	placedOrder := event.NewOrderPlaced(order)
	if err := o.publisher.Publish(txCtx, placedOrder); err != nil {
		return nil, err
	}

	if err := txCtx.Commit(); err != nil {
		return nil, err
	}

	return &command.PlacedOrder{ID: order.ID}, nil
}
