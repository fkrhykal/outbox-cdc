package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fkrhykal/outbox-cdc/data"
	"github.com/fkrhykal/outbox-cdc/internal/inventory/command"
	"github.com/fkrhykal/outbox-cdc/internal/inventory/entity"
	"github.com/fkrhykal/outbox-cdc/internal/inventory/event"
	"github.com/fkrhykal/outbox-cdc/internal/inventory/repository"
	"github.com/fkrhykal/outbox-cdc/internal/inventory/service"
	"github.com/fkrhykal/outbox-cdc/internal/outbox"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReservationService_PlaceItemReservation_Success(t *testing.T) {
	// Setup mocks
	txManager := data.NewMockTxManager[any](t)
	productRepo := repository.NewMockProductRepository(t)
	reservationRepo := repository.NewMockReservationRepository(t)
	outboxRepo := outbox.NewMockOutboxRepository(t)
	reservationService := service.NewReservationService(txManager, productRepo, reservationRepo, outboxRepo)

	testCtx := context.Background()
	txCtx := data.NewMockTxContext[any](t)
	productID := uuid.New()
	reservationKey := uuid.New()

	cmd := &command.PlaceReservation{
		ProductID:      productID,
		ReservationKey: reservationKey,
		EstimatedPrice: 1000,
		Quantity:       2,
	}

	product := &entity.Product{
		ID:    productID,
		Name:  "Test Product",
		Stock: 10,
		Price: 1000,
	}

	// Setup expectations
	txManager.EXPECT().Begin(testCtx).Return(txCtx, nil)
	productRepo.EXPECT().FindByIDLockForUpdate(txCtx, productID).Return(product, nil)
	productRepo.EXPECT().UpdateStock(txCtx, product).Return(nil)
	reservationRepo.EXPECT().Save(txCtx, mock.AnythingOfType("*entity.Reservation")).Return(nil)
	outboxRepo.EXPECT().SaveEvent(txCtx, mock.AnythingOfType("*event.ReservationPlaced")).Return(nil)
	txCtx.EXPECT().Commit().Return(nil)
	txCtx.EXPECT().Rollback().Return(nil)

	// Execute
	err := reservationService.PlaceReservation(testCtx, cmd)

	// Assertions
	assert.NoError(t, err)

	// Verify
	txManager.AssertExpectations(t)
	productRepo.AssertExpectations(t)
	reservationRepo.AssertExpectations(t)
	outboxRepo.AssertExpectations(t)
}

func TestReservationService_PlaceItemReservation_ProductNotFound(t *testing.T) {
	// Setup mocks
	txManager := data.NewMockTxManager[any](t)
	productRepo := repository.NewMockProductRepository(t)
	reservationRepo := repository.NewMockReservationRepository(t)
	outboxRepo := outbox.NewMockOutboxRepository(t)
	reservationService := service.NewReservationService(txManager, productRepo, reservationRepo, outboxRepo)

	testCtx := context.Background()
	txCtx := data.NewMockTxContext[any](t)
	productID := uuid.New()

	cmd := &command.PlaceReservation{
		ProductID:      productID,
		ReservationKey: uuid.New(),
		EstimatedPrice: 1000,
		Quantity:       2,
	}

	// Setup expectations
	txManager.EXPECT().Begin(testCtx).Return(txCtx, nil)
	productRepo.EXPECT().FindByIDLockForUpdate(txCtx, productID).Return(nil, nil)
	outboxRepo.EXPECT().SaveEvent(txCtx, mock.AnythingOfType("*event.ProductNotFound")).Return(nil)
	txCtx.EXPECT().Rollback().Return(nil)

	// Execute
	err := reservationService.PlaceReservation(testCtx, cmd)

	// Assertions
	var productNotFound *event.ProductNotFound
	assert.ErrorAs(t, err, &productNotFound)
	assert.Equal(t, productID, productNotFound.ProductID)

	// Verify
	txManager.AssertExpectations(t)
	productRepo.AssertExpectations(t)
	outboxRepo.AssertExpectations(t)
}

func TestReservationService_PlaceItemReservation_InsuficientStock(t *testing.T) {
	// Setup mocks
	txManager := data.NewMockTxManager[any](t)
	productRepo := repository.NewMockProductRepository(t)
	reservationRepo := repository.NewMockReservationRepository(t)
	outboxRepo := outbox.NewMockOutboxRepository(t)
	reservationService := service.NewReservationService(txManager, productRepo, reservationRepo, outboxRepo)

	testCtx := context.Background()
	txCtx := data.NewMockTxContext[any](t)
	productID := uuid.New()

	cmd := &command.PlaceReservation{
		ProductID:      productID,
		ReservationKey: uuid.New(),
		EstimatedPrice: 1000,
		Quantity:       20, // More than stock
	}

	product := &entity.Product{
		ID:    productID,
		Name:  "Test Product",
		Stock: 10,
		Price: 1000,
	}

	// Setup expectations
	txManager.EXPECT().Begin(testCtx).Return(txCtx, nil)
	productRepo.EXPECT().FindByIDLockForUpdate(txCtx, productID).Return(product, nil)
	outboxRepo.EXPECT().SaveEvent(txCtx, mock.AnythingOfType("*event.InsuficientStock")).Return(nil)
	txCtx.EXPECT().Rollback().Return(nil)

	// Execute
	err := reservationService.PlaceReservation(testCtx, cmd)

	// Assertions
	var insufficientStock *event.InsuficientStock
	assert.ErrorAs(t, err, &insufficientStock)

	// Verify
	txManager.AssertExpectations(t)
	productRepo.AssertExpectations(t)
	outboxRepo.AssertExpectations(t)
}

func TestReservationService_PlaceItemReservation_MismatchedPrice(t *testing.T) {
	// Setup mocks
	txManager := data.NewMockTxManager[any](t)
	productRepo := repository.NewMockProductRepository(t)
	reservationRepo := repository.NewMockReservationRepository(t)
	outboxRepo := outbox.NewMockOutboxRepository(t)
	reservationService := service.NewReservationService(txManager, productRepo, reservationRepo, outboxRepo)

	testCtx := context.Background()
	txCtx := data.NewMockTxContext[any](t)
	productID := uuid.New()

	cmd := &command.PlaceReservation{
		ProductID:      productID,
		ReservationKey: uuid.New(),
		EstimatedPrice: 500, // Different from product price
		Quantity:       2,
	}

	product := &entity.Product{
		ID:    productID,
		Name:  "Test Product",
		Stock: 10,
		Price: 1000,
	}

	// Setup expectations
	txManager.EXPECT().Begin(testCtx).Return(txCtx, nil)
	productRepo.EXPECT().FindByIDLockForUpdate(txCtx, productID).Return(product, nil)
	outboxRepo.EXPECT().SaveEvent(txCtx, mock.AnythingOfType("*event.MismatchedPrice")).Return(nil)
	txCtx.EXPECT().Rollback().Return(nil)

	// Execute
	err := reservationService.PlaceReservation(testCtx, cmd)

	// Assertions
	var mismatchedPrice *event.MismatchedPrice
	assert.ErrorAs(t, err, &mismatchedPrice)

	// Verify
	txManager.AssertExpectations(t)
	productRepo.AssertExpectations(t)
	outboxRepo.AssertExpectations(t)
}

func TestReservationService_PlaceItemReservation_TxBeginFailure(t *testing.T) {
	// Setup mocks
	txManager := data.NewMockTxManager[any](t)
	productRepo := repository.NewMockProductRepository(t)
	reservationRepo := repository.NewMockReservationRepository(t)
	outboxRepo := outbox.NewMockOutboxRepository(t)
	reservationService := service.NewReservationService(txManager, productRepo, reservationRepo, outboxRepo)

	testCtx := context.Background()
	cmd := &command.PlaceReservation{
		ProductID: uuid.New(),
	}

	// Setup expectations
	txManager.EXPECT().Begin(testCtx).Return(nil, errors.New("tx begin failed"))

	// Execute
	err := reservationService.PlaceReservation(testCtx, cmd)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to begin transaction")

	// Verify
	txManager.AssertExpectations(t)
}
