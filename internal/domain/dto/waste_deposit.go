package dto

import (
	"time"

	"github.com/google/uuid"
)

type DepositRequest struct {
	Name         string  `json:"name" validate:"required"`
	WasteType    string  `json:"waste_type" validate:"required"`
	WasteWeight  float64 `json:"waste_weight" validate:"required"`
	PickupMethod string  `json:"pickup_method" validate:"required,oneof=pickup dropoff"`
}

type DepositResponse struct {
	DepositId    uuid.UUID `json:"deposit_id"`
	Name         string    `json:"name"`
	WasteType    string    `json:"waste_type"`
	WasteWeight  float64   `json:"weight"`
	Reward       float64   `json:"reward"`
	PickupMethod string    `json:"pickup_method"`
	Status       string    `json:"status"`
	PickupDate   time.Time `json:"pickup_date,omitempty"`
}
