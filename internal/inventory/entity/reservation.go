package entity

import "github.com/google/uuid"

type Reservation struct {
	ID             uuid.UUID
	ReservationKey uuid.UUID
	ProductID      uuid.UUID
	PriceLevel     int
	Quantity       int
}
