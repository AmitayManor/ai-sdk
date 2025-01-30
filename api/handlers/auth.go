package handlers

import (
	"api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/supabase-community/supabase-go"
)

type AuthHandler struct {
	supabaseClient *supabase.Client
}

func NewAuthHandler(supabaseClient *supabase.Client) *AuthHandler {
	return &AuthHandler{
		supabaseClient: supabaseClient,
	}
}
func (h *AuthHandler) SignUp(c *fiber.Ctx) error {

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": models.ErrInvalidRequestStatus.Error,
		})
	}

	_, err := h.supabaseClient.Auth.SignInWithEmailPassword(input.Email, input.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "signup failed",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Signup successful. Please check your email for verification.",
	})
}

func (h *AuthHandler) SignIn(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	authResponse, err := h.supabaseClient.Auth.SignInWithEmailPassword(input.Email, input.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Return the access token
	return c.JSON(fiber.Map{
		"token": authResponse.AccessToken,
	})
}
