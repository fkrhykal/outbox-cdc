package entity

import "github.com/google/uuid"

type ItemReservation struct {
	ID        uuid.UUID
	ProductID uuid.UUID
	Quantity  int
}
