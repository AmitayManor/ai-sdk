package middleware

import (
	"api/config"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid authorization header",
			})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		client := config.GetSupabaseClient()
		client.Auth.WithToken(token)
		user, err := client.Auth.GetUser()

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		c.Locals("user_id", user.ID)

		return c.Next()
	}
}
