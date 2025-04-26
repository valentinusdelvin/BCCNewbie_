package repository

import (
	"hackfest-uc/internal/domain/dto"
	"hackfest-uc/internal/domain/entity"

	"gorm.io/gorm"
)

type InterPaymentRepository interface {
	CreatePayment(payment entity.Payment) error
	UpdatePaymentStatus(tx *gorm.DB, status string, orderID string) error
	GetInvoice(orderID string) (dto.Payment, error)
}

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) InterPaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

func (pr *PaymentRepository) CreatePayment(payment entity.Payment) error {
	err := pr.db.Create(&payment).Error
	if err != nil {
		return err
	}
	return nil
}

func (pr *PaymentRepository) UpdatePaymentStatus(tx *gorm.DB, status string, orderID string) error {
	return tx.Model(&entity.Payment{}).
		Where("order_id = ?", orderID).
		Update("status", status).
		Error
}

func (pr *PaymentRepository) GetInvoice(orderID string) (dto.Payment, error) {
	var invoice dto.Payment
	err := pr.db.Model(entity.Payment{}).Where("order_id = ?", orderID).First(&invoice).Error
	return invoice, err
}
