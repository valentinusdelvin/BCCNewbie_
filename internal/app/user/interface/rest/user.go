package rest

import (
	"hackfest-uc/internal/app/user/usecase"
	"hackfest-uc/internal/domain/dto"
	"hackfest-uc/internal/middleware"
	"hackfest-uc/internal/validation"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	usecase    usecase.UserUsecaseItf
	validator  validation.InputValidation
	middleware middleware.MiddlewareItf
}

func NewUserHandler(routerGroup fiber.Router, validator validation.InputValidation, userUsecase usecase.UserUsecaseItf, middleware middleware.MiddlewareItf) {
	UserHandler := UserHandler{
		usecase:    userUsecase,
		validator:  validator,
		middleware: middleware,
	}

	routerGroup = routerGroup.Group("/users")
	routerGroup.Post("/register", UserHandler.Register)
	routerGroup.Post("/login", UserHandler.Login)

	protectedGroup := routerGroup.Group("/users", middleware.Authentication)
	protectedGroup.Get("/", UserHandler.GetProfile)
	protectedGroup.Put("/", UserHandler.UpdateProfile)
	protectedGroup.Put("/picture", UserHandler.UpdateProfilePicture)
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

	profile, err := h.usecase.UpdateProfile(userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile",
		})
	}

	return c.JSON(profile)
}

func (h UserHandler) UpdateProfilePicture(c *fiber.Ctx) error {
	userID, ok := c.Locals("userId").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Get file from form data
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Image is required",
		})
	}

	// Validate file size (max 2MB)
	if file.Size > 2<<20 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Image too large (max 2MB)",
		})
	}

	profile, err := h.usecase.UpdateProfilePicture(userID, file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update profile picture",
			"details": err.Error(),
		})
	}

	return c.JSON(profile)
}
