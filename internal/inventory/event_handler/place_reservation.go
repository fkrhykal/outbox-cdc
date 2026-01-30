package eventhandler

import (
	"context"

	"github.com/fkrhykal/outbox-cdc/internal/inventory/command"
	"github.com/fkrhykal/outbox-cdc/internal/messaging"
	"github.com/fkrhykal/outbox-cdc/internal/order/event"
)

func PlaceReservationHandler(handler command.PlaceReservationHandler) messaging.EventHandler[event.OrderPlaced] {
	return func(ctx context.Context, event event.OrderPlaced) error {
		err := handler.PlaceReservation(ctx, &command.PlaceReservation{
			ReservationKey: event.OrderID,
			ProductID:      event.ProductID,
			Quantity:       event.Quantity,
			EstimatedPrice: event.EstimatedPrice,
		})
		if _, ok := err.(messaging.FailureEvent); ok {
			return nil
		}
		return err
	}
}
