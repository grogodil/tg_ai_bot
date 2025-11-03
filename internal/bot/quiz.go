package bot

import (
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var quizQuestions = []struct {
    Question string
    Answer   string
}{
    {"Столица Франции?", "Париж"},
    {"2 + 2 = ?", "4"},
    {"Основатель OpenAI?", "Илон Маск"},
}

func SendQuizQuestion(bot *tgbotapi.BotAPI, chatID int64, state *UserState) error {
    if state.QuizIndex >= len(quizQuestions) {
        msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Квест окончен! Результат: %d из %d верных.", state.QuizScore, len(quizQuestions)))
        _, err := bot.Send(msg)
        return err
    }
    q := quizQuestions[state.QuizIndex].Question
    msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Вопрос %d: %s", state.QuizIndex+1, q))
    _, err := bot.Send(msg)
    return err
}

func HandleQuizAnswer(bot *tgbotapi.BotAPI, chatID int64, state *UserState, text string) error {
    correctAnswer := quizQuestions[state.QuizIndex].Answer
    var respMsg string
    if text == correctAnswer {
        state.QuizScore++
        respMsg = "✅ Верно!"
    } else {
        respMsg = fmt.Sprintf("❌ Неверно! Правильный ответ: %s", correctAnswer)
    }
    state.QuizIndex++

    if _, err := bot.Send(tgbotapi.NewMessage(chatID, respMsg)); err != nil {
        return err
    }
    return SendQuizQuestion(bot, chatID, state)
}