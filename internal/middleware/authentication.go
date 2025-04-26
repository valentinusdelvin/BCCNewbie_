package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) Authentication(ctx *fiber.Ctx) error {
	authHeader := ctx.GetReqHeaders()["Authorization"]

	if authHeader == nil {
		ctx.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	if len(authHeader) < 1 {
		ctx.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		return nil
	}
	bearerToken := authHeader[0]

	if bearerToken == "" {
		ctx.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		return nil
	}

	token := strings.Split(bearerToken, " ")[1]
	fmt.Println(token)

	id, isAdmin, err := m.jwt.ValidateToken(token)
	if err != nil {
		ctx.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   err.Error(),
		})
		return nil
	}

	ctx.Locals("userId", id)
	ctx.Locals("isAdmin", isAdmin)
	return ctx.Next()
}
