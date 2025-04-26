package middleware

import (
	"hackfest-uc/internal/infra/jwt"

	"github.com/gofiber/fiber/v2"
)

type MiddlewareI interface {
	Authentication(ctx *fiber.Ctx) error
	Authorization(ctx *fiber.Ctx) error
}

type Middleware struct {
	jwt *jwt.JWT
}

func NewMiddleware(jwt *jwt.JWT) MiddlewareI {
	return &Middleware{
		jwt: jwt,
	}
}
