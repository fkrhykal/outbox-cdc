package service

import (
	"context"
	"fmt"

	"github.com/fkrhykal/outbox-cdc/data"
	"github.com/fkrhykal/outbox-cdc/internal/inventory/command"
	"github.com/fkrhykal/outbox-cdc/internal/inventory/event"
	"github.com/fkrhykal/outbox-cdc/internal/inventory/repository"
	"github.com/fkrhykal/outbox-cdc/internal/outbox"
	"github.com/google/uuid"
)

var _ command.PlaceReservationHandler = (*ReservationService[any])(nil)

type ReservationService[T any] struct {
	txManager             data.TxManager[T]
	productRepository     repository.ProductRepository
	reservationRepository repository.ReservationRepository
	outboxRepository      outbox.OutboxRepository
}

func NewReservationService[T any](
	txManager data.TxManager[T],
	productRepository repository.ProductRepository,
	reservationRepository repository.ReservationRepository,
	outboxRepository outbox.OutboxRepository,
) *ReservationService[T] {
	return &ReservationService[T]{
		txManager:             txManager,
		productRepository:     productRepository,
		reservationRepository: reservationRepository,
		outboxRepository:      outboxRepository,
	}
}

// PlaceProductReservation implements [command.PlaceReservationHandler].
func (rs *ReservationService[T]) PlaceReservation(ctx context.Context, cmd *command.PlaceReservation) error {
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
		productNotFound := &event.ProductNotFound{
			ID:             uuid.New(),
			ReservationKey: cmd.ReservationKey,
			ProductID:      cmd.ProductID,
		}
		if err := rs.outboxRepository.SaveEvent(txCtx, productNotFound); err != nil {
			return err
		}
		return productNotFound
	}

	reservation, reservationFailed := product.Reserve(cmd.ReservationKey, cmd.EstimatedPrice, cmd.Quantity)
	if reservationFailed != nil {
		if err := rs.outboxRepository.SaveEvent(txCtx, reservationFailed); err != nil {
			return fmt.Errorf("failed to save failure event: %w", err)
		}
		return reservationFailed
	}

	if err := rs.productRepository.UpdateStock(txCtx, product); err != nil {
		return fmt.Errorf("failed to update product stock: %w", err)
	}

	if err := rs.reservationRepository.Save(txCtx, reservation); err != nil {
		return fmt.Errorf("failed to save item reservation: %w", err)
	}

	reservationPlaced := &event.ReservationPlaced{
		ID:            uuid.New(),
		ProductID:     reservation.ProductID,
		ReservationID: reservation.ID,
		PriceLevel:    reservation.PriceLevel,
		Quantity:      reservation.Quantity,
	}

	if err := rs.outboxRepository.SaveEvent(txCtx, reservationPlaced); err != nil {
		return fmt.Errorf("failed to save reservation placed event: %w", err)
	}

	if err := txCtx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
