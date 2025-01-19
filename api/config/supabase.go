package config

import (
	"errors"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
	"os"
)

var SupabaseClient *supabase.Client

func InitSupabase() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_ANNON_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		return errors.New("missing required environment variables")
	}

	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		return err
	}

	SupabaseClient = client

	return nil
}

func GetSupabaseClient() *supabase.Client {
	return SupabaseClient
}
