package entity

import (
	"time"

	"github.com/google/uuid"
)

type WasteDeposit struct {
	DepositId    uuid.UUID `json:"deposit_id" gorm:"type:char(36);primaryKey"`
	UserId       uuid.UUID `json:"user_id" gorm:"type:char(36);not null;index"`
	Name         string    `json:"name" gorm:"type:varchar(255)"`
	WasteType    string    `json:"waste_type" gorm:"type:enum('Limbah Organik Basah','Limbah Organik Kering','Limbah Campuran')"`
	WasteWeight  float64   `json:"weight" gorm:"type:decimal(10,2)"`
	Reward       float64   `json:"reward" gorm:"type:decimal(10,2)"`
	PickupMethod string    `json:"pickup_method" gorm:"type:enum('Pick-Up','Drop-Off')"`
	Status       string    `json:"status" gorm:"type:enum('Pending','Scheduled','Completed');default:'Pending'"`
	PickupDate   time.Time `json:"pickup_date" gorm:"type:datetime"`

	User User `gorm:"foreignKey:UserId"`
}
