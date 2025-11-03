package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken    string
	OpenRouterApiKey string
}

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found")
	}

	tgToken := os.Getenv("TELEGRAM_TOKEN")
	openRouterKey := os.Getenv("OPENROUTER_API_KEY")
	if openRouterKey == "" {
		log.Fatal("OPENROUTER_API_KEY не установлена")
	}

	if tgToken == "" || openRouterKey == "" {
		log.Fatal("Не установлены TELEGRAM_TOKEN или OPENROUTER_API_KEY")
	}

	return Config{
		TelegramToken:    tgToken,
		OpenRouterApiKey: openRouterKey,
	}
}
