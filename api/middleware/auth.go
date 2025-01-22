package middleware

import (
	"api/config"
	"api/models"
	"github.com/gofiber/fiber/v2"
	"strings"
	"sync"
	"time"
)

type TokenBlacklist struct {
	tokens map[string]time.Time
	mutex  sync.RWMutex
}

var blacklist = &TokenBlacklist{
	tokens: make(map[string]time.Time),
}

func (tb *TokenBlacklist) IsBlackListed(token string) bool {
	tb.mutex.RLock()
	defer tb.mutex.RUnlock()

	expiryTime, exists := tb.tokens[token]
	if !exists {
		return false
	}

	if time.Now().After(expiryTime) {
		tb.mutex.RUnlock()
		tb.mutex.Lock()

		delete(tb.tokens, token)

		tb.mutex.Unlock()
		tb.mutex.RLock()
		return false
	}
	return true
}

func (tb *TokenBlacklist) AddToBlacklist(token string, duration time.Duration) {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	tb.tokens[token] = time.Now().Add(duration)
}

func ValidateToken(token string) (*models.User, error) {
	if blacklist.IsBlackListed(token) {
		return nil, models.ErrUnauthorized
	}

	client := config.GetSupabaseClient()
	client.Auth.WithToken(token)

	user, err := client.Auth.GetUser()
	if err != nil {
		return nil, models.ErrUnauthorized
	}

	isAdmin := false
	if adminValue, ok := user.AppMetadata["is_admin"]; ok {
		isAdmin, _ = adminValue.(bool)
	}

	appUser := &models.User{
		Email:    user.Email,
		IsAdmin:  isAdmin,
		IsActive: true,
	}

	return appUser, nil
}

func Protected() fiber.Handler {

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		user, err := ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": models.ErrUnauthorized.Error,
			})
		}

		c.Locals("token", token)
		c.Locals("user", user)

		return c.Next()
	}
}

func AdminOnly() fiber.Handler {

	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*models.User)
		if !ok || user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": models.ErrUnaithenticated.Error(),
			})
		}

		if !user.IsAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": models.ErrNotAdmin.Error(),
			})
		}
		return c.Next()
	}
}

func RateLimiter(request int, duration time.Duration) fiber.Handler {
	type client struct {
		count    int
		lastSeen time.Time
	}

	var (
		clients = make(map[string]*client)
		mu      sync.RWMutex
	)

	go func() {
		for {
			time.Sleep(duration)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > duration {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *fiber.Ctx) error {
		ip := c.IP()

		mu.Lock()
		defer mu.Unlock()

		if clients[ip] == nil {
			clients[ip] = &client{count: 1, lastSeen: time.Now()}
			return c.Next()
		}

		if time.Since(clients[ip].lastSeen) > duration {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
			})
		}

		clients[ip].count++
		clients[ip].lastSeen = time.Now()
		return c.Next()
	}
}
