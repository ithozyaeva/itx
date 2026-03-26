package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"ithozyeva/config"
	"ithozyeva/internal/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const summarizeMessageLimit = 200

type openAIRequest struct {
	Model    string          `json:"model"`
	Messages []openAIMessage `json:"messages"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	Choices []struct {
		Message openAIMessage `json:"message"`
	} `json:"choices"`
}

func (b *TelegramBot) handleSummarizeCommand(message *tgbotapi.Message) {
	// Удаляем команду из чата
	deleteMsg := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
	b.bot.Request(deleteMsg)

	if config.CFG.OpenAIKey == "" {
		b.SendDirectMessage(message.From.ID, "Суммаризация не настроена: отсутствует OPENAI_API_KEY.")
		return
	}

	// Получаем последние сообщения из чата
	repo := repository.NewChatActivityRepository()
	messages, err := repo.GetRecentMessages(message.Chat.ID, summarizeMessageLimit)
	if err != nil {
		log.Printf("Error fetching messages for summarize: %v", err)
		b.SendDirectMessage(message.From.ID, "Ошибка при получении сообщений.")
		return
	}

	if len(messages) == 0 {
		b.SendDirectMessage(message.From.ID, "Нет сообщений для суммаризации в этом чате.")
		return
	}

	// Формируем текст для суммаризации (в хронологическом порядке)
	var sb strings.Builder
	for i := len(messages) - 1; i >= 0; i-- {
		m := messages[i]
		name := m.TelegramUsername
		if name == "" {
			name = m.TelegramFirstName
		}
		sb.WriteString(fmt.Sprintf("[%s] %s: %s\n", m.SentAt.Format("02.01 15:04"), name, m.MessageText))
	}

	// Отправляем уведомление пользователю
	b.SendDirectMessage(message.From.ID, fmt.Sprintf("⏳ Суммаризирую последние %d сообщений из чата <b>%s</b>...", len(messages), message.Chat.Title))

	// Вызываем OpenAI API
	summary, err := callOpenAI(sb.String())
	if err != nil {
		log.Printf("Error calling OpenAI for summarize: %v", err)
		b.SendDirectMessage(message.From.ID, "Ошибка при обращении к AI.")
		return
	}

	result := fmt.Sprintf("📋 <b>Суммаризация чата %s</b>\n(%d сообщений)\n\n%s", message.Chat.Title, len(messages), summary)
	b.SendDirectMessage(message.From.ID, result)
}

func callOpenAI(chatLog string) (string, error) {
	baseURL := config.CFG.OpenAIBaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	model := config.CFG.OpenAIModel
	if model == "" {
		model = "gpt-4o-mini"
	}

	reqBody := openAIRequest{
		Model: model,
		Messages: []openAIMessage{
			{
				Role: "system",
				Content: `Ты — помощник, который суммаризирует переписки из Telegram-чатов IT-сообщества.
Твоя задача — кратко и структурировано изложить основные темы обсуждений, ключевые мнения и выводы.
Формат: маркированный список основных тем с краткими пояснениями.
Пиши на русском языке. Будь лаконичен.`,
			},
			{
				Role:    "user",
				Content: "Суммаризируй эту переписку:\n\n" + chatLog,
			},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.CFG.OpenAIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI API error %d: %s", resp.StatusCode, string(respBody))
	}

	var result openAIResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return result.Choices[0].Message.Content, nil
}
