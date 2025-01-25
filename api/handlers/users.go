package handlers

import (
	"api/config"
	"api/models"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
)

type UserHandler struct {
	supabaseClient *supabase.Client
}

func NewUserHandler(supabaseClient *supabase.Client) *UserHandler {
	return &UserHandler{
		supabaseClient: supabaseClient,
	}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user models.User

	// Verify body is correct
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid body request",
		})
	}

	// Check there is no users as the new one
	dbClient := config.GetDBClient()
	_, count, err := dbClient.From("users").
		Select("id, email", "exact", false).
		Eq("email", user.Email).
		Execute()

	if err == nil && count > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": models.ErrUserAlreadyExists.Error,
		})
	}

	// Create the new user
	_, _, err = dbClient.From("users").
		Insert(user, false, "", "representation", "exact").
		Execute()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		IsAdmin:  user.IsAdmin,
		IsActive: user.IsActive,
	})
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	id, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": models.ErrUserNotFound.Error,
		})
	}
	dbClient := config.GetDBClient()
	result, count, err := dbClient.From("users").
		Select("*", "exact", false).
		Eq("id", id.String()).
		Execute()

	if err != nil || count == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": models.ErrUserNotFound.Error,
		})
	}

	var user models.User
	if err := json.Unmarshal(result, &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse response",
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		IsAdmin:  user.IsAdmin,
		IsActive: user.IsActive,
	})
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	// pagination
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	dbClient := config.GetDBClient()
	res, count, err := dbClient.From("users").
		Select("*", "exact", false).
		Range(offset, offset+limit-1, "").
		Execute()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	var users []models.User
	if err := json.Unmarshal(res, &users); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to parse response",
		})
	}

	response := make([]models.UserResponse, len(users))
	for i, user := range users {
		response[i] = models.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			IsAdmin:  user.IsAdmin,
			IsActive: user.IsActive,
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"users": response,
		"total": count,
		"page":  page,
		"limit": limit,
	})
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	id, err := uuid.Parse(userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": models.ErrInvalidCredentials.Error,
		})
	}

	var updateData models.User
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Body request",
		})
	}

	dbClient := config.GetDBClient()
	_, count, err := dbClient.From("users").
		Select("*", "exact", false).
		Eq("id", id.String()).
		Execute()

	if err != nil || count == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": models.ErrUserNotFound.Error,
		})
	}

	_, _, err = dbClient.From("users").
		Update(updateData, "representation", "exact").
		Eq("id", id.String()).
		Execute()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update user",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	id, err := uuid.Parse(userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": models.ErrInvalidCredentials.Error,
		})
	}

	// using Soft delete by turning Is_Active to False instead of Hard deleting
	dbClient := config.GetDBClient()
	updateData := map[string]interface{}{"is_active": false}
	_, _, err = dbClient.From("users").
		Update(updateData, "representation", "exact").
		Eq("id", id.String()).
		Execute()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete user",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *UserHandler) UpdateLoginAttempts(c *fiber.Ctx) error {
	userId := c.Params("id")
	id, err := uuid.Parse(userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": models.ErrInvalidCredentials.Error,
		})
	}

	dbClient := config.GetDBClient()
	res, count, err := dbClient.From("users").
		Select("failed_login_attempts", "exact", false).
		Eq("id", id.String()).
		Execute()

	if err != nil || count == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": models.ErrUserNotFound.Error,
		})
	}

	var user models.User
	if err := json.Unmarshal(res, &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to parse response",
		})
	}

	updateData := map[string]interface{}{
		"failed_login_attempts": user.FailedLoginAttempts + 1,
		"is_active":             user.FailedLoginAttempts < 4,
	}

	_, _, err = dbClient.From("users").
		Update(updateData, "representation", "exact").
		Eq("id", id.String()).
		Execute()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update login attempts",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *UserHandler) ResetLoginAttempts(c *fiber.Ctx) error {
	userId := c.Params("id")
	id, err := uuid.Parse(userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": models.ErrInvalidCredentials.Error,
		})
	}

	updateDate := map[string]interface{}{
		"failed_login_attempts": 0,
		"last_login":            "NOW()",
	}

	dbClient := config.GetDBClient()
	_, _, err = dbClient.From("users").
		Update(updateDate, "representation", "exact").
		Eq("id", id.String()).
		Execute()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to reset login attempts",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// i think that i need to change the c.Locals('user') to 'token'
func (h *UserHandler) AdminUpdateUser(c *fiber.Ctx) error {
	reqUserData := c.Locals("user")
	reqUser, ok := reqUserData.(models.User)
	if !ok || !reqUser.IsAdmin {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": models.ErrUnauthorized.Error,
		})
	}

	return h.UpdateUser(c)
}

func (h *UserHandler) SignIn(c *fiber.Ctx) error {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&credentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	authResponse, err := h.supabaseClient.Auth.SignInWithEmailPassword(credentials.Email, credentials.Password)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Invalid credentials",
			"details": err.Error,
		})
	}

	return c.JSON(fiber.Map{
		"access_token": authResponse.AccessToken,
		"token_type":   "Bearer",
	})
}
