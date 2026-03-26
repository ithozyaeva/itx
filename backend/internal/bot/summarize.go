package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"ithozyeva/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const summarizeMessageLimit = 200

var openAIClient = &http.Client{Timeout: 120 * time.Second}

// fallbackModels — цепочка моделей для retry: от лучшей к запасным.
var fallbackModels = []string{
	"Qwen/Qwen3-235B-A22B-Instruct-2507", // 235B, лучшая для суммаризации (20.7/61)
	"t-tech/T-pro-it-2.1",                 // T-Pro, хороший русский (бесплатная)
	"zai-org/GLM-4.7",                     // GLM 4.7 (бесплатная)
	"ai-sage/GigaChat3-10B-A1.8B",         // GigaChat лёгкая (12.2/12.2)
	"zai-org/GLM-4.7-Flash",               // GLM Flash (бесплатная)
}

type openAIRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
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
	messages, err := b.chatActivityService.GetRecentMessages(message.Chat.ID, summarizeMessageLimit)
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

	// Вызываем AI с retry по fallback-моделям
	summary, usedModel, err := callOpenAIWithRetry(sb.String())
	if err != nil {
		log.Printf("Error calling OpenAI for summarize (all models failed): %v", err)
		b.SendDirectMessage(message.From.ID, "Ошибка: все AI-модели недоступны.")
		return
	}

	result := fmt.Sprintf("📋 <b>Суммаризация чата %s</b>\n(%d сообщений, модель: %s)\n\n%s", message.Chat.Title, len(messages), usedModel, summary)
	b.SendDirectMessage(message.From.ID, result)
}

// callOpenAIWithRetry пробует модели по цепочке: сначала из конфига, потом fallback.
func callOpenAIWithRetry(chatLog string) (summary string, model string, err error) {
	models := buildModelChain()

	var lastErr error
	for _, m := range models {
		summary, err := callOpenAI(chatLog, m)
		if err == nil {
			return summary, m, nil
		}
		log.Printf("Model %s failed: %v, trying next...", m, err)
		lastErr = err
	}

	return "", "", fmt.Errorf("all %d models failed, last error: %w", len(models), lastErr)
}

// buildModelChain строит цепочку моделей: конфиг-модель первая, затем fallback (без дублей).
func buildModelChain() []string {
	primary := config.CFG.OpenAIModel
	if primary == "" {
		return fallbackModels
	}

	chain := []string{primary}
	for _, m := range fallbackModels {
		if m != primary {
			chain = append(chain, m)
		}
	}
	return chain
}

func callOpenAI(chatLog string, model string) (string, error) {
	baseURL := config.CFG.OpenAIBaseURL
	if baseURL == "" {
		baseURL = "https://foundation-models.api.cloud.ru/v1"
	}

	reqBody := openAIRequest{
		Model:       model,
		MaxTokens:   2500,
		Temperature: 0.5,
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

	resp, err := openAIClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
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
