package routes

import (
	"ithozyeva/internal/handler"
	"ithozyeva/internal/middleware"
	"ithozyeva/internal/models"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, redisClient *redis.Client) {
	SetupPublicRoutes(app, db)
	SetupAdminRoutes(app, db, redisClient)
	SetupPlatformRoutes(app, db)
}
func SetupPublicRoutes(app *fiber.App, db *gorm.DB) {
	// Инициализация сервисов и репозиториев
	telegramAuthHandler, err := handler.NewTelegramAuthHandler()
	if err != nil {
		log.Fatalf("Failed to create TelegramAuthHandler: %v", err)
	}

	api := app.Group("/api")

	// Маршруты для авторизации через Telegram
	auth := api.Group("/auth")
	auth.Post("/telegram/refresh", telegramAuthHandler.RefreshToken)
	auth.Post("/telegram", telegramAuthHandler.Authenticate)
	auth.Post("/telegram-from-bot", telegramAuthHandler.HandleBotMessage)

	mentorHandler := handler.NewMentorHandler()
	api.Get("/mentors", mentorHandler.GetAllWithRelations)

	// Маршруты для профессиональных тегов
	profTagHandler := handler.NewProfTagsHandler()
	api.Get("/profTags", profTagHandler.Search)

	// Маршруты для участников
	memberHandler := handler.NewMembersHandler()
	api.Get("/members", memberHandler.Search)

	// Маршруты для отзывов на услуги
	reviewOnServiceHandler := handler.NewReviewOnServiceHandler()
	api.Get("/reviews-on-service", reviewOnServiceHandler.GetReviewsWithMentorInfo)

	// Маршруты для отзывов о сообществе
	reviewHandler := handler.NewReviewOnCommunityHandler()
	api.Get("/review-on-community", reviewHandler.GetApproved)

	eventsHandler := handler.NewEventsHandler()
	api.Get("/events/old", eventsHandler.GetOld)
	api.Get("/events/next", eventsHandler.GetNext)
	api.Get("/events/ics", eventsHandler.GetICSFile)

	// Маршруты для словарей
	dictionaryHandler := handler.NewDictionaryHandler()
	api.Get("/dictionaries", dictionaryHandler.GetDictionaries)
}

func SetupAdminRoutes(app *fiber.App, db *gorm.DB, redisClient *redis.Client) {
	authMiddleware := middleware.NewAuthMiddleware(db)

	// Защищенные маршруты
	protected := app.Group("/api/admin", authMiddleware.RequireAuth)

	// Маршруты для статистики
	statsHandler := handler.NewStatsHandler()
	protected.Get("/stats", statsHandler.GetStats)
	protected.Get("/stats/charts", statsHandler.GetChartStats)

	// Маршруты для менторов
	mentorHandler := handler.NewMentorHandler()
	mentors := protected.Group("/mentors", authMiddleware.RequirePermission(models.PermissionCanViewAdminMentors))
	mentors.Get("/", mentorHandler.GetAllWithRelations)
	mentors.Get("/:id", mentorHandler.GetById)
	mentors.Post("/", authMiddleware.RequirePermission(models.PermissionCanEditAdminMentors), mentorHandler.Create)
	mentors.Put("/:id", authMiddleware.RequirePermission(models.PermissionCanEditAdminMentors), mentorHandler.Update)
	mentors.Delete("/:id", authMiddleware.RequirePermission(models.PermissionCanEditAdminMentors), mentorHandler.Delete)
	mentors.Post("/review", mentorHandler.AddReviewToService)
	mentors.Patch("/:id/order", authMiddleware.RequirePermission(models.PermissionCanEditAdminMentors), mentorHandler.UpdateOrder)
	mentors.Get("/:id/services", mentorHandler.GetServices)

	// Здесь будут защищенные маршруты

	// Маршруты для профессиональных тегов
	profTagHandler := handler.NewProfTagsHandler()
	profTags := protected.Group("/profTags")
	profTags.Get("/", profTagHandler.Search)
	profTags.Get("/:id", profTagHandler.GetById)
	profTags.Post("/", profTagHandler.Create)
	profTags.Put("/", profTagHandler.Update)
	profTags.Delete("/:id", profTagHandler.Delete)

	// Маршруты для участников
	memberHandler := handler.NewMembersHandler()
	members := protected.Group("/members", authMiddleware.RequirePermission(models.PermissionCanViewAdminMembers))
	members.Get("/", memberHandler.Search)
	members.Post("/", authMiddleware.RequirePermission(models.PermissionCanEditAdminMembers), memberHandler.Create)
	members.Get("/:id", memberHandler.GetById)
	members.Put("/:id", authMiddleware.RequirePermission(models.PermissionCanEditAdminMembers), memberHandler.Update)
	members.Delete("/:id", authMiddleware.RequirePermission(models.PermissionCanEditAdminMembers), memberHandler.Delete)

	protected.Get("/me/permissions", memberHandler.GetPermissions)
	// Маршруты для отзывов о сообществе
	reviewHandler := handler.NewReviewOnCommunityHandler()
	reviews := protected.Group("/reviews", authMiddleware.RequirePermission(models.PermissionCanViewAdminReviews))
	reviews.Post("/:id/approve", authMiddleware.RequirePermission(models.PermissionCanApprovedAdminReviews), reviewHandler.Approve)
	reviews.Get("/", reviewHandler.GetAllWithAuthor)
	reviews.Post("/", authMiddleware.RequirePermission(models.PermissionCanEditAdminReviews), reviewHandler.CreateReview)
	reviews.Get("/:id", reviewHandler.GetById)
	reviews.Patch("/:id", authMiddleware.RequirePermission(models.PermissionCanEditAdminReviews), reviewHandler.Update)
	reviews.Delete("/:id", authMiddleware.RequirePermission(models.PermissionCanEditAdminReviews), reviewHandler.Delete)

	// Маршруты для отзывов на услуги
	reviewOnServiceHandler := handler.NewReviewOnServiceHandler()
	reviewsOnService := protected.Group("/reviews-on-service", authMiddleware.RequirePermission(models.PermissionCanViewAdminMentorsReview))
	reviewsOnService.Get("/", reviewOnServiceHandler.Search)
	reviewsOnService.Get("/:id", reviewOnServiceHandler.GetById)
	reviewsOnService.Post("/", authMiddleware.RequirePermission(models.PermissionCanEditAdminMentorsReview), reviewOnServiceHandler.CreateReview)
	reviewsOnService.Post("/:id/approve", authMiddleware.RequirePermission(models.PermissionCanApproveAdminMentorsReview), reviewOnServiceHandler.Approve)
	reviewsOnService.Patch("/:id", authMiddleware.RequirePermission(models.PermissionCanEditAdminMentorsReview), reviewOnServiceHandler.Update)
	reviewsOnService.Delete("/:id", authMiddleware.RequirePermission(models.PermissionCanEditAdminMentorsReview), reviewOnServiceHandler.Delete)

	// Маршруты для ивентов
	eventHandler := handler.NewEventsHandler()
	events := protected.Group("/events", authMiddleware.RequirePermission(models.PermissionCanViewAdminEvents))
	events.Get("/", eventHandler.Search)
	events.Get("/:id", eventHandler.GetById)
	events.Post("/", authMiddleware.RequirePermission(models.PermissionCanEditAdminEvents), eventHandler.Create)
	events.Put("/:id", authMiddleware.RequirePermission(models.PermissionCanEditAdminEvents), eventHandler.Update)
	events.Delete("/:id", authMiddleware.RequirePermission(models.PermissionCanEditAdminEvents), eventHandler.Delete)
	resumeHandler := handler.NewResumeHandler()
	resumes := protected.Group("/resumes", authMiddleware.RequirePermission(models.PermissionCanViewAdminResumes))
	resumes.Get("/", resumeHandler.AdminList)
	resumes.Get("/download", resumeHandler.AdminDownload)
	resumes.Get("/:id", resumeHandler.AdminGet)

	// Маршруты для массовых операций
	bulkHandler := handler.NewBulkHandler()
	bulk := protected.Group("/bulk")
	bulk.Post("/events/delete", authMiddleware.RequirePermission(models.PermissionCanEditAdminEvents), bulkHandler.BulkDeleteEvents)
	bulk.Post("/mentors/delete", authMiddleware.RequirePermission(models.PermissionCanEditAdminMentors), bulkHandler.BulkDeleteMentors)
	bulk.Post("/members/delete", authMiddleware.RequirePermission(models.PermissionCanEditAdminMembers), bulkHandler.BulkDeleteMembers)
	bulk.Post("/reviews/delete", authMiddleware.RequirePermission(models.PermissionCanEditAdminReviews), bulkHandler.BulkDeleteReviews)
	bulk.Post("/reviews/approve", authMiddleware.RequirePermission(models.PermissionCanApprovedAdminReviews), bulkHandler.BulkApproveReviews)
	bulk.Post("/mentors-reviews/delete", authMiddleware.RequirePermission(models.PermissionCanEditAdminMentorsReview), bulkHandler.BulkDeleteMentorsReviews)
	bulk.Post("/mentors-reviews/approve", authMiddleware.RequirePermission(models.PermissionCanApproveAdminMentorsReview), bulkHandler.BulkApproveServiceReviews)

	// Маршруты для баллов (админ)
	pointsHandler := handler.NewPointsHandler()
	points := protected.Group("/points", authMiddleware.RequirePermission(models.PermissionCanViewAdminPoints))
	points.Get("/", pointsHandler.AdminSearch)
	points.Post("/", authMiddleware.RequirePermission(models.PermissionCanEditAdminPoints), pointsHandler.AdminAward)
	points.Delete("/:id", authMiddleware.RequirePermission(models.PermissionCanEditAdminPoints), pointsHandler.AdminDelete)

	// Маршруты для журнала действий
	auditLogHandler := handler.NewAuditLogHandler()
	protected.Get("/audit-logs", authMiddleware.RequirePermission(models.PermissionCanViewAdminAuditLogs), auditLogHandler.Search)

	// Маршруты для реферальных ссылок (админ)
	referalHandler := handler.NewReferalLinkHandler()
	adminReferals := protected.Group("/admin-referals")
	adminReferals.Get("/", referalHandler.AdminSearch)
	adminReferals.Get("/:id", referalHandler.AdminGetById)
	adminReferals.Delete("/:id", referalHandler.AdminDelete)

	// Маршруты для активности чатов
	chatActivityHandler := handler.NewChatActivityHandler()
	chatActivity := protected.Group("/chat-activity")
	chatActivity.Get("/stats", chatActivityHandler.GetStats)
	chatActivity.Get("/chart", chatActivityHandler.GetChart)
	chatActivity.Get("/top-users", chatActivityHandler.GetTopUsers)
	chatActivity.Get("/chats", chatActivityHandler.GetChats)
	chatActivity.Get("/user-stats", chatActivityHandler.GetUserStats)
	chatActivity.Get("/export", chatActivityHandler.ExportCSV)

	// Маршруты для заданий чатов (админ)
	chatQuestHandler := handler.NewChatQuestHandler()
	chatQuests := protected.Group("/chat-quests")
	chatQuests.Get("/", chatQuestHandler.GetAll)
	chatQuests.Post("/", chatQuestHandler.Create)
	chatQuests.Put("/:id", chatQuestHandler.Update)
	chatQuests.Delete("/:id", chatQuestHandler.Delete)

	// Маршруты для сезонов (админ)
	seasonHandler := handler.NewSeasonHandler()
	adminSeasons := protected.Group("/seasons")
	adminSeasons.Get("/", seasonHandler.GetAll)
	adminSeasons.Post("/", seasonHandler.Create)
	adminSeasons.Post("/:id/finish", seasonHandler.Finish)

	// Маршруты для розыгрышей (админ)
	raffleHandler := handler.NewRaffleHandler()
	adminRaffles := protected.Group("/raffles")
	adminRaffles.Get("/", raffleHandler.GetAllAdmin)
	adminRaffles.Post("/", raffleHandler.Create)
	adminRaffles.Delete("/:id", raffleHandler.Delete)

	// Маршруты для мини-игр (админ)
	casinoHandler := handler.NewCasinoHandler()
	adminCasino := protected.Group("/minigames")
	adminCasino.Get("/stats", casinoHandler.GetAdminStats)
	adminCasino.Get("/bets", casinoHandler.GetAdminBets)

	// Маршруты для тегов ивентов
	eventTagHandler := handler.NewEventTagHandler()
	eventTags := protected.Group("/eventTags")
	eventTags.Get("/", eventTagHandler.Search)
	eventTags.Get("/:id", eventTagHandler.GetById)
	eventTags.Post("/", authMiddleware.RequirePermission(models.PermissionCanEditAdminEvents), eventTagHandler.Create)
	eventTags.Put("/:id", authMiddleware.RequirePermission(models.PermissionCanEditAdminEvents), eventTagHandler.Update)
	eventTags.Delete("/:id", authMiddleware.RequirePermission(models.PermissionCanEditAdminEvents), eventTagHandler.Delete)

	// Маршруты для подписок
	if redisClient != nil {
		subscriptionHandler := handler.NewSubscriptionHandler(redisClient)
		subs := protected.Group("/subscriptions", authMiddleware.RequirePermission(models.PermissionCanViewAdminSubscriptions))
		subs.Get("/stats", subscriptionHandler.GetStats)
		subs.Get("/tiers", subscriptionHandler.GetTiers)
		subs.Get("/chats", subscriptionHandler.GetChats)
		subs.Get("/users", subscriptionHandler.GetUsers)
		subs.Get("/users/:id", subscriptionHandler.GetUser)
		subs.Put("/users/:id/override", authMiddleware.RequirePermission(models.PermissionCanEditAdminSubscriptions), subscriptionHandler.SetOverride)
		subs.Delete("/users/:id/override", authMiddleware.RequirePermission(models.PermissionCanEditAdminSubscriptions), subscriptionHandler.ClearOverride)
		subs.Delete("/users/:id/access/:chatId", authMiddleware.RequirePermission(models.PermissionCanEditAdminSubscriptions), subscriptionHandler.RevokeAccess)
	}
}

func SetupPlatformRoutes(app *fiber.App, db *gorm.DB) {
	authMiddleware := middleware.NewAuthMiddleware(db)

	// Защищенные маршруты
	protected := app.Group("/api/platform", authMiddleware.RequireTGAuth)

	// Маршруты для отзывов о сообществе
	reviewHandler := handler.NewReviewOnCommunityHandler()
	reviews := protected.Group("/reviews")
	reviews.Post("/add", reviewHandler.AddReview)
	reviews.Get("/my", reviewHandler.GetMyReviews)
	reviews.Patch("/:id", reviewHandler.UpdateMyReview)
	reviews.Delete("/:id", reviewHandler.DeleteMyReview)

	// Маршруты для участников
	memberHandler := handler.NewMembersHandler()
	members := protected.Group("/members")
	members.Get("/me", memberHandler.Me)
	members.Get("/:id", memberHandler.GetPublicProfile)
	members.Patch("/me", memberHandler.UpdateProfile)
	members.Post("/me/avatar", memberHandler.UploadAvatar)

	// Маршруты для менторов
	mentorsHandler := handler.NewMentorHandler()
	mentors := protected.Group("/mentors")
	mentors.Get("/:id", mentorsHandler.GetById)
	mentors.Post("/:id/reviews", mentorsHandler.AddReviewFromPlatform)

	mentorsMe := protected.Group("/mentors/me")
	mentorsMe.Post("/update-info", authMiddleware.RequirePermission(models.PermissionCanEditPlatformMentors), mentorsHandler.UpdateInfo)
	mentorsMe.Post("/update-prof-tags", authMiddleware.RequirePermission(models.PermissionCanEditPlatformMentors), mentorsHandler.UpdateProfTags)
	mentorsMe.Post("/update-services", authMiddleware.RequirePermission(models.PermissionCanEditPlatformMentors), mentorsHandler.UpdateServices)
	mentorsMe.Post("/update-contacts", authMiddleware.RequirePermission(models.PermissionCanEditPlatformMentors), mentorsHandler.UpdateContacts)

	// Маршруты для ивентов
	eventHandler := handler.NewEventsHandler()
	events := protected.Group("/events")
	events.Get("/", eventHandler.Search)
	events.Get("/:id", eventHandler.GetById)
	events.Post("/apply", eventHandler.AddMember)
	events.Post("/decline", eventHandler.RemoveMember)

	// Маршурты для таблицы рефералов
	referalsHandler := handler.NewReferalLinkHandler()
	referals := protected.Group("/referals")
	referals.Get("/", referalsHandler.Search)
	referals.Post("/add-link", referalsHandler.AddLink)
	referals.Put("/update-link", referalsHandler.UpdateLink)
	referals.Delete("/delete-link", referalsHandler.DeleteLink)
	referals.Post("/track-conversion", referalsHandler.TrackConversion)

	resumeHandler := handler.NewResumeHandler()
	resumes := protected.Group("/resumes")
	resumes.Post("/", resumeHandler.Upload)
	resumes.Get("/me", resumeHandler.ListMy)
	resumes.Get("/:id/download", resumeHandler.DownloadMy)
	resumes.Patch("/:id", resumeHandler.UpdateMy)
	resumes.Delete("/:id", resumeHandler.DeleteMy)

	// Маршруты для уведомлений
	notificationHandler := handler.NewNotificationHandler()
	notifications := protected.Group("/notifications")
	notifications.Get("/", notificationHandler.GetMy)
	notifications.Get("/unread-count", notificationHandler.GetUnreadCount)
	notifications.Patch("/:id/read", notificationHandler.MarkAsRead)
	notifications.Post("/read-all", notificationHandler.MarkAllAsRead)

	// Маршруты для баллов
	pointsHandler := handler.NewPointsHandler()
	points := protected.Group("/points")
	points.Get("/me", pointsHandler.GetMyPoints)
	points.Get("/leaderboard", pointsHandler.GetLeaderboard)

	// Маршруты для заданий чатов (платформа)
	chatQuestHandler := handler.NewChatQuestHandler()
	chatQuests := protected.Group("/chat-quests")
	chatQuests.Get("/all", chatQuestHandler.GetAllForMember)
	chatQuests.Get("/active", chatQuestHandler.GetActive)

	// Маршруты для достижений
	achievementHandler := handler.NewAchievementHandler()
	achievements := protected.Group("/achievements")
	achievements.Get("/me", achievementHandler.GetMyAchievements)
	achievements.Get("/member/:id", achievementHandler.GetMemberAchievements)

	// Маршруты для биржи заданий
	taskExchangeHandler := handler.NewTaskExchangeHandler()
	tasks := protected.Group("/tasks")
	tasks.Get("/", taskExchangeHandler.Search)
	tasks.Get("/:id", taskExchangeHandler.GetById)
	tasks.Post("/", taskExchangeHandler.Create)
	tasks.Put("/:id", taskExchangeHandler.Update)
	tasks.Post("/:id/assign", taskExchangeHandler.Assign)
	tasks.Post("/:id/unassign", taskExchangeHandler.Unassign)
	tasks.Delete("/:id/assignees/:memberId", taskExchangeHandler.RemoveAssignee)
	tasks.Post("/:id/done", taskExchangeHandler.MarkDone)
	tasks.Post("/:id/approve", authMiddleware.RequirePermission(models.PermissionCanApprovePlatformTasks), taskExchangeHandler.Approve)
	tasks.Post("/:id/reject", authMiddleware.RequirePermission(models.PermissionCanApprovePlatformTasks), taskExchangeHandler.Reject)
	tasks.Delete("/:id", taskExchangeHandler.Delete)

	// Маршруты для настроек уведомлений
	notifSettingsHandler := handler.NewNotificationSettingsHandler()
	notifSettings := protected.Group("/notification-settings")
	notifSettings.Get("/", notifSettingsHandler.GetMy)
	notifSettings.Patch("/", notifSettingsHandler.UpdateMy)

	// Маршруты для хайлайтов из чатов
	highlightHandler := handler.NewChatHighlightHandler()
	highlights := protected.Group("/highlights")
	highlights.Get("/recent", highlightHandler.GetRecent)
	highlights.Get("/", highlightHandler.Search)

	// Маршруты для барахолки
	marketplaceHandler := handler.NewMarketplaceHandler()
	marketplace := protected.Group("/marketplace")
	marketplace.Get("/", marketplaceHandler.Search)
	marketplace.Get("/:id", marketplaceHandler.GetById)
	marketplace.Post("/", marketplaceHandler.Create)
	marketplace.Patch("/:id", marketplaceHandler.Update)
	marketplace.Post("/:id/request-purchase", marketplaceHandler.RequestPurchase)
	marketplace.Post("/:id/cancel-purchase", marketplaceHandler.CancelPurchase)
	marketplace.Post("/:id/sold", marketplaceHandler.MarkSold)
	marketplace.Delete("/:id", marketplaceHandler.Delete)

	// Маршруты для стены благодарностей
	kudosHandler := handler.NewKudosHandler()
	kudos := protected.Group("/kudos")
	kudos.Get("/", kudosHandler.GetRecent)
	kudos.Post("/", kudosHandler.Send)

	// Маршруты для сезонов
	seasonHandler := handler.NewSeasonHandler()
	seasons := protected.Group("/seasons")
	seasons.Get("/", seasonHandler.GetAll)
	seasons.Get("/active", seasonHandler.GetActive)
	seasons.Get("/:id/leaderboard", seasonHandler.GetLeaderboard)

	// Маршруты для розыгрышей
	raffleHandler := handler.NewRaffleHandler()
	raffles := protected.Group("/raffles")
	raffles.Get("/", raffleHandler.GetAll)
	raffles.Post("/:id/buy", raffleHandler.BuyTickets)

	// Маршруты для казино
	casinoHandler := handler.NewCasinoHandler()
	casino := protected.Group("/minigames")
	casino.Post("/coin-flip", casinoHandler.PlayCoinFlip)
	casino.Post("/dice-roll", casinoHandler.PlayDiceRoll)
	casino.Post("/wheel", casinoHandler.PlayWheel)
	casino.Get("/history", casinoHandler.GetHistory)
	casino.Get("/feed", casinoHandler.GetGlobalFeed)
	casino.Get("/stats", casinoHandler.GetStats)

	// Маршруты для гильдий
	guildHandler := handler.NewGuildHandler()
	guilds := protected.Group("/guilds")
	guilds.Get("/", guildHandler.GetAll)
	guilds.Post("/", guildHandler.Create)
	guilds.Put("/:id", guildHandler.Update)
	guilds.Delete("/:id", guildHandler.Delete)
	guilds.Post("/:id/join", guildHandler.Join)
	guilds.Post("/:id/leave", guildHandler.Leave)
	guilds.Get("/:id/members", guildHandler.GetMembers)

	// Маршруты для статистики профиля
	profileStatsHandler := handler.NewProfileStatsHandler()
	profileStats := protected.Group("/profile-stats")
	profileStats.Get("/me", profileStatsHandler.GetMyStats)
	profileStats.Get("/:id", profileStatsHandler.GetMemberStats)

	// SSE (Server-Sent Events) для real-time обновлений
	sseHandler := handler.NewSSEHandler()
	protected.Get("/sse", sseHandler.Stream)
}
