package rest

import (
	"hackfest-uc/internal/app/payment/usecase"
	"hackfest-uc/internal/domain/entity"
	"hackfest-uc/internal/middleware"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	usecase    usecase.InterPaymentUsecase
	middleware middleware.MiddlewareItf
}

func NewPaymentHandler(routerGroup fiber.Router, paymentUsecase usecase.InterPaymentUsecase, middleware middleware.MiddlewareItf) {
	paymentHandler := PaymentHandler{
		usecase:    paymentUsecase,
		middleware: middleware,
	}

	routerGroup = routerGroup.Group("/payments")
	routerGroup.Post("/purchase/:id", paymentHandler.Purchase)
}

func (p *PaymentHandler) Purchase(ctx *fiber.Ctx) error {
	user := ctx.Locals("userId")
	if user == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "failed to get user",
		})
	}

	parsedID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to parse product ID",
		})
	}

	amount := ctx.Query("amount")
	int_amount, err := strconv.Atoi(amount)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to parse amount",
		})
	}

	payment := entity.Payment{
		OrderID:   uuid.New(),
		UserID:    user.(uuid.UUID),
		ProductID: parsedID,
		Amount:    uint64(int_amount),
	}

	paymentLink, err := p.usecase.Purchase(payment)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to Purchase",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"payment_link": paymentLink})
}
