package rest

import (
	"hackfest-uc/internal/app/user/usecase"
	"hackfest-uc/internal/domain/dto"
	"hackfest-uc/internal/middleware"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	usecase    usecase.UserUsecaseItf
	validator  *validator.Validate
	middleware middleware.MiddlewareItf
}

func NewUserHandler(routerGroup fiber.Router, userUsecase usecase.UserUsecaseItf, validator validator.Validate, middleware middleware.MiddlewareItf) {
	UserHandler := UserHandler{
		usecase:    userUsecase,
		validator:  &validator,
		middleware: middleware,
	}

	routerGroup = routerGroup.Group("/users")
	routerGroup.Post("/register", UserHandler.Register)
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

	if err := h.validator.Struct(register); err != nil {
		log.Printf("Validation error: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
		})
	}

	createdUser, err := h.usecase.Register(register)
	if err != nil {
		errorMap := fiber.Map{"general": err.Error()}
		status := fiber.StatusBadRequest

		if strings.Contains(err.Error(), "email already exists") {
			status = fiber.StatusConflict
			errorMap = fiber.Map{
				"email": "Email already exists",
			}
		}
		if strings.Contains(err.Error(), "username already exists") {
			status = fiber.StatusConflict
			errorMap = fiber.Map{
				"username": "Username already exists",
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
		"data":    createdUser,
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

	if err := h.validator.Struct(login); err != nil {
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
