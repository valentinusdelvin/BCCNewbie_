package repository

import (
	"hackfest-uc/internal/domain/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserMySQLItf interface {
	Create(user entity.User) error
	FindByEmail(email string) (entity.User, error)
	FindByUsername(username string) (entity.User, error)
	GetById(id uuid.UUID) (entity.User, error)
	Update(user entity.User) error
	UpdateProfilePicture(userID uuid.UUID, imageURL string) error
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

func (r UserMySQL) GetById(id uuid.UUID) (entity.User, error) {
	var user entity.User
	err := r.db.First(&user, "user_id = ?", id).Error
	return user, err
}

func (r UserMySQL) Update(user entity.User) error {
	return r.db.Model(user).Updates(map[string]interface{}{
		"first_name":  user.FirstName,
		"last_name":   user.LastName,
		"profile_pic": user.ProfilePic,
	}).Error
}

func (r UserMySQL) UpdateProfilePicture(userId uuid.UUID, imageURL string) error {
	return r.db.Model(&entity.User{}).
		Where("user_id = ?", userId).
		Update("profile_pic", imageURL).Error
}
