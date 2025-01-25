package handlers

import (
	"api/config"
	"api/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/supabase-community/gotrue-go/types"
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

	var input types.SignupRequest

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": models.ErrInvalidRequestStatus.Error,
		})
	}

	authResp, err := h.supabaseClient.Auth.Signup(input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "signup failed",
		})
	}

	user := models.User{
		ID:       uuid.MustParse(authResp.User.ID.String()),
		Email:    authResp.User.Email,
		IsActive: true,
	}

	dbClient := config.GetDBClient()
	_, _, err = dbClient.From("users").
		Insert(user, false, "", "representation", "exact").
		Execute()

	fmt.Printf("Insert attempt result: %v\n", err)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "registration failed",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Signup successful. Please check your email for verification.",
		"user": models.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			IsActive: user.IsActive,
			IsAdmin:  user.IsAdmin,
		},
	})
}
