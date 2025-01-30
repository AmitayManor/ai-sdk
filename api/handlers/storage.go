package handlers

import (
	"github.com/gofiber/fiber/v2"
	storage_go "github.com/supabase-community/storage-go"
)

type StorageHandler struct {
	storageClient *storage_go.Client
}

func NewStorageHandler(storageClient *storage_go.Client) *StorageHandler {
	return &StorageHandler{
		storageClient: storageClient,
	}
}

func (h *StorageHandler) GetGeneratedImage(c *fiber.Ctx) error {
	imageID := c.Params("id")

	data, err := h.storageClient.DownloadFile("generated-images", imageID)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Image not found",
		})
	}

	return c.Send(data)
}
