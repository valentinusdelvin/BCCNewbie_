package usecase

import (
	"hackfest-uc/internal/app/user/repository"
	"hackfest-uc/internal/infra/jwt"
)

type UserUsecaseItf interface{}

type UserUsecase struct {
	userRepo repository.UserMySQLItf
	jwt      jwt.JWT
}

func NewUserUsecase(userRepo repository.UserMySQLItf, jwt jwt.JWT) UserUsecaseItf {
	return &UserUsecase{
		userRepo: userRepo,
		jwt:      jwt,
	}
}
