package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/supabase-community/postgrest-go"
	storage_go "github.com/supabase-community/storage-go"
	"github.com/supabase-community/supabase-go"
	"os"
	"strings"
)

var SupabaseClient *supabase.Client

func InitSupabase() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	// Add debug logging
	fmt.Printf("Initializing Supabase client...\n")

	// The URL should include the full path
	supabaseURL := os.Getenv("SUPABASE_URL")
	if !strings.HasPrefix(supabaseURL, "https://") {
		supabaseURL = "https://" + supabaseURL
	}

	supabaseKey := os.Getenv("SUPABASE_ANON_KEY")

	fmt.Printf("Supabase URL: %s\n", supabaseURL)
	fmt.Printf("Supabase Key length: %d\n", len(supabaseKey))

	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		return fmt.Errorf("failed to create Supabase client: %w", err)
	}

	SupabaseClient = client
	fmt.Printf("Supabase client initialized successfully\n")

	return nil
}
func GetSupabaseClient() *supabase.Client {
	return SupabaseClient
}

var dbClient *postgrest.Client

func InitPostgres() error {
	supabaseUrl := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_ANON_KEY")

	headers := map[string]string{
		"apikey":        supabaseKey,
		"Authorization": fmt.Sprintf("Bearer %s", supabaseKey),
		"Content-Type":  "application/json",
		"Prefer":        "return=minimal",
	}

	client := postgrest.NewClient(fmt.Sprintf("%s/rest/v1", supabaseUrl), "public", headers)

	fmt.Printf("DB Client created with URL: %s\n", supabaseUrl)

	if client.ClientError != nil {
		return fmt.Errorf("failed to create Postgres client: %w", client.ClientError)
	}

	dbClient = client
	return nil
}

func GetDBClient() *postgrest.Client {
	return dbClient
}

var storageClient *storage_go.Client

func InitStorage() error {
	fmt.Printf("Initializing Storage client...\n")

	storageUrl := os.Getenv("STORAGE_URL")
	supabaseKey := os.Getenv("SUPABASE_ANON_KEY")

	headers := map[string]string{
		"apikey": supabaseKey,
	}

	// Initialize storage client
	client := storage_go.NewClient(storageUrl, supabaseKey, headers)
	if client == nil {
		return fmt.Errorf("failed to create Storage client is empty")
	}

	storageClient = client
	fmt.Printf("Storage client initialized successfully\n")
	return nil
}

func GetStorageClient() *storage_go.Client {
	return storageClient
}
