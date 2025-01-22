package middleware

import (
	"api/config"
	"github.com/gofiber/fiber/v2"
	"strings"
)

// TODO: Test implementation

func Protected() fiber.Handler {

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		client := config.GetSupabaseClient()

		client.Auth.WithToken(token)

		_, err := client.Auth.GetUser()
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		c.Locals("token", token)

		return c.Next()
	}
}

func AdminOnly() fiber.Handler {

	return func(c *fiber.Ctx) error {

		token := c.Locals("token")

		if token == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Not Authenticated",
			})
		}

		client := config.GetSupabaseClient()
		user, err := client.Auth.GetUser()
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		isAdmin, ok := user.AppMetadata["is_admin"].(bool)
		if !ok || !isAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Admin access required",
			})
		}

		//userData, ok := token.(map[string]interface{})
		//if !ok {
		//	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		//		"error": "Invalid token data format",
		//	})
		//}
		//
		//isAuthorized, ok := userData["is_authorized"].(bool)
		//if !ok || !isAuthorized {
		//	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		//		"error": "Admin access required",
		//	})
		//}

		return c.Next()

	}
}
