package validation

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name string
}

func TestValidatorRegistry(t *testing.T) {
	registry := &ValidatorRegistry{}

	t.Run("success validation", func(t *testing.T) {
		validator := func(ctx context.Context, s TestStruct, validationErrors *Errors) UnexpectedError {
			if s.Name == "" {
				validationErrors.Append(Error{
					Field:   "name",
					Message: "name cannot be empty",
				})
			}
			return nil
		}

		registry.Register(Coerce(validator))

		err := registry.Validate(context.Background(), TestStruct{Name: "valid name"})
		assert.NoError(t, err)
	})

	t.Run("failure validation", func(t *testing.T) {
		err := registry.Validate(context.Background(), TestStruct{Name: ""})
		assert.Error(t, err)

		validationErrors, ok := err.(Errors)
		assert.True(t, ok, "error should be of type ValidationErrors")
		assert.Len(t, validationErrors, 1)
		assert.Equal(t, "name", validationErrors[0].Field)
		assert.Equal(t, "name cannot be empty", validationErrors[0].Message)

		assert.Contains(t, err.Error(), `[{"name":"name","message":"name cannot be empty","meta":null}]`)
	})

	t.Run("unsupported type", func(t *testing.T) {
		err := registry.Validate(context.Background(), "some string")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported value: no validation for type")
	})
}
