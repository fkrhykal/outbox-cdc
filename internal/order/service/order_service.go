package service

import (
	"context"
	"fmt"
	"time"

	"github.com/fkrhykal/outbox-cdc/data"
	command "github.com/fkrhykal/outbox-cdc/internal/order/comand"
	"github.com/fkrhykal/outbox-cdc/internal/order/entity"
	"github.com/fkrhykal/outbox-cdc/internal/order/event"
	"github.com/fkrhykal/outbox-cdc/internal/order/repository"
	"github.com/fkrhykal/outbox-cdc/internal/outbox"
	"github.com/fkrhykal/outbox-cdc/internal/validation"

	"github.com/google/uuid"
)

var _ command.PlaceOrderHandler = (*OrderService[any])(nil)

type OrderService[T any] struct {
	validator         validation.Validator
	txManager         data.TxManager[T]
	orderRepository   repository.OrderRepository
	outboxPersistence outbox.OutboxRepository
}

func NewOrderService[T any](
	validator validation.Validator,
	txManager data.TxManager[T],
	orderRepository repository.OrderRepository,
	outboxRepository outbox.OutboxRepository,
) *OrderService[T] {
	return &OrderService[T]{
		validator:         validator,
		txManager:         txManager,
		orderRepository:   orderRepository,
		outboxPersistence: outboxRepository,
	}
}

// PlaceOrder creates a new order and publishes an OrderPlaced event.
// It uses the outbox pattern to ensure that the event is published if and only if the order is saved.
func (o *OrderService[T]) PlaceOrder(ctx context.Context, cmd *command.PlaceOrder) (*command.PlacedOrder, error) {
	if err := o.validator.Validate(ctx, cmd); err != nil {
		return nil, err
	}
	txCtx, err := o.txManager.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer txCtx.Rollback()

	order := &entity.Order{
		ID:             uuid.New(),
		ProductID:      cmd.ProductID,
		Quantity:       cmd.Quantity,
		EstimatedPrice: cmd.EstimatedPrice,
		PlacedAt:       time.Now(),
	}
	if err := o.orderRepository.Save(txCtx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	orderPlaced := event.NewOrderPlaced(order)
	outboxRecord, err := outbox.Event(orderPlaced)
	if err != nil {
		return nil, fmt.Errorf("failed to map event to outbox: %w", err)
	}

	if err := o.outboxPersistence.Save(txCtx, outboxRecord); err != nil {
		return nil, fmt.Errorf("failed to save place order event: %w", err)
	}

	if err := txCtx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &command.PlacedOrder{ID: order.ID}, nil
}
