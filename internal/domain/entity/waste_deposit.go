package entity

import (
	"time"

	"github.com/google/uuid"
)

type WasteDeposit struct {
	DepositId    uuid.UUID `json:"deposit_id" gorm:"type:char(36);primaryKey"`
	UserId       uuid.UUID `json:"user_id" gorm:"type:char(36);not null;index"`
	Name         string    `json:"name" gorm:"type:varchar(255);not null"`
	WasteType    string    `json:"waste_type" gorm:"type:enum('Limbah Organik Basah','Limbah Organik Kering','Limbah Campuran');not null"`
	WasteWeight  float64   `json:"waste_weight" gorm:"type:decimal(10,2);not null"`
	Reward       float64   `json:"reward" gorm:"type:decimal(10,2);not null"`
	PickupMethod string    `json:"pickup_method" gorm:"type:enum('Pick-Up','Drop-Off');not null"`
	Status       string    `json:"status" gorm:"type:enum('Pending','Scheduled','Completed');default:'Completed'"`
	PickupDate   time.Time `json:"pickup_date" gorm:"type:datetime;not null"`

	User User `gorm:"foreignKey:UserId"`
}
