package dto

import "github.com/google/uuid"

type Payment struct {
	OrderID   uuid.UUID
	UserID    uuid.UUID
	ProductID uuid.UUID
}
