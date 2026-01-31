package entity

import (
	"testing"

	"github.com/fkrhykal/outbox-cdc/internal/inventory/event"
	"github.com/fkrhykal/outbox-cdc/internal/messaging"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProduct_Reserve_Success(t *testing.T) {
	product := &Product{
		ID:    uuid.New(),
		Name:  "Test Product",
		Stock: 10,
		Price: 100,
	}

	reservationKey := uuid.New()
	estimatedPrice := 100
	quantity := 3

	reservation, err := product.Reserve(reservationKey, estimatedPrice, quantity)

	assert.NoError(t, err)
	assert.NotNil(t, reservation)
	assert.Equal(t, reservationKey, reservation.ID)
	assert.Equal(t, product.ID, reservation.ProductID)
	assert.Equal(t, quantity, reservation.Quantity)
	assert.Equal(t, 7, product.Stock) // Stock should be reduced
}

func TestProduct_Reserve_MismatchedPrice(t *testing.T) {
	product := &Product{
		ID:    uuid.New(),
		Name:  "Test Product",
		Stock: 10,
		Price: 100,
	}

	reservationKey := uuid.New()
	estimatedPrice := 150 // Different from product price
	quantity := 3

	reservation, err := product.Reserve(reservationKey, estimatedPrice, quantity)

	assert.Error(t, err)
	assert.Nil(t, reservation)
	assert.IsType(t, &event.MismatchedPrice{}, err)
	assert.Implements(t, (*messaging.Event)(nil), err)

	mismatchedPriceErr := err.(*event.MismatchedPrice)
	assert.Equal(t, reservationKey, mismatchedPriceErr.ReservationKey)
	assert.Equal(t, product.ID, mismatchedPriceErr.ProductID)
	assert.Equal(t, product.Price, mismatchedPriceErr.ActualPrice)
	assert.Equal(t, estimatedPrice, mismatchedPriceErr.EstimedPrice)

	// Stock should not be reduced
	assert.Equal(t, 10, product.Stock)
}

func TestProduct_Reserve_OutOfStock(t *testing.T) {
	product := &Product{
		ID:    uuid.New(),
		Name:  "Test Product",
		Stock: 2,
		Price: 100,
	}

	reservationKey := uuid.New()
	estimatedPrice := 100
	quantity := 5 // More than available stock

	reservation, err := product.Reserve(reservationKey, estimatedPrice, quantity)

	assert.Error(t, err)
	assert.Nil(t, reservation)
	assert.IsType(t, &event.InsuficientStock{}, err)
	assert.Implements(t, (*messaging.Event)(nil), err)

	outOfStockErr := err.(*event.InsuficientStock)
	assert.Equal(t, reservationKey, outOfStockErr.ReservationKey)
	assert.Equal(t, product.ID, outOfStockErr.ProductID)
	assert.Equal(t, product.Stock, outOfStockErr.AvailableStock)
	assert.Equal(t, quantity, outOfStockErr.RequestedQuantity)

	// Stock should not be reduced
	assert.Equal(t, 2, product.Stock)
}

func TestMismatchedPrice_Error(t *testing.T) {
	err := &event.MismatchedPrice{
		ID:             uuid.New(),
		ReservationKey: uuid.New(),
		ProductID:      uuid.New(),
		ActualPrice:    100,
		EstimedPrice:   150,
	}

	expected := "mismatched price: expected 150 found 100"
	assert.Equal(t, expected, err.Error())
}

func TestInsuficientStock_Error(t *testing.T) {
	err := &event.InsuficientStock{
		ID:                uuid.New(),
		ReservationKey:    uuid.New(),
		ProductID:         uuid.New(),
		AvailableStock:    2,
		RequestedQuantity: 5,
	}

	expected := "insuficient stock: requested 5 but only 2 available"
	assert.Equal(t, expected, err.Error())
}

func TestMismatchedPrice_EventImplementation(t *testing.T) {
	productID := uuid.New()
	eventID := uuid.New()
	reservationKey := uuid.New()

	event := &event.MismatchedPrice{
		ID:             eventID,
		ReservationKey: reservationKey,
		ProductID:      productID,
		ActualPrice:    100,
		EstimedPrice:   150,
	}

	assert.Equal(t, productID.String(), event.AggregateID())
	assert.Equal(t, "product", event.AggregateType())
	assert.Equal(t, eventID.String(), event.EventID())
	assert.Equal(t, "mismatched_price", event.EventType())
}

func TestInsuficientStock_EventImplementation(t *testing.T) {
	productID := uuid.New()
	eventID := uuid.New()
	reservationKey := uuid.New()

	event := &event.InsuficientStock{
		ID:                eventID,
		ReservationKey:    reservationKey,
		ProductID:         productID,
		AvailableStock:    2,
		RequestedQuantity: 5,
	}

	assert.Equal(t, productID.String(), event.AggregateID())
	assert.Equal(t, "product", event.AggregateType())
	assert.Equal(t, eventID.String(), event.EventID())
	assert.Equal(t, "insuficient_stock", event.EventType())
}
