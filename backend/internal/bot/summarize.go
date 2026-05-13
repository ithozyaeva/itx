package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"ithozyeva/config"
	"ithozyeva/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	summarizeDefaultLimit = 200
	summarizeDailyLimit   = 5
)

var openAIClient = &http.Client{Timeout: 120 * time.Second}

// fallbackModels — цепочка моделей для retry: от лучшей к запасным.
var fallbackModels = []string{
	"Qwen/Qwen3-235B-A22B-Instruct-2507",
	"zai-org/GLM-4.7",
	"ai-sage/GigaChat3-10B-A1.8B",
	"zai-org/GLM-4.7-Flash",
	"t-tech/T-pro-it-2.0",
	"t-tech/T-pro-it-2.1",
}

// userSummarizeCount — лимит суммаризаций на пользователя в день
var (
	userSummarizeCount = make(map[int64]map[string]int) // userID -> date -> count
	summarizeMu        sync.Mutex
)

func checkAndIncrementLimit(userID int64) bool {
	summarizeMu.Lock()
	defer summarizeMu.Unlock()

	today := time.Now().Format("2006-01-02")
	if userSummarizeCount[userID] == nil {
		userSummarizeCount[userID] = make(map[string]int)
	}

	if userSummarizeCount[userID][today] >= summarizeDailyLimit {
		return false
	}
	userSummarizeCount[userID][today]++

	// Очищаем старые даты
	for date := range userSummarizeCount[userID] {
		if date != today {
			delete(userSummarizeCount[userID], date)
		}
	}
	return true
}

func getRemainingLimit(userID int64) int {
	summarizeMu.Lock()
	defer summarizeMu.Unlock()

	today := time.Now().Format("2006-01-02")
	if userSummarizeCount[userID] == nil {
		return summarizeDailyLimit
	}
	used := userSummarizeCount[userID][today]
	remaining := summarizeDailyLimit - used
	if remaining < 0 {
		return 0
	}
	return remaining
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

// handleSummarizeCommand — /summarize [N|day|week|3d]
func (b *TelegramBot) handleSummarizeCommand(message *tgbotapi.Message) {
	deleteMsg := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
	b.bot.Request(deleteMsg)

	if config.CFG.OpenAIKey == "" {
		b.SendDirectMessage(message.From.ID, "Суммаризация не настроена: отсутствует OPENAI_API_KEY.")
		return
	}

	if !checkAndIncrementLimit(message.From.ID) {
		b.SendDirectMessage(message.From.ID, fmt.Sprintf("Лимит суммаризаций исчерпан (%d/%d в день). Попробуйте завтра.", summarizeDailyLimit, summarizeDailyLimit))
		return
	}

	arg := strings.TrimSpace(message.CommandArguments())
	messages, label, err := b.fetchMessages(message.Chat.ID, arg)
	if err != nil {
		log.Printf("Error fetching messages for summarize: %v", err)
		b.SendDirectMessage(message.From.ID, "Ошибка при получении сообщений.")
		return
	}

	if len(messages) == 0 {
		b.SendDirectMessage(message.From.ID, "Нет сообщений для суммаризации за указанный период.")
		return
	}

	remaining := getRemainingLimit(message.From.ID)
	b.SendDirectMessage(message.From.ID, fmt.Sprintf("⏳ Суммаризирую %d сообщений (%s) из чата <b>%s</b>...\nОсталось запросов: %d/%d",
		len(messages), label, html.EscapeString(message.Chat.Title), remaining, summarizeDailyLimit))

	var sb strings.Builder
	for i := len(messages) - 1; i >= 0; i-- {
		m := messages[i]
		name := m.TelegramUsername
		if name == "" {
			name = m.TelegramFirstName
		}
		sb.WriteString(fmt.Sprintf("[%s] %s: %s\n", m.SentAt.Format("02.01 15:04"), name, m.MessageText))
	}

	summary, usedModel, err := callOpenAIWithRetry(sb.String())
	if err != nil {
		log.Printf("Error calling OpenAI for summarize (all models failed): %v", err)
		b.SendDirectMessage(message.From.ID, "Ошибка: все AI-модели недоступны.")
		return
	}

	result := fmt.Sprintf("📋 <b>Суммаризация чата %s</b>\n(%d сообщений, %s, модель: %s)\n\n%s",
		html.EscapeString(message.Chat.Title), len(messages), label, usedModel, html.EscapeString(summary))
	b.SendDirectMessage(message.From.ID, result)
}

// fetchMessages возвращает сообщения в зависимости от аргумента:
// "" — последние 200
// "day" — за сутки
// "week" — за неделю
// "3d" — за 3 дня
// число (50, 100, 500) — последние N (макс 1000)
func (b *TelegramBot) fetchMessages(chatID int64, arg string) ([]models.ChatMessage, string, error) {
	switch strings.ToLower(arg) {
	case "day", "today", "сегодня", "день":
		since := time.Now().Add(-24 * time.Hour)
		msgs, err := b.chatActivityService.GetMessagesSince(chatID, since)
		return msgs, "за сутки", err

	case "week", "неделя":
		since := time.Now().Add(-7 * 24 * time.Hour)
		msgs, err := b.chatActivityService.GetMessagesSince(chatID, since)
		return msgs, "за неделю", err

	case "3d", "3дня":
		since := time.Now().Add(-3 * 24 * time.Hour)
		msgs, err := b.chatActivityService.GetMessagesSince(chatID, since)
		return msgs, "за 3 дня", err

	default:
		if n, err := strconv.Atoi(arg); err == nil && n > 0 {
			if n > 1000 {
				n = 1000
			}
			msgs, err := b.chatActivityService.GetRecentMessages(chatID, n)
			return msgs, fmt.Sprintf("последние %d", n), err
		}

		msgs, err := b.chatActivityService.GetRecentMessages(chatID, summarizeDefaultLimit)
		return msgs, fmt.Sprintf("последние %d", summarizeDefaultLimit), err
	}
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
Пиши на русском языке. Будь лаконичен.
ВАЖНО: Используй HTML-разметку для форматирования (Telegram HTML). Доступные теги: <b>жирный</b>, <i>курсив</i>, <code>код</code>. НЕ используй Markdown (**, __, #, -). Для списков используй символ • в начале строки.`,
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
