package bot

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"strconv"
	"strings"
	"time"

	"ithozyeva/config"
	"ithozyeva/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const msgAccessRevoked = "Уровень вашей подписки изменился. Доступ к некоторым чатам был отозван."

// strPtr returns a pointer to s if non-empty, nil otherwise.
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// notifyUserOfSyncResult sends links for granted chats and revocation message for revoked ones.
func (b *TelegramBot) notifyUserOfSyncResult(userID int64, result *service.SyncResult) {
	if len(result.Granted) > 0 {
		b.sendSubscriptionLinks(userID, result)
	}
	if len(result.Revoked) > 0 {
		b.SendDirectMessage(userID, msgAccessRevoked)
	}
}

// isAdmin checks if a Telegram user ID belongs to a member with the ADMIN role.
func (b *TelegramBot) isAdmin(userID int64) bool {
	return b.member.IsAdminByTelegramID(userID)
}

// subscriptionAdminID is the only Telegram user allowed to manage subscriptions via bot.
const subscriptionAdminID int64 = 931916742

// isSubscriptionAdmin checks if the user is allowed to manage subscriptions.
func (b *TelegramBot) isSubscriptionAdmin(userID int64) bool {
	return userID == subscriptionAdminID
}

// --- Telegram API helpers for subscription system ---

// isChatMember checks if a user is in a specific chat via Telegram API.
func (b *TelegramBot) isChatMember(chatID, userID int64) bool {
	member, err := b.bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: userID,
		},
	})
	if err != nil {
		log.Printf("Failed to check membership: chat=%d user=%d: %v", chatID, userID, err)
		return false
	}
	return isActiveMemberStatus(member.Status)
}

// createOneTimeInviteLink creates a single-use invite link for a chat.
func (b *TelegramBot) createOneTimeInviteLink(chatID int64) (string, error) {
	link, err := b.bot.Request(tgbotapi.CreateChatInviteLinkConfig{
		ChatConfig: tgbotapi.ChatConfig{ChatID: chatID},
		MemberLimit: 1,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create invite link for chat %d: %w", chatID, err)
	}

	var result struct {
		InviteLink string `json:"invite_link"`
	}
	if err := parseAPIResponse(link, &result); err != nil {
		return "", err
	}
	return result.InviteLink, nil
}

// kickFromChat kicks a user by ban+unban.
func (b *TelegramBot) kickFromChat(chatID, userID int64) {
	_, err := b.bot.Request(tgbotapi.BanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
	})
	if err != nil {
		log.Printf("Failed to ban user %d from chat %d: %v", userID, chatID, err)
		return
	}

	time.Sleep(1 * time.Second)

	_, err = b.bot.Request(tgbotapi.UnbanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
		OnlyIfBanned: true,
	})
	if err != nil {
		log.Printf("Failed to unban user %d from chat %d: %v", userID, chatID, err)
	}
}

// botCheckFunc returns a closure for the subscription service.
func (b *TelegramBot) botCheckFunc() func(int64, int64) bool {
	return func(chatID, userID int64) bool {
		return b.subscriptionService.IsMember(chatID, userID, b.isChatMember)
	}
}

// createInviteLinkFunc returns a closure for creating invite links.
func (b *TelegramBot) createInviteLinkFunc() func(int64) (string, error) {
	return b.createOneTimeInviteLink
}

// kickUserFunc returns a closure for kicking users.
func (b *TelegramBot) kickUserFunc() func(int64, int64) {
	return b.kickFromChat
}

// --- User commands ---

// handleSubCommand checks subscription and grants access.
func (b *TelegramBot) handleSubCommand(message *tgbotapi.Message) {
	user := message.From
	result, err := b.subscriptionService.OnboardUser(
		user.ID, strPtr(user.UserName), user.FirstName+" "+user.LastName,
		b.botCheckFunc(), b.createInviteLinkFunc(), b.kickUserFunc(),
	)
	if err != nil {
		log.Printf("Error onboarding user %d: %v", user.ID, err)
		b.sendMessage(message.Chat.ID, "Произошла ошибка. Попробуйте позже.")
		return
	}

	if result.EffectiveTierID == nil {
		b.sendMessage(message.Chat.ID,
			"У вас нет активной подписки.\n\n"+
				"Подпишитесь через Boosty или Tribute, затем вернитесь и нажмите /sub снова.")
		return
	}

	if len(result.Granted) > 0 {
		b.sendSubscriptionLinks(message.Chat.ID, result)
	} else {
		b.sendMessage(message.Chat.ID,
			"Ваша подписка активна! У вас уже есть доступ ко всем чатам.\n"+
				"Используйте /substatus для просмотра.")
	}
}

// handleSubStatusCommand shows current subscription status.
func (b *TelegramBot) handleSubStatusCommand(message *tgbotapi.Message) {
	user, err := b.subscriptionService.GetUser(message.From.ID)
	if err != nil {
		b.sendMessage(message.Chat.ID, "Вы не зарегистрированы. Используйте /sub для начала.")
		return
	}

	tierID := user.EffectiveTierID()
	tierName := "Нет"
	if tierID != nil {
		tier, err := b.subscriptionService.GetTier(*tierID)
		if err == nil {
			tierName = tier.Name
		}
	}

	access, _ := b.subscriptionService.GetActiveAccess(message.From.ID)

	text := fmt.Sprintf("<b>Статус подписки</b>\n\n"+
		"Тир: %s\n"+
		"Активных чатов: %d\n", tierName, len(access))

	if len(access) > 0 {
		text += "\nДоступные чаты:\n"
		for _, a := range access {
			chat, err := b.subscriptionService.GetChat(a.ChatID)
			if err == nil {
				text += fmt.Sprintf("  • %s\n", chat.Title)
			}
		}
	}

	b.SendDirectMessage(message.Chat.ID, text)
}

// subscriptionDeepLink returns a t.me link that launches the bot with /start sub.
func (b *TelegramBot) subscriptionDeepLink() string {
	return fmt.Sprintf("https://t.me/%s?start=sub", b.bot.Self.UserName)
}

// postAnchorWelcome posts a welcome message in the anchor chat with a button
// that deep-links to the bot so the user can receive DM invite links.
func (b *TelegramBot) postAnchorWelcome(chatID int64, user *tgbotapi.User) {
	mention := "@" + user.UserName
	if user.UserName == "" {
		name := strings.TrimSpace(user.FirstName + " " + user.LastName)
		if name == "" {
			name = "друг"
		}
		mention = fmt.Sprintf("<a href=\"tg://user?id=%d\">%s</a>", user.ID, name)
	}

	text := fmt.Sprintf(
		"%s, добро пожаловать! Нажмите кнопку ниже, чтобы получить доступ к остальным чатам.",
		mention)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.DisableNotification = true
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Получить доступ", b.subscriptionDeepLink()),
		),
	)
	if _, err := b.bot.Send(msg); err != nil {
		log.Printf("Failed to post anchor welcome in chat %d: %v", chatID, err)
	}
}

// sendSubscriptionLinks sends invite links grouped by category as a single HTML message.
func (b *TelegramBot) sendSubscriptionLinks(chatID int64, result *service.SyncResult) {
	type entry struct {
		title string
		link  string
	}
	groups := make(map[string][]entry)
	groupEmoji := make(map[string]string)
	var order []string

	for _, g := range result.Granted {
		chat, err := b.subscriptionService.GetChat(g.ChatID)
		title := fmt.Sprintf("Chat %d", g.ChatID)
		cat := "Прочее"
		emoji := "💬"
		if err == nil {
			title = chat.Title
			if chat.Category != nil && *chat.Category != "" {
				cat = *chat.Category
			}
			if chat.Emoji != nil && *chat.Emoji != "" {
				emoji = *chat.Emoji
			}
		}
		if _, exists := groups[cat]; !exists {
			order = append(order, cat)
			groupEmoji[cat] = emoji
		}
		groups[cat] = append(groups[cat], entry{title: title, link: g.Link})
	}

	text := fmt.Sprintf("Подписка подтверждена! Доступно чатов: <b>%d</b>\n", len(result.Granted))
	for _, cat := range order {
		text += fmt.Sprintf("\n%s <b>%s</b>\n", groupEmoji[cat], cat)
		for _, e := range groups[cat] {
			text += fmt.Sprintf("• <a href=\"%s\">%s</a>\n", e.link, e.title)
		}
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true
	if _, err := b.bot.Send(msg); err != nil {
		log.Printf("Failed to send subscription links to %d: %v", chatID, err)
	}
}

// --- Chat member event handlers ---

// handleChatMemberUpdated reacts to user join/leave in anchor chats.
func (b *TelegramBot) handleChatMemberUpdated(update *tgbotapi.ChatMemberUpdated) {
	chat, err := b.subscriptionService.GetChat(update.Chat.ID)
	if err != nil || chat.AnchorForTierID == nil {
		return // Not an anchor chat
	}

	userID := update.NewChatMember.User.ID
	oldActive := isActiveMemberStatus(update.OldChatMember.Status)
	newActive := isActiveMemberStatus(update.NewChatMember.Status)

	if oldActive == newActive {
		return
	}

	log.Printf("Anchor chat member change: chat=%d user=%d %s->%s",
		update.Chat.ID, userID, update.OldChatMember.Status, update.NewChatMember.Status)

	// Invalidate membership cache
	b.subscriptionService.InvalidateMemberCache(update.Chat.ID, userID)

	// Ensure user exists and sync access
	tgUser := update.NewChatMember.User
	usernamePtr := strPtr(tgUser.UserName)
	result, err := b.subscriptionService.OnboardUser(
		userID, usernamePtr, tgUser.FirstName+" "+tgUser.LastName,
		b.botCheckFunc(), b.createInviteLinkFunc(), b.kickUserFunc(),
	)
	if err != nil {
		return
	}

	// Post welcome in anchor only on join (not on leave/kick).
	// Telegram не даёт боту писать первым в ЛС: сообщение в anchor с deep-link
	// — единственный способ заставить юзера один раз стартануть бота.
	if newActive && !oldActive {
		b.postAnchorWelcome(update.Chat.ID, tgUser)
	}

	b.notifyUserOfSyncResult(userID, result)
}

// handleMyChatMemberUpdated handles bot being added/removed from chats.
func (b *TelegramBot) handleMyChatMemberUpdated(update *tgbotapi.ChatMemberUpdated) {
	chat := update.Chat
	if chat.Type != "group" && chat.Type != "supergroup" && chat.Type != "channel" {
		return
	}

	newStatus := update.NewChatMember.Status
	oldStatus := update.OldChatMember.Status

	// Bot added to chat
	if !isActiveMemberStatus(oldStatus) && isActiveMemberStatus(newStatus) {
		title := chat.Title
		if title == "" {
			title = fmt.Sprintf("Chat %d", chat.ID)
		}
		err := b.subscriptionService.UpsertChat(chat.ID, title, chat.Type)
		if err != nil {
			log.Printf("Failed to register chat %d: %v", chat.ID, err)
			return
		}
		log.Printf("Bot added to chat %d (%s), registered in DB", chat.ID, title)

		// Бот должен быть администратором, чтобы получать chat_member updates
		// и создавать invite-ссылки. Если добавили не админом — предупреждаем.
		if newStatus != "administrator" {
			log.Printf("WARNING: bot is not an administrator in chat %d (status=%s) — chat_member updates and invite links will not work", chat.ID, newStatus)
		}

		// Notify only the subscription admin (не всех админов платформы).
		// Всё идёт под HTML parse-mode — placeholder-ы и пользовательский ввод экранируем.
		addedBy := ""
		if update.From.UserName != "" {
			addedBy = fmt.Sprintf("\nДобавил: @%s", html.EscapeString(update.From.UserName))
		} else if update.From.ID != 0 {
			addedBy = fmt.Sprintf("\nДобавил: id=%d", update.From.ID)
		}
		b.SendDirectMessage(subscriptionAdminID, fmt.Sprintf(
			"Бот добавлен в чат:\nID: <code>%d</code>\nНазвание: %s%s\n\n"+
				"Настройте роль через:\n"+
				"/subaddchat %d &lt;tier_slug&gt; — content чат\n"+
				"/subaddchat %d &lt;tier_slug&gt; anchor — anchor чат",
			chat.ID, html.EscapeString(title), addedBy, chat.ID, chat.ID))
	}

	// Bot removed from chat
	if isActiveMemberStatus(oldStatus) && !isActiveMemberStatus(newStatus) {
		b.subscriptionService.DeleteChat(chat.ID)
		log.Printf("Bot removed from chat %d, deleted from DB", chat.ID)
	}
}

func isActiveMemberStatus(status string) bool {
	return status == "creator" || status == "administrator" || status == "member" || status == "restricted"
}

// --- Periodic subscription checker ---

func (b *TelegramBot) startSubscriptionChecker() {
	interval := time.Duration(config.CFG.SubscriptionCheckIntervalHours) * time.Hour
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		b.subscriptionService.PeriodicCheck(
			b.botCheckFunc(),
			b.createInviteLinkFunc(),
			b.kickUserFunc(),
			b.notifyUserOfSyncResult,
			50*time.Millisecond,
		)
	}
}

// --- Admin commands ---

func (b *TelegramBot) handleSubTiersCommand(message *tgbotapi.Message) {
	if !b.isSubscriptionAdmin(message.From.ID) {
		return
	}
	tiers, err := b.subscriptionService.GetAllTiers()
	if err != nil {
		b.sendMessage(message.Chat.ID, "Ошибка получения тиров.")
		return
	}

	tierCounts, _ := b.subscriptionService.CountAllUsersByTier()
	text := "<b>Тиры подписок:</b>\n\n"
	for _, t := range tiers {
		text += fmt.Sprintf("Level %d: <b>%s</b> (%s) — %d пользователей\n", t.Level, t.Name, t.Slug, tierCounts[t.ID])
	}

	b.SendDirectMessage(message.Chat.ID, text)
}

func (b *TelegramBot) handleSubChatsCommand(message *tgbotapi.Message) {
	if !b.isSubscriptionAdmin(message.From.ID) {
		return
	}
	chats, err := b.subscriptionService.GetAllChats()
	if err != nil || len(chats) == 0 {
		b.sendMessage(message.Chat.ID, "Нет зарегистрированных чатов.")
		return
	}

	text := "<b>Зарегистрированные чаты:</b>\n\n"
	for _, c := range chats {
		role := ""
		if c.AnchorForTierID != nil {
			tier, err := b.subscriptionService.GetTier(*c.AnchorForTierID)
			if err == nil {
				role = fmt.Sprintf(" [ANCHOR → %s]", tier.Name)
			}
		}
		text += fmt.Sprintf("<code>%d</code> — %s%s\n", c.ID, c.Title, role)
	}

	b.SendDirectMessage(message.Chat.ID, text)
}

func (b *TelegramBot) handleSubAddChatCommand(message *tgbotapi.Message) {
	if !b.isSubscriptionAdmin(message.From.ID) {
		return
	}

	args := strings.Fields(message.Text)
	// /subaddchat <chat_id> <tier_slug> [anchor]
	if len(args) < 3 {
		b.sendMessage(message.Chat.ID,
			"Использование:\n"+
				"/subaddchat <chat_id> <tier_slug> — content чат\n"+
				"/subaddchat <chat_id> <tier_slug> anchor — anchor чат")
		return
	}

	chatID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		b.sendMessage(message.Chat.ID, "Неверный chat_id.")
		return
	}

	tier, err := b.subscriptionService.GetTierBySlug(args[2])
	if err != nil {
		b.sendMessage(message.Chat.ID, fmt.Sprintf("Тир '%s' не найден.", args[2]))
		return
	}

	isAnchor := len(args) > 3 && args[3] == "anchor"

	// Ensure chat exists in DB, fetch real title from Telegram
	title := fmt.Sprintf("Chat %d", chatID)
	chatType := "supergroup"
	chatConfig := tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: chatID}}
	if tgChat, err := b.bot.GetChat(chatConfig); err == nil {
		if tgChat.Title != "" {
			title = tgChat.Title
		}
		if tgChat.Type != "" {
			chatType = tgChat.Type
		}
	} else if chat, err := b.subscriptionService.GetChat(chatID); err == nil {
		title = chat.Title
		chatType = chat.ChatType
	}
	b.subscriptionService.UpsertChat(chatID, title, chatType)

	if isAnchor {
		b.subscriptionService.SetAnchor(chatID, &tier.ID)
		b.SendDirectMessage(message.Chat.ID, fmt.Sprintf(
			"Чат <code>%d</code> установлен как <b>anchor</b> для тира %s.", chatID, tier.Name))
	} else {
		b.subscriptionService.AddChatToTier(chatID, tier.ID)
		b.SendDirectMessage(message.Chat.ID, fmt.Sprintf(
			"Чат <code>%d</code> добавлен как <b>content</b> для тира %s.", chatID, tier.Name))
	}
}

func (b *TelegramBot) handleSubSetAnchorCommand(message *tgbotapi.Message) {
	if !b.isSubscriptionAdmin(message.From.ID) {
		return
	}

	args := strings.Fields(message.Text)
	if len(args) < 3 {
		b.sendMessage(message.Chat.ID, "Использование: /subsetanchor <chat_id> <tier_slug|clear>")
		return
	}

	chatID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		b.sendMessage(message.Chat.ID, "Неверный chat_id.")
		return
	}

	if args[2] == "clear" {
		b.subscriptionService.SetAnchor(chatID, nil)
		b.SendDirectMessage(message.Chat.ID, fmt.Sprintf("Anchor снят с чата <code>%d</code>.", chatID))
		return
	}

	tier, err := b.subscriptionService.GetTierBySlug(args[2])
	if err != nil {
		b.sendMessage(message.Chat.ID, fmt.Sprintf("Тир '%s' не найден.", args[2]))
		return
	}

	b.subscriptionService.SetAnchor(chatID, &tier.ID)
	b.SendDirectMessage(message.Chat.ID, fmt.Sprintf(
		"Чат <code>%d</code> теперь anchor для тира %s.", chatID, tier.Name))
}

func (b *TelegramBot) handleSubRemoveChatCommand(message *tgbotapi.Message) {
	if !b.isSubscriptionAdmin(message.From.ID) {
		return
	}

	args := strings.Fields(message.Text)
	if len(args) < 2 {
		b.sendMessage(message.Chat.ID, "Использование: /subremovechat <chat_id>")
		return
	}

	chatID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		b.sendMessage(message.Chat.ID, "Неверный chat_id.")
		return
	}

	b.subscriptionService.DeleteChat(chatID)
	b.SendDirectMessage(message.Chat.ID, fmt.Sprintf("Чат <code>%d</code> удалён.", chatID))
}

func (b *TelegramBot) handleSubUsersCommand(message *tgbotapi.Message) {
	if !b.isSubscriptionAdmin(message.From.ID) {
		return
	}

	args := strings.Fields(message.Text)
	page := 0
	if len(args) > 1 {
		if p, err := strconv.Atoi(args[1]); err == nil {
			page = p
		}
	}

	pageSize := 20
	total, _ := b.subscriptionService.CountAllUsers()
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if totalPages == 0 {
		totalPages = 1
	}
	users, _ := b.subscriptionService.GetPaginatedUsers(page*pageSize, pageSize)

	text := fmt.Sprintf("<b>Пользователи</b> (стр. %d/%d, всего: %d)\n\n", page+1, totalPages, total)
	for _, u := range users {
		tierInfo := ""
		if u.ManualTierID != nil {
			tierInfo = fmt.Sprintf(" [manual:%d]", *u.ManualTierID)
		} else if u.ResolvedTierID != nil {
			tierInfo = fmt.Sprintf(" [tier:%d]", *u.ResolvedTierID)
		}
		active := "+"
		if !u.IsActive {
			active = "-"
		}
		usernameStr := "?"
		if u.Username != nil {
			usernameStr = *u.Username
		}
		text += fmt.Sprintf("%s <code>%d</code> @%s%s\n", active, u.ID, usernameStr, tierInfo)
	}

	b.SendDirectMessage(message.Chat.ID, text)
}

func (b *TelegramBot) handleSubUserInfoCommand(message *tgbotapi.Message) {
	if !b.isSubscriptionAdmin(message.From.ID) {
		return
	}

	args := strings.Fields(message.Text)
	if len(args) < 2 {
		b.sendMessage(message.Chat.ID, "Использование: /subuserinfo <user_id>")
		return
	}

	userID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		b.sendMessage(message.Chat.ID, "Неверный user_id.")
		return
	}

	user, err := b.subscriptionService.GetUser(userID)
	if err != nil {
		b.sendMessage(message.Chat.ID, "Пользователь не найден.")
		return
	}

	effTierName := "Нет"
	effTierID := user.EffectiveTierID()
	if effTierID != nil {
		if tier, err := b.subscriptionService.GetTier(*effTierID); err == nil {
			effTierName = tier.Name
		}
	}

	access, _ := b.subscriptionService.GetActiveAccess(userID)
	usernameStr := "?"
	if user.Username != nil {
		usernameStr = *user.Username
	}
	lastCheck := "никогда"
	if user.LastCheckAt != nil {
		lastCheck = user.LastCheckAt.Format("2006-01-02 15:04")
	}

	text := fmt.Sprintf("<b>Инфо о пользователе</b>\n\n"+
		"ID: <code>%d</code>\n"+
		"Username: @%s\n"+
		"Имя: %s\n"+
		"Активен: %v\n"+
		"Resolved tier: %v\n"+
		"Manual tier: %v\n"+
		"Effective tier: %s\n"+
		"Последняя проверка: %s\n"+
		"Активных чатов: %d\n",
		user.ID, usernameStr, user.FullName, user.IsActive,
		user.ResolvedTierID, user.ManualTierID, effTierName,
		lastCheck, len(access))

	b.SendDirectMessage(message.Chat.ID, text)
}

func (b *TelegramBot) handleSubOverrideCommand(message *tgbotapi.Message) {
	if !b.isSubscriptionAdmin(message.From.ID) {
		return
	}

	args := strings.Fields(message.Text)
	if len(args) < 3 {
		b.sendMessage(message.Chat.ID, "Использование: /suboverride <user_id> <tier_slug|clear>")
		return
	}

	userID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		b.sendMessage(message.Chat.ID, "Неверный user_id.")
		return
	}

	if _, err := b.subscriptionService.GetUser(userID); err != nil {
		b.sendMessage(message.Chat.ID, "Пользователь не найден.")
		return
	}

	if args[2] == "clear" {
		b.subscriptionService.SetManualTier(userID, nil)
		b.subscriptionService.AddAudit(userID, "manual_override", map[string]interface{}{
			"tier": nil, "by": message.From.ID,
		})
		b.SendDirectMessage(message.Chat.ID, fmt.Sprintf(
			"Ручной тир снят для пользователя <code>%d</code>.", userID))
	} else {
		tier, err := b.subscriptionService.GetTierBySlug(args[2])
		if err != nil {
			b.sendMessage(message.Chat.ID, fmt.Sprintf("Тир '%s' не найден.", args[2]))
			return
		}
		b.subscriptionService.SetManualTier(userID, &tier.ID)
		b.subscriptionService.AddAudit(userID, "manual_override", map[string]interface{}{
			"tier_id": tier.ID, "tier_slug": tier.Slug, "by": message.From.ID,
		})
		b.SendDirectMessage(message.Chat.ID, fmt.Sprintf(
			"Пользователь <code>%d</code> установлен на тир %s.", userID, tier.Name))
	}

	// Re-sync
	result, err := b.subscriptionService.CheckAndSyncUser(
		userID, b.botCheckFunc(), b.createInviteLinkFunc(), b.kickUserFunc(),
	)
	if err == nil {
		b.notifyUserOfSyncResult(userID, result)
	}
}

func (b *TelegramBot) handleSubCheckAllCommand(message *tgbotapi.Message) {
	if !b.isSubscriptionAdmin(message.From.ID) {
		return
	}

	b.sendMessage(message.Chat.ID, "Запуск полной проверки подписок...")

	go func() {
		b.subscriptionService.PeriodicCheck(
			b.botCheckFunc(),
			b.createInviteLinkFunc(),
			b.kickUserFunc(),
			b.notifyUserOfSyncResult,
			50*time.Millisecond,
		)
		b.SendDirectMessage(message.Chat.ID, "Проверка подписок завершена.")
	}()
}

func (b *TelegramBot) handleSubStatsCommand(message *tgbotapi.Message) {
	if !b.isSubscriptionAdmin(message.From.ID) {
		return
	}

	total, _ := b.subscriptionService.CountAllUsers()
	tiers, _ := b.subscriptionService.GetAllTiers()

	tierCounts, _ := b.subscriptionService.CountAllUsersByTier()
	text := fmt.Sprintf("<b>Статистика подписок</b>\n\nВсего пользователей: %d\n\n", total)
	for _, t := range tiers {
		text += fmt.Sprintf("%s: %d\n", t.Name, tierCounts[t.ID])
	}

	b.SendDirectMessage(message.Chat.ID, text)
}

// handleSubPinCommand posts and pins a welcome message in an anchor chat
// so existing members can click the button to receive DM invite links.
func (b *TelegramBot) handleSubPinCommand(message *tgbotapi.Message) {
	if !b.isSubscriptionAdmin(message.From.ID) {
		return
	}

	args := strings.Fields(message.Text)
	if len(args) < 2 {
		b.sendMessage(message.Chat.ID, "Использование: /subpin <anchor_chat_id>")
		return
	}

	chatID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		b.sendMessage(message.Chat.ID, "Неверный chat_id.")
		return
	}

	chat, err := b.subscriptionService.GetChat(chatID)
	if err != nil || chat.AnchorForTierID == nil {
		b.sendMessage(message.Chat.ID, "Чат не зарегистрирован как anchor.")
		return
	}

	text := "Добро пожаловать! Нажмите кнопку ниже, чтобы получить доступ к остальным чатам по вашей подписке."
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Получить доступ", b.subscriptionDeepLink()),
		),
	)
	sent, err := b.bot.Send(msg)
	if err != nil {
		b.sendMessage(message.Chat.ID, fmt.Sprintf("Не удалось запостить: %v", err))
		return
	}

	pinCfg := tgbotapi.PinChatMessageConfig{
		ChatID:              chatID,
		MessageID:           sent.MessageID,
		DisableNotification: false,
	}
	if _, err := b.bot.Request(pinCfg); err != nil {
		b.sendMessage(message.Chat.ID, fmt.Sprintf("Запостил, но не удалось закрепить: %v", err))
		return
	}

	b.SendDirectMessage(message.Chat.ID, fmt.Sprintf("Сообщение запощено и закреплено в чате <code>%d</code>.", chatID))
}

// parseAPIResponse parses the Telegram API response into the target struct.
func parseAPIResponse(resp *tgbotapi.APIResponse, target interface{}) error {
	if resp == nil || resp.Result == nil {
		return fmt.Errorf("empty API response")
	}
	data, err := resp.Result.MarshalJSON()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}
