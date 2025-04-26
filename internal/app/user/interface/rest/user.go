package rest

import (
	"fmt"
	"hackfest-uc/internal/app/user/usecase"
	"hackfest-uc/internal/domain/dto"
	"hackfest-uc/internal/infra/supabase"
	"hackfest-uc/internal/middleware"
	"hackfest-uc/internal/validation"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	usecase    usecase.UserUsecaseItf
	validator  validation.InputValidation
	middleware middleware.MiddlewareItf
	supabase   supabase.SupabaseItf
}

func NewUserHandler(routerGroup fiber.Router, validator validation.InputValidation, userUsecase usecase.UserUsecaseItf, middleware middleware.MiddlewareItf, supabase supabase.SupabaseItf) {
	UserHandler := UserHandler{
		usecase:    userUsecase,
		validator:  validator,
		middleware: middleware,
		supabase:   supabase,
	}

	routerGroup.Post("/register", UserHandler.Register)
	routerGroup.Post("/login", UserHandler.Login)

	userProtected := routerGroup.Group("/users", middleware.Authentication)
	userProtected.Get("/", UserHandler.GetProfile)
	userProtected.Patch("/", UserHandler.UpdateProfile)
	userProtected.Patch("/picture", UserHandler.UpdateProfilePicture)
}

func (h *UserHandler) Register(ctx *fiber.Ctx) error {
	var register dto.Register

	if err := ctx.BodyParser(&register); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"success": false,
				"message": "Invalid request",
			})
	}

	if err := h.validator.Validate(register); err != nil {
		log.Printf("Validation error: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
		})
	}

	_, err := h.usecase.Register(register)
	if err != nil {
		errorMap := fiber.Map{"general": err.Error()}
		status := fiber.StatusBadRequest

		if strings.Contains(err.Error(), "email already exists") {
			status = fiber.StatusConflict
			errorMap = fiber.Map{
				"email": "Email already exists",
			}
		}

		return ctx.Status(status).JSON(fiber.Map{
			"success": false,
			"message": "Registration failed",
			"errors":  errorMap,
		})
	}
	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "User registered successfully",
	})
}

func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	var login dto.Login

	if err := ctx.BodyParser(&login); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"success": false,
				"message": "Invalid request",
			})
	}

	if err := h.validator.Validate(login); err != nil {
		log.Printf("Validation error: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
		})
	}

	token, err := h.usecase.Login(login)
	if err != nil {
		errorMap := fiber.Map{"general": err.Error()}
		status := fiber.StatusBadRequest

		if strings.Contains(err.Error(), "invalid email or password") {
			status = fiber.StatusUnauthorized
			errorMap = fiber.Map{
				"email":    "Invalid email or password",
				"password": "Invalid email or password",
			}
		}

		return ctx.Status(status).JSON(fiber.Map{
			"success": false,
			"message": "Login failed",
			"errors":  errorMap,
		})
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    token,
	})
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("userId").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	profile, err := h.usecase.GetProfile(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get profile",
		})
	}

	return c.JSON(profile)
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("userId").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req dto.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validator.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed",
		})
	}

	if err := h.validator.Validator.Struct(req); err != nil {
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			field := strings.ToLower(err.Field())
			validationErrors[field] = fmt.Sprintf("%s is invalid", field)
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Validation failed",
			"details": validationErrors,
		})
	}

	updatedProfile, err := h.usecase.UpdateProfile(userID, req)
	if err != nil {
		log.Printf("Failed to update profile: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to update profile",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Profile updated successfully",
		"data":    updatedProfile,
	})
}

func (h *UserHandler) UpdateProfilePicture(c *fiber.Ctx) error {
	userID, ok := c.Locals("userId").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid user ID",
		})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Image file is required",
		})
	}

	if file.Size > 2<<20 { // 2MB limit
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Image too large (max 2MB)",
		})
	}

	ext := filepath.Ext(file.Filename)
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	if !allowedExtensions[strings.ToLower(ext)] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Only JPG/JPEG/PNG files are allowed",
		})
	}

	if err := h.usecase.UpdateProfilePicture(userID, file); err != nil {
		log.Printf("Failed to update profile picture: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to update profile picture",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Profile picture updated successfully",
	})
}
