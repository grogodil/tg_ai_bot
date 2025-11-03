package bot

import (
	"context"
	"log"
	"time"

	"tg_ai_bot/internal/ai"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Start(telegramToken string, openRouter *ai.Client) error {
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		return err
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		go func(update tgbotapi.Update) {
			if update.Message != nil {
				chatID := update.Message.Chat.ID
				userState := GetUserState(chatID)

				if update.Message.IsCommand() {
					switch update.Message.Command() {
					case "start":
						userState.QuizIndex = 0
						userState.QuizScore = 0
						userState.WaitingForAI = false
						msg := tgbotapi.NewMessage(chatID, "Выберите действие:")
						msg.ReplyMarkup = GetMainMenu()
						_, err := bot.Send(msg)
						if err != nil {
							log.Println("Ошибка отправки меню:", err)
						}
					}
					return
				}

				if userState.WaitingForAI {
					ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
					defer cancel()

					response, err := openRouter.Ask(ctx, update.Message.Text)
					if err != nil {
						bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при обращении к ИИ"))
					} else {
						bot.Send(tgbotapi.NewMessage(chatID, response))
					}
					userState.WaitingForAI = false
					return
				}

				if userState.QuizIndex < len(quizQuestions) {
					if err := HandleQuizAnswer(bot, chatID, userState, update.Message.Text); err != nil {
						log.Println("Ошибка обработки квиза:", err)
					}
					return
				}

			} else if update.CallbackQuery != nil {
				var chatID int64
				if update.CallbackQuery.Message != nil {
					chatID = update.CallbackQuery.Message.Chat.ID
				} else {
					chatID = update.CallbackQuery.From.ID
				}
				userState := GetUserState(chatID)

				switch update.CallbackQuery.Data {
				case "quiz":
					userState.QuizIndex = 0
					userState.QuizScore = 0
					userState.WaitingForAI = false
					SendQuizQuestion(bot, chatID, userState)
					callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Начинаем квест!")
					if _, err := bot.Request(callback); err != nil {
						log.Println("Ошибка callback:", err)
					}
				case "ai":
					userState.WaitingForAI = true
					msg := tgbotapi.NewMessage(chatID, "Введите ваш вопрос для ИИ:")
					bot.Send(msg)
					callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Жду вопрос")
					if _, err := bot.Request(callback); err != nil {
						log.Println("Ошибка callback:", err)
					}
				case "about":
					msg := tgbotapi.NewMessage(chatID, "🤖 Бот на Go с квизом и бесплатным ИИ OpenRouter.ai")
					bot.Send(msg)
					callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Информация")
					if _, err := bot.Request(callback); err != nil {
						log.Println("Ошибка callback:", err)
					}
				}
			}
		}(update)
	}
	return nil
}
