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

type ModelHandler struct {
	dbClient *postgrest.Client
}

func NewModelHandler() *ModelHandler {
	return &ModelHandler{
		dbClient: config.GetDBClient(),
	}
}

func (h *ModelHandler) CreateModel(c *fiber.Ctx) error {
	var model models.AIModel
	if err := c.BodyParser(&model); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := utils.ValidateModelMetadata(&model); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Metadata are not valid",
		})
	}

	if err := utils.VerifyHuggingfaceModel(model.HuggingfaceID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid HuggingFace model ID",
		})
	}

	model.ID = uuid.New()
	model.CreatedAt = time.Now()
	model.IsActive = true
	model.FunctionURL = utils.GenerateEdgeFunctionURL(model.ModelType, model.HuggingfaceID)

	_, _, err := h.dbClient.From("ai_models").
		Insert(model, false, "", "representation", "exact").
		Execute()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create model",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model)
}

func (h *ModelHandler) GetModel(c *fiber.Ctx) error {
	modelID := c.Params("id")
	id, err := uuid.Parse(modelID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid model ID",
		})
	}

	result, count, err := h.dbClient.From("ai_models").
		Select("*", "exact", false).
		Eq("id", id.String()).
		Execute()

	if err != nil || count == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": models.ErrModelNotFound.Error,
		})
	}

	var model models.AIModel
	if err := json.Unmarshal(result, &model); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse response",
		})
	}

	return c.JSON(model)
}

func (h *ModelHandler) ListModels(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	result, count, err := h.dbClient.From("ai_models").
		Select("*", "exact", false).
		Range(offset, offset+limit-1, "").
		Execute()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch models",
		})
	}

	var resModels []models.AIModel
	if err := json.Unmarshal(result, &resModels); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse response",
		})
	}

	return c.JSON(fiber.Map{
		"models": resModels,
		"total":  count,
		"page":   page,
		"limit":  limit,
	})
}

func (h *ModelHandler) UpdateModel(c *fiber.Ctx) error {
	modelID := c.Params("id")
	id, err := uuid.Parse(modelID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid model ID",
		})
	}

	var updateData models.AIModel
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := utils.ValidateModelMetadata(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Metadata are not valid",
		})
	}

	_, count, err := h.dbClient.From("ai_models").
		Select("*", "exact", false).
		Eq("id", id.String()).
		Execute()

	if err != nil || count == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": models.ErrModelNotFound.Error,
		})
	}

	updateData.FunctionURL = utils.GenerateEdgeFunctionURL(updateData.ModelType, updateData.HuggingfaceID)

	_, _, err = h.dbClient.From("ai_models").
		Update(updateData, "representation", "exact").
		Eq("id", id.String()).
		Execute()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update model",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *ModelHandler) DeleteModel(c *fiber.Ctx) error {
	modelID := c.Params("id")
	id, err := uuid.Parse(modelID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid model ID",
		})
	}

	updateData := map[string]interface{}{"is_active": false}
	_, _, err = h.dbClient.From("ai_models").
		Update(updateData, "representation", "exact").
		Eq("id", id.String()).
		Execute()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete model",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
