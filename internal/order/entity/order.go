package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID             uuid.UUID
	ItemID         int
	Quantity       int
	EstimatedPrice int
	PlacedAt       time.Time
}
