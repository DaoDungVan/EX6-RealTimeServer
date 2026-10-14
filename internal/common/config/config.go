package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	godotenv.Load()
}

func GetEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
