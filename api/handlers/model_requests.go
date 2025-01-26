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

type RequestHandler struct {
	dbClient *postgrest.Client
}

func NewRequestHandler() *RequestHandler {
	return &RequestHandler{
		dbClient: config.GetDBClient(),
	}
}

func (h *RequestHandler) CreateRequest(c *fiber.Ctx) error {
	var request models.ModelRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user := c.Locals("user").(*models.User)
	request.UserID = user.ID

	result, count, err := h.dbClient.From("ai_models").
		Select("*", "exact", false).
		Eq("id", request.ModelID.String()).
		Eq("is_active", "true").
		Execute()

	if err != nil || count == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": models.ErrModelNotFound.Error,
		})
	}

	var model models.AIModel
	if err := json.Unmarshal(result, &model); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse model data",
		})
	}

	request.ID = uuid.New()
	request.CreatedAt = time.Now()
	request.Status = "PENDING"

	_, _, err = h.dbClient.From("model_requests").
		Insert(request, false, "", "representation", "exact").
		Execute()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create request",
		})
	}

	go utils.ProcessModelRequest(request, model)

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"request_id": request.ID,
		"status":     "PENDING",
	})
}

func (h *RequestHandler) GetRequest(c *fiber.Ctx) error {
	requestID := c.Params("id")
	id, err := uuid.Parse(requestID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request ID",
		})
	}

	user := c.Locals("user").(*models.User)

	result, count, err := h.dbClient.From("model_requests").
		Select("*", "exact", false).
		Eq("id", id.String()).
		Execute()

	if err != nil || count == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Request not found",
		})
	}

	var request models.ModelRequest
	if err := json.Unmarshal(result, &request); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse response",
		})
	}

	if !user.IsAdmin && request.UserID != user.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Not authorized to access this request",
		})
	}

	return c.JSON(request)
}

func (h *RequestHandler) ListRequests(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	user := c.Locals("user").(*models.User)

	var result []byte
	var count int64
	var err error

	// Admins can see all requests, users only see their own
	if user.IsAdmin {
		result, count, err = h.dbClient.From("model_requests").
			Select("*", "exact", false).
			Range(offset, offset+limit-1, "").
			Execute()
	} else {
		result, count, err = h.dbClient.From("model_requests").
			Select("*", "exact", false).
			Eq("user_id", user.ID.String()).
			Range(offset, offset+limit-1, "").
			Execute()
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch requests",
		})
	}

	var requests []models.ModelRequest
	if err := json.Unmarshal(result, &requests); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse response",
		})
	}

	return c.JSON(fiber.Map{
		"requests": requests,
		"total":    count,
		"page":     page,
		"limit":    limit,
	})
}
