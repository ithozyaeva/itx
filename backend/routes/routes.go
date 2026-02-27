package routes

import (
	"ithozyeva/internal/handler"
	"ithozyeva/internal/middleware"
	"ithozyeva/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	SetupPublicRoutes(app, db)
	SetupAdminRoutes(app, db)
	SetupPlatformRoutes(app, db)
}
func SetupPublicRoutes(app *fiber.App, db *gorm.DB) {
	// Инициализация сервисов и репозиториев
	telegramAuthHandler := handler.NewTelegramAuthHandler()

	api := app.Group("/api")

	// Маршруты для авторизации через Telegram
	auth := api.Group("/auth")
	auth.Post("/telegram/refresh", telegramAuthHandler.RefreshToken)
	auth.Post("/telegram", telegramAuthHandler.Authenticate)
	auth.Post("/telegram-from-bot", telegramAuthHandler.HandleBotMessage)

	userHandler := handler.NewUserHandler()
	// Маршруты для аутентификации в админ панели
	// TODO: Рассмотреть вариант авторизации в админке через тг + роль admin у member (подумать о его переименование в user, или же оставить тотальное разделение между public и admin зонами)
	auth.Post("/login", userHandler.Login)
	auth.Post("/refresh", userHandler.RefreshToken)

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

func SetupAdminRoutes(app *fiber.App, db *gorm.DB) {
	authMiddleware := middleware.NewAuthMiddleware(db)

	// Защищенные маршруты
	protected := app.Group("/api/admin", authMiddleware.RequireAuth)

	// Маршруты для менторов
	mentorHandler := handler.NewMentorHandler()
	mentors := protected.Group("/mentors", authMiddleware.RequirePermission(models.PermissionCanViewAdminMentors))
	mentors.Get("/", mentorHandler.GetAllWithRelations)
	mentors.Get("/:id", mentorHandler.GetById)
	mentors.Post("/", authMiddleware.RequirePermission(models.PermissionCanEditAdminMentors), mentorHandler.Create)
	mentors.Put("/:id", authMiddleware.RequirePermission(models.PermissionCanEditAdminMentors), mentorHandler.Update)
	mentors.Delete("/:id", authMiddleware.RequirePermission(models.PermissionCanEditAdminMentors), mentorHandler.Delete)
	mentors.Post("/review", mentorHandler.AddReviewToService)
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

	// Маршруты для тегов ивентов
	eventTagHandler := handler.NewEventTagHandler()
	eventTags := protected.Group("/eventTags")
	eventTags.Get("/", eventTagHandler.Search)
	eventTags.Get("/:id", eventTagHandler.GetById)
	eventTags.Post("/", authMiddleware.RequirePermission(models.PermissionCanEditAdminEvents), eventTagHandler.Create)
	eventTags.Put("/:id", authMiddleware.RequirePermission(models.PermissionCanEditAdminEvents), eventTagHandler.Update)
	eventTags.Delete("/:id", authMiddleware.RequirePermission(models.PermissionCanEditAdminEvents), eventTagHandler.Delete)
}

func SetupPlatformRoutes(app *fiber.App, db *gorm.DB) {
	authMiddleware := middleware.NewAuthMiddleware(db)

	// Защищенные маршруты
	protected := app.Group("/api/platform", authMiddleware.RequireTGAuth)

	// Маршруты для отзывов о сообществе
	reviewHandler := handler.NewReviewOnCommunityHandler()
	reviews := protected.Group("/reviews")
	reviews.Post("/add", reviewHandler.AddReview)

	// Маршруты для участников
	memberHandler := handler.NewMembersHandler()
	members := protected.Group("/members")
	members.Get("/me", memberHandler.Me)
	members.Patch("/me", memberHandler.UpdateProfile)

	// Маршруты для ментора
	mentorsHandler := handler.NewMentorHandler()
	mentorsMe := protected.Group("/mentors/me")
	mentorsMe.Post("/update-info", authMiddleware.RequirePermission(models.PermissionCanEditPlatformMentors), mentorsHandler.UpdateInfo)
	mentorsMe.Post("/update-prof-tags", authMiddleware.RequirePermission(models.PermissionCanEditPlatformMentors), mentorsHandler.UpdateProfTags)
	mentorsMe.Post("/update-services", authMiddleware.RequirePermission(models.PermissionCanEditPlatformMentors), mentorsHandler.UpdateServices)
	mentorsMe.Post("/update-contacts", authMiddleware.RequirePermission(models.PermissionCanEditPlatformMentors), mentorsHandler.UpdateContacts)

	// Маршруты для ивентов
	eventHandler := handler.NewEventsHandler()
	events := protected.Group("/events")
	events.Get("/", eventHandler.Search)
	events.Post("/apply", eventHandler.AddMember)
	events.Post("/decline", eventHandler.RemoveMember)

	// Маршурты для таблицы рефералов
	referalsHandler := handler.NewReferalLinkHandler()
	referals := protected.Group("/referals")
	referals.Get("/", referalsHandler.Search)
	referals.Post("/add-link", referalsHandler.AddLink)
	referals.Put("/update-link", referalsHandler.UpdateLink)
	referals.Delete("/delete-link", referalsHandler.DeleteLink)

	resumeHandler := handler.NewResumeHandler()
	resumes := protected.Group("/resumes")
	resumes.Post("/", resumeHandler.Upload)
	resumes.Get("/me", resumeHandler.ListMy)
	resumes.Patch("/:id", resumeHandler.UpdateMy)
	resumes.Delete("/:id", resumeHandler.DeleteMy)
}
