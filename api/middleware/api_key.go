package middleware

import (
	"api/config"
	"api/models"
	"api/utils"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"strings"
	"sync"
	"time"
)

type keyUsage struct {
	count    int
	lastSeen time.Time
}

var (
	usageMap = make(map[string]*keyUsage)
	usageMux sync.RWMutex
)

func init() {
	// Clean up usage tracking periodically
	go func() {
		for {
			time.Sleep(15 * time.Minute)
			cleanupUsage()
		}
	}()
}

func cleanupUsage() {
	usageMux.Lock()
	defer usageMux.Unlock()

	for key, usage := range usageMap {
		if time.Since(usage.lastSeen) > time.Hour {
			delete(usageMap, key)
		}
	}
}

func ValidateAPIKey() fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			apiKey = strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		}

		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "API key required",
			})
		}

		if !utils.ValidateKeyFormat(apiKey) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid API key format",
			})
		}

		keyHash := utils.HashAPIKey(apiKey)

		dbClient := config.GetDBClient()
		result, count, err := dbClient.From("api_keys").
			Select("*", "exact", false).
			Eq("key_hash", keyHash).
			Eq("is_active", "true").
			Execute()

		if err != nil || count == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid API key",
			})
		}

		var keys []models.APIKey
		if err := json.Unmarshal(result, &keys); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to process API key",
			})
		}

		key := keys[0]

		usageMux.Lock()
		if usageMap[keyHash] == nil {
			usageMap[keyHash] = &keyUsage{
				count:    1,
				lastSeen: time.Now(),
			}
		} else {
			if time.Since(usageMap[keyHash].lastSeen) <= time.Minute {
				if usageMap[keyHash].count >= key.RateLimit {
					usageMux.Unlock()
					return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
						"error": "Rate limit exceeded",
					})
				}
				usageMap[keyHash].count++
			} else {
				usageMap[keyHash].count = 1
			}
			usageMap[keyHash].lastSeen = time.Now()
		}
		usageMux.Unlock()

		// Update last_used timestamp in database
		updateData := map[string]interface{}{
			"last_used": time.Now(),
		}

		_, _, err = dbClient.From("api_keys").
			Update(updateData, "representation", "exact").
			Eq("id", key.ID.String()).
			Execute()

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update API key usage",
			})
		}

		// Store key info for handlers
		c.Locals("api_key", key)
		return c.Next()
	}
}

func RequireAPIKey() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals("api_key") == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Valid API key required",
			})
		}
		return c.Next()
	}
}
