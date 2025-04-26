package rest

import (
	"hackfest-uc/internal/app/user/usecase"
	"hackfest-uc/internal/middleware"

	"github.com/go-playground/validator"
)

type UserHandler struct {
	usecase    usecase.UserUsecaseItf
	validator  validator.Validate
	middleware middleware.MiddlewareI
}

func NewUserHandler() {
	UserHandler := UserHandler{}
}
