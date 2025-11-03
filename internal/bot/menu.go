package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func GetMainMenu() tgbotapi.InlineKeyboardMarkup {
    buttons := [][]tgbotapi.InlineKeyboardButton{
        {
            tgbotapi.NewInlineKeyboardButtonData("🎯 Пройти квест", "quiz"),
        },
        {
            tgbotapi.NewInlineKeyboardButtonData("🤖 Получить совет от ИИ", "ai"),
        },
        {
            tgbotapi.NewInlineKeyboardButtonData("ℹ️ О боте", "about"),
        },
    }
    return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}