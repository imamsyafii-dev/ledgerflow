package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
}

func Load() (Config, error) {
	// Saat development, API dijalankan dari folder backend,
	// sedangkan .env berada di root repository.
	_ = godotenv.Load("../.env")

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	return Config{
		DatabaseURL: databaseURL,
	}, nil
}
