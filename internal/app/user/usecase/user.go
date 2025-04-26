package usecase

import (
	"errors"
	"fmt"
	"hackfest-uc/internal/app/user/repository"
	"hackfest-uc/internal/domain/dto"
	"hackfest-uc/internal/domain/entity"
	"hackfest-uc/internal/infra/jwt"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecaseItf interface {
	Register(register dto.Register) (entity.User, error)
}

type UserUsecase struct {
	userRepo repository.UserMySQLItf
	validate *validator.Validate
	jwt      jwt.JWT
}

func NewUserUsecase(userRepo repository.UserMySQLItf, jwt jwt.JWT) UserUsecaseItf {
	return &UserUsecase{
		userRepo: userRepo,
		validate: validator.New(),
		jwt:      jwt,
	}
}

func (u UserUsecase) Register(register dto.Register) (entity.User, error) {
	var user entity.User

	if err := u.validate.Struct(register); err != nil {
		return entity.User{}, fmt.Errorf("validation error: %w", err)
	}

	if _, err := u.userRepo.FindByEmail(register.Email); err == nil {
		return entity.User{}, errors.New("Email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		return entity.User{}, err
	}

	user = entity.User{
		UserId:    uuid.New(),
		Email:     register.Email,
		FirstName: register.FirstName,
		LastName:  register.LastName,
		Password:  string(hashedPassword),
	}

	err = u.userRepo.Create(user)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}
