package main

import (
	"api/config"
	"log"
)

func main() {

	if err := config.InitSupabase(); err != nil {
		log.Fatal("Failed to initialized Supabase: %v", err)

	}
}
