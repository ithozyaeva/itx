package bot

import (
	"fmt"
	"html"
	"ithozyeva/database"
	"ithozyeva/internal/models"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// autoQuestPollInterval — как часто бот сканирует БД на новые
// авто-сгенерированные квесты, которые ещё не опубликованы в чатах.
// API создаёт записи в 11:00 MSK; минута задержки в публикации не критична.
const autoQuestPollInterval = time.Minute

// startAutoQuestPoster — фоновая горутина, постящая в Telegram-чаты
// уведомления о новых авто-сгенерированных квестах. Pull-based: API живёт
// на основном сервере (РФ), бот — на NL (Telegram заблочен из РФ), поэтому
// прямой RPC между ними дороже, чем минутный polling общей БД.
func (b *TelegramBot) startAutoQuestPoster() {
	ticker := time.NewTicker(autoQuestPollInterval)
	defer ticker.Stop()

	b.postPendingAutoQuests()
	for range ticker.C {
		b.postPendingAutoQuests()
	}
}

// postPendingAutoQuests читает квесты с auto_generated=true и пустым
// notification_posted_at и публикует их в соответствующие чаты. Помечаем
// notification_posted_at до отправки в TG, чтобы при сбое до Telegram не
// флудить чат повторными попытками — лучше потерять один пост, чем спамить.
func (b *TelegramBot) postPendingAutoQuests() {
	var quests []models.ChatQuest
	err := database.DB.
		Where("auto_generated = ? AND notification_posted_at IS NULL AND chat_id IS NOT NULL AND ends_at > NOW()", true).
		Order("created_at ASC").
		Limit(20).
		Find(&quests).Error
	if err != nil {
		log.Printf("auto-quest-poster: load pending error: %v", err)
		return
	}
	if len(quests) == 0 {
		return
	}

	for _, q := range quests {
		b.postOneAutoQuest(q)
	}
}

// postOneAutoQuest публикует одно уведомление и помечает квест отправленным.
// Порядок «mark first, then send» — намеренный: при падении SendMessage
// между чатом и сетью повторного спама не будет (см. комментарий выше).
func (b *TelegramBot) postOneAutoQuest(q models.ChatQuest) {
	if q.ChatID == nil {
		return
	}

	now := time.Now()
	res := database.DB.Model(&models.ChatQuest{}).
		Where("id = ? AND notification_posted_at IS NULL", q.Id).
		Update("notification_posted_at", now)
	if res.Error != nil {
		log.Printf("auto-quest-poster: mark posted error quest=%d: %v", q.Id, res.Error)
		return
	}
	if res.RowsAffected == 0 {
		// Гонка с другим инстансом бота (или повторный заход в текущем тикере) —
		// кто-то уже взял этот квест. Не отправляем повторно.
		return
	}

	// Description формируется генератором с подстановкой названия чата,
	// которое может содержать <, >, & — экранируем перед HTML-режимом TG,
	// иначе сообщение упадёт с can't parse entities.
	text := fmt.Sprintf(
		"🔥 <b>Новое задание для чата</b>\n\n%s\n\nЦель: <b>%d сообщений за 24 часа</b>.\nНаграда: <b>+%d баллов</b> каждому, кто выполнит.",
		html.EscapeString(q.Description), q.TargetCount, q.PointsReward,
	)
	msg := tgbotapi.NewMessage(*q.ChatID, text)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true

	if _, err := b.bot.Send(msg); err != nil {
		log.Printf("auto-quest-poster: send error quest=%d chat=%d: %v", q.Id, *q.ChatID, err)
		return
	}
	log.Printf("auto-quest-poster: posted quest=%d to chat=%d", q.Id, *q.ChatID)
}
