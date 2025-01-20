package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID                  uuid.UUID  `json:"id"`
	Email               string     `json:"email"`
	CreatedAt           time.Time  `json:"created_at"`
	LastLogin           *time.Time `json:"last_login"`
	IsAdmin             bool       `json:"is_admin"`
	IsActive            bool       `json:"is_active"`
	FailedLoginAttempts int        `json:"failed_login_attempts"`
}

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	IsAdmin  bool      `json:"is_admin"`
	IsActive bool      `json:"is_active"`
}

type APIKey struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	KeyHash   string     `json:"key_hash"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	LastUsed  *time.Time `json:"last_update"`
	IsActive  bool       `json:"is_active"`
	RateLimit int        `json:"rate_limit"`
}

type AIModel struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	ModelType     string    `json:"model_type"`
	Version       string    `json:"version"`
	HuggingfaceID string    `json:"huggingface_id"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
}

type ModelRequest struct {
	ID          uuid.UUID              `json:"id"`
	UserID      uuid.UUID              `json:"user_id"`
	ModelID     uuid.UUID              `json:"model_id"`
	APIKeyID    uuid.UUID              `json:"api_key_id"`
	CreatedAt   time.Time              `json:"created_at"`
	CompletedAt *time.Time             `json:"completed_at"`
	Status      string                 `json:"status"`
	InputData   map[string]interface{} `json:"input_data"`
	OutputData  map[string]interface{} `json:"output_data"`
	ErrorMsg    string                 `json:"error_msg"`
	TokenUsed   int                    `json:"token_used"`
}
