package main

import (
	"ithozyeva/config"
	"ithozyeva/database"
	"ithozyeva/internal/bot"
	"ithozyeva/internal/service"
	"ithozyeva/routes"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Загружаем конфигурацию
	config.LoadConfig()

	// Подключаемся к базе данных
	database.SetupDatabase()

	// Создаем экземпляр Fiber
	app := fiber.New(fiber.Config{
		AppName: "ITX API",
	})

	// Добавляем middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Telegram-User-Token",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Настраиваем маршруты
	routes.SetupRoutes(app, database.DB)

	// Запускаем фоновую задачу для автозамораживания реферальных ссылок
	go func() {
		referalSvc := service.NewReferalLinkService()
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		// Выполняем сразу при старте
		if count, err := referalSvc.ExpireLinks(); err != nil {
			log.Printf("Error expiring referral links: %v", err)
		} else if count > 0 {
			log.Printf("Expired %d referral links", count)
		}

		for range ticker.C {
			if count, err := referalSvc.ExpireLinks(); err != nil {
				log.Printf("Error expiring referral links: %v", err)
			} else if count > 0 {
				log.Printf("Expired %d referral links", count)
			}
		}
	}()

	// Запускаем фоновую задачу для начисления баллов за события и бонусов активности
	go func() {
		pointsSvc := service.NewPointsService()
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()

		// Выполняем сразу при старте
		pointsSvc.AwardPointsForPastEvents()
		pointsSvc.AwardActivityBonuses()

		for range ticker.C {
			pointsSvc.AwardPointsForPastEvents()
			pointsSvc.AwardActivityBonuses()
		}
	}()

	// Запускаем Telegram бота в отдельной горутине
	go func() {
		telegramBot, err := bot.NewTelegramBot()
		if err != nil {
			log.Printf("Error creating bot: %v", err)
			return
		}

		// Устанавливаем глобальный экземпляр бота
		bot.SetGlobalBot(telegramBot)

		log.Println("Telegram bot started successfully")
		telegramBot.Start()
	}()

	// Запускаем сервер
	log.Printf("Server starting on port %s", config.CFG.Port)
	if err := app.Listen(":" + config.CFG.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
