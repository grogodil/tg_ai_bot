package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Client struct {
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

type requestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type responseBody struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func (c *Client) Ask(ctx context.Context, question string) (string, error) {
	if c.apiKey == "" {
		return "", errors.New("API ключ OpenRouter не установлен")
	}

	body := requestBody{
		Model: "nvidia/nemotron-nano-9b-v2:free", // подставляем нужную модель
		Messages: []Message{
			{Role: "system", Content: "Ты полезный помощник."},
			{Role: "user", Content: question},
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Ошибка HTTP запроса на OpenRouter: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("OpenRouter API вернул ошибку: %s - %s", resp.Status, string(bodyBytes))
		return "", fmt.Errorf("OpenRouter error: %s", resp.Status)
	}

	var respData responseBody
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", err
	}

	if len(respData.Choices) == 0 {
		return "", errors.New("пустой ответ от OpenRouter")
	}

	return respData.Choices[0].Message.Content, nil
}
