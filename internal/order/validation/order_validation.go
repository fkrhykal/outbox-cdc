package validation

import (
	"context"

	command "github.com/fkrhykal/outbox-cdc/internal/order/comand"
	"github.com/fkrhykal/outbox-cdc/internal/validation"
	"github.com/google/uuid"
)

func ValidatePlaceOrderCommand() validation.ValidatorFn[*command.PlaceOrder] {
	return func(ctx context.Context, cmd *command.PlaceOrder, errors *validation.Errors) validation.UnexpectedError {
		if cmd.ProductID == uuid.Nil {
			errors.Append(validation.Error{
				Field:   "product_id",
				Message: "product_id is required",
			})
		}
		if cmd.EstimatedPrice <= 0 {
			errors.Append(validation.Error{
				Field:   "estimated_price",
				Message: "estimated_price must be greater than zero",
			})
		}
		if cmd.Quantity <= 0 {
			errors.Append(validation.Error{
				Field:   "quantity",
				Message: "quantity must be greater than zero",
			})
		}
		return nil
	}
}
