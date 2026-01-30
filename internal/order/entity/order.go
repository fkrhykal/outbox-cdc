package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID             uuid.UUID
	ProductID      uuid.UUID
	Quantity       int
	EstimatedPrice int
	PlacedAt       time.Time
}
