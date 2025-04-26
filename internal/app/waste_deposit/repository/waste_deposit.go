package repository

import(
    "gorm.io/gorm"
)

type WasteDepositMySQLItf interface {}

type WasteDepositMySQL struct {
    db *gorm.DB
}

func NewWasteDepositMySQL(db *gorm.DB) WasteDepositMySQLItf {
    return &WasteDepositMySQL{db}
}
