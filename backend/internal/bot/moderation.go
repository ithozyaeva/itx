package bot

import (
	"fmt"
	"html"
	"log"
	"strconv"
	"strings"
	"time"

	"ithozyeva/internal/models"
	"ithozyeva/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Параметры голосования по умолчанию.
const (
	votebanRequiredVotes = 5
	votebanWindowSeconds = 15 * 60 // 15 минут — окно сбора голосов
	// votebanKickSeconds — длительность санкции (BanChatMember с UntilDate).
	// В БД хранится в колонке voteban.mute_seconds (имя оставлено по
	// совместимости с уже отгруженной миграцией T1-модерации).
	votebanKickSeconds = 60 * 60 // 1 час

	// Anti-abuse cooldowns / стаж голосующего.
	votebanCooldownChatSeconds      = 5 * 60        // не чаще одного /voteban в чате раз в 5 мин (любым target'ом)
	votebanCooldownInitiatorSeconds = 30 * 60       // не чаще /voteban от одного инициатора в чате раз в 30 мин
	voterMinActivityWindow          = 7 * 24 * time.Hour
	// Юзер может голосовать только если за последние 7 дней написал хотя бы
	// одно сообщение в этом чате — отсекает «свежезашедших проходящих мимо».
	voterMinMessages = 1

	cleanupDefaultPeriod = 24 * time.Hour
	cleanupMaxPeriod     = 7 * 24 * time.Hour
	cleanupBatchSleep    = 35 * time.Millisecond // Telegram ~30 удалений/сек

	moderationWatcherTick = time.Minute
)

// --- Permissions ---

// isChatAdmin checks via Telegram API whether userID is creator/admin in chatID.
func (b *TelegramBot) isChatAdmin(chatID, userID int64) bool {
	member, err := b.bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: userID,
		},
	})
	if err != nil {
		log.Printf("isChatAdmin: GetChatMember failed chat=%d user=%d: %v", chatID, userID, err)
		return false
	}
	return member.Status == "creator" || member.Status == "administrator"
}

// canModerate true для super-admin платформы, любого ADMIN в БД и админов конкретного чата.
func (b *TelegramBot) canModerate(chatID, userID int64) bool {
	if b.isSubscriptionAdmin(userID) {
		return true
	}
	if b.isAdmin(userID) {
		return true
	}
	return b.isChatAdmin(chatID, userID)
}

// --- Argument parsing ---

// commandArgs возвращает аргументы команды без самой команды (с поддержкой @bot_name).
func commandArgs(message *tgbotapi.Message) []string {
	return strings.Fields(message.CommandArguments())
}

// parseTargetFromArg достаёт telegram user id из строки: либо число, либо @username
// (находим в chat_messages для текущего чата). Возвращает (id, displayName, ok).
func (b *TelegramBot) parseTargetFromArg(chatID int64, arg string) (int64, string, bool) {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		return 0, "", false
	}
	if id, err := strconv.ParseInt(arg, 10, 64); err == nil {
		return id, fmt.Sprintf("id=%d", id), true
	}
	username := strings.TrimPrefix(arg, "@")
	if username == "" {
		return 0, "", false
	}
	id, err := b.chatActivityService.LookupUserIDByUsername(chatID, username)
	if err != nil || id == 0 {
		return 0, "", false
	}
	return id, "@" + username, true
}

// targetDisplay — html-имя для упоминания в служебных сообщениях.
func targetDisplay(user *tgbotapi.User) string {
	if user == nil {
		return "пользователь"
	}
	if user.UserName != "" {
		return "@" + html.EscapeString(user.UserName)
	}
	name := strings.TrimSpace(user.FirstName + " " + user.LastName)
	if name == "" {
		name = fmt.Sprintf("id=%d", user.ID)
	}
	return fmt.Sprintf("<a href=\"tg://user?id=%d\">%s</a>", user.ID, html.EscapeString(name))
}

// --- /ban ---

func (b *TelegramBot) handleBanCommand(message *tgbotapi.Message) {
	if message.Chat.Type != "group" && message.Chat.Type != "supergroup" {
		return
	}
	if !b.canModerate(message.Chat.ID, message.From.ID) {
		return
	}
	if message.ReplyToMessage == nil {
		b.replyAndAutoDelete(message, "Используйте /ban в ответ на сообщение нарушителя. Опционально: /ban 1h, /ban 1d.")
		return
	}
	target := message.ReplyToMessage.From
	if target == nil || target.IsBot {
		return
	}
	if b.canModerate(message.Chat.ID, target.ID) {
		b.replyAndAutoDelete(message, "Нельзя забанить администратора.")
		return
	}

	args := commandArgs(message)
	var duration time.Duration
	if len(args) > 0 {
		d, err := service.ParseHumanDuration(args[0])
		if err != nil {
			b.replyAndAutoDelete(message, fmt.Sprintf("Не понял длительность: %v. Примеры: 30m, 1h, 1d.", err))
			return
		}
		duration = d
	}

	until := int64(0)
	var expiresAt *time.Time
	if duration > 0 {
		t := time.Now().Add(duration)
		until = t.Unix()
		expiresAt = &t
	}

	if _, err := b.bot.Request(tgbotapi.BanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: message.Chat.ID,
			UserID: target.ID,
		},
		UntilDate: until,
	}); err != nil {
		log.Printf("/ban: BanChatMember failed chat=%d user=%d: %v", message.Chat.ID, target.ID, err)
		b.replyAndAutoDelete(message, "Telegram отказался банить (нужны права администратора с правом блокировки).")
		return
	}

	durSec := int(duration.Seconds())
	durPtr := &durSec
	if duration == 0 {
		durPtr = nil
	}
	if err := b.moderationService.LogAction(&models.ModerationAction{
		ChatID:          message.Chat.ID,
		TargetUserID:    target.ID,
		ActorUserID:     message.From.ID,
		Action:          models.ModerationActionBan,
		DurationSeconds: durPtr,
		ExpiresAt:       expiresAt,
	}); err != nil {
		log.Printf("/ban: log failed: %v", err)
	}

	durStr := service.FormatDurationHuman(duration)
	b.sendChatHTML(message.Chat.ID, fmt.Sprintf("⛔ %s забанен (%s).", targetDisplay(target), durStr))
	b.tryDelete(message.Chat.ID, message.MessageID)
}

// --- /unban ---

func (b *TelegramBot) handleUnbanCommand(message *tgbotapi.Message) {
	if message.Chat.Type != "group" && message.Chat.Type != "supergroup" {
		return
	}
	if !b.canModerate(message.Chat.ID, message.From.ID) {
		return
	}

	var targetID int64
	display := ""
	if message.ReplyToMessage != nil && message.ReplyToMessage.From != nil {
		targetID = message.ReplyToMessage.From.ID
		display = targetDisplay(message.ReplyToMessage.From)
	} else {
		args := commandArgs(message)
		if len(args) == 0 {
			b.replyAndAutoDelete(message, "Использование: /unban в ответ на сообщение, или /unban @username, или /unban <user_id>.")
			return
		}
		id, d, ok := b.parseTargetFromArg(message.Chat.ID, args[0])
		if !ok {
			b.replyAndAutoDelete(message, "Не нашёл пользователя. Передайте user_id или @username из этого чата.")
			return
		}
		targetID = id
		display = html.EscapeString(d)
	}

	if _, err := b.bot.Request(tgbotapi.UnbanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: message.Chat.ID,
			UserID: targetID,
		},
		OnlyIfBanned: true,
	}); err != nil {
		log.Printf("/unban: UnbanChatMember failed chat=%d user=%d: %v", message.Chat.ID, targetID, err)
		b.replyAndAutoDelete(message, "Telegram отказался разбанить.")
		return
	}

	if err := b.moderationService.LogAction(&models.ModerationAction{
		ChatID:       message.Chat.ID,
		TargetUserID: targetID,
		ActorUserID:  message.From.ID,
		Action:       models.ModerationActionUnban,
	}); err != nil {
		log.Printf("/unban: log failed: %v", err)
	}

	b.sendChatHTML(message.Chat.ID, fmt.Sprintf("✅ %s разбанен.", display))
	b.tryDelete(message.Chat.ID, message.MessageID)
}

// --- /mute ---

// restrictPermissionsMuted — запрет писать/слать медиа на время мута.
// Все CanXxx=false; RestrictChatMember в Telegram трактует «не указано»
// как «отнять» — поэтому дополнительные права тоже всё равно убираются.
func restrictPermissionsMuted() *tgbotapi.ChatPermissions {
	return &tgbotapi.ChatPermissions{}
}

// restrictPermissionsAllow — обратно «всё можно» (для unmute).
func restrictPermissionsAllow() *tgbotapi.ChatPermissions {
	return &tgbotapi.ChatPermissions{
		CanSendMessages:       true,
		CanSendMediaMessages:  true,
		CanSendPolls:          true,
		CanSendOtherMessages:  true,
		CanAddWebPagePreviews: true,
		CanInviteUsers:        true,
	}
}

func (b *TelegramBot) handleMuteCommand(message *tgbotapi.Message) {
	if message.Chat.Type != "group" && message.Chat.Type != "supergroup" {
		return
	}
	if !b.canModerate(message.Chat.ID, message.From.ID) {
		return
	}
	if message.ReplyToMessage == nil {
		b.replyAndAutoDelete(message, "Используйте /mute в ответ на сообщение. Опционально: /mute 30m, /mute 1h.")
		return
	}
	target := message.ReplyToMessage.From
	if target == nil || target.IsBot {
		return
	}
	if b.canModerate(message.Chat.ID, target.ID) {
		b.replyAndAutoDelete(message, "Нельзя замутить администратора.")
		return
	}

	args := commandArgs(message)
	var duration time.Duration
	if len(args) > 0 {
		d, err := service.ParseHumanDuration(args[0])
		if err != nil {
			b.replyAndAutoDelete(message, fmt.Sprintf("Не понял длительность: %v.", err))
			return
		}
		duration = d
	}

	until := int64(0)
	var expiresAt *time.Time
	if duration > 0 {
		t := time.Now().Add(duration)
		until = t.Unix()
		expiresAt = &t
	}

	if err := b.muteUserInChat(message.Chat.ID, target.ID, until); err != nil {
		log.Printf("/mute: failed chat=%d user=%d: %v", message.Chat.ID, target.ID, err)
		b.replyAndAutoDelete(message, "Не удалось замутить.")
		return
	}

	durSec := int(duration.Seconds())
	durPtr := &durSec
	if duration == 0 {
		durPtr = nil
	}
	if err := b.moderationService.LogAction(&models.ModerationAction{
		ChatID:          message.Chat.ID,
		TargetUserID:    target.ID,
		ActorUserID:     message.From.ID,
		Action:          models.ModerationActionMute,
		DurationSeconds: durPtr,
		ExpiresAt:       expiresAt,
	}); err != nil {
		log.Printf("/mute: log failed: %v", err)
	}

	b.sendChatHTML(message.Chat.ID, fmt.Sprintf("🔇 %s замучен (%s).", targetDisplay(target), service.FormatDurationHuman(duration)))
	b.tryDelete(message.Chat.ID, message.MessageID)
}

func (b *TelegramBot) muteUserInChat(chatID, userID int64, untilUnix int64) error {
	_, err := b.bot.Request(tgbotapi.RestrictChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
		Permissions: restrictPermissionsMuted(),
		UntilDate:   untilUnix,
	})
	return err
}

// --- /cleanup ---

func (b *TelegramBot) handleCleanupCommand(message *tgbotapi.Message) {
	if message.Chat.Type != "group" && message.Chat.Type != "supergroup" {
		return
	}
	if !b.canModerate(message.Chat.ID, message.From.ID) {
		return
	}

	var targetID int64
	var display string
	if message.ReplyToMessage != nil && message.ReplyToMessage.From != nil {
		targetID = message.ReplyToMessage.From.ID
		display = targetDisplay(message.ReplyToMessage.From)
	} else {
		args := commandArgs(message)
		if len(args) == 0 {
			b.replyAndAutoDelete(message, "Использование: /cleanup в ответ на сообщение [период], или /cleanup @username [период]. Период по умолчанию — 24h.")
			return
		}
		id, d, ok := b.parseTargetFromArg(message.Chat.ID, args[0])
		if !ok {
			b.replyAndAutoDelete(message, "Не нашёл пользователя в этом чате.")
			return
		}
		targetID = id
		display = html.EscapeString(d)
	}

	period := cleanupDefaultPeriod
	args := commandArgs(message)
	// Если в reply, period — args[0]. Если без reply, period — args[1].
	periodArg := ""
	if message.ReplyToMessage != nil {
		if len(args) > 0 {
			periodArg = args[0]
		}
	} else {
		if len(args) > 1 {
			periodArg = args[1]
		}
	}
	if periodArg != "" {
		d, err := service.ParseHumanDuration(periodArg)
		if err != nil {
			b.replyAndAutoDelete(message, fmt.Sprintf("Не понял период: %v.", err))
			return
		}
		if d <= 0 {
			b.replyAndAutoDelete(message, "Период должен быть больше нуля.")
			return
		}
		if d > cleanupMaxPeriod {
			b.replyAndAutoDelete(message, "Период не больше 7 дней.")
			return
		}
		period = d
	}

	since := time.Now().Add(-period)
	ids, err := b.moderationService.MessagesForCleanup(message.Chat.ID, targetID, since)
	if err != nil {
		log.Printf("/cleanup: query failed: %v", err)
		b.replyAndAutoDelete(message, "Ошибка при поиске сообщений.")
		return
	}

	// Ack — чтобы команда не висела молча; финальную сводку шлём по завершении.
	b.tryDelete(message.Chat.ID, message.MessageID)

	if len(ids) == 0 {
		b.sendChatHTML(message.Chat.ID, fmt.Sprintf("🧹 У %s нет сообщений за %s.", display, service.FormatDurationHuman(period)))
		return
	}

	go b.runCleanup(message.Chat.ID, targetID, message.From.ID, display, period, ids)
}

func (b *TelegramBot) runCleanup(chatID, targetID, actorID int64, display string, period time.Duration, ids []int) {
	deleted := 0
	failed := 0
	successIDs := make([]int, 0, len(ids))
	for _, mid := range ids {
		if _, err := b.bot.Request(tgbotapi.NewDeleteMessage(chatID, mid)); err != nil {
			// Сообщение могло быть уже удалено вручную — это не ошибка для нас.
			if !strings.Contains(err.Error(), "message to delete not found") {
				failed++
				log.Printf("cleanup: delete msg=%d in chat=%d failed: %v", mid, chatID, err)
			}
			continue
		}
		deleted++
		successIDs = append(successIDs, mid)
		time.Sleep(cleanupBatchSleep)
	}
	if _, err := b.moderationService.DeleteCleanedMessages(chatID, successIDs); err != nil {
		log.Printf("cleanup: drop chat_messages rows failed: %v", err)
	}

	if err := b.moderationService.LogActionWithMeta(&models.ModerationAction{
		ChatID:       chatID,
		TargetUserID: targetID,
		ActorUserID:  actorID,
		Action:       models.ModerationActionCleanup,
	}, map[string]interface{}{
		"period":  service.FormatDurationHuman(period),
		"deleted": deleted,
		"failed":  failed,
		"matched": len(ids),
	}); err != nil {
		log.Printf("cleanup: log failed: %v", err)
	}

	skipped := len(ids) - deleted - failed
	summary := fmt.Sprintf("🧹 У %s удалено %d/%d сообщений за %s.",
		display, deleted, len(ids), service.FormatDurationHuman(period))
	if failed > 0 {
		summary += fmt.Sprintf(" Не удалось: %d.", failed)
	}
	if skipped > 0 {
		summary += fmt.Sprintf(" Пропущено: %d.", skipped)
	}
	summary += "\n<i>Сообщения, отправленные до включения этой функции, удалить нельзя — у бота нет их Telegram-id.</i>"
	b.sendChatHTML(chatID, summary)
}

// --- /voteban ---

func (b *TelegramBot) handleVotebanCommand(message *tgbotapi.Message) {
	if message.Chat.Type != "group" && message.Chat.Type != "supergroup" {
		return
	}
	// Любой участник чата может запустить голосование. Не пускаем только в
	// тех чатах, которые мы не отслеживаем (бот не админ или соседний чат).
	if !b.chatActivityService.IsTrackedChat(message.Chat.ID) {
		return
	}
	if message.ReplyToMessage == nil {
		b.replyAndAutoDelete(message, "Используйте /voteban в ответ на сообщение нарушителя.")
		return
	}
	target := message.ReplyToMessage.From
	if target == nil || target.IsBot {
		return
	}
	if target.ID == message.From.ID {
		b.replyAndAutoDelete(message, "Самому на себя голосование не нужно.")
		return
	}
	if b.canModerate(message.Chat.ID, target.ID) {
		b.replyAndAutoDelete(message, "Нельзя начать голосование на администратора.")
		return
	}
	if !b.isChatMember(message.Chat.ID, message.From.ID) {
		// На случай, когда юзер вышел/был удалён, но успел отправить команду.
		return
	}

	// Cooldown по чату — защита от спама голосований подряд.
	if last, _ := b.moderationService.LatestVotebanCreatedInChat(message.Chat.ID); last != nil {
		if remain := time.Duration(votebanCooldownChatSeconds)*time.Second - time.Since(*last); remain > 0 {
			b.replyAndAutoDelete(message, fmt.Sprintf(
				"В чате уже было голосование недавно. Попробуйте через %s.", service.FormatDurationHuman(remain.Round(time.Second))))
			return
		}
	}
	// Cooldown по инициатору — защита от одного активного троля.
	if last, _ := b.moderationService.LatestVotebanCreatedByInitiator(message.Chat.ID, message.From.ID); last != nil {
		if remain := time.Duration(votebanCooldownInitiatorSeconds)*time.Second - time.Since(*last); remain > 0 {
			b.replyAndAutoDelete(message, fmt.Sprintf(
				"Ваш предыдущий /voteban был недавно. Подождите %s.", service.FormatDurationHuman(remain.Round(time.Second))))
			return
		}
	}

	triggerID := message.ReplyToMessage.MessageID
	triggerPtr := &triggerID
	chatTitle := message.Chat.Title

	// Сначала отправляем poll-сообщение — нам нужен его MessageID для записи.
	pollText := b.formatVotebanPoll(target, message.From, models.VotebanTally{}, votebanRequiredVotes, time.Duration(votebanWindowSeconds)*time.Second)
	pollMsg := tgbotapi.NewMessage(message.Chat.ID, pollText)
	pollMsg.ParseMode = "HTML"
	pollMsg.DisableWebPagePreview = true
	pollMsg.ReplyToMessageID = triggerID
	// Кнопки добавим после получения id записи (они содержат voteban_id).
	sent, err := b.bot.Send(pollMsg)
	if err != nil {
		log.Printf("voteban: send poll failed: %v", err)
		return
	}

	vb, err := b.moderationService.StartVoteban(service.VotebanStartParams{
		ChatID:           message.Chat.ID,
		ChatTitle:        chatTitle,
		TargetUserID:     target.ID,
		TargetUsername:   target.UserName,
		TargetFirstName:  target.FirstName,
		InitiatorUserID:  message.From.ID,
		TriggerMessageID: triggerPtr,
		PollMessageID:    sent.MessageID,
		RequiredVotes:    votebanRequiredVotes,
		MuteSeconds:      votebanKickSeconds, // длительность kick-санкции (БД-колонка из T1)
		WindowSeconds:    votebanWindowSeconds,
	})
	if err != nil {
		// Уже идёт голосование — удаляем только что отправленный poll, оставляем существующий.
		if err == service.ErrVotebanAlreadyOpen {
			b.tryDelete(message.Chat.ID, sent.MessageID)
			b.replyAndAutoDelete(message, "На этого участника уже идёт голосование.")
			return
		}
		log.Printf("voteban: start failed: %v", err)
		b.tryDelete(message.Chat.ID, sent.MessageID)
		return
	}

	// Обновляем poll-сообщение, добавив кнопки с реальным id.
	editMarkup := tgbotapi.NewEditMessageReplyMarkup(message.Chat.ID, sent.MessageID, b.votebanKeyboard(vb.Id, models.VotebanTally{}))
	if _, err := b.bot.Send(editMarkup); err != nil {
		log.Printf("voteban: edit markup failed: %v", err)
	}

	// Инициатор автоматически голосует «за» — один клик меньше.
	if res, err := b.moderationService.CastVote(vb.Id, message.From.ID, models.VotebanVoteFor); err == nil {
		b.refreshVotebanMessage(vb, res.Tally)
		if res.Threshold {
			b.finalizeVotebanPassed(vb)
		}
	}

	b.tryDelete(message.Chat.ID, message.MessageID)
}

func (b *TelegramBot) formatVotebanPoll(target, initiator *tgbotapi.User, tally models.VotebanTally, required int, window time.Duration) string {
	return fmt.Sprintf(
		"⚖️ <b>Голосование за кик</b>\n\n"+
			"Кого: %s\n"+
			"Кто запустил: %s\n"+
			"Окно: %s · нужно «за»: %d · санкция: кик на %s\n\n"+
			"Голосуют участники чата (с активностью за последние 7 дней). Цель не голосует.",
		targetDisplay(target),
		targetDisplay(initiator),
		service.FormatDurationHuman(window),
		required,
		service.FormatDurationHuman(time.Duration(votebanKickSeconds)*time.Second),
	) + b.formatVotebanTally(tally, required)
}

func (b *TelegramBot) formatVotebanTally(tally models.VotebanTally, required int) string {
	return fmt.Sprintf("\n\n✅ За: %d/%d   ❌ Против: %d", tally.For, required, tally.Against)
}

func (b *TelegramBot) votebanKeyboard(votebanID int64, tally models.VotebanTally) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("✅ За (%d)", tally.For), fmt.Sprintf("vb:%d:up", votebanID)),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("❌ Против (%d)", tally.Against), fmt.Sprintf("vb:%d:down", votebanID)),
		),
	)
}

// refreshVotebanMessage обновляет текст и кнопки poll-сообщения с актуальной раскладкой.
func (b *TelegramBot) refreshVotebanMessage(vb *models.Voteban, tally models.VotebanTally) {
	target := &tgbotapi.User{
		ID:        vb.TargetUserID,
		UserName:  vb.TargetUsername,
		FirstName: vb.TargetFirstName,
	}
	initiator := &tgbotapi.User{ID: vb.InitiatorUserID}
	window := time.Until(vb.ExpiresAt)
	if window < 0 {
		window = 0
	}
	text := b.formatVotebanPoll(target, initiator, tally, vb.RequiredVotes, window)

	edit := tgbotapi.NewEditMessageTextAndMarkup(vb.ChatID, vb.PollMessageID, text, b.votebanKeyboard(vb.Id, tally))
	edit.ParseMode = "HTML"
	edit.DisableWebPagePreview = true
	if _, err := b.bot.Send(edit); err != nil {
		log.Printf("voteban: refresh poll failed: %v", err)
	}
}

// handleVotebanCallback обрабатывает нажатие на ✅/❌ под голосованием.
func (b *TelegramBot) handleVotebanCallback(callback *tgbotapi.CallbackQuery) {
	parts := strings.Split(callback.Data, ":")
	if len(parts) != 3 {
		return
	}
	votebanID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return
	}
	var vote int16
	switch parts[2] {
	case "up":
		vote = models.VotebanVoteFor
	case "down":
		vote = models.VotebanVoteAgainst
	default:
		return
	}

	vb, err := b.moderationService.GetVoteban(votebanID)
	if err != nil || vb == nil {
		b.answerCallbackQuery(callback.ID, "Голосование не найдено.")
		return
	}
	if vb.Status != models.VotebanStatusOpen {
		b.answerCallbackQuery(callback.ID, "Голосование уже закрыто.")
		return
	}
	if !b.isChatMember(vb.ChatID, callback.From.ID) {
		b.answerCallbackQuery(callback.ID, "Голосовать могут только участники чата.")
		return
	}
	// Стаж: голосует только тот, кто реально пишет в этом чате. Цель и
	// инициатор автоматически проходят (инициатор уже проверен на cooldown,
	// цель отсекается чуть ниже в CastVote → ErrVoteSelfTarget).
	if callback.From.ID != vb.TargetUserID && callback.From.ID != vb.InitiatorUserID {
		count, _ := b.chatActivityService.CountUserMessagesInChatSince(
			vb.ChatID, callback.From.ID, time.Now().Add(-voterMinActivityWindow))
		if count < int64(voterMinMessages) {
			b.answerCallbackQuery(callback.ID, "Голосовать могут активные участники чата за последние 7 дней.")
			return
		}
	}

	res, err := b.moderationService.CastVote(votebanID, callback.From.ID, vote)
	if err != nil {
		switch err {
		case service.ErrVoteSelfTarget:
			b.answerCallbackQuery(callback.ID, "Цель голосования не может голосовать.")
		case service.ErrVotebanClosed:
			b.answerCallbackQuery(callback.ID, "Голосование уже закрыто.")
		default:
			log.Printf("voteban: cast failed: %v", err)
			b.answerCallbackQuery(callback.ID, "Ошибка.")
		}
		return
	}

	if !res.Changed {
		b.answerCallbackQuery(callback.ID, "Ваш голос уже учтён.")
	} else {
		b.answerCallbackQuery(callback.ID, "Голос принят.")
	}
	b.refreshVotebanMessage(vb, res.Tally)

	if res.Threshold {
		b.finalizeVotebanPassed(vb)
	}
}

// finalizeVotebanPassed применяет санкцию: kick (BanChatMember с UntilDate) на
// vb.MuteSeconds + удаление триггер-сообщения. Идемпотентно: повторный вызов
// после успешной финализации ничего не делает (запись уже не open).
//
// Имя поля MuteSeconds в БД сохранено по совместимости с T1-миграцией,
// фактически в нём лежит длительность санкции.
func (b *TelegramBot) finalizeVotebanPassed(vb *models.Voteban) {
	ok, err := b.moderationService.FinalizeVoteban(vb.Id, models.VotebanStatusPassed)
	if err != nil {
		log.Printf("voteban: finalize-passed failed: %v", err)
		return
	}
	if !ok {
		return
	}

	until := time.Now().Add(time.Duration(vb.MuteSeconds) * time.Second)
	if _, err := b.bot.Request(tgbotapi.BanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: vb.ChatID,
			UserID: vb.TargetUserID,
		},
		UntilDate: until.Unix(),
	}); err != nil {
		log.Printf("voteban: ban failed chat=%d user=%d: %v", vb.ChatID, vb.TargetUserID, err)
	}
	if vb.TriggerMessageID != nil {
		b.tryDelete(vb.ChatID, *vb.TriggerMessageID)
	}

	dur := time.Duration(vb.MuteSeconds) * time.Second
	durSec := vb.MuteSeconds
	expiresAt := until
	_ = b.moderationService.LogActionWithMeta(&models.ModerationAction{
		ChatID:          vb.ChatID,
		TargetUserID:    vb.TargetUserID,
		ActorUserID:     0,
		Action:          models.ModerationActionVotebanKick,
		DurationSeconds: &durSec,
		ExpiresAt:       &expiresAt,
	}, map[string]interface{}{
		"voteban_id": vb.Id,
		"initiator":  vb.InitiatorUserID,
	})

	tally, _ := b.moderationService.CountVotes(vb.Id)
	target := &tgbotapi.User{ID: vb.TargetUserID, UserName: vb.TargetUsername, FirstName: vb.TargetFirstName}
	text := fmt.Sprintf("⚖️ Голосование завершено: %s кикнут из чата на %s (✅ %d / ❌ %d). Авто-возврат после истечения срока.",
		targetDisplay(target), service.FormatDurationHuman(dur), tally.For, tally.Against)
	edit := tgbotapi.NewEditMessageText(vb.ChatID, vb.PollMessageID, text)
	edit.ParseMode = "HTML"
	if _, err := b.bot.Send(edit); err != nil {
		log.Printf("voteban: edit final passed failed: %v", err)
	}
}

// finalizeVotebanFailed закрывает голосование без санкций (истечение окна).
func (b *TelegramBot) finalizeVotebanFailed(vb *models.Voteban) {
	ok, err := b.moderationService.FinalizeVoteban(vb.Id, models.VotebanStatusFailed)
	if err != nil {
		log.Printf("voteban: finalize-failed failed: %v", err)
		return
	}
	if !ok {
		return
	}

	tally, _ := b.moderationService.CountVotes(vb.Id)
	target := &tgbotapi.User{ID: vb.TargetUserID, UserName: vb.TargetUsername, FirstName: vb.TargetFirstName}
	text := fmt.Sprintf("⚖️ Голосование закрыто: голосов недостаточно (✅ %d / ❌ %d). %s остаётся в чате.",
		tally.For, tally.Against, targetDisplay(target))
	edit := tgbotapi.NewEditMessageText(vb.ChatID, vb.PollMessageID, text)
	edit.ParseMode = "HTML"
	if _, err := b.bot.Send(edit); err != nil {
		log.Printf("voteban: edit final failed: %v", err)
	}
}

// startVotebanWatcher финализирует протёкшие голосования и шлёт алерты
// «срок санкции истёк» для ban/mute/voteban_kick. Telegram сам снимает ban
// при наступлении until_date — нам нужно только уведомить чат.
func (b *TelegramBot) startVotebanWatcher() {
	ticker := time.NewTicker(moderationWatcherTick)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()

		// 1) Финализация голосований по истечении окна.
		expired, err := b.moderationService.ListExpiredOpenVotebans(now)
		if err != nil {
			log.Printf("voteban-watcher: list failed: %v", err)
		} else {
			for i := range expired {
				vb := expired[i]
				tally, _ := b.moderationService.CountVotes(vb.Id)
				if tally.For >= vb.RequiredVotes {
					b.finalizeVotebanPassed(&vb)
				} else {
					b.finalizeVotebanFailed(&vb)
				}
			}
		}

		// 2) Алерты «срок санкции истёк» для ban/mute/voteban.
		actions, err := b.moderationService.ListExpiredUnnotifiedActions(now)
		if err != nil {
			log.Printf("expiry-watcher: list failed: %v", err)
			continue
		}
		for i := range actions {
			b.notifyActionExpired(&actions[i])
		}
	}
}

// notifyActionExpired шлёт сообщение в чат «срок санкции истёк» и помечает
// запись expired_notified_at, чтобы не дублировать.
func (b *TelegramBot) notifyActionExpired(action *models.ModerationAction) {
	if action.ChatID == 0 {
		// Глобальные баны (chat_id = 0) — отдельный flow в #294 (там же добавим алерты).
		_ = b.moderationService.MarkActionExpiredNotified(action.Id)
		return
	}
	verb := "разблокирован"
	switch action.Action {
	case models.ModerationActionBan, models.ModerationActionVotebanKick, models.ModerationActionVotebanMute:
		verb = "снова в чате"
	case models.ModerationActionMute:
		verb = "может писать снова"
	}
	target := &tgbotapi.User{ID: action.TargetUserID}
	b.sendChatHTML(action.ChatID, fmt.Sprintf("⏰ Срок санкции истёк — %s %s.", targetDisplay(target), verb))

	if err := b.moderationService.MarkActionExpiredNotified(action.Id); err != nil {
		log.Printf("expiry-watcher: mark notified failed action=%d: %v", action.Id, err)
	}
}

// --- helpers ---

// sendChatHTML отправляет HTML-сообщение в чат, без preview.
func (b *TelegramBot) sendChatHTML(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true
	if _, err := b.bot.Send(msg); err != nil {
		log.Printf("sendChatHTML failed chat=%d: %v", chatID, err)
	}
}

// tryDelete удаляет сообщение, ошибки только логирует.
func (b *TelegramBot) tryDelete(chatID int64, messageID int) {
	if _, err := b.bot.Request(tgbotapi.NewDeleteMessage(chatID, messageID)); err != nil {
		// «message to delete not found» — нормальный исход (уже удалено).
		if !strings.Contains(err.Error(), "message to delete not found") {
			log.Printf("tryDelete chat=%d msg=%d: %v", chatID, messageID, err)
		}
	}
}

// replyAndAutoDelete отвечает на команду подсказкой и через 15 сек удаляет
// и подсказку, и саму команду — чтобы не засорять чат.
func (b *TelegramBot) replyAndAutoDelete(message *tgbotapi.Message, text string) {
	reply := tgbotapi.NewMessage(message.Chat.ID, text)
	reply.ReplyToMessageID = message.MessageID
	sent, err := b.bot.Send(reply)
	if err != nil {
		log.Printf("replyAndAutoDelete send failed: %v", err)
		return
	}
	go func() {
		time.Sleep(15 * time.Second)
		b.tryDelete(message.Chat.ID, sent.MessageID)
		b.tryDelete(message.Chat.ID, message.MessageID)
	}()
}
