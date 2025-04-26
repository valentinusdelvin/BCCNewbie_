package entity

import "github.com/google/uuid"

type Payment struct {
	OrderID     uuid.UUID
	UserID      uuid.UUID `gorm:"foreignkey:ID;references:users"`
	ProductID   uuid.UUID
	ProductName string
	Amount      uint64
	Price       uint64
	FinalPrice  uint64
	PaymentLink string
	Status      string
}
