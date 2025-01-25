package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
)

const (
	keyPrefix = "sk_"
	keyLength = 32
)

func GenerateAPIKey() (string, error) {
	randomBytes := make([]byte, keyLength)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	randomStr := base64.URLEncoding.EncodeToString(randomBytes)

	keyUUID := uuid.New().String()

	apiKey := fmt.Sprintf("%s%s_%s", keyPrefix, keyUUID, randomStr)

	return apiKey, nil
}

func HashAPIKey(apiKey string) string {
	hasher := sha256.New()
	hasher.Write([]byte(apiKey))
	return hex.EncodeToString(hasher.Sum(nil))
}

func ValidateKeyFormat(apiKey string) bool {
	if len(apiKey) < len(keyPrefix)+37 || !startWith(apiKey, keyPrefix) {
		return false
	}

	parts := split(apiKey[len(keyPrefix):], "_")

	_, err := uuid.Parse(parts[0])
	return err == nil
}

func startWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

func split(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if i+1 <= len(s) && s[i:i+1] == sep {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		result = append(result, s[start:])
	}
	return result
}
