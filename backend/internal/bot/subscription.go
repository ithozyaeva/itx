package bot

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"ithozyeva/config"
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const msgAccessRevoked = "Уровень вашей подписки изменился. Доступ к некоторым чатам был отозван."

// resolveUserTier возвращает имя тира и флаг «подписчик» для юзера.
// Используется при онбординге в боте — определяет, какой welcome-wizard
// показать (полный для подписчиков, прогрев с тарифами для UNSUBSCRIBER).
// Источник правды — subscription_users.EffectiveTierID().
func (b *TelegramBot) resolveUserTier(userID int64) (tierName string, isSubscriber bool) {
	user, err := b.subscriptionService.GetSubscriptionUser(userID)
	if err != nil || user == nil {
		return "", false
	}
	tierID := user.EffectiveTierID()
	if tierID == nil {
		return "", false
	}
	tier, err := b.subscriptionService.GetTier(*tierID)
	if err != nil || tier == nil {
		// Запись о тире есть, но название не достали — всё равно подписчик.
		return "", true
	}
	return tier.Name, true
}

// sendTariffsMessage показывает публичные тарифы с кнопками на Boosty.
// Источник цен/ссылок — subscription_tiers.is_public=true (см. service.GetPublicTiers).
// Используется в welcome-wizard для UNSUBSCRIBER и доступен подписчикам через
// callback wiz:tariffs (например, чтобы апгрейдиться).
func (b *TelegramBot) sendTariffsMessage(chatID int64) {
	tiers, err := b.subscriptionService.GetPublicTiers()
	if err != nil || len(tiers) == 0 {
		b.sendMessage(chatID, "Не удалось загрузить тарифы. Напишите админу: /support")
		return
	}

	var text strings.Builder
	text.WriteString("<b>Тарифы IT-ХОЗЯЕВА</b>\n\n")
	for _, t := range tiers {
		text.WriteString("💎 <b>")
		text.WriteString(html.EscapeString(t.Name))
		text.WriteString("</b>")
		if t.Price > 0 {
			text.WriteString(fmt.Sprintf(" — %d ₽/мес", t.Price))
		}
		if t.Description != "" {
			text.WriteString(". ")
			text.WriteString(html.EscapeString(t.Description))
		}
		text.WriteString("\n")
	}
	text.WriteString("\nПосле оплаты нажмите /sub — бот выдаст инвайты в чаты по вашему тиру.")

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(tiers)+1)
	for _, t := range tiers {
		if t.BoostyURL == "" {
			continue
		}
		label := fmt.Sprintf("%s — %d ₽", t.Name, t.Price)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(label, t.BoostyURL),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✅ Я оплатил → проверить", "wiz:sub"),
	))

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	if _, err := b.bot.Send(msg); err != nil {
		log.Printf("sendTariffsMessage: send failed for chat %d: %v", chatID, err)
	}
}

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

// subscriptionAdminID — единственный Telegram-пользователь, которому разрешено
// управлять подписками через бот. Источник правды — config.CFG.SuperAdminTelegramID
// (env SUPER_ADMIN_TELEGRAM_ID), с fallback на исторический id для безопасности,
// если конфиг ещё не инициализирован.
func subscriptionAdminID() int64 {
	if config.CFG != nil && config.CFG.SuperAdminTelegramID != 0 {
		return config.CFG.SuperAdminTelegramID
	}
	return 931916742
}

// isSubscriptionAdmin checks if the user is allowed to manage subscriptions.
func (b *TelegramBot) isSubscriptionAdmin(userID int64) bool {
	return userID == subscriptionAdminID()
}

// --- Telegram API helpers for subscription system ---

// isChatMember checks if a user is in a specific chat via Telegram API.
//
// Возвращает (bool, error). Ошибку API (rate-limit/таймаут/network)
// пропускаем наверх отдельно от «честного not-member», чтобы вызывающий
// код мог отличить «юзер вышел» от «не дозвонились» — это решает проблему
// ложного понижения тира при rate-limit'е (кончалось reverse-кик'ом из
// master-only чатов на следующем периодике).
func (b *TelegramBot) isChatMember(chatID, userID int64) (bool, error) {
	member, err := b.bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: userID,
		},
	})
	if err != nil {
		log.Printf("Failed to check membership: chat=%d user=%d: %v", chatID, userID, err)
		return false, err
	}
	return isActiveMemberStatus(member.Status), nil
}

// createOneTimeInviteLink creates a single-use invite link for a chat.
func (b *TelegramBot) createOneTimeInviteLink(chatID int64) (string, error) {
	return b.createInviteLinkWithLimit(chatID, 1)
}

// createInviteLinkWithLimit создаёт invite-link с заданным member_limit.
// memberLimit=1 — эквивалент старой createOneTimeInviteLink.
// Для массовых рассылок шлём одну ссылку с limit=len(users), чтобы не
// упираться в Telegram rate-limit на createChatInviteLink (~20/мин на чат).
// memberLimit=0 означает ссылку без ограничения (до 99999 юзеров).
func (b *TelegramBot) createInviteLinkWithLimit(chatID int64, memberLimit int) (string, error) {
	link, err := b.bot.Request(tgbotapi.CreateChatInviteLinkConfig{
		ChatConfig:  tgbotapi.ChatConfig{ChatID: chatID},
		MemberLimit: memberLimit,
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
// Управляется фича-флагом SUBSCRIPTION_AUTO_KICK_ENABLED: если он выключен,
// бот возвращает false, и вызывающий код (CheckAndSyncUser) пропускает
// запись revoke в БД и нотификацию — это полноценный dry-run всей цепочки,
// а не только Telegram-вызова. Возвращаем true только если кик реально
// прошёл (или хотя бы был отправлен запрос).
func (b *TelegramBot) kickFromChat(chatID, userID int64) bool {
	if config.CFG == nil || !config.CFG.SubscriptionAutoKickEnabled {
		log.Printf("[auto-kick disabled] dry-run: would kick user %d from chat %d", userID, chatID)
		return false
	}

	_, err := b.bot.Request(tgbotapi.BanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
	})
	if err != nil {
		log.Printf("Failed to ban user %d from chat %d: %v", userID, chatID, err)
		return false
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
	return true
}

// botCheckFunc returns a closure for the subscription service.
func (b *TelegramBot) botCheckFunc() service.MemberCheckFunc {
	return func(chatID, userID int64) (bool, error) {
		return b.subscriptionService.IsMember(chatID, userID, b.isChatMember)
	}
}

// uncachedCheckFunc — прямой getChatMember без Redis-кэша. Нужен sweep-у
// и dry-run, иначе после первого прохода они будут читать стейл-данные
// до 5 минут (TTL membership-кэша).
func (b *TelegramBot) uncachedCheckFunc() service.MemberCheckFunc {
	return b.isChatMember
}

// createInviteLinkFunc returns a closure for creating invite links.
func (b *TelegramBot) createInviteLinkFunc() func(int64) (string, error) {
	return b.createOneTimeInviteLink
}

// kickUserFunc returns a closure for kicking users.
// Bool-возвращаемое значение — флаг «реально ли произошёл kick». false при
// SUBSCRIPTION_AUTO_KICK_ENABLED=false: вызывающий код по этому флагу
// решает, делать ли запись revoke в БД (см. CheckAndSyncUser).
func (b *TelegramBot) kickUserFunc() func(int64, int64) bool {
	return b.kickFromChat
}

// --- User commands ---

// handleSubCommand checks subscription and grants access.
func (b *TelegramBot) handleSubCommand(message *tgbotapi.Message) {
	user := message.From

	// Глобально заблокированных не пускаем дальше — иначе /sub попытается
	// раздать им invite-ссылки, а через минуту их снова кикнут наши же
	// модерационные хуки. Лучше один внятный отказ.
	if active, gb, _ := b.moderationService.IsGloballyBanned(user.ID); active && gb != nil {
		text := "Доступ ограничен."
		if gb.Reason != nil && *gb.Reason != "" {
			text += "\nПричина: " + *gb.Reason
		}
		if gb.ExpiresAt != nil {
			text += "\nДо: " + gb.ExpiresAt.Format("2006-01-02 15:04")
		}
		b.sendMessage(message.Chat.ID, text)
		return
	}

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

	// Раньше /sub показывал только что granted. Если у юзера уже есть
	// access на часть чатов (например, они были добавлены в тир раньше),
	// эти чаты пропускались — казалось, что /sub «не видит» ИИ-чат или
	// ещё что-то, хотя на деле он просто не в списке «новых». Теперь
	// показываем полный актуальный список — как в /substatus.
	b.handleSubStatusCommand(message)
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

	// Собираем доступные юзеру чаты по двум источникам:
	//  1. Записи в subscription_user_chat_access — явно выданный доступ.
	//  2. Чаты, привязанные к его тиру (или ниже) — доступны даже если
	//     access-запись ещё не проставлена (например, чат добавили недавно,
	//     рассылка не прошла или юзер заблокировал бота).
	// Объединение даёт полный ответ на вопрос «что мне доступно по подписке».
	unique := map[int64]models.SubscriptionChat{}
	access, _ := b.subscriptionService.GetActiveAccess(message.From.ID)
	for _, a := range access {
		if chat, err := b.subscriptionService.GetChat(a.ChatID); err == nil {
			unique[a.ChatID] = *chat
		}
	}
	if tierID != nil {
		if tier, err := b.subscriptionService.GetTier(*tierID); err == nil {
			chats, _ := b.subscriptionService.GetChatsForTierLevel(tier.Level)
			for _, c := range chats {
				unique[c.ID] = c
			}
		}
	}

	text := fmt.Sprintf("<b>Статус подписки</b>\n\n"+
		"Тир: %s\n"+
		"Доступных чатов: %d\n", tierName, len(unique))

	if len(unique) > 0 {
		items := make([]chatListItem, 0, len(unique))
		for _, chat := range unique {
			// Одноразовая ссылка на каждый чат: хранить их смысла нет, в БД
			// ссылки не держим (см. комментарий в /substatus v1).
			link, linkErr := b.createOneTimeInviteLink(chat.ID)
			if linkErr != nil {
				log.Printf("substatus: failed to create invite link for chat %d: %v", chat.ID, linkErr)
			}
			items = append(items, chatListItem{chat: chat, link: link})
		}
		text += formatChatsGrouped(items)
		text += "\n<i>Во все сразу вступать не обязательно — Telegram после нескольких " +
			"подряд вступлений просит подождать. Заходи, куда хочется.</i>"
	}

	b.SendDirectMessage(message.Chat.ID, text)
}

// handleMyGroupsCommand shows chats available to the user's tier that they
// haven't joined yet, with fresh one-time invite links for each.
func (b *TelegramBot) handleMyGroupsCommand(message *tgbotapi.Message) {
	userID := message.From.ID
	user, err := b.subscriptionService.GetUser(userID)
	if err != nil {
		b.sendMessage(message.Chat.ID, "Вы не зарегистрированы. Используйте /sub, чтобы начать.")
		return
	}

	effectiveTierID := user.EffectiveTierID()
	if effectiveTierID == nil {
		b.sendMessage(message.Chat.ID,
			"У вас нет активной подписки. Используйте /sub, чтобы получить доступ к чатам.")
		return
	}

	tier, err := b.subscriptionService.GetTier(*effectiveTierID)
	if err != nil {
		log.Printf("mygroups: failed to get tier %d: %v", *effectiveTierID, err)
		b.sendMessage(message.Chat.ID, "Не удалось получить ваш тир. Попробуйте позже.")
		return
	}

	chats, err := b.subscriptionService.GetChatsForTierLevel(tier.Level)
	if err != nil {
		log.Printf("mygroups: failed to list chats for level %d: %v", tier.Level, err)
		b.sendMessage(message.Chat.ID, "Не удалось получить список чатов.")
		return
	}

	if len(chats) == 0 {
		b.SendDirectMessage(message.Chat.ID,
			"По вашему тиру <b>"+html.EscapeString(tier.Name)+"</b> пока нет подключённых чатов.")
		return
	}

	// Показываем все доступные чаты, а не только «куда ещё не вступил».
	// Для каждого делаем одноразовую invite-ссылку; рядом с теми, где юзер
	// уже состоит, ставим ✅ — так человек видит полный scope своей подписки
	// и может проверить, что нигде не пропустил.
	items := make([]chatListItem, 0, len(chats))
	for _, chat := range chats {
		link, linkErr := b.createOneTimeInviteLink(chat.ID)
		if linkErr != nil {
			log.Printf("mygroups: invite-link failed for chat %d: %v", chat.ID, linkErr)
		}
		// Здесь ✅-маркер декоративный (показываем «уже состоит»). При ошибке
		// API трактуем как «не member»: пользователь увидит чат без галочки —
		// безопаснее, чем показать неверное «вы уже там».
		isMember, _ := b.subscriptionService.IsMember(chat.ID, userID, b.isChatMember)
		items = append(items, chatListItem{
			chat:     chat,
			link:     link,
			isMember: isMember,
		})
	}

	text := fmt.Sprintf(
		"<b>Доступные чаты по подписке (%s):</b>\n",
		html.EscapeString(tier.Name))
	text += formatChatsGrouped(items)
	text += "\n<i>✅ — чат, в котором вы уже состоите. Во все сразу вступать " +
		"не обязательно — Telegram ограничивает подряд идущие вступления, " +
		"так что выбирай, что тебе интересно.</i>"

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

// chatListItem — строка в форматируемом списке: чат с, возможно, заранее
// сгенерированной invite-ссылкой и отметкой «юзер уже там» для /mygroups.
type chatListItem struct {
	chat     models.SubscriptionChat
	link     string
	isMember bool
}

// formatChatsGrouped группирует items по Category (NULL → «Прочее»),
// сортирует категории по MAX(priority) DESC, внутри категории — по title.
// Результат — HTML-строка с заголовками-категориями и пунктами-списком.
// Передавайте прегенерированные link-и (для юзерских команд — одноразовые);
// если link пуст, выводится просто название.
func formatChatsGrouped(items []chatListItem) string {
	const fallbackCategory = "Прочее"
	const fallbackEmoji = "💬"

	type group struct {
		items       []chatListItem
		emoji       string
		maxPriority int
	}
	groups := make(map[string]*group)
	for _, it := range items {
		cat := fallbackCategory
		emoji := fallbackEmoji
		if it.chat.Category != nil && *it.chat.Category != "" {
			cat = *it.chat.Category
		}
		if it.chat.Emoji != nil && *it.chat.Emoji != "" {
			emoji = *it.chat.Emoji
		}
		g, ok := groups[cat]
		if !ok {
			g = &group{emoji: emoji, maxPriority: it.chat.Priority}
			groups[cat] = g
		} else if g.emoji == fallbackEmoji && emoji != fallbackEmoji {
			// Если в группе появился чат с явно заданным emoji — используем его.
			g.emoji = emoji
		}
		if it.chat.Priority > g.maxPriority {
			g.maxPriority = it.chat.Priority
		}
		g.items = append(g.items, it)
	}

	order := make([]string, 0, len(groups))
	for cat := range groups {
		order = append(order, cat)
	}
	sort.Slice(order, func(i, j int) bool {
		gi, gj := groups[order[i]], groups[order[j]]
		if gi.maxPriority != gj.maxPriority {
			return gi.maxPriority > gj.maxPriority
		}
		return order[i] < order[j]
	})

	var sb strings.Builder
	for _, cat := range order {
		g := groups[cat]
		sort.Slice(g.items, func(i, j int) bool {
			return g.items[i].chat.Title < g.items[j].chat.Title
		})
		sb.WriteString(fmt.Sprintf("\n%s <b>%s</b>\n", g.emoji, html.EscapeString(cat)))
		for _, it := range g.items {
			title := html.EscapeString(it.chat.Title)
			prefix := "• "
			if it.isMember {
				prefix = "• ✅ "
			}
			if it.link != "" {
				sb.WriteString(fmt.Sprintf("%s<a href=\"%s\">%s</a>\n", prefix, it.link, title))
			} else {
				sb.WriteString(fmt.Sprintf("%s%s\n", prefix, title))
			}
		}
	}
	return sb.String()
}

// sendSubscriptionLinks sends invite links grouped by category as a single HTML message.
func (b *TelegramBot) sendSubscriptionLinks(chatID int64, result *service.SyncResult) {
	items := make([]chatListItem, 0, len(result.Granted))
	for _, g := range result.Granted {
		chat, err := b.subscriptionService.GetChat(g.ChatID)
		if err != nil {
			// Не нашли чат в БД — формируем заглушку, чтобы всё равно отдать ссылку.
			chat = &models.SubscriptionChat{ID: g.ChatID, Title: fmt.Sprintf("Chat %d", g.ChatID)}
		}
		items = append(items, chatListItem{chat: *chat, link: g.Link})
	}

	text := fmt.Sprintf("Подписка подтверждена! Доступно чатов: <b>%d</b>\n", len(result.Granted))
	text += formatChatsGrouped(items)
	text += "\n<i>Необязательно вступать во все сразу — Telegram после нескольких подряд " +
		"вступлений просит подождать. Выбирай чаты, которые тебе интересны; " +
		"остальные всегда под рукой в /mygroups.</i>"

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true
	if _, err := b.bot.Send(msg); err != nil {
		log.Printf("Failed to send subscription links to %d: %v", chatID, err)
	}
}

// --- Chat member event handlers ---

// handleChatMemberUpdated reacts to user join/leave events. Anchor-чаты
// триггерят полный sync (пересчёт тира + рассылка инвайтов в content-чаты),
// content-чаты — лёгкое обновление subscription_user_chat_access (без
// инвайтов и киков): без этого таблица access отражала бы только то, что
// бот сам выдал ссылкой, и периодик не видел бы реально сидящих в чатах.
func (b *TelegramBot) handleChatMemberUpdated(update *tgbotapi.ChatMemberUpdated) {
	chat, err := b.subscriptionService.GetChat(update.Chat.ID)
	if err != nil {
		return // Not a known subscription chat
	}

	userID := update.NewChatMember.User.ID
	oldActive := isActiveMemberStatus(update.OldChatMember.Status)
	newActive := isActiveMemberStatus(update.NewChatMember.Status)
	if oldActive == newActive {
		return
	}

	tgUser := update.NewChatMember.User
	usernamePtr := strPtr(tgUser.UserName)
	fullName := strings.TrimSpace(tgUser.FirstName + " " + tgUser.LastName)

	if chat.AnchorForTierID != nil {
		log.Printf("Anchor chat member change: chat=%d user=%d %s->%s",
			update.Chat.ID, userID, update.OldChatMember.Status, update.NewChatMember.Status)
		b.subscriptionService.InvalidateMemberCache(update.Chat.ID, userID)

		result, err := b.subscriptionService.OnboardUser(
			userID, usernamePtr, fullName,
			b.botCheckFunc(), b.createInviteLinkFunc(), b.kickUserFunc(),
		)
		if err != nil {
			return
		}
		// Раньше постили персональное "@user, добро пожаловать, нажми кнопку"
		// прямо в anchor-чат — это видели все участники, и люди жаловались на
		// спам. Теперь молчим: приветствие/инвайты уходят в ЛС через
		// notifyUserOfSyncResult (если юзер стартовал бота). Для тех, кто
		// ещё не нажимал /start, в anchor-чате есть закреплённое сообщение
		// с deep-link — его ставит pinAnchorWelcome при установке anchor.
		b.notifyUserOfSyncResult(userID, result)
		return
	}

	// Content-чат: только фиксируем факт членства. Тир пересчитает
	// ближайший PeriodicCheck (раз в N часов) — здесь дёргать его
	// преждевременно дорого (придётся снова обходить все anchor-чаты).
	//
	// isManual: добавил админ-человек, не сам юзер по invite-link и не сам бот.
	// Эта пометка защищает запись в subscription_user_chat_access от revoke
	// в periodic — чтобы юзер, добавленный админом за заслуги, не выпал на
	// следующем тикере, потому что чат не входит в его подписочный тир.
	//
	// update.From — value-тип (tgbotapi.User), при отсутствии в payload'е
	// будет zero-value с ID=0; такие случаи не считаем manual (fail-safe).
	isManual := update.From.ID != 0 &&
		update.From.ID != userID &&
		update.From.ID != b.bot.Self.ID
	log.Printf("Content chat member change: chat=%d user=%d %s->%s manual=%v",
		update.Chat.ID, userID, update.OldChatMember.Status, update.NewChatMember.Status, isManual)
	if newActive {
		if err := b.subscriptionService.SyncContentJoin(userID, update.Chat.ID, usernamePtr, fullName, isManual); err != nil {
			log.Printf("SyncContentJoin failed user=%d chat=%d: %v", userID, update.Chat.ID, err)
		}
		if isManual {
			b.notifyManualAdd(update.NewChatMember.User, &update.From, update.Chat.ID, chat.Title)
		}
	} else {
		if err := b.subscriptionService.SyncContentLeave(userID, update.Chat.ID); err != nil {
			log.Printf("SyncContentLeave failed user=%d chat=%d: %v", userID, update.Chat.ID, err)
		}
	}
}

// notifyManualAdd шлёт админу подтверждение, что юзер был добавлен в чат
// вручную и его access помечен как manual (periodic не вышибет). Без этого
// admin не имел бы immediate feedback, что auto-detect сработал — пришлось
// бы лезть в БД проверять is_manual.
func (b *TelegramBot) notifyManualAdd(addedUser *tgbotapi.User, addedBy *tgbotapi.User, chatID int64, chatTitle string) {
	if addedUser == nil {
		return
	}
	if chatTitle == "" {
		chatTitle = fmt.Sprintf("chat %d", chatID)
	}
	by := "?"
	if addedBy != nil {
		by = targetDisplay(addedBy)
	}
	text := fmt.Sprintf(
		"🛡 %s (<code>%d</code>) добавлен вручную в <b>%s</b>.\nДобавил: %s\nДоступ помечен как manual — periodic не тронет.",
		targetDisplay(addedUser), addedUser.ID, html.EscapeString(chatTitle), by,
	)
	b.SendDirectMessage(subscriptionAdminID(), text)
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

		// Сразу начинаем трекать активность в этом чате, чтобы новый чат
		// автоматически появился на дашборде «Активность чатов» без ручных миграций.
		if err := b.chatActivityService.AddTrackedChat(chat.ID, title, chat.Type); err != nil {
			log.Printf("Failed to add tracked chat %d: %v", chat.ID, err)
		}

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
		b.SendDirectMessage(subscriptionAdminID(), fmt.Sprintf(
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

		// Снимаем чат с отслеживания активности. История сообщений в chat_messages
		// остаётся — дашборд покажет только активные источники.
		if err := b.chatActivityService.RemoveTrackedChat(chat.ID); err != nil {
			log.Printf("Failed to deactivate tracked chat %d: %v", chat.ID, err)
		}
	}
}

func isActiveMemberStatus(status string) bool {
	return status == "creator" || status == "administrator" || status == "member" || status == "restricted"
}

// --- Periodic subscription checker ---

func (b *TelegramBot) startSubscriptionChecker() {
	interval := time.Duration(config.CFG.SubscriptionCheckIntervalHours) * time.Hour
	tierTicker := time.NewTicker(interval)
	defer tierTicker.Stop()

	// Sweep реального членства — отдельный 24-часовой тикер: проход
	// долгий (≈10 мин на 250 юзеров × 41 чат), а тиры мы перепроверяем
	// каждые 4 часа. На отдельном тикере sweep не блокирует tier-check.
	sweepTicker := time.NewTicker(24 * time.Hour)
	defer sweepTicker.Stop()

	for {
		select {
		case <-tierTicker.C:
			results := b.subscriptionService.PeriodicCheck(
				b.botCheckFunc(),
				b.createInviteLinkFunc(),
				b.kickUserFunc(),
				50*time.Millisecond,
			)
			b.sendPeriodicCheckReport(subscriptionAdminID(), "Периодическая проверка подписок", results)
		case <-sweepTicker.C:
			stats, err := b.subscriptionService.SweepRealMembership(b.uncachedCheckFunc(), 50*time.Millisecond)
			if err != nil {
				log.Printf("Daily membership sweep failed: %v", err)
			} else {
				log.Printf("Daily membership sweep done: scanned=%d created=%d checks=%d failed=%d granted=%d revoked=%d",
					stats.UsersScanned, stats.UsersCreated, stats.ChecksPerformed,
					stats.ChecksFailed, stats.AccessGranted, stats.AccessRevoked)
			}
		}
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
		// Сразу ставим единое закреплённое приветствие с deep-link в бота —
		// иначе новый участник якорного чата не знает, куда идти (писать
		// боту первым нельзя — Telegram запрещает; нужен его click).
		// Ошибка пина не критична: чат уже записан как anchor, админ
		// может вызвать /subpin руками.
		if err := b.pinAnchorWelcome(chatID); err != nil {
			log.Printf("subaddchat: failed to pin anchor welcome for chat %d: %v", chatID, err)
			b.SendDirectMessage(message.Chat.ID, fmt.Sprintf(
				"Чат <code>%d</code> установлен как <b>anchor</b> для тира %s, но закрепить welcome не удалось: %v.\n"+
					"Вызовите /subpin %d вручную.", chatID, tier.Name, err, chatID))
			return
		}
		b.SendDirectMessage(message.Chat.ID, fmt.Sprintf(
			"Чат <code>%d</code> установлен как <b>anchor</b> для тира %s, welcome-сообщение закреплено.",
			chatID, tier.Name))
		return
	}

	b.subscriptionService.AddChatToTier(chatID, tier.ID)
	b.SendDirectMessage(message.Chat.ID, fmt.Sprintf(
		"Чат <code>%d</code> добавлен как <b>content</b> для тира %s.", chatID, tier.Name))

	// Рассылаем уведомления всем пользователям, у которых эффективный тир
	// достаточного уровня и ещё нет доступа в этот чат. Делаем асинхронно,
	// чтобы команда /subaddchat не блокировалась на N Telegram-запросах.
	go b.notifyNewChatAccess(chatID, title, tier.Level, message.Chat.ID)
}

// notifyNewChatAccess выдаёт доступ в chatID всем пользователям с нужным
// уровнем тира (у кого его ещё нет) и рассылает им invite-ссылку в ЛС.
// В конце отправляет супер-админу сводку по рассылке.
//
// Создаём ОДНУ invite-ссылку с member_limit = len(users), а не по штуке
// на каждого — Telegram лимитирует createChatInviteLink ~20/мин на чат,
// и раньше рассылка на 50+ юзеров массово падала с «Too Many Requests».
func (b *TelegramBot) notifyNewChatAccess(chatID int64, chatTitle string, tierLevel int, adminChatID int64) {
	users, err := b.subscriptionService.GetEligibleUsersWithoutAccessForChat(chatID, tierLevel)
	if err != nil {
		log.Printf("notifyNewChatAccess: failed to fetch eligible users for chat %d: %v", chatID, err)
		b.SendDirectMessage(adminChatID, fmt.Sprintf("Не удалось собрать список получателей для чата <code>%d</code>: %v", chatID, err))
		return
	}
	if len(users) == 0 {
		b.SendDirectMessage(adminChatID, fmt.Sprintf("Рассылка не нужна: в чате <code>%d</code> уже все с подходящим тиром.", chatID))
		return
	}

	link, err := b.createInviteLinkWithLimit(chatID, len(users))
	if err != nil {
		log.Printf("notifyNewChatAccess: failed to create shared invite link for chat %d: %v", chatID, err)
		b.SendDirectMessage(adminChatID, fmt.Sprintf(
			"Не удалось создать invite-ссылку для чата <code>%d</code>: %v", chatID, err))
		return
	}

	titleEscaped := html.EscapeString(chatTitle)
	text := fmt.Sprintf(
		"🆕 Вам открыт новый чат по вашей подписке:\n\n<b>%s</b>\n\n<a href=\"%s\">Перейти в чат</a>",
		titleEscaped, link)
	delivered, skipped, failed := 0, 0, 0

	for _, user := range users {
		msg := tgbotapi.NewMessage(user.ID, text)
		msg.ParseMode = "HTML"
		msg.DisableWebPagePreview = true
		if _, err := b.bot.Send(msg); err != nil {
			// Forbidden = пользователь не нажимал /start боту. Не считаем это
			// ошибкой, просто skip — ссылка shared, новая «не протухает».
			if strings.Contains(err.Error(), "Forbidden") {
				skipped++
			} else {
				log.Printf("notifyNewChatAccess: DM failed to user %d: %v", user.ID, err)
				failed++
			}
			continue
		}
		if err := b.subscriptionService.GrantAccess(user.ID, chatID, false); err != nil {
			log.Printf("notifyNewChatAccess: grant failed for user %d chat %d: %v", user.ID, chatID, err)
			failed++
			continue
		}
		delivered++
		// Небольшая пауза между отправками — Telegram limits ~30 msg/sec
		// в разные чаты. 50ms с запасом.
		time.Sleep(50 * time.Millisecond)
	}

	b.SendDirectMessage(adminChatID, fmt.Sprintf(
		"Рассылка по чату <code>%d</code>: доставлено %d, пропущено (бот заблокирован) %d, ошибок %d.",
		chatID, delivered, skipped, failed))
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
	if err := b.pinAnchorWelcome(chatID); err != nil {
		log.Printf("subsetanchor: failed to pin anchor welcome for chat %d: %v", chatID, err)
		b.SendDirectMessage(message.Chat.ID, fmt.Sprintf(
			"Чат <code>%d</code> теперь anchor для тира %s, но закрепить welcome не удалось: %v.\n"+
				"Вызовите /subpin %d вручную.", chatID, tier.Name, err, chatID))
		return
	}
	b.SendDirectMessage(message.Chat.ID, fmt.Sprintf(
		"Чат <code>%d</code> теперь anchor для тира %s, welcome-сообщение закреплено.",
		chatID, tier.Name))
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
		results := b.subscriptionService.PeriodicCheck(
			b.botCheckFunc(),
			b.createInviteLinkFunc(),
			b.kickUserFunc(),
			50*time.Millisecond,
		)
		b.sendPeriodicCheckReport(message.Chat.ID, "Полная проверка подписок", results)
	}()
}

// sendPeriodicCheckReport форматирует и шлёт админу сводку по результату
// PeriodicCheck. PeriodicCheck — массовый автоматический проход; юзерам
// нотификации не уходят (это спам), вместо этого админ видит, что бот
// сделал и руками решает, что с этим делать. Сообщения дробятся по 3500
// символов — Telegram режет на 4096.
//
// Skipped — юзеры, которых периодик пропустил из-за ошибок Telegram API
// (rate-limit, таймаут, network). Раньше это терялось в логах и сводка
// «Изменений нет» выглядела как success, хотя сервис мог быть деградирован.
func (b *TelegramBot) sendPeriodicCheckReport(adminID int64, title string, report service.PeriodicReport) {
	skipNote := ""
	if len(report.Skipped) > 0 {
		// Показываем до 10 ID, чтобы не раздувать сообщение; полный список
		// в логах пода.
		shown := report.Skipped
		extra := ""
		if len(shown) > 10 {
			shown = shown[:10]
			extra = fmt.Sprintf(" (показаны первые 10 из %d)", len(report.Skipped))
		}
		ids := make([]string, len(shown))
		for i, id := range shown {
			ids[i] = fmt.Sprintf("<code>%d</code>", id)
		}
		skipNote = fmt.Sprintf("\n\n⚠️ Пропущено юзеров из-за ошибок Telegram: %d%s\n%s",
			len(report.Skipped), extra, strings.Join(ids, ", "))
	}

	if len(report.Changed) == 0 {
		summary := fmt.Sprintf("<b>%s</b>\n\nИзменений нет.", html.EscapeString(title))
		if report.UsersTotal > 0 {
			summary += fmt.Sprintf(" Просмотрено юзеров: %d.", report.UsersTotal)
		}
		b.SendDirectMessage(adminID, summary+skipNote)
		return
	}

	chats, _ := b.subscriptionService.GetAllChats()
	chatTitle := make(map[int64]string, len(chats))
	for _, c := range chats {
		chatTitle[c.ID] = c.Title
	}

	var totalGrant, totalRevoke int
	for _, r := range report.Changed {
		totalGrant += len(r.Granted)
		totalRevoke += len(r.Revoked)
	}

	header := fmt.Sprintf(
		"<b>%s</b>\n\n"+
			"Юзеров с изменениями: %d (всего просмотрено: %d)\n"+
			"Всего grant: %d\n"+
			"Всего revoke: %d%s\n\n",
		html.EscapeString(title), len(report.Changed), report.UsersTotal,
		totalGrant, totalRevoke, skipNote,
	)

	lines := make([]string, 0, len(report.Changed))
	for _, r := range report.Changed {
		uname := ""
		if u, err := b.subscriptionService.GetUser(r.UserID); err == nil && u.Username != nil && *u.Username != "" {
			uname = "@" + *u.Username
		}
		line := fmt.Sprintf("<code>%d</code> %s — grant=%d revoke=%d",
			r.UserID, html.EscapeString(uname), len(r.Granted), len(r.Revoked))
		if len(r.Revoked) > 0 {
			names := make([]string, 0, len(r.Revoked))
			for _, cid := range r.Revoked {
				t := chatTitle[cid]
				if t == "" {
					t = fmt.Sprintf("%d", cid)
				}
				names = append(names, html.EscapeString(t))
			}
			line += "\n  ↳ revoke: " + strings.Join(names, ", ")
		}
		if len(r.Granted) > 0 {
			names := make([]string, 0, len(r.Granted))
			for _, g := range r.Granted {
				t := chatTitle[g.ChatID]
				if t == "" {
					t = fmt.Sprintf("%d", g.ChatID)
				}
				names = append(names, html.EscapeString(t))
			}
			line += "\n  ↳ grant: " + strings.Join(names, ", ")
		}
		lines = append(lines, line)
	}

	body := strings.Join(lines, "\n\n")
	const maxLen = 3500
	if len(header)+len(body) <= maxLen {
		b.SendDirectMessage(adminID, header+body)
		return
	}
	b.SendDirectMessage(adminID, header+"(детали ниже частями)")
	var buf strings.Builder
	for _, l := range lines {
		if buf.Len()+len(l)+2 > maxLen {
			b.SendDirectMessage(adminID, buf.String())
			buf.Reset()
		}
		if buf.Len() > 0 {
			buf.WriteString("\n\n")
		}
		buf.WriteString(l)
	}
	if buf.Len() > 0 {
		b.SendDirectMessage(adminID, buf.String())
	}
}

// handleSubMemberSweepCommand — backfill реального членства. Идёт по
// subscription_users ∪ chat_messages.distinct(user_id) × всем
// subscription_chats и приводит subscription_user_chat_access в
// соответствие реальности по getChatMember. Долгая команда (~10 мин на
// 250 юзеров × 41 чат с rateDelay=50ms), запускается в фоне.
func (b *TelegramBot) handleSubMemberSweepCommand(message *tgbotapi.Message) {
	if !b.isSubscriptionAdmin(message.From.ID) {
		return
	}

	b.sendMessage(message.Chat.ID, "Запуск sweep реального членства... это займёт несколько минут.")

	go func() {
		stats, err := b.subscriptionService.SweepRealMembership(b.uncachedCheckFunc(), 50*time.Millisecond)
		if err != nil {
			b.SendDirectMessage(message.Chat.ID, fmt.Sprintf("Sweep упал: %v", err))
			return
		}
		failedNote := ""
		if stats.ChecksFailed > 0 {
			failedNote = fmt.Sprintf("\n⚠️ Упало getChatMember: %d (см. логи — обычно rate-limit или бот без прав)",
				stats.ChecksFailed)
		}
		b.SendDirectMessage(message.Chat.ID, fmt.Sprintf(
			"<b>Sweep завершён</b>\n\n"+
				"Просканировано юзеров: %d\n"+
				"Заведено новых в subscription_users: %d\n"+
				"Запросов getChatMember: %d\n"+
				"Выдано access (новых): %d\n"+
				"Снято access (revoke): %d%s",
			stats.UsersScanned, stats.UsersCreated, stats.ChecksPerformed,
			stats.AccessGranted, stats.AccessRevoked, failedNote,
		))
	}()
}

// handleSubKickDryCommand — прогон PeriodicCheck в режиме «без действий».
// Считает, кого бот удалил бы из чатов прямо сейчас (по текущему состоянию
// БД и реальному членству в anchor-чатах). Ничего в БД и в Telegram не
// меняется. Используется как пред-чек перед включением SUBSCRIPTION_AUTO_KICK_ENABLED.
func (b *TelegramBot) handleSubKickDryCommand(message *tgbotapi.Message) {
	if !b.isSubscriptionAdmin(message.From.ID) {
		return
	}

	b.sendMessage(message.Chat.ID, "Запуск dry-run проверки... ничего реально не кикается.")

	go func() {
		results, err := b.subscriptionService.DryRunPeriodicCheck(b.uncachedCheckFunc(), 50*time.Millisecond)
		if err != nil {
			b.SendDirectMessage(message.Chat.ID, fmt.Sprintf("Dry-run упал: %v", err))
			return
		}

		// Сворачиваем chat_id → title заранее одним запросом, чтобы в
		// отчёте показывать человекочитаемые названия.
		chats, _ := b.subscriptionService.GetAllChats()
		chatTitle := make(map[int64]string, len(chats))
		for _, c := range chats {
			chatTitle[c.ID] = c.Title
		}

		var totalRevoke, usersWithRevoke, totalGrant int
		var lines []string
		for _, r := range results {
			if len(r.WouldRevoke) == 0 && len(r.WouldGrant) == 0 {
				continue
			}
			totalRevoke += len(r.WouldRevoke)
			totalGrant += len(r.WouldGrant)
			if len(r.WouldRevoke) > 0 {
				usersWithRevoke++
			}

			uname := ""
			if r.Username != nil && *r.Username != "" {
				uname = "@" + *r.Username
			}
			line := fmt.Sprintf("<code>%d</code> %s — revoke=%d grant=%d",
				r.UserID, html.EscapeString(uname), len(r.WouldRevoke), len(r.WouldGrant))
			if len(r.WouldRevoke) > 0 {
				names := make([]string, 0, len(r.WouldRevoke))
				for _, cid := range r.WouldRevoke {
					t := chatTitle[cid]
					if t == "" {
						t = fmt.Sprintf("%d", cid)
					}
					names = append(names, html.EscapeString(t))
				}
				line += "\n  ↳ из: " + strings.Join(names, ", ")
			}
			lines = append(lines, line)
		}

		header := fmt.Sprintf("<b>Dry-run проверки подписок</b>\n\n"+
			"Просканировано юзеров: %d\n"+
			"Юзеров под кик: %d (всего access-revoke действий: %d)\n"+
			"Access-grant действий: %d\n\n",
			len(results), usersWithRevoke, totalRevoke, totalGrant)

		// Telegram режет сообщения на 4096 символов — режем сами по строкам.
		body := strings.Join(lines, "\n\n")
		const maxLen = 3500
		if len(body) <= maxLen {
			b.SendDirectMessage(message.Chat.ID, header+body)
			return
		}
		b.SendDirectMessage(message.Chat.ID, header+"(полный отчёт ниже частями)")
		var buf strings.Builder
		for _, l := range lines {
			if buf.Len()+len(l)+2 > maxLen {
				b.SendDirectMessage(message.Chat.ID, buf.String())
				buf.Reset()
			}
			if buf.Len() > 0 {
				buf.WriteString("\n\n")
			}
			buf.WriteString(l)
		}
		if buf.Len() > 0 {
			b.SendDirectMessage(message.Chat.ID, buf.String())
		}
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

// pinAnchorWelcome публикует и закрепляет единое welcome-сообщение с
// deep-link в бот — чтобы новый участник якорного чата мог одним нажатием
// стартнуть /start и получить инвайты. Используется и бот-командой /subpin,
// и автоматически при установке anchor через /subaddchat … anchor.
func (b *TelegramBot) pinAnchorWelcome(chatID int64) error {
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
		return fmt.Errorf("send: %w", err)
	}
	pinCfg := tgbotapi.PinChatMessageConfig{
		ChatID:              chatID,
		MessageID:           sent.MessageID,
		DisableNotification: false,
	}
	if _, err := b.bot.Request(pinCfg); err != nil {
		return fmt.Errorf("pin: %w", err)
	}
	return nil
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

	if err := b.pinAnchorWelcome(chatID); err != nil {
		b.sendMessage(message.Chat.ID, fmt.Sprintf("Не удалось закрепить: %v", err))
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
