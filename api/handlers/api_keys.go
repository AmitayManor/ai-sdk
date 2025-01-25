package handlers

import (
	"api/config"
	"api/models"
	"api/utils"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/supabase-community/postgrest-go"
	"time"
)

type APIKeyHandler struct {
	dbClient *postgrest.Client
}

func NewAPIKeyHandler() *APIKeyHandler {
	return &APIKeyHandler{
		dbClient: config.GetDBClient(),
	}
}

func (h *APIKeyHandler) CreateKey(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	var input struct {
		Name      string `json:"name"`
		RateLimit int    `json:"rate_limit"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	apiKey, err := utils.GenerateAPIKey()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate new api key",
		})
	}

	keyHash := utils.HashAPIKey(apiKey)

	newKey := models.APIKey{
		ID:        uuid.New(),
		UserID:    user.ID,
		KeyHash:   keyHash,
		Name:      input.Name,
		CreatedAt: time.Now(),
		IsActive:  true,
		RateLimit: input.RateLimit,
	}

	_, _, err = h.dbClient.From("api_keys").
		Insert(newKey, false, "", "representation", "exact").
		Execute()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create API key",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"key": apiKey,
		"id":  newKey.ID,
	})
}

func (h *APIKeyHandler) ListKeys(c *fiber.Ctx) error {
	user := c.Locals("users").(*models.User)

	res, count, err := h.dbClient.From("api_keys").
		Select("id, name, created_at, last_used, is_active, rate_limit", "exact", false).
		Eq("user_id", user.ID.String()).
		Execute()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch API key",
		})
	}

	var keys []models.APIKey
	if err := json.Unmarshal(res, &keys); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to parse response",
		})
	}

	return c.JSON(fiber.Map{
		"keys":  keys,
		"total": count,
	})
}

func (h *APIKeyHandler) DeactivateKey(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	keyID := c.Params("id")

	id, err := uuid.Parse(keyID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid key ID",
		})
	}

	_, count, err := h.dbClient.From("api_keys").
		Select("id", "exact", false).
		Eq("id", id.String()).
		Eq("user_id", user.ID.String()).
		Execute()

	if err != nil || count == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "API key not found",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "API Key not found",
		})
	}

	updateData := map[string]interface{}{
		"is_active": false,
	}

	_, _, err = h.dbClient.From("api_keys").
		Update(updateData, "representation", "exact").
		Eq("id", id.String()).
		Execute()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed deactivate API key",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *APIKeyHandler) UpdateKey(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	keyID := c.Params("id")

	id, err := uuid.Parse(keyID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid key ID",
		})
	}

	var input struct {
		Name      string `json:"name"`
		RateLimit int    `json:"rate_limit"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	_, count, err := h.dbClient.From("api_keys").
		Select("id", "exact", false).
		Eq("id", id.String()).
		Eq("user_id", user.ID.String()).
		Execute()

	if err != nil || count == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "API key not found",
		})
	}
	updateData := map[string]interface{}{
		"name":       input.Name,
		"rate_limit": input.RateLimit,
	}

	_, _, err = h.dbClient.From("api_keys").
		Update(updateData, "representation", "exact").
		Eq("id", id.String()).
		Execute()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update API key",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
