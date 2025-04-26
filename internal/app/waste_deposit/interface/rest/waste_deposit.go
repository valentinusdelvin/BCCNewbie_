package rest

import (
	"hackfest-uc/internal/app/waste_deposit/usecase"
	"hackfest-uc/internal/domain/dto"
	"hackfest-uc/internal/middleware"
	"hackfest-uc/internal/validation"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WasteDepositHandler struct {
	wasteDepositUsecase usecase.WasteDepositUsecaseItf
	validator           validation.InputValidation
	middleware          middleware.MiddlewareItf
}

func NewWasteDepositHandler(routerGroup fiber.Router, wasteDepositUsecase usecase.WasteDepositUsecaseItf, middleware middleware.MiddlewareItf, validator validation.InputValidation) {
	WasteDepositHandler := WasteDepositHandler{
		wasteDepositUsecase: wasteDepositUsecase,
		middleware:          middleware,
		validator:           validator,
	}

	routerGroup = routerGroup.Group("/deposits", middleware.Authentication)
	routerGroup.Post("/", WasteDepositHandler.CreateDeposit)
	routerGroup.Get("/history", WasteDepositHandler.GetUserDepositHistory)
	routerGroup.Get("/reward", WasteDepositHandler.GetUserReward)

}

func (h WasteDepositHandler) CreateDeposit(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(uuid.UUID)
	if !ok {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Invalid user ID format",
		})
	}

	var req dto.DepositRequest
	if err := ctx.BodyParser(&req); err != nil {
		if strings.Contains(err.Error(), "berat limbah") {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Validation failed",
				"error":   err.Error(),
			})
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
			"error":   err.Error(),
		})
	}

	deposit, err := h.wasteDepositUsecase.CreateDeposit(userId, req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create deposit",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(deposit)
}

func (h WasteDepositHandler) GetUserDepositHistory(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(uuid.UUID)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}

	deposit, err := h.wasteDepositUsecase.GetUserDeposits(userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get deposit history",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(deposit)
}

func (h WasteDepositHandler) GetUserReward(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals("userId").(uuid.UUID)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}

	deposit, err := h.wasteDepositUsecase.GetUserReward(userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get reward history",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(deposit)
}
