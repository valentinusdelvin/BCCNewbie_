package entity

import "github.com/google/uuid"

type User struct {
	UserId   uuid.UUID `json:"user_id" gorm:"type:char(36);primaryKey"`
	Email    string    `json:"email" gorm:"type:varchar(255);unique"`
	Username string    `json:"username" gorm:"type:varchar(255);unique"`
	Password string    `json:"password" gorm:"type:varchar(255);unique"`
}
