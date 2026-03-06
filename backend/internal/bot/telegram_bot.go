package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"ithozyeva/config"
	"ithozyeva/database"
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	globalBot *TelegramBot
	botMutex  sync.RWMutex
)

// getEventLocation возвращает *time.Location для таймзоны события.
// Поддерживает формат "UTC", "UTC+3", "UTC-5" и т.д.
func getEventLocation(timezone string) *time.Location {
	if timezone == "" || timezone == "UTC" {
		return time.UTC
	}

	// Парсим "UTC+3" или "UTC-5"
	if strings.HasPrefix(timezone, "UTC") {
		offsetStr := timezone[3:] // "+3" или "-5"
		if offsetStr == "" {
			return time.UTC
		}
		hours, err := strconv.Atoi(offsetStr)
		if err != nil {
			log.Printf("Warning: failed to parse timezone %q, falling back to UTC", timezone)
			return time.UTC
		}
		return time.FixedZone(timezone, hours*3600)
	}

	// Пробуем как IANA таймзону (на будущее)
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Printf("Warning: failed to load timezone %q: %v, falling back to UTC", timezone, err)
		return time.UTC
	}
	return loc
}

// formatEventDateStr форматирует дату события с учётом его таймзоны
func formatEventDateStr(eventDate time.Time, timezone string) string {
	loc := getEventLocation(timezone)
	dateInTz := eventDate.In(loc)
	return dateInTz.Format("02.01.2006 в 15:04")
}

// formatTimezoneLabel возвращает человекочитаемую метку таймзоны
func formatTimezoneLabel(timezone string) string {
	if timezone == "" || timezone == "UTC" {
		return "UTC"
	}
	return timezone
}

// GetGlobalBot возвращает глобальный экземпляр бота
func GetGlobalBot() *TelegramBot {
	botMutex.RLock()
	defer botMutex.RUnlock()
	return globalBot
}

// SetGlobalBot устанавливает глобальный экземпляр бота
func SetGlobalBot(b *TelegramBot) {
	botMutex.Lock()
	defer botMutex.Unlock()
	globalBot = b

	// Регистрируем callback для отправки Telegram DM из service-слоя (избегаем circular import)
	service.SendTelegramDMFunc = func(chatID int64, text string) {
		b.SendDirectMessage(chatID, text)
	}
}

type TelegramBot struct {
	bot                         *tgbotapi.BotAPI
	tg_service                  *service.TelegramService
	member                      *service.MemberService
	eventAlertSubscription      *service.EventAlertSubscriptionService
	eventService                *service.EventsService
	chatActivityService         *service.ChatActivityService
	notificationSettingsService *service.NotificationSettingsService
	chatHighlightService        *service.ChatHighlightService
}

func NewTelegramBot() (*TelegramBot, error) {

	botToken := config.CFG.TelegramToken
	if botToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, fmt.Errorf("error creating bot: %v", err)
	}

	tg_service, err := service.NewTelegramService()
	if err != nil {
		return nil, fmt.Errorf("error creating telegram service: %v", err)
	}

	member_service := service.NewMemberService()
	eventAlertSubscriptionService := service.NewEventAlertSubscriptionService()
	eventService := service.NewEventsService()

	chatActivityService := service.NewChatActivityService()
	notificationSettingsService := service.NewNotificationSettingsService()
	chatHighlightService := service.NewChatHighlightService()

	return &TelegramBot{
		bot:                         bot,
		tg_service:                  tg_service,
		member:                      member_service,
		eventAlertSubscription:      eventAlertSubscriptionService,
		eventService:                eventService,
		chatActivityService:         chatActivityService,
		notificationSettingsService: notificationSettingsService,
		chatHighlightService:        chatHighlightService,
	}, nil
}

func (b *TelegramBot) Start() {
	// Register bot commands menu
	b.registerCommands()

	// Start birthday checker
	go b.startBirthdayChecker()

	// Start event alerts scheduler
	go b.startEventAlertsScheduler()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		// Обработка callback кнопок
		if update.CallbackQuery != nil {
			b.handleCallbackQuery(update.CallbackQuery)
			continue
		}

		if update.Message == nil {
			continue
		}

		// Трекинг активности чатов — для каждого сообщения (асинхронно, чтобы не блокировать обработку)
		go b.chatActivityService.TrackMessage(update.Message)

		// Обработка новых участников чата
		if update.Message.NewChatMembers != nil {
			for _, newMember := range update.Message.NewChatMembers {
				b.handleNewChatMember(update.Message.Chat.ID, &newMember)
			}
			continue
		}

		// Команда /chatid — отправляет ID чата владельцу и удаляет сообщение
		if update.Message.IsCommand() && update.Message.Command() == "chatid" {
			b.handleChatIDCommand(update.Message)
			continue
		}

		// Команда /highlight — сохранение сообщения как хайлайт
		if update.Message.IsCommand() && update.Message.Command() == "highlight" {
			b.handleHighlightCommand(update.Message)
			continue
		}

		// Команда /whois — информация об участнике в групповых чатах
		if update.Message.IsCommand() && update.Message.Command() == "whois" {
			b.handleWhoisCommand(update.Message)
			continue
		}

		// Бот отвечает только в личных сообщениях
		if update.Message.Chat.Type != "private" {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				b.handleStartCommand(update.Message)
			case "mypoints":
				b.handleMyPointsCommand(update.Message)
			case "events":
				b.handleEventsCommand(update.Message)
			case "help":
				b.handleHelpCommand(update.Message)
			}
		}
	}
}

func (b *TelegramBot) registerCommands() {
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Авторизация на платформе"},
		{Command: "mypoints", Description: "Мои баллы"},
		{Command: "events", Description: "Ближайшие события"},
		{Command: "whois", Description: "Кто этот участник"},
		{Command: "help", Description: "Помощь"},
	}
	cfg := tgbotapi.NewSetMyCommands(commands...)
	if _, err := b.bot.Request(cfg); err != nil {
		log.Printf("Error registering bot commands: %v", err)
	}
}

func (b *TelegramBot) handleMyPointsCommand(message *tgbotapi.Message) {
	member, err := b.member.GetByTelegramID(message.From.ID)
	if err != nil {
		b.sendMessage(message.Chat.ID, "Вы не зарегистрированы на платформе. Используйте /start для авторизации.")
		return
	}

	pointsSvc := service.NewPointsService()
	balance, err := pointsSvc.GetBalance(member.Id)
	if err != nil {
		b.sendMessage(message.Chat.ID, "Ошибка при получении баллов.")
		return
	}

	text := fmt.Sprintf("Ваш баланс: %d баллов", balance)
	b.sendMessage(message.Chat.ID, text)
}

func (b *TelegramBot) handleEventsCommand(message *tgbotapi.Message) {
	events, err := b.eventService.GetUpcomingEvents(3)
	if err != nil || len(events) == 0 {
		b.sendMessage(message.Chat.ID, "Ближайших событий не найдено.")
		return
	}

	var builder strings.Builder
	builder.WriteString("<b>Ближайшие события:</b>\n\n")
	for _, event := range events {
		dateStr := formatEventDateStr(event.Date, event.Timezone)
		tzLabel := formatTimezoneLabel(event.Timezone)
		builder.WriteString(fmt.Sprintf("📆 <b>%s</b>\n%s (%s)\n\n", event.Title, dateStr, tzLabel))
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, builder.String())
	msg.ParseMode = "HTML"
	b.bot.Send(msg)
}

func (b *TelegramBot) handleHelpCommand(message *tgbotapi.Message) {
	text := "Доступные команды:\n" +
		"/start - Авторизация на платформе\n" +
		"/mypoints - Посмотреть баланс баллов\n" +
		"/events - Ближайшие события\n" +
		"/whois - Кто этот участник (ответьте на сообщение или /whois @username)\n" +
		"/help - Помощь"
	b.sendMessage(message.Chat.ID, text)
}

// handleWhoisCommand показывает информацию об участнике сообщества
func (b *TelegramBot) handleWhoisCommand(message *tgbotapi.Message) {
	var member *models.Member
	var err error

	// Способ 1: ответ на сообщение
	if message.ReplyToMessage != nil {
		member, err = b.member.GetByTelegramID(message.ReplyToMessage.From.ID)
	} else {
		// Способ 2: /whois @username
		args := strings.TrimSpace(message.CommandArguments())
		if args == "" {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Ответьте на сообщение командой /whois или укажите username: /whois @username")
			msg.ReplyToMessageID = message.MessageID
			b.bot.Send(msg)
			return
		}
		username := strings.TrimPrefix(args, "@")
		member, err = b.member.GetByUsername(username)
	}

	if err != nil || member == nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Участник не найден на платформе.")
		msg.ReplyToMessageID = message.MessageID
		b.bot.Send(msg)
		return
	}

	// Формируем карточку участника
	var builder strings.Builder

	name := strings.TrimSpace(fmt.Sprintf("%s %s", member.FirstName, member.LastName))
	if name == "" {
		name = member.Username
	}
	builder.WriteString(fmt.Sprintf("👤 <b>%s</b>", name))
	if member.Username != "" {
		builder.WriteString(fmt.Sprintf(" (@%s)", member.Username))
	}
	builder.WriteString("\n")

	// Роли
	if len(member.Roles) > 0 {
		roleLabels := map[models.Role]string{
			models.MemberRoleSubscriber:   "Участник",
			models.MemberRoleMentor:       "Ментор",
			models.MemberRoleAdmin:        "Админ",
			models.MemberRoleEventMaker:   "Организатор событий",
			models.MemberRoleUnsubscriber: "Неактивный",
		}
		var roles []string
		for _, role := range member.Roles {
			if label, ok := roleLabels[role]; ok {
				roles = append(roles, label)
			}
		}
		if len(roles) > 0 {
			builder.WriteString(fmt.Sprintf("🏷 %s\n", strings.Join(roles, ", ")))
		}
	}

	// Био
	if member.Bio != "" {
		builder.WriteString(fmt.Sprintf("\n📝 %s\n", member.Bio))
	}

	// Баллы
	pointsSvc := service.NewPointsService()
	balance, pointsErr := pointsSvc.GetBalance(member.Id)
	if pointsErr == nil {
		builder.WriteString(fmt.Sprintf("\n⭐ Баллы: %d", balance))
	}

	// Менторская информация
	mentor, mentorErr := b.member.GetMentor(member.Id)
	if mentorErr == nil && mentor != nil {
		if mentor.Occupation != "" {
			builder.WriteString(fmt.Sprintf("\n💼 %s", mentor.Occupation))
		}
		if mentor.Experience != "" {
			builder.WriteString(fmt.Sprintf("\n📊 Опыт: %s", mentor.Experience))
		}
		if len(mentor.ProfTags) > 0 {
			var tags []string
			for _, tag := range mentor.ProfTags {
				tags = append(tags, tag.Title)
			}
			builder.WriteString(fmt.Sprintf("\n🔧 %s", strings.Join(tags, ", ")))
		}
	}

	// Ссылка на профиль
	platformURL := config.CFG.PublicDomain
	if platformURL != "" {
		if !strings.HasPrefix(platformURL, "http://") && !strings.HasPrefix(platformURL, "https://") {
			platformURL = "https://" + platformURL
		}
		builder.WriteString(fmt.Sprintf("\n\n🔗 <a href=\"%s/members/%d\">Профиль на платформе</a>", platformURL, member.Id))
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, builder.String())
	msg.ParseMode = "HTML"
	msg.ReplyToMessageID = message.MessageID
	msg.DisableWebPagePreview = true
	b.bot.Send(msg)
}

// handleChatIDCommand отправляет ID чата владельцу в ЛС и удаляет команду
func (b *TelegramBot) handleChatIDCommand(message *tgbotapi.Message) {
	// Удаляем сообщение с командой
	deleteMsg := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
	b.bot.Request(deleteMsg)

	// Отправляем ID чата владельцу в ЛС
	const ownerChatID int64 = 931916742
	text := fmt.Sprintf("Chat ID: <code>%d</code>\nTitle: %s\nType: %s", message.Chat.ID, message.Chat.Title, message.Chat.Type)
	msg := tgbotapi.NewMessage(ownerChatID, text)
	msg.ParseMode = "HTML"
	b.bot.Send(msg)
}

// SendDirectMessage отправляет личное сообщение пользователю по chatID
func (b *TelegramBot) SendDirectMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	if _, err := b.bot.Send(msg); err != nil {
		log.Printf("Error sending direct message to %d: %v", chatID, err)
	}
}

func (b *TelegramBot) startBirthdayChecker() {
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
		if now.After(next) {
			next = next.Add(24 * time.Hour)
		}
		time.Sleep(time.Until(next))

		b.checkBirthdays()
	}
}

func (b *TelegramBot) checkBirthdays() {
	birthdays, err := b.member.GetTodayBirthdays()
	if err != nil {
		log.Printf("Error checking birthdays: %v", err)
		return
	}

	if len(birthdays) == 0 {
		return
	}

	// Get random congratulations
	congrats := []string{
		"🎉 С днем рождения! Желаю счастья, здоровья и успехов!",
		"🎂 Поздравляю с днем рождения! Пусть каждый день будет наполнен радостью!",
		"🎊 С днем рождения! Пусть все мечты становятся реальностью!",
		"🎈 С днем рождения! Желаю удачи во всех начинаниях!",
		"🎁 Поздравляю с днем рождения! Пусть жизнь будет полна приятных сюрпризов!",
	}
	randomCongrats := congrats[rand.Intn(len(congrats))]

	// Mention all users with birthdays
	mentions := make([]string, len(birthdays))
	for i, username := range birthdays {
		mentions[i] = fmt.Sprintf("@%s", username)
	}
	mentionText := strings.Join(mentions, " ")

	// Send birthday message
	message := fmt.Sprintf("%s\n%s", mentionText, randomCongrats)
	b.sendMessage(config.CFG.TelegramMainChatID, message)
}

func (b *TelegramBot) handleStartCommand(message *tgbotapi.Message) {
	log.Printf("Received /start command from user %d with args: %s", message.From.ID, message.CommandArguments())

	// Получаем аргументы команды
	args := strings.Split(message.CommandArguments(), " ")
	if len(args) == 0 || args[0] == "" {
		log.Printf("No arguments provided for /start command")
		b.sendMessage(message.Chat.ID, "Пожалуйста, используйте кнопку авторизации на сайте")
		return
	}

	// Первый аргумент - URL для перенаправления
	redirectUrl := config.CFG.PublicDomain
	log.Printf("Redirect URL before processing: %s", redirectUrl)

	if !strings.HasPrefix(redirectUrl, "http://") && !strings.HasPrefix(redirectUrl, "https://") {
		redirectUrl = "http://" + redirectUrl
	}
	log.Printf("Final redirect URL: %s", redirectUrl)

	// Генерируем токен для пользователя
	token := b.tg_service.GenerateAuthToken(message.From.ID)

	log.Printf("Generated token for user %d: %s", message.From.ID, token)

	// Формируем URL для перенаправления с токеном
	authUrl := fmt.Sprintf("%s?token=%s", redirectUrl, token)
	log.Printf("Final auth URL: %s", authUrl)

	sendAuthToBackend(token, message.From)

	// Отправляем сообщение с кнопкой для авторизации
	msg := tgbotapi.NewMessage(message.Chat.ID, "Нажмите кнопку ниже для авторизации")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Авторизоваться", authUrl),
		),
	)

	if _, err := b.bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	} else {
		log.Printf("Successfully sent auth button to user %d", message.From.ID)
	}
}

// handleHighlightCommand сохраняет сообщение как хайлайт
func (b *TelegramBot) handleHighlightCommand(message *tgbotapi.Message) {
	if message.ReplyToMessage == nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Ответьте на сообщение командой /highlight, чтобы сохранить его как хайлайт.")
		msg.ReplyToMessageID = message.MessageID
		b.bot.Send(msg)
		return
	}

	reply := message.ReplyToMessage
	if reply.Text == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Можно сохранить только текстовые сообщения.")
		msg.ReplyToMessageID = message.MessageID
		b.bot.Send(msg)
		return
	}

	// Ищем member по telegram ID автора сообщения
	var memberID *int64
	member, err := b.member.GetByTelegramID(reply.From.ID)
	if err == nil && member != nil {
		memberID = &member.Id
	}

	highlight := &models.ChatHighlight{
		ChatID:           message.Chat.ID,
		MessageID:        reply.MessageID,
		AuthorTelegramID: reply.From.ID,
		AuthorUsername:    reply.From.UserName,
		AuthorFirstName:  reply.From.FirstName,
		MessageText:      reply.Text,
		HighlightedBy:    message.From.ID,
		MemberID:         memberID,
	}

	_, err = b.chatHighlightService.Create(highlight)
	if err != nil {
		log.Printf("Error saving highlight: %v", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Ошибка при сохранении хайлайта.")
		msg.ReplyToMessageID = message.MessageID
		b.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "⭐ Сообщение сохранено как хайлайт!")
	msg.ReplyToMessageID = message.MessageID
	b.bot.Send(msg)
}

// handleNewChatMember приветствует новых участников в основном чате
func (b *TelegramBot) handleNewChatMember(chatID int64, user *tgbotapi.User) {
	if chatID != config.CFG.TelegramMainChatID {
		return
	}
	if user.IsBot {
		return
	}

	name := user.FirstName
	if user.UserName != "" {
		name = "@" + user.UserName
	}

	text := fmt.Sprintf("👋 Приветствуем <b>%s</b> в IT-Хозяевах!\n\n"+
		"🌐 <a href=\"https://ithozyaeva.ru/platform\">Платформа</a> — здесь всё самое интересное!",
		name)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true
	if _, err := b.bot.Send(msg); err != nil {
		log.Printf("Error sending welcome message: %v", err)
	}
}

func (b *TelegramBot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := b.bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func (b *TelegramBot) SendEventAlert(telegramID int64, event *models.Event, isInitial bool) error {
	now := time.Now()
	timeUntilEvent := event.Date.Sub(now)
	messageText := b.formatEventAlert(event, isInitial, timeUntilEvent)

	msg := tgbotapi.NewMessage(telegramID, messageText)
	msg.ParseMode = "HTML"

	if isInitial {
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("✅ Приду", fmt.Sprintf("event_attend:%d", event.Id)),
				tgbotapi.NewInlineKeyboardButtonData("❌ Не приду", fmt.Sprintf("event_decline:%d", event.Id)),
			),
		)
		msg.ReplyMarkup = keyboard
	}

	_, err := b.bot.Send(msg)
	return err
}

func (b *TelegramBot) formatEventAlert(event *models.Event, isInitial bool, timeUntilEvent time.Duration) string {
	var builder strings.Builder

	if event.ExclusiveChatID != nil && *event.ExclusiveChatID != 0 {
		label := event.ExclusiveChatTitle
		if label == "" {
			label = "Эксклюзив"
		}
		builder.WriteString(fmt.Sprintf("👑 <b>%s</b>\n", label))
	}

	if isInitial {
		builder.WriteString("⭐ <b>Новое событие!</b>\n\n")
	} else if timeUntilEvent <= 1*time.Minute && timeUntilEvent > -2*time.Minute {
		builder.WriteString("🚀 <b>Событие началось!</b>\n\n")
	} else {
		timeRemaining := b.formatTimeRemaining(timeUntilEvent)
		builder.WriteString(fmt.Sprintf("📌 <b>Напоминание о событии</b>%s\n\n", timeRemaining))
	}

	builder.WriteString(fmt.Sprintf("<b>%s</b>\n", event.Title))

	if event.Description != "" {
		builder.WriteString(fmt.Sprintf("\n%s\n", event.Description))
	}

	dateStr := formatEventDateStr(event.Date, event.Timezone)
	tzLabel := formatTimezoneLabel(event.Timezone)
	builder.WriteString(fmt.Sprintf("\n📆 <b>Дата:</b> %s (%s)\n", dateStr, tzLabel))

	if len(event.Hosts) > 0 {
		builder.WriteString("\n👥 <b>Спикеры:</b>\n")
		for _, host := range event.Hosts {
			name := strings.TrimSpace(fmt.Sprintf("%s %s", host.FirstName, host.LastName))
			if name == "" {
				name = host.Username
			}

			if host.Username != "" {
				builder.WriteString(fmt.Sprintf("• %s (@%s)\n", name, host.Username))
			} else {
				builder.WriteString(fmt.Sprintf("• %s\n", name))
			}
		}
	}

	if event.PlaceType == models.EventOnline {
		builder.WriteString(fmt.Sprintf("\n🔗 <b>Ссылка:</b> %s\n", event.Place))
	} else {
		place := event.Place
		if event.CustomPlaceType != "" {
			place = event.CustomPlaceType + ", " + event.Place
		}
		builder.WriteString(fmt.Sprintf("\n📍 <b>Место:</b> %s\n", place))
	}

	// Добавляем информацию о повторениях
	if event.IsRepeating && event.RepeatPeriod != nil {
		builder.WriteString("\n🔄 <b>Повторяющееся событие:</b> ")
		interval := 1
		if event.RepeatInterval != nil {
			interval = *event.RepeatInterval
		}

		periodLabels := map[string]string{
			"DAILY":   "день",
			"WEEKLY":  "неделя",
			"MONTHLY": "месяц",
			"YEARLY":  "год",
		}

		periodLabel := periodLabels[*event.RepeatPeriod]
		if periodLabel == "" {
			periodLabel = strings.ToLower(*event.RepeatPeriod)
		}

		if interval == 1 {
			builder.WriteString(fmt.Sprintf("каждый %s", periodLabel))
		} else {
			builder.WriteString(fmt.Sprintf("каждые %d %s", interval, b.pluralizePeriod(interval, periodLabel)))
		}

		if event.RepeatEndDate != nil {
			loc := getEventLocation(event.Timezone)
			endDateStr := event.RepeatEndDate.In(loc).Format("02.01.2006")
			builder.WriteString(fmt.Sprintf(" до %s", endDateStr))
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

func (b *TelegramBot) formatTimeRemaining(timeUntilEvent time.Duration) string {
	if timeUntilEvent <= 0 {
		return " (событие началось)"
	}

	days := int(timeUntilEvent.Hours()) / 24
	hours := int(timeUntilEvent.Hours()) % 24
	minutes := int(timeUntilEvent.Minutes()) % 60

	var parts []string
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d %s", days, b.pluralize(days, "день", "дня", "дней")))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d %s", hours, b.pluralize(hours, "час", "часа", "часов")))
	}
	if minutes > 0 && days == 0 {
		parts = append(parts, fmt.Sprintf("%d %s", minutes, b.pluralize(minutes, "минута", "минуты", "минут")))
	}

	if len(parts) > 0 {
		return fmt.Sprintf(" (до события осталось %s)", strings.Join(parts, " "))
	}

	return ""
}

func (b *TelegramBot) pluralize(n int, one, few, many string) string {
	if n%10 == 1 && n%100 != 11 {
		return one
	}
	if n%10 >= 2 && n%10 <= 4 && (n%100 < 10 || n%100 >= 20) {
		return few
	}
	return many
}

func (b *TelegramBot) pluralizePeriod(n int, period string) string {
	forms := map[string][]string{
		"день":   {"дня", "дней"},
		"неделя": {"недели", "недель"},
		"месяц":  {"месяца", "месяцев"},
		"год":    {"года", "лет"},
	}

	if forms[period] == nil {
		return period
	}

	if n%10 == 1 && n%100 != 11 {
		return period
	}
	if n%10 >= 2 && n%10 <= 4 && (n%100 < 10 || n%100 >= 20) {
		return forms[period][0]
	}
	return forms[period][1]
}

// handleCallbackQuery обрабатывает нажатия на callback кнопки
func (b *TelegramBot) handleCallbackQuery(callback *tgbotapi.CallbackQuery) {
	data := callback.Data
	userID := callback.From.ID

	// Парсим callback data
	if strings.HasPrefix(data, "event_attend:") {
		eventIdStr := strings.TrimPrefix(data, "event_attend:")
		var eventId int64
		fmt.Sscanf(eventIdStr, "%d", &eventId)

		// Получаем пользователя по telegram_id
		member, err := b.member.GetByTelegramID(userID)
		if err != nil {
			log.Printf("Error getting member by telegram ID %d: %v", userID, err)
			b.answerCallbackQuery(callback.ID, "Ошибка: пользователь не найден")
			return
		}

		// Обновляем подписку на SUBSCRIBED
		_, err = b.eventAlertSubscription.UpdateSubscriptionStatus(eventId, member.Id, models.EventAlertStatusSubscribed)
		if err != nil {
			log.Printf("Error updating subscription status: %v", err)
			b.answerCallbackQuery(callback.ID, "Ошибка при обновлении подписки")
			return
		}

		// Синхронизируем участие: добавляем в event_members
		_, err = b.eventService.AddMember(int(eventId), int(member.Id))
		if err != nil {
			log.Printf("Error adding member %d to event %d: %v", member.Id, eventId, err)
		}

		b.answerCallbackQuery(callback.ID, "Отлично! Вы записаны на мероприятие")

		// Обновляем сообщение, убирая кнопки
		editMsg := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, callback.Message.Text)
		editMsg.ParseMode = "HTML"
		b.bot.Send(editMsg)

	} else if strings.HasPrefix(data, "event_decline:") {
		eventIdStr := strings.TrimPrefix(data, "event_decline:")
		var eventId int64
		fmt.Sscanf(eventIdStr, "%d", &eventId)

		// Получаем пользователя по telegram_id
		member, err := b.member.GetByTelegramID(userID)
		if err != nil {
			log.Printf("Error getting member by telegram ID %d: %v", userID, err)
			b.answerCallbackQuery(callback.ID, "Ошибка: пользователь не найден")
			return
		}

		// Обновляем подписку на UNSUBSCRIBED
		_, err = b.eventAlertSubscription.UpdateSubscriptionStatus(eventId, member.Id, models.EventAlertStatusUnsubscribed)
		if err != nil {
			log.Printf("Error updating subscription status: %v", err)
			b.answerCallbackQuery(callback.ID, "Ошибка при обновлении подписки")
			return
		}

		// Синхронизируем участие: убираем из event_members
		_, err = b.eventService.RemoveMember(int(eventId), int(member.Id))
		if err != nil {
			log.Printf("Error removing member %d from event %d: %v", member.Id, eventId, err)
		}

		b.answerCallbackQuery(callback.ID, "Вы отписаны от мероприятия")

		// Обновляем сообщение, убирая кнопки
		editMsg := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, callback.Message.Text)
		editMsg.ParseMode = "HTML"
		b.bot.Send(editMsg)
	}
}

// answerCallbackQuery отвечает на callback query
func (b *TelegramBot) answerCallbackQuery(callbackID string, text string) {
	callbackConfig := tgbotapi.NewCallback(callbackID, text)
	if _, err := b.bot.Request(callbackConfig); err != nil {
		log.Printf("Error answering callback query: %v", err)
	}
}

// getNotificationSettingsMap получает настройки уведомлений для списка участников
func (b *TelegramBot) getNotificationSettingsMap(members []models.Member) map[int64]*models.NotificationSettings {
	memberIds := make([]int64, 0, len(members))
	for _, m := range members {
		if m.TelegramID != 0 {
			memberIds = append(memberIds, m.Id)
		}
	}
	settingsMap, err := b.notificationSettingsService.GetByMemberIds(memberIds)
	if err != nil {
		log.Printf("Error getting notification settings: %v", err)
		return make(map[int64]*models.NotificationSettings)
	}
	return settingsMap
}

// SendInitialEventAlerts отправляет инициализирующие алерты всем подписанным пользователям
func (b *TelegramBot) SendInitialEventAlerts(event *models.Event) error {
	members, err := b.member.GetSubscribedMembersWithTelegram()
	if err != nil {
		return fmt.Errorf("error getting subscribed members: %v", err)
	}

	// Для эксклюзивных событий — алерты только участникам указанного чата
	var exclusiveMemberIDs map[int64]bool
	if event.ExclusiveChatID != nil && *event.ExclusiveChatID != 0 {
		chatActivitySvc := service.NewChatActivityService()
		memberIDs, chatErr := chatActivitySvc.GetMemberIDsByChatID(*event.ExclusiveChatID)
		if chatErr != nil {
			log.Printf("Error getting exclusive chat members: %v", chatErr)
		} else {
			exclusiveMemberIDs = make(map[int64]bool, len(memberIDs))
			for _, id := range memberIDs {
				exclusiveMemberIDs[id] = true
			}
			log.Printf("Exclusive event %d: sending alerts to %d chat members", event.Id, len(exclusiveMemberIDs))
		}
	}

	settingsMap := b.getNotificationSettingsMap(members)

	for _, member := range members {
		if member.TelegramID == 0 {
			continue
		}

		// Фильтруем по эксклюзивному чату
		if exclusiveMemberIDs != nil && !exclusiveMemberIDs[member.Id] {
			continue
		}

		// Проверяем настройки уведомлений
		if s, ok := settingsMap[member.Id]; ok && (s.MuteAll || !s.NewEvents) {
			continue
		}

		_, err := b.eventAlertSubscription.CreateSubscription(event.Id, member.Id)
		if err != nil {
			log.Printf("Error creating subscription for member %d: %v", member.Id, err)
			continue
		}

		err = b.SendEventAlert(member.TelegramID, event, true)
		if err != nil {
			if strings.Contains(err.Error(), "chat not found") {
				continue
			}
			log.Printf("Error sending event alert to user %d: %v", member.TelegramID, err)
			continue
		}
	}

	return nil
}

func (b *TelegramBot) SendRepeatingEventAlert(event *models.Event, alertType string) error {
	members, err := b.eventAlertSubscription.GetSubscribedMembersForEvent(event.Id)
	if err != nil {
		return fmt.Errorf("error getting subscribed members for event: %v", err)
	}

	settingsMap := b.getNotificationSettingsMap(members)

	for _, member := range members {
		if member.TelegramID == 0 {
			continue
		}

		// Проверяем настройки уведомлений по типу алерта
		if s, ok := settingsMap[member.Id]; ok && s.MuteAll {
			continue
		}
		if s, ok := settingsMap[member.Id]; ok {
			switch alertType {
			case "first":
				if !s.RemindWeek {
					continue
				}
			case "second":
				if !s.RemindDay {
					continue
				}
			case "third":
				if !s.RemindHour {
					continue
				}
			case "start":
				if !s.EventStart {
					continue
				}
			default:
				log.Printf("Unknown alertType %q for event %d", alertType, event.Id)
			}
		}

		err = b.SendEventAlert(member.TelegramID, event, false)
		if err != nil {
			if strings.Contains(err.Error(), "chat not found") {
				continue
			}
			log.Printf("Error sending repeating event alert to user %d: %v", member.TelegramID, err)
			continue
		}
	}

	return nil
}

// SendEventUpdateAlert отправляет уведомление об изменении события всем подписанным пользователям
func (b *TelegramBot) SendEventUpdateAlert(event *models.Event) error {
	members, err := b.eventAlertSubscription.GetSubscribedMembersForEvent(event.Id)
	if err != nil {
		return fmt.Errorf("error getting subscribed members for event: %v", err)
	}

	settingsMap := b.getNotificationSettingsMap(members)

	for _, member := range members {
		if member.TelegramID == 0 {
			continue
		}

		// Проверяем настройки уведомлений
		if s, ok := settingsMap[member.Id]; ok && (s.MuteAll || !s.EventUpdates) {
			continue
		}

		messageText := b.formatEventUpdateAlert(event)
		msg := tgbotapi.NewMessage(member.TelegramID, messageText)
		msg.ParseMode = "HTML"

		_, err = b.bot.Send(msg)
		if err != nil {
			if strings.Contains(err.Error(), "chat not found") {
				continue
			}
			log.Printf("Error sending event update alert to user %d: %v", member.TelegramID, err)
			continue
		}
	}

	return nil
}

// SendEventCancelAlert отправляет уведомление об отмене события всем подписанным пользователям
func (b *TelegramBot) SendEventCancelAlert(event *models.Event) error {
	members, err := b.eventAlertSubscription.GetSubscribedMembersForEvent(event.Id)
	if err != nil {
		return fmt.Errorf("error getting subscribed members for event: %v", err)
	}

	settingsMap := b.getNotificationSettingsMap(members)

	for _, member := range members {
		if member.TelegramID == 0 {
			continue
		}

		// Проверяем настройки уведомлений
		if s, ok := settingsMap[member.Id]; ok && (s.MuteAll || !s.EventCancelled) {
			continue
		}

		messageText := fmt.Sprintf("❌ <b>Событие отменено!</b>\n\n<b>%s</b>\n\nСобытие было отменено организаторами.", event.Title)
		msg := tgbotapi.NewMessage(member.TelegramID, messageText)
		msg.ParseMode = "HTML"

		_, err = b.bot.Send(msg)
		if err != nil {
			if strings.Contains(err.Error(), "chat not found") {
				continue
			}
			log.Printf("Error sending event cancel alert to user %d: %v", member.TelegramID, err)
			continue
		}
	}

	return nil
}

// formatEventUpdateAlert форматирует сообщение об изменении события
func (b *TelegramBot) formatEventUpdateAlert(event *models.Event) string {
	var builder strings.Builder

	builder.WriteString("📝 <b>Событие изменено!</b>\n\n")
	builder.WriteString(fmt.Sprintf("<b>%s</b>\n", event.Title))

	if event.Description != "" {
		builder.WriteString(fmt.Sprintf("\n%s\n", event.Description))
	}

	dateStr := formatEventDateStr(event.Date, event.Timezone)
	tzLabel := formatTimezoneLabel(event.Timezone)
	builder.WriteString(fmt.Sprintf("\n📆 <b>Дата:</b> %s (%s)\n", dateStr, tzLabel))

	if len(event.Hosts) > 0 {
		builder.WriteString("\n👥 <b>Спикеры:</b>\n")
		for _, host := range event.Hosts {
			name := strings.TrimSpace(fmt.Sprintf("%s %s", host.FirstName, host.LastName))
			if name == "" {
				name = host.Username
			}

			if host.Username != "" {
				builder.WriteString(fmt.Sprintf("• %s (@%s)\n", name, host.Username))
			} else {
				builder.WriteString(fmt.Sprintf("• %s\n", name))
			}
		}
	}

	if event.PlaceType == models.EventOnline {
		builder.WriteString(fmt.Sprintf("\n🔗 <b>Ссылка:</b> %s\n", event.Place))
	} else {
		place := event.Place
		if event.CustomPlaceType != "" {
			place = event.CustomPlaceType + ", " + event.Place
		}
		builder.WriteString(fmt.Sprintf("\n📍 <b>Место:</b> %s\n", place))
	}

	// Добавляем информацию о повторениях
	if event.IsRepeating && event.RepeatPeriod != nil {
		builder.WriteString("\n🔄 <b>Повторяющееся событие:</b> ")
		interval := 1
		if event.RepeatInterval != nil {
			interval = *event.RepeatInterval
		}

		periodLabels := map[string]string{
			"DAILY":   "день",
			"WEEKLY":  "неделя",
			"MONTHLY": "месяц",
			"YEARLY":  "год",
		}

		periodLabel := periodLabels[*event.RepeatPeriod]
		if periodLabel == "" {
			periodLabel = strings.ToLower(*event.RepeatPeriod)
		}

		if interval == 1 {
			builder.WriteString(fmt.Sprintf("каждый %s", periodLabel))
		} else {
			builder.WriteString(fmt.Sprintf("каждые %d %s", interval, b.pluralizePeriod(interval, periodLabel)))
		}

		if event.RepeatEndDate != nil {
			loc := getEventLocation(event.Timezone)
			endDateStr := event.RepeatEndDate.In(loc).Format("02.01.2006")
			builder.WriteString(fmt.Sprintf(" до %s", endDateStr))
		}
		builder.WriteString("\n")
	}

	builder.WriteString("\n💡 <i>Пожалуйста, проверьте актуальную информацию о событии</i>")

	return builder.String()
}

func (b *TelegramBot) startEventAlertsScheduler() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		b.checkAndSendEventAlerts()
	}
}

func (b *TelegramBot) checkAndSendEventAlerts() {
	now := time.Now()
	futureEvents, err := b.eventService.GetFutureEvents(now.Add(-1 * time.Minute))
	if err != nil {
		log.Printf("Error getting future events for alerts: %v", err)
		return
	}

	for _, event := range futureEvents {
		b.checkReminderAlert(&event, now)

		// Для повторяющихся событий проверяем все будущие повторения
		if event.IsRepeating && event.RepeatPeriod != nil {
			b.checkRepeatingEventOccurrences(&event, now)
		} else {
			// Для обычных событий проверяем только исходную дату
			b.checkRepeatingAlerts(&event, now)
		}
	}
}

// getNextOccurrence вычисляет следующее повторение события после указанной даты
func (b *TelegramBot) getNextOccurrence(event *models.Event, after time.Time) *time.Time {
	if !event.IsRepeating || event.RepeatPeriod == nil {
		return nil
	}

	interval := 1
	if event.RepeatInterval != nil {
		interval = *event.RepeatInterval
	}

	// Проверяем, не истекло ли событие
	if event.RepeatEndDate != nil && after.After(*event.RepeatEndDate) {
		return nil
	}

	// Начинаем с исходной даты события
	currentDate := event.Date

	// Если исходная дата уже прошла, вычисляем следующее повторение
	if currentDate.Before(after) || currentDate.Equal(after) {
		switch *event.RepeatPeriod {
		case "DAILY":
			daysSinceStart := int(after.Sub(currentDate).Hours() / 24)
			nextOccurrenceDays := ((daysSinceStart / interval) + 1) * interval
			currentDate = currentDate.AddDate(0, 0, nextOccurrenceDays)
		case "WEEKLY":
			weeksSinceStart := int(after.Sub(currentDate).Hours() / (24 * 7))
			nextOccurrenceWeeks := ((weeksSinceStart / interval) + 1) * interval
			currentDate = currentDate.AddDate(0, 0, nextOccurrenceWeeks*7)
		case "MONTHLY":
			monthsSinceStart := 0
			tempDate := currentDate
			for tempDate.Before(after) || tempDate.Equal(after) {
				tempDate = tempDate.AddDate(0, interval, 0)
				if tempDate.Before(after) || tempDate.Equal(after) {
					monthsSinceStart++
				}
			}
			currentDate = currentDate.AddDate(0, (monthsSinceStart+1)*interval, 0)
		case "YEARLY":
			yearsSinceStart := 0
			tempDate := currentDate
			for tempDate.Before(after) || tempDate.Equal(after) {
				tempDate = tempDate.AddDate(interval, 0, 0)
				if tempDate.Before(after) || tempDate.Equal(after) {
					yearsSinceStart++
				}
			}
			currentDate = currentDate.AddDate((yearsSinceStart+1)*interval, 0, 0)
		}
	}

	// Проверяем ограничения по дате окончания
	if event.RepeatEndDate != nil && currentDate.After(*event.RepeatEndDate) {
		return nil
	}

	return &currentDate
}

// checkRepeatingEventOccurrences проверяет и отправляет алерты для всех будущих повторений события
func (b *TelegramBot) checkRepeatingEventOccurrences(event *models.Event, now time.Time) {
	// Получаем следующее повторение события
	nextOccurrence := b.getNextOccurrence(event, now)
	if nextOccurrence == nil {
		return
	}

	// Создаем временное событие с датой следующего повторения для проверки алертов
	tempEvent := *event
	tempEvent.Date = *nextOccurrence
	b.checkRepeatingAlerts(&tempEvent, now)
}

func (b *TelegramBot) getReminderInterval() time.Duration {
	return time.Duration(config.CFG.AlertReminderIntervalMinutes) * time.Minute
}

func (b *TelegramBot) checkReminderAlert(event *models.Event, now time.Time) {
	subscriptions, err := b.eventAlertSubscription.GetPendingSubscriptionsForEvent(event.Id)
	if err != nil {
		log.Printf("Error getting pending subscriptions: %v", err)
		return
	}

	reminderInterval := b.getReminderInterval()

	for _, subscription := range subscriptions {
		if subscription.ReminderSentAt != nil {
			if subscription.ReminderSentAt.Add(reminderInterval).Before(now) {
				_, err := b.eventAlertSubscription.UpdateSubscriptionStatus(
					subscription.EventId,
					subscription.MemberId,
					models.EventAlertStatusUnsubscribed,
				)
				if err != nil {
					log.Printf("Error unsubscribing after reminder: %v", err)
				}
			}
			continue
		}

		timeSinceCreation := now.Sub(subscription.CreatedAt)
		if timeSinceCreation >= reminderInterval {
			member, err := b.member.GetById(subscription.MemberId)
			if err != nil || member.TelegramID == 0 {
				continue
			}

			err = b.SendEventAlert(member.TelegramID, event, true)
			if err != nil {
				if strings.Contains(err.Error(), "chat not found") {
					continue
				}
				log.Printf("Error sending reminder alert to user %d: %v", member.TelegramID, err)
				continue
			}

			reminderTime := now
			subscription.ReminderSentAt = &reminderTime
			_, err = b.eventAlertSubscription.CreateOrUpdate(&subscription)
			if err != nil {
				log.Printf("Error updating subscription reminder time: %v", err)
			}
		}
	}
}

func (b *TelegramBot) getAlertIntervals() (alertFirst, alertSecond, alertThird time.Duration) {
	return time.Duration(config.CFG.AlertReminderFirstIntervalMinutes) * time.Minute,
		time.Duration(config.CFG.AlertReminderSecondIntervalMinutes) * time.Minute,
		time.Duration(config.CFG.AlertReminderThirdIntervalMinutes) * time.Minute
}

func (b *TelegramBot) checkRepeatingAlerts(event *models.Event, now time.Time) {
	eventTime := event.Date
	timeUntilEvent := eventTime.Sub(now)

	alertFirst, alertSecond, alertThird := b.getAlertIntervals()

	eventLocation := getEventLocation(event.Timezone)
	nowInMoscow := now.In(eventLocation)

	scheduledHour := config.CFG.AlertScheduledHour
	scheduledMinute := config.CFG.AlertScheduledMinute

	shouldSend := false
	var alertType string

	if timeUntilEvent <= 1*time.Minute && timeUntilEvent > -2*time.Minute {
		alertType = "start"
		shouldSend = true
	} else if timeUntilEvent <= alertThird && timeUntilEvent > 1*time.Minute {
		alertType = "third"
		shouldSend = true
	} else if timeUntilEvent <= alertSecond && timeUntilEvent > alertThird {
		if nowInMoscow.Hour() == scheduledHour && nowInMoscow.Minute() == scheduledMinute {
			alertType = "second"
			shouldSend = true
		}
	} else if timeUntilEvent <= alertFirst && timeUntilEvent > alertSecond {
		if nowInMoscow.Hour() == scheduledHour && nowInMoscow.Minute() == scheduledMinute {
			alertType = "first"
			shouldSend = true
		}
	}

	if shouldSend {
		if event.LastRepeatingAlertSentAt != nil {
			if alertType == "start" {
				timeSinceLastAlert := now.Sub(*event.LastRepeatingAlertSentAt)
				if timeSinceLastAlert < 2*time.Minute {
					return
				}
			} else {
				lastSentDay := event.LastRepeatingAlertSentAt.Day()
				lastSentMonth := event.LastRepeatingAlertSentAt.Month()
				lastSentYear := event.LastRepeatingAlertSentAt.Year()
				currentDay := now.Day()
				currentMonth := now.Month()
				currentYear := now.Year()

				if lastSentDay == currentDay && lastSentMonth == currentMonth && lastSentYear == currentYear {
					return
				}
			}
		}

		log.Printf("Sending repeating alert for event %d, type: %s, timeUntilEvent: %v", event.Id, alertType, timeUntilEvent)
		if err := b.SendRepeatingEventAlert(event, alertType); err != nil {
			log.Printf("Error sending repeating alert: %v", err)
			return
		}

		if err := database.DB.Model(&models.Event{}).
			Where("id = ?", event.Id).
			Update("last_repeating_alert_sent_at", now).Error; err != nil {
			log.Printf("Error updating event last alert sent time: %v", err)
		}
	}
}

type AuthRequest struct {
	Token     string      `json:"token"`
	UserID    int64       `json:"user_id"`
	Username  string      `json:"username"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Role      models.Role `json:"role"`
}

func sendAuthToBackend(token string, user *tgbotapi.User) {
	isSubcriber, err := CheckUserInChat(user.ID)
	if err != nil {
		log.Println("Ошибка проверки пользователя в чате:", err)
	}
	var role models.Role

	if isSubcriber {
		role = models.MemberRoleSubscriber
	} else {
		role = models.MemberRoleUnsubscriber
	}

	data := AuthRequest{
		Token:     token,
		UserID:    user.ID,
		Username:  user.UserName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      role,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Ошибка сериализации JSON:", err)
		return
	}

	url := fmt.Sprintf("%s/api/auth/telegram-from-bot", config.CFG.BackendDomain)

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Ошибка создания запроса:", err)
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Bot-Secret", config.CFG.BotSharedSecret)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Println("Ошибка отправки запроса:", err)
		return
	}
	defer resp.Body.Close()

	log.Println("Ответ от Fiber:", resp.Status)
}

func CheckUserInChat(userID int64) (bool, error) {
	telegramApiUrl := fmt.Sprintf("https://api.telegram.org/bot%s/getChatMember?chat_id=%d&user_id=%d", config.CFG.TelegramToken, config.CFG.TelegramMainChatID, userID)

	resp, err := http.Get(telegramApiUrl)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	var result struct {
		Ok     bool `json:"ok"`
		Result struct {
			Status string `json:"status"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	if !result.Ok {
		return false, fmt.Errorf("telegram API error")
	}

	switch result.Result.Status {
	case "member", "administrator", "creator":
		return true, nil
	default:
		return false, nil
	}
}
