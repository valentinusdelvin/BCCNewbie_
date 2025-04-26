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
	Login(login dto.Login) (string, error)
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
		FirstName: register.FirstName,
		LastName:  register.LastName,
		Email:     register.Email,
		Password:  string(hashedPassword),
	}

	err = u.userRepo.Create(user)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (u UserUsecase) Login(login dto.Login) (string, error) {
	var user entity.User

	if err := u.validate.Struct(login); err != nil {
		return "", fmt.Errorf("validation error: %w", err)
	}

	user, err := u.userRepo.FindByEmail(login.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := u.jwt.GenerateToken(user.UserId)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}
