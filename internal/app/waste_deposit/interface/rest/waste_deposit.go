package rest

import (
	"hackfest-uc/internal/app/waste_deposit/usecase"
	"hackfest-uc/internal/domain/dto"
	"hackfest-uc/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WasteDepositHandler struct {
	wasteDepositUsecase usecase.WasteDepositUsecaseItf
	middleware          middleware.MiddlewareItf
}

func NewWasteDepositHandler(routerGroup fiber.Router, wasteDepositUsecase usecase.WasteDepositUsecaseItf, middleware middleware.MiddlewareItf) {
	WasteDepositHandler := WasteDepositHandler{
		wasteDepositUsecase: wasteDepositUsecase,
		middleware:          middleware,
	}

	routerGroup = routerGroup.Group("/deposits", middleware.Authentication)
	routerGroup.Post("/", WasteDepositHandler.CreateDeposit)

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
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
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
