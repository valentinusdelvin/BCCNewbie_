package repository

import (
	"hackfest-uc/internal/domain/entity"

	"gorm.io/gorm"
)

type UserMySQLItf interface {
	Create(user entity.User) error
	FindByEmail(email string) (entity.User, error)
	FindByUsername(username string) (entity.User, error)
}

type UserMySQL struct {
	db *gorm.DB
}

func NewUserMySQL(db *gorm.DB) UserMySQLItf {
	return &UserMySQL{db}
}

func (r UserMySQL) Create(user entity.User) error {
	return r.db.Create(user).Error
}

func (r UserMySQL) FindByEmail(email string) (entity.User, error) {
	var user entity.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (r UserMySQL) FindByUsername(username string) (entity.User, error) {
	var user entity.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
