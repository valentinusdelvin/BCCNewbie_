package dto

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type Register struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Password  string `json:"password" validate:"required,min=8"`
}

type Login struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type ErrorInputResponse struct {
	FieldName string `json:"fieldName"`
	Message   string `json:"message"`
}

type ProfileResponse struct {
	UserID     uuid.UUID `json:"user_id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	ProfilePic string    `json:"profile_pic,omitempty"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"first_name" validate:"omitempty,min=2"`
	LastName  string `json:"last_name" validate:"omitempty,min=2"`
}

type UpdateProfilePictureRequest struct {
	Image *multipart.FileHeader `form:"image" validate:"required"`
}
