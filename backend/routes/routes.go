package routes

import (
	"ithozyeva/internal/handler"
	"ithozyeva/internal/middleware"
	"ithozyeva/internal/models"
	"ithozyeva/internal/service"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, redisClient *redis.Client) {
	SetupPublicRoutes(app, db)
	SetupInternalRoutes(app, db, redisClient)
	SetupAdminRoutes(app, db, redisClient)
	SetupPlatformRoutes(app, db, redisClient)
}

// SetupInternalRoutes — server-to-server API для соседних сервисов.
// Защита: shared secret в заголовке X-Internal-Secret (RequireInternalSecret).
// CORS НЕ выступает защитой — секрет проверяется на каждом запросе через
// constant-time compare.
func SetupInternalRoutes(app *fiber.App, db *gorm.DB, redisClient *redis.Client) {
	authMiddleware := middleware.NewAuthMiddleware(db)

	internal := app.Group("/api/internal", authMiddleware.RequireInternalSecret)

	if redisClient != nil {
		subscriptionHandler := handler.NewSubscriptionHandler(redisClient)
		internal.Get("/subscription/:tg_id", subscriptionHandler.GetInternalUserSubscription)
	}
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
	auth.Post("/telegram-webapp", telegramAuthHandler.AuthenticateWebApp)
	auth.Post("/telegram-from-bot", telegramAuthHandler.HandleBotMessage)

	mentorHandler := handler.NewMentorHandler()
	api.Get("/mentors", mentorHandler.GetAllWithRelationsPublic)

	// Маршруты для профессиональных тегов
	profTagHandler := handler.NewProfTagsHandler()
	api.Get("/profTags", profTagHandler.Search)

	// Маршруты для участников
	memberHandler := handler.NewMembersHandler()
	api.Get("/members", memberHandler.SearchPublic)

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

	// Маршруты для розыгрышей (админ)
	raffleHandler := handler.NewRaffleHandler()
	adminRaffles := protected.Group("/raffles")
	adminRaffles.Get("/", raffleHandler.GetAllAdmin)
	adminRaffles.Post("/", raffleHandler.Create)
	adminRaffles.Delete("/:id", raffleHandler.Delete)

	// Геймификация (админ): пул дейликов и шаблоны челленджей
	adminDailyTasks := protected.Group("/daily-tasks", authMiddleware.RequirePermission(models.PermissionCanViewAdminPoints))
	adminDailyTaskHandler := handler.NewAdminDailyTaskHandler()
	adminDailyTasks.Get("/", adminDailyTaskHandler.List)
	adminDailyTasks.Get("/sets", adminDailyTaskHandler.RecentSets)
	adminDailyTasks.Post("/", authMiddleware.RequirePermission(models.PermissionCanEditAdminPoints), adminDailyTaskHandler.Create)
	adminDailyTasks.Put("/:id", authMiddleware.RequirePermission(models.PermissionCanEditAdminPoints), adminDailyTaskHandler.Update)
	adminDailyTasks.Delete("/:id", authMiddleware.RequirePermission(models.PermissionCanEditAdminPoints), adminDailyTaskHandler.Delete)

	adminChallenges := protected.Group("/challenges", authMiddleware.RequirePermission(models.PermissionCanViewAdminPoints))
	adminChallengeHandler := handler.NewAdminChallengeHandler()
	adminChallenges.Get("/", adminChallengeHandler.ListTemplates)
	adminChallenges.Get("/instances", adminChallengeHandler.ListInstances)
	adminChallenges.Post("/", authMiddleware.RequirePermission(models.PermissionCanEditAdminPoints), adminChallengeHandler.CreateTemplate)
	adminChallenges.Put("/:id", authMiddleware.RequirePermission(models.PermissionCanEditAdminPoints), adminChallengeHandler.UpdateTemplate)
	adminChallenges.Delete("/:id", authMiddleware.RequirePermission(models.PermissionCanEditAdminPoints), adminChallengeHandler.DeleteTemplate)

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
		subs.Get("/chats/resolve/:id", subscriptionHandler.ResolveChat)
		subs.Get("/chats/:id", subscriptionHandler.GetChatDetail)
		subs.Post("/chats", authMiddleware.RequireSuperAdmin, subscriptionHandler.CreateChat)
		subs.Put("/chats/:id", authMiddleware.RequireSuperAdmin, subscriptionHandler.UpdateChat)
		subs.Delete("/chats/:id", authMiddleware.RequireSuperAdmin, subscriptionHandler.DeleteChat)
		subs.Get("/users", subscriptionHandler.GetUsers)
		subs.Get("/users/:id", subscriptionHandler.GetUser)
		subs.Put("/users/:id/override", authMiddleware.RequirePermission(models.PermissionCanEditAdminSubscriptions), subscriptionHandler.SetOverride)
		subs.Delete("/users/:id/override", authMiddleware.RequirePermission(models.PermissionCanEditAdminSubscriptions), subscriptionHandler.ClearOverride)
		subs.Delete("/users/:id/access/:chatId", authMiddleware.RequirePermission(models.PermissionCanEditAdminSubscriptions), subscriptionHandler.RevokeAccess)
	}

	// Маршруты для обратной связи (NPS)
	feedbackHandler := handler.NewFeedbackHandler()
	feedback := protected.Group("/feedback", authMiddleware.RequirePermission(models.PermissionCanViewAdminFeedback))
	feedback.Get("/", feedbackHandler.AdminList)

	// Маршруты модерации (списки санкций/глобальных банов/голосований + снятие).
	// Действия выполняет бот через Redis pub/sub — backend в РФ, TG API заблокирован.
	if redisClient != nil {
		moderationHandler := handler.NewModerationHandler(redisClient)
		moderation := protected.Group("/moderation", authMiddleware.RequirePermission(models.PermissionCanViewAdminModeration))
		moderation.Get("/sanctions", moderationHandler.GetActiveSanctions)
		moderation.Get("/actions", moderationHandler.GetRecentActions)
		moderation.Get("/global-bans", moderationHandler.GetGlobalBans)
		moderation.Get("/votebans", moderationHandler.GetOpenVotebans)
		moderation.Post("/sanctions/:id/revoke",
			authMiddleware.RequirePermission(models.PermissionCanEditAdminModeration),
			moderationHandler.RevokeSanction)
		moderation.Delete("/global-bans/:user_id",
			authMiddleware.RequirePermission(models.PermissionCanEditAdminModeration),
			moderationHandler.RevokeGlobalBan)
		moderation.Post("/votebans/:id/cancel",
			authMiddleware.RequirePermission(models.PermissionCanEditAdminModeration),
			moderationHandler.CancelVoteban)
	}
}

func SetupPlatformRoutes(app *fiber.App, db *gorm.DB, redisClient *redis.Client) {
	authMiddleware := middleware.NewAuthMiddleware(db)

	// protected — все авторизованные через Telegram, без проверки подписки.
	// Сюда попадают «прогревные» эндпоинты, которые видит UNSUBSCRIBER:
	// профиль, тарифы, FAQ, отзывы, базовый список менторов.
	protected := app.Group("/api/platform", authMiddleware.RequireTGAuth)

	// subscribed — только для активных подписчиков (EffectiveTierID != nil).
	// Гейтится через RequireSubscription за SUBSCRIPTION_GATE_ENABLED флагом.
	subscribed := app.Group("/api/platform", authMiddleware.RequireTGAuth, authMiddleware.RequireSubscription)

	// Shared CommentService — единая точка работы со всеми комментами
	// платформы. Visibility-чекеры per-entity_type инкапсулируют
	// специфику доступа: AI-материал и event — только видимость самой
	// сущности (любой подписчик имеет доступ к разделу — гейт на уровне
	// /ai-materials и /events).
	aiMaterialSvc := service.NewAIMaterialService()
	eventsSvc := service.NewEventsService()
	commentSvc := service.NewCommentService(map[models.CommentEntityType]service.EntityVisibilityChecker{
		models.CommentEntityAIMaterial: service.AIMaterialVisibilityChecker(aiMaterialSvc),
		models.CommentEntityEvent:      service.EventVisibilityChecker(eventsSvc),
	})
	commentHandler := handler.NewCommentHandler(commentSvc)

	// --- protected (доступно UNSUBSCRIBER'ам для прогрева) ---

	// Отзывы о сообществе — нужно UNSUBSCRIBER'у тоже, чтобы оставить отклик.
	reviewHandler := handler.NewReviewOnCommunityHandler()
	reviews := protected.Group("/reviews")
	reviews.Post("/add", reviewHandler.AddReview)
	reviews.Get("/my", reviewHandler.GetMyReviews)
	reviews.Patch("/:id", reviewHandler.UpdateMyReview)
	reviews.Delete("/:id", reviewHandler.DeleteMyReview)

	// Профиль (включая редактирование). UNSUBSCRIBER должен иметь возможность
	// заполнить аватар/био заранее — после оплаты данные не теряются.
	memberHandler := handler.NewMembersHandler()
	members := protected.Group("/members")
	members.Get("/me", memberHandler.Me)
	members.Patch("/me", memberHandler.UpdateProfile)
	members.Post("/me/avatar", memberHandler.UploadAvatar)

	// Менторы — список и детали read-only открыты UNSUBSCRIBER'у (витрина).
	// Действия с менторами (написать отзыв, контакт) требуют подписки.
	mentorsHandler := handler.NewMentorHandler()
	mentors := protected.Group("/mentors")
	mentors.Get("/:id", mentorsHandler.GetByIdPublic)

	// Публичная карточка участника — открыта всем авторизованным. Юзер,
	// перешедший по ссылке /members/:id из /whois в боте, не должен
	// упираться в гейт подписки и улетать на /tariffs прямо с профиля.
	members.Get("/:id", memberHandler.GetPublicProfile)

	// Уведомления и их настройки — нужны и UNSUBSCRIBER'у (например, подписка
	// истекла → push-уведомление, или анонс снижения цены).
	notificationHandler := handler.NewNotificationHandler()
	notifications := protected.Group("/notifications")
	notifications.Get("/", notificationHandler.GetMy)
	notifications.Get("/unread-count", notificationHandler.GetUnreadCount)
	notifications.Patch("/:id/read", notificationHandler.MarkAsRead)
	notifications.Post("/read-all", notificationHandler.MarkAllAsRead)

	notifSettingsHandler := handler.NewNotificationSettingsHandler()
	notifSettings := protected.Group("/notification-settings")
	notifSettings.Get("/", notifSettingsHandler.GetMy)
	notifSettings.Patch("/", notifSettingsHandler.UpdateMy)

	// Публичные тарифы для /tariffs и прогрева в боте.
	if redisClient != nil {
		subscriptionHandler := handler.NewSubscriptionHandler(redisClient)
		protected.Get("/subscriptions/tiers", subscriptionHandler.PublicTiers)
	}

	// Обратная связь о платформе — должен мочь оставить кто угодно.
	feedbackHandler := handler.NewFeedbackHandler()
	protected.Post("/feedback", feedbackHandler.Create)

	// --- subscribed (только для подписчиков) ---

	// Контакт ментора (отзыв) — write-эндпоинт.
	subscribedMentors := subscribed.Group("/mentors")
	subscribedMentors.Post("/:id/reviews", mentorsHandler.AddReviewFromPlatform)

	// Менторские админ-эндпоинты (только для самих менторов с пермишеном).
	mentorsMe := subscribed.Group("/mentors/me")
	mentorsMe.Post("/update-info", authMiddleware.RequirePermission(models.PermissionCanEditPlatformMentors), mentorsHandler.UpdateInfo)
	mentorsMe.Post("/update-prof-tags", authMiddleware.RequirePermission(models.PermissionCanEditPlatformMentors), mentorsHandler.UpdateProfTags)
	mentorsMe.Post("/update-services", authMiddleware.RequirePermission(models.PermissionCanEditPlatformMentors), mentorsHandler.UpdateServices)
	mentorsMe.Post("/update-contacts", authMiddleware.RequirePermission(models.PermissionCanEditPlatformMentors), mentorsHandler.UpdateContacts)

	// События
	eventHandler := handler.NewEventsHandler()
	events := subscribed.Group("/events")
	events.Get("/", eventHandler.Search)
	events.Get("/:id", eventHandler.GetById)
	events.Post("/apply", eventHandler.AddMember)
	events.Post("/decline", eventHandler.RemoveMember)
	// Комменты к событиям — открыты любому подписчику (как остальные
	// /events). Гейт по master+ не требуется, в отличие от AI-материалов.
	events.Get("/:id/comments", commentHandler.ListForEntity(models.CommentEntityEvent))
	events.Post("/:id/comments", commentHandler.CreateForEntity(models.CommentEntityEvent))

	// Индивидуальные операции над комментами — отдельная группа /comments/:id.
	// Доступна на subscribed (любой подписчик), потому что включает комменты
	// к event'ам. Доступ к конкретному комменту контролируется визибилити
	// родительской сущности через visibility-checker внутри CommentService.
	comments := subscribed.Group("/comments")
	comments.Patch("/:id", commentHandler.Update)
	comments.Delete("/:id", commentHandler.Delete)
	comments.Post("/:id/like", commentHandler.ToggleLike)
	comments.Post("/:id/hidden", commentHandler.SetHidden)

	// Реферальные ссылки
	referalsHandler := handler.NewReferalLinkHandler()
	referals := subscribed.Group("/referals")
	referals.Get("/", referalsHandler.Search)
	referals.Post("/add-link", referalsHandler.AddLink)
	referals.Put("/update-link", referalsHandler.UpdateLink)
	referals.Delete("/delete-link", referalsHandler.DeleteLink)
	referals.Post("/track-conversion", referalsHandler.TrackConversion)

	// Резюме
	resumeHandler := handler.NewResumeHandler()
	resumes := subscribed.Group("/resumes")
	resumes.Post("/", resumeHandler.Upload)
	resumes.Get("/me", resumeHandler.ListMy)
	resumes.Get("/:id/download", resumeHandler.DownloadMy)
	resumes.Patch("/:id", resumeHandler.UpdateMy)
	resumes.Delete("/:id", resumeHandler.DeleteMy)

	// Баллы и лидерборд
	pointsHandler := handler.NewPointsHandler()
	points := subscribed.Group("/points")
	points.Get("/me", pointsHandler.GetMyPoints)
	points.Get("/leaderboard", pointsHandler.GetLeaderboard)

	// Чат-квесты
	chatQuestHandler := handler.NewChatQuestHandler()
	chatQuests := subscribed.Group("/chat-quests")
	chatQuests.Get("/all", chatQuestHandler.GetAllForMember)
	chatQuests.Get("/active", chatQuestHandler.GetActive)

	// Достижения
	achievementHandler := handler.NewAchievementHandler()
	achievements := subscribed.Group("/achievements")
	achievements.Get("/me", achievementHandler.GetMyAchievements)
	achievements.Get("/member/:id", achievementHandler.GetMemberAchievements)

	// Биржа заданий
	taskExchangeHandler := handler.NewTaskExchangeHandler()
	tasks := subscribed.Group("/tasks")
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

	// Хайлайты из чатов
	highlightHandler := handler.NewChatHighlightHandler()
	highlights := subscribed.Group("/highlights")
	highlights.Get("/recent", highlightHandler.GetRecent)
	highlights.Get("/", highlightHandler.Search)

	// AI-материалы — открыты любому подписчику. Раньше были master+,
	// но раздел оказался полезным как точка притяжения для всей платной
	// аудитории, поэтому перенесли с tierMaster на subscribed.
	aiMaterialHandler := handler.NewAIMaterialHandler()
	aiMaterials := subscribed.Group("/ai-materials")
	aiMaterials.Get("/", aiMaterialHandler.Search)
	aiMaterials.Get("/tags", aiMaterialHandler.TopTags)
	aiMaterials.Post("/", aiMaterialHandler.Create)
	aiMaterials.Get("/:id", aiMaterialHandler.GetByID)
	// PUT, не PATCH — UpdateAIMaterialRequest требует все поля и валидируется
	// целиком; partial-update семантика не поддерживается.
	aiMaterials.Put("/:id", aiMaterialHandler.Update)
	aiMaterials.Delete("/:id", aiMaterialHandler.Delete)
	aiMaterials.Post("/:id/hidden", aiMaterialHandler.SetHidden)
	aiMaterials.Post("/:id/like", aiMaterialHandler.ToggleLike)
	aiMaterials.Post("/:id/bookmark", aiMaterialHandler.ToggleBookmark)

	// Comments к AI-материалу — list/create на /<entity>/:id/comments,
	// действия над комментом — на /comments/:id (см. блок ниже).
	aiMaterials.Get("/:id/comments", commentHandler.ListForEntity(models.CommentEntityAIMaterial))
	aiMaterials.Post("/:id/comments", commentHandler.CreateForEntity(models.CommentEntityAIMaterial))

	// Барахолка
	marketplaceHandler := handler.NewMarketplaceHandler()
	marketplace := subscribed.Group("/marketplace")
	marketplace.Get("/", marketplaceHandler.Search)
	marketplace.Get("/:id", marketplaceHandler.GetById)
	marketplace.Post("/", marketplaceHandler.Create)
	marketplace.Patch("/:id", marketplaceHandler.Update)
	marketplace.Post("/:id/request-purchase", marketplaceHandler.RequestPurchase)
	marketplace.Post("/:id/cancel-purchase", marketplaceHandler.CancelPurchase)
	marketplace.Post("/:id/sold", marketplaceHandler.MarkSold)
	marketplace.Delete("/:id", marketplaceHandler.Delete)

	// Стена благодарностей
	kudosHandler := handler.NewKudosHandler()
	kudos := subscribed.Group("/kudos")
	kudos.Get("/", kudosHandler.GetRecent)
	kudos.Post("/", kudosHandler.Send)

	// Розыгрыши
	raffleHandler := handler.NewRaffleHandler()
	raffles := subscribed.Group("/raffles")
	raffles.Get("/", raffleHandler.GetAll)
	raffles.Get("/daily/today", raffleHandler.DailyToday)
	raffles.Post("/:id/buy", raffleHandler.BuyTickets)

	// Казино
	casinoHandler := handler.NewCasinoHandler()
	casino := subscribed.Group("/minigames")
	casino.Post("/coin-flip", casinoHandler.PlayCoinFlip)
	casino.Post("/dice-roll", casinoHandler.PlayDiceRoll)
	casino.Post("/wheel", casinoHandler.PlayWheel)
	casino.Get("/history", casinoHandler.GetHistory)
	casino.Get("/feed", casinoHandler.GetGlobalFeed)
	casino.Get("/stats", casinoHandler.GetStats)

	// Статистика профиля
	profileStatsHandler := handler.NewProfileStatsHandler()
	profileStats := subscribed.Group("/profile-stats")
	profileStats.Get("/me", profileStatsHandler.GetMyStats)
	profileStats.Get("/:id", profileStatsHandler.GetMemberStats)

	// Геймификация: ежедневный check-in, дейлики, стрики
	dailiesHandler := handler.NewDailiesHandler()
	dailies := subscribed.Group("/dailies")
	dailies.Get("/today", dailiesHandler.Today)
	dailies.Post("/check-in", dailiesHandler.CheckIn)
	streak := subscribed.Group("/streak")
	streak.Get("/me", dailiesHandler.MyStreak)

	// Челленджи (еженедельные + ежемесячные)
	challengesHandler := handler.NewChallengesHandler()
	challenges := subscribed.Group("/challenges")
	challenges.Get("/", challengesHandler.GetMyChallenges)

	// SSE — реал-тайм обновления (events, casino, и т.п. — премиум-функции).
	sseHandler := handler.NewSSEHandler()
	subscribed.Get("/sse", sseHandler.Stream)
}
