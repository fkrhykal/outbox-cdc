package service

import (
	"context"
	"fmt"

	"github.com/fkrhykal/outbox-cdc/data"
	"github.com/fkrhykal/outbox-cdc/internal/inventory/command"
	"github.com/fkrhykal/outbox-cdc/internal/inventory/entity"
	"github.com/fkrhykal/outbox-cdc/internal/inventory/event"
	"github.com/fkrhykal/outbox-cdc/internal/inventory/repository"
	"github.com/fkrhykal/outbox-cdc/internal/messaging"
)

var _ command.PlaceItemReservationHandler = (*ReservationService[any])(nil)

type ReservationService[T any] struct {
	txManager                 data.TxManager[T]
	productRepository         repository.ProductRepository
	itemReservationRepository repository.ItemReservationRepository
	publisher                 messaging.EventPublisher[messaging.Event]
}

// PlaceItemReservation implements [command.PlaceItemReservationHandler].
func (rs *ReservationService[T]) PlaceItemReservation(ctx context.Context, cmd *command.PlaceItemReservation) error {
	txCtx, err := rs.txManager.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer txCtx.Rollback()

	product, err := rs.productRepository.FindByIDLockForUpdate(txCtx, cmd.ProductID)
	if err != nil {
		return fmt.Errorf("failed to retrieve product with id %s: %w", cmd.ProductID, err)
	}
	if product == nil {
		return fmt.Errorf("product with id %s doesn't exist: %w", cmd.ProductID, err)
	}

	if product.Price < cmd.EstimatedPrice {
		return fmt.Errorf("")
	}
	if product.Stock < cmd.Quantity {
		return fmt.Errorf("")
	}

	product.Stock -= cmd.Quantity
	if err := rs.productRepository.UpdateStock(txCtx, product.ID, product.Stock); err != nil {
		return fmt.Errorf("failed to update product stock: %w", err)
	}

	itemReservation := &entity.ItemReservation{
		ID:        cmd.ReservationKey,
		ProductID: product.ID,
		Quantity:  cmd.Quantity,
	}
	if err := rs.itemReservationRepository.Save(txCtx, itemReservation); err != nil {
		return fmt.Errorf("failed to save item reservation: %w", err)
	}

	itemReserved := event.ItemReserved{}
	if err := rs.publisher.Publish(txCtx, &itemReserved); err != nil {
		return fmt.Errorf("failed to published item reserved event: %w", err)
	}

	if err := txCtx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
