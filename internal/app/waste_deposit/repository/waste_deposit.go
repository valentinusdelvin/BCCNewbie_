package repository

import (
	"hackfest-uc/internal/domain/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WasteDepositMySQLItf interface {
	Create(deposit entity.WasteDeposit) error
	GetByUserId(userId uuid.UUID) ([]entity.WasteDeposit, error)
}

type WasteDepositMySQL struct {
	db *gorm.DB
}

func NewWasteDepositMySQL(db *gorm.DB) WasteDepositMySQLItf {
	return &WasteDepositMySQL{db}
}

func (r WasteDepositMySQL) Create(deposit entity.WasteDeposit) error {
	if err := r.db.Create(&deposit).Error; err != nil {
		return err
	}
	return nil
}

func (r WasteDepositMySQL) GetByUserId(userId uuid.UUID) ([]entity.WasteDeposit, error) {
	var deposits []entity.WasteDeposit
	err := r.db.Where("user_id = ?", userId).Find(&deposits).Error
	return deposits, err
}
