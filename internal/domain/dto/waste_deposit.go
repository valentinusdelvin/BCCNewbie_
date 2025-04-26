package dto

import (
	"time"

	"github.com/google/uuid"
)

const (
	WasteTypeOrganicWet = "Limbah Organik Basah"
	WasteTypeOrganicDry = "Limbah Organik Kering"
	WasteTypeMixed      = "Limbah Campuran"
	PickupMethodPickup  = "Pick-Up"
	PickupMethodDropoff = "Drop-Off"
)

type DepositRequest struct {
	Name         string  `json:"name" validate:"required"`
	WasteType    string  `json:"waste_type" validate:"required,oneof=Limbah Organik Basah Limbah Organik Kering Limbah Campuran"`
	WasteWeight  float64 `json:"weight" validate:"required,gt=0"`
	PickupMethod string  `json:"pickup_method" validate:"required,oneof=Pick-Up Drop-Off"`
}

type DepositResponse struct {
	DepositId    uuid.UUID `json:"deposit_id"`
	Name         string    `json:"name"`
	WasteType    string    `json:"waste_type"`
	WasteWeight  float64   `json:"waste_weight"`
	Reward       float64   `json:"reward"`
	PickupMethod string    `json:"pickup_method"`
	Status       string    `json:"status"`
	PickupDate   time.Time `json:"pickup_date"`
}
