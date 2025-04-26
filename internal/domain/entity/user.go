package entity

import "github.com/google/uuid"

type User struct {
	UserId     uuid.UUID `json:"user_id" gorm:"type:char(36);primaryKey"`
	Email      string    `json:"email" gorm:"type:varchar(255);unique"`
	FirstName  string    `json:"first_name" gorm:"type:varchar(255)"`
	LastName   string    `json:"last_name" gorm:"type:varchar(255)"`
	Password   string    `json:"password" gorm:"type:varchar(255)"`
	ProfilePic string    `json:"profile_pic,omitempty" gorm:"type:text"`
}
