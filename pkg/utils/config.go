package utils

import (
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		GetLogger().Sugar().Error("Error loading .env file")
	}
}
