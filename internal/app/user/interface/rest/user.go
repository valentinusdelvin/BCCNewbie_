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
		"message": "User registered succesfully",
	})

}
