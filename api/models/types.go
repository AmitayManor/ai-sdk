package models

import (
	"github.com/google/uuid"
	"time"
)

type ModelRequest struct {
	ID             uuid.UUID              `json:"id"`
	CreatedAt      time.Time              `json:"created_at"`
	UserID         uuid.UUID              `json:"user_id"`
	Status         string                 `json:"status"`
	InputData      map[string]interface{} `json:"input_data"`
	OutputData     map[string]interface{} `json:"output_data"`
	ErrorMsg       string                 `json:"error_msg"`
	TokenUsed      int                    `json:"token_used"`
	TokenCount     int                    `json:"token_count"`
	ProcessingTime time.Duration          `json:"processing_time"`
	ModelType      string                 `json:"model_type"`
}
