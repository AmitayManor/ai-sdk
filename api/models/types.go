package models

import (
	"time"
)

type User struct {
	ID                  string
	Email               string
	CreatedAt           time.Time
	LastLogin           *time.Time
	IsAdmin             bool
	IsActive            bool
	FailedLoginAttempts int
}
type APIKey struct {
	ID        string
	UserID    string
	KeyHash   string
	Name      string
	CreatedAt time.Time
	LastUsed  *time.Time
	IsActive  bool
	RateLimit int
}

type AIModel struct {
	ID            string
	Name          string
	ModelType     string
	Version       string
	HuggingfaceID string
	IsActive      bool
	CreatedAt     time.Time
}

type ModelRequest struct {
	ID          string
	UserID      string
	ModelID     string
	APIKeyID    string
	CreatedAt   time.Time
	CompletedAt *time.Time
	Status      string
	InputData   map[string]interface{}
	OutputData  map[string]interface{}
	ErrorMsg    string
	TokenUsed   int
}

