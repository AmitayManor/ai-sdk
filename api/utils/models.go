package utils

import (
	"api/models"
	"errors"
	"fmt"
	"net/http"
)

func ValidateModelMetadata(model *models.AIModel) error {
	if model.Name == "" {
		return errors.New("model names is required")
	}

	if model.ModelType == "" {
		return errors.New("model type is required")
	}

	if model.Version == "" {
		return errors.New("model version is required")
	}

	if model.HuggingfaceID == "" {
		return errors.New("huggingface ID is required")
	}
	return nil
}

func VerifyHuggingfaceModel(modelID string) error {
	url := fmt.Sprintf("https://huggingface.co/api/models/%s", modelID)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("huggingface model not found")
	}

	return nil
}

func GenerateEdgeFunctionURL(modelType, modelID string) string {
	return fmt.Sprintf("https://tpuhjnicfmhvgoufjvvn.supabase.co/functions/v1/models/%s/%s",
		modelType, modelID)
}
