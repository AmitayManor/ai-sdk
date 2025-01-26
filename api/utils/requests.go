package utils

import (
	"api/config"
	"api/models"
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/supabase-community/postgrest-go"
	"net/http"
	"time"
)

func ProcessModelRequest(req models.ModelRequest, model models.AIModel) {
	dbClient := config.GetDBClient()

	payload := map[string]interface{}{
		"input":      req.InputData,
		"request_id": req.ID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		updateRequestStatus(req.ID, "FAILED", "Failed to prepare request", dbClient)
		return
	}

	resp, err := http.Post(model.FunctionURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		updateRequestStatus(req.ID, "FAILED", "Failed to call model endpoint", dbClient)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		updateRequestStatus(req.ID, "FAILED", "Failed to parse response", dbClient)
		return
	}

	now := time.Now()
	updateDate := map[string]interface{}{
		"status":       "COMPLETED",
		"completed_at": now,
		"output_data":  result,
	}

	_, _, err = dbClient.From("model_requests").
		Update(updateDate, "representation", "excat").
		Eq("id", req.ID.String()).
		Execute()

	if err != nil {
		updateRequestStatus(req.ID, "FAILED", "Failed to store results", dbClient)
	}
}

func updateRequestStatus(requestID uuid.UUID, status string, errorMsg string, dbClient *postgrest.Client) {
	updateData := map[string]interface{}{
		"status":       status,
		"error_msg":    errorMsg,
		"completed_at": time.Now(),
	}

	dbClient.From("model_requests").
		Update(updateData, "representation", "exact").
		Eq("id", requestID.String()).
		Execute()
}
