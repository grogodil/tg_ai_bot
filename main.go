package main

import (
	"log"
	"tg_ai_bot/config"
	"tg_ai_bot/internal/ai"
	"tg_ai_bot/internal/bot"
)

func main() {
	cfg := config.LoadConfig()
	openRouterClient := ai.NewClient(cfg.OpenRouterApiKey)

	if err := bot.Start(cfg.TelegramToken, openRouterClient); err != nil {
		log.Fatalf("Ошибка запуска бота: %v", err)
	}
}
