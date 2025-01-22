package config

import (
	"fmt"
	"github.com/joho/godotenv"
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
