package validation

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
)

type Error struct {
	Field   string            `json:"name"`
	Message string            `json:"message"`
	Meta    map[string]string `json:"meta"`
}

type Errors []Error

func (e *Errors) Append(err Error) {
	*e = append(*e, err)
}

func (e Errors) IsError() bool {
	return len(e) > 0
}

func (e Errors) Error() string {
	message, _ := json.Marshal(e)
	return string(message)
}

type Validator interface {
	Validate(ctx context.Context, v any) error
}

var _ Validator = (*ValidatorRegistry)(nil)

type UnexpectedError interface {
	error
}

type ValidatorFn[T any] func(ctx context.Context, v T, errors *Errors) UnexpectedError

type ValidatorRegistry struct {
	registry map[reflect.Type]ValidatorFn[any]
}

func NewValidatorRegistry() *ValidatorRegistry {
	return &ValidatorRegistry{
		registry: make(map[reflect.Type]ValidatorFn[any]),
	}
}

// Validate implements [Validator].
func (r *ValidatorRegistry) Validate(ctx context.Context, v any) error {
	t := reflect.TypeOf(v)
	if r.registry == nil {
		r.registry = make(map[reflect.Type]ValidatorFn[any])
	}
	validate, ok := r.registry[t]
	if !ok {
		return fmt.Errorf("unsupported value: no validation for type %q", t)
	}

	validationErrors := make(Errors, 0)

	if err := validate(ctx, v, &validationErrors); err != nil {
		return fmt.Errorf("failed to validate value %q: %w", t, err)
	}

	if validationErrors.IsError() {
		return validationErrors
	}

	return nil
}

func (r *ValidatorRegistry) Register(t reflect.Type, validator ValidatorFn[any]) {
	if r.registry == nil {
		r.registry = make(map[reflect.Type]ValidatorFn[any])
	}
	r.registry[t] = validator
}

func Coerce[T any](validator ValidatorFn[T]) (reflect.Type, ValidatorFn[any]) {
	return reflect.TypeFor[T](), func(ctx context.Context, v any, errors *Errors) UnexpectedError {
		return validator(ctx, v.(T), errors)
	}
}
