package usecase

import (
	"errors"
	"fmt"
	"hackfest-uc/internal/app/user/repository"
	"hackfest-uc/internal/domain/dto"
	"hackfest-uc/internal/domain/entity"
	"hackfest-uc/internal/infra/jwt"
	"hackfest-uc/internal/infra/supabase"
	"hackfest-uc/internal/validation"
	"mime/multipart"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecaseItf interface {
	Register(register dto.Register) (entity.User, error)
	Login(login dto.Login) (string, error)
	GetProfile(userId uuid.UUID) (dto.ProfileResponse, error)
	UpdateProfile(userId uuid.UUID, req dto.UpdateProfileRequest) (dto.ProfileResponse, error)
	UpdateProfilePicture(userID uuid.UUID, file *multipart.FileHeader) (dto.ProfileResponse, error)
}

type UserUsecase struct {
	userRepo  repository.UserMySQLItf
	validator validation.InputValidation
	jwt       jwt.JWT
	supabase  supabase.SupabaseItf
}

func NewUserUsecase(userRepo repository.UserMySQLItf, jwt jwt.JWT, validator validation.InputValidation, supabase supabase.SupabaseItf) UserUsecaseItf {
	return &UserUsecase{
		userRepo:  userRepo,
		validator: validator,
		jwt:       jwt,
		supabase:  supabase,
	}
}

func (u UserUsecase) Register(register dto.Register) (entity.User, error) {
	var user entity.User

	validationErrors := u.validator.Validate(register)
	if len(validationErrors) > 0 {
		return entity.User{}, fmt.Errorf("validation failed: %v", validationErrors)
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

func (u UserUsecase) Login(login dto.Login) (string, error) {
	var user entity.User
	validationErrors := u.validator.Validate(login)
	if len(validationErrors) > 0 {
		return "", fmt.Errorf("validation failed: %v", validationErrors)
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

func (u UserUsecase) GetProfile(userId uuid.UUID) (dto.ProfileResponse, error) {
	user, err := u.userRepo.GetById(userId)
	if err != nil {
		return dto.ProfileResponse{}, err
	}

	return dto.ProfileResponse{
		UserID:     user.UserId,
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		ProfilePic: user.ProfilePic,
	}, nil
}

func (u UserUsecase) UpdateProfile(userId uuid.UUID, req dto.UpdateProfileRequest) (dto.ProfileResponse, error) {
	user, err := u.userRepo.GetById(userId)
	if err != nil {
		return dto.ProfileResponse{}, err
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}

	if err := u.userRepo.Update(user); err != nil {
		return dto.ProfileResponse{}, err
	}

	return u.GetProfile(userId)
}

func (u UserUsecase) UpdateProfilePicture(userID uuid.UUID, file *multipart.FileHeader) (dto.ProfileResponse, error) {
	imageURL, err := u.supabase.Upload(file)
	if err != nil {
		return dto.ProfileResponse{}, err
	}

	if err := u.userRepo.UpdateProfilePicture(userID, imageURL); err != nil {
		return dto.ProfileResponse{}, err
	}

	return u.GetProfile(userID)
}
