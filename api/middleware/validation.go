package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/supabase-community/gotrue-go/types"
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

func ValidateSignUp() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input types.SignupRequest

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		if !emailRegex.MatchString(input.Email) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid email format",
			})
		}

		if len(input.Password) < 8 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid password format",
			})
		}

		return c.Next()
	}
}
