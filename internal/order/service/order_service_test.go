package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fkrhykal/outbox-cdc/data"
	command "github.com/fkrhykal/outbox-cdc/internal/order/comand"
	"github.com/fkrhykal/outbox-cdc/internal/order/repository"
	"github.com/fkrhykal/outbox-cdc/internal/order/service"
	"github.com/fkrhykal/outbox-cdc/internal/outbox"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrderService_PlaceOrder_Success(t *testing.T) {
	// Setup mocks
	txManager := data.NewMockTxManager[any](t)
	orderRepository := repository.NewMockOrderRepository(t)
	outboxRepository := outbox.NewMockOutboxRepository(t)
	orderService := service.NewOrderService(txManager, orderRepository, outboxRepository)

	testCtx := context.Background()
	txCtx := data.NewMockTxContext[any](t)
	productID := uuid.New()

	cmd := &command.PlaceOrder{
		ProductID:      productID,
		EstimatedPrice: 1000,
		Quantity:       2,
	}

	// Setup expectations for the mocks
	txManager.EXPECT().Begin(testCtx).Return(txCtx, nil)
	txCtx.EXPECT().Commit().Return(nil)
	txCtx.EXPECT().Rollback().Return(nil)

	orderRepository.EXPECT().Save(txCtx, mock.AnythingOfType("*entity.Order")).Return(nil)
	outboxRepository.EXPECT().Save(txCtx, mock.AnythingOfType("*outbox.Outbox")).Return(nil)

	// Execute the test
	placedOrder, err := orderService.PlaceOrder(testCtx, cmd)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, placedOrder)

	// Verify all mock expectations were met
	txManager.AssertExpectations(t)
	txCtx.AssertExpectations(t)
	orderRepository.AssertExpectations(t)
	outboxRepository.AssertExpectations(t)
}

func TestOrderService_PlaceOrder_TransactionBeginFailure(t *testing.T) {
	// Setup mocks
	txManager := data.NewMockTxManager[any](t)
	orderRepository := repository.NewMockOrderRepository(t)
	outboxRepository := outbox.NewMockOutboxRepository(t)
	orderService := service.NewOrderService(txManager, orderRepository, outboxRepository)

	testCtx := context.Background()
	cmd := &command.PlaceOrder{
		ProductID:      uuid.New(),
		EstimatedPrice: 1000,
		Quantity:       2,
	}

	// Setup expectations - transaction begin fails
	txManager.EXPECT().Begin(testCtx).Return(nil, errors.New("transaction begin failed"))

	// Execute the test
	placedOrder, err := orderService.PlaceOrder(testCtx, cmd)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, placedOrder)
	assert.Contains(t, err.Error(), "failed to start transaction")

	// Verify only the transaction manager was called
	txManager.AssertExpectations(t)
	orderRepository.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	outboxRepository.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
}

func TestOrderService_PlaceOrder_OrderSaveFailure(t *testing.T) {
	// Setup mocks
	txManager := data.NewMockTxManager[any](t)
	orderRepository := repository.NewMockOrderRepository(t)
	outboxRepository := outbox.NewMockOutboxRepository(t)
	orderService := service.NewOrderService(txManager, orderRepository, outboxRepository)

	testCtx := context.Background()
	txCtx := data.NewMockTxContext[any](t)
	cmd := &command.PlaceOrder{
		ProductID:      uuid.New(),
		EstimatedPrice: 1000,
		Quantity:       2,
	}

	// Setup expectations
	txManager.EXPECT().Begin(testCtx).Return(txCtx, nil)
	orderRepository.EXPECT().Save(txCtx, mock.AnythingOfType("*entity.Order")).Return(errors.New("order save failed"))
	txCtx.EXPECT().Rollback().Return(nil)

	// Execute the test
	placedOrder, err := orderService.PlaceOrder(testCtx, cmd)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, placedOrder)
	assert.Contains(t, err.Error(), "failed to save order")

	// Verify expectations
	txManager.AssertExpectations(t)
	txCtx.AssertExpectations(t)
	orderRepository.AssertExpectations(t)
	outboxRepository.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
}

func TestOrderService_PlaceOrder_EventMappingFailure(t *testing.T) {
	// Setup mocks
	txManager := data.NewMockTxManager[any](t)
	orderRepository := repository.NewMockOrderRepository(t)
	outboxRepository := outbox.NewMockOutboxRepository(t)
	orderService := service.NewOrderService(txManager, orderRepository, outboxRepository)

	testCtx := context.Background()
	txCtx := data.NewMockTxContext[any](t)
	cmd := &command.PlaceOrder{
		ProductID:      uuid.New(),
		EstimatedPrice: 1000,
		Quantity:       2,
	}

	// Setup expectations - event mapping fails
	txManager.EXPECT().Begin(testCtx).Return(txCtx, nil)
	orderRepository.EXPECT().Save(txCtx, mock.AnythingOfType("*entity.Order")).Return(nil)
	outboxRepository.EXPECT().Save(txCtx, mock.AnythingOfType("*outbox.Outbox")).Return(errors.New("event mapping failed"))
	txCtx.EXPECT().Rollback().Return(nil)

	// Execute the test
	placedOrder, err := orderService.PlaceOrder(testCtx, cmd)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, placedOrder)
	assert.Contains(t, err.Error(), "failed to save place order event")

	// Verify expectations
	txManager.AssertExpectations(t)
	txCtx.AssertExpectations(t)
	orderRepository.AssertExpectations(t)
	outboxRepository.AssertExpectations(t)
}

func TestOrderService_PlaceOrder_TransactionCommitFailure(t *testing.T) {
	// Setup mocks
	txManager := data.NewMockTxManager[any](t)
	orderRepository := repository.NewMockOrderRepository(t)
	outboxRepository := outbox.NewMockOutboxRepository(t)
	orderService := service.NewOrderService(txManager, orderRepository, outboxRepository)

	testCtx := context.Background()
	txCtx := data.NewMockTxContext[any](t)
	cmd := &command.PlaceOrder{
		ProductID:      uuid.New(),
		EstimatedPrice: 1000,
		Quantity:       2,
	}

	// Setup expectations
	txManager.EXPECT().Begin(testCtx).Return(txCtx, nil)
	orderRepository.EXPECT().Save(txCtx, mock.AnythingOfType("*entity.Order")).Return(nil)
	outboxRepository.EXPECT().Save(txCtx, mock.AnythingOfType("*outbox.Outbox")).Return(nil)
	txCtx.EXPECT().Commit().Return(errors.New("commit failed"))
	txCtx.EXPECT().Rollback().Return(nil)

	// Execute the test
	placedOrder, err := orderService.PlaceOrder(testCtx, cmd)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, placedOrder)
	assert.Contains(t, err.Error(), "failed to commit transaction")

	// Verify expectations
	txManager.AssertExpectations(t)
	txCtx.AssertExpectations(t)
	orderRepository.AssertExpectations(t)
	outboxRepository.AssertExpectations(t)
}
