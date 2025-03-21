package routes

import (
	"ithozyeva/internal/handler"
	"ithozyeva/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Маршруты для менторов
	mentorHandler := handler.NewMentorHandler()
	mentors := app.Group("/mentors")
	mentors.Get("/", mentorHandler.GetAllWithRelations)
	mentors.Get("/:id", mentorHandler.GetById)
	mentors.Post("/", middleware.Protected(), mentorHandler.Create)
	mentors.Put("/:id", middleware.Protected(), mentorHandler.Update)
	mentors.Delete("/:id", middleware.Protected(), mentorHandler.Delete)
	mentors.Post("/findByTag", mentorHandler.FindByTag)
	mentors.Post("/review", middleware.Protected(), mentorHandler.AddReviewToService)
	mentors.Get("/:id/services", mentorHandler.GetServices)

	// Маршруты для профессиональных тегов
	profTagHandler := handler.NewProfTagsHandler()
	profTags := app.Group("/profTags")
	profTags.Get("/", profTagHandler.Search)
	profTags.Get("/:id", profTagHandler.GetById)
	profTags.Post("/", middleware.Protected(), profTagHandler.Create)
	profTags.Put("/", middleware.Protected(), profTagHandler.Update)
	profTags.Delete("/:id", middleware.Protected(), profTagHandler.Delete)

	// Маршруты для участников
	memberHandler := handler.NewMembersHandler()
	members := app.Group("/members")
	members.Get("/", memberHandler.Search)
	members.Get("/:id", memberHandler.GetById)
	members.Post("/", middleware.Protected(), memberHandler.Create)
	members.Put("/", middleware.Protected(), memberHandler.Update)
	members.Delete("/:id", middleware.Protected(), memberHandler.Delete)

	// Маршруты для отзывов о сообществе
	reviewHandler := handler.NewReviewOnCommunityHandler()
	reviews := app.Group("/reviews")
	reviews.Get("/", reviewHandler.GetAllWithAuthor)
	reviews.Post("/", middleware.Protected(), reviewHandler.AddReview)
	reviews.Get("/:id", reviewHandler.GetById)
	reviews.Put("/:id", middleware.Protected(), reviewHandler.Update)
	reviews.Delete("/:id", middleware.Protected(), reviewHandler.Delete)

	// Маршруты для отзывов на услуги
	reviewOnServiceHandler := handler.NewReviewOnServiceHandler()
	reviewsOnService := app.Group("/reviews-on-service")
	reviewsOnService.Get("/", reviewOnServiceHandler.GetReviewsWithMentorInfo)
	reviewsOnService.Get("/:id", reviewOnServiceHandler.GetById)
	reviewsOnService.Post("/", middleware.Protected(), reviewOnServiceHandler.CreateReview)
	reviewsOnService.Put("/:id", middleware.Protected(), reviewOnServiceHandler.Update)
	reviewsOnService.Delete("/:id", middleware.Protected(), reviewOnServiceHandler.Delete)

	// Маршруты для аутентификации
	userHandler := handler.NewUserHandler()
	auth := app.Group("/auth")
	auth.Post("/login", userHandler.Login)
	auth.Post("/refresh", userHandler.RefreshToken)
}
