package main

import (
	"context"
	"ithozyeva/config"
	"ithozyeva/database"
	"ithozyeva/internal/bot"
	"ithozyeva/internal/service"
	"ithozyeva/internal/utils"
	"ithozyeva/routes"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Загружаем конфигурацию
	config.LoadConfig()

	log.Printf("Starting in %s mode", config.CFG.AppMode)

	// Подключаемся к базе данных
	database.SetupDatabase()

	// Инициализируем глобальный S3 клиент для presigned URL
	utils.InitGlobalS3()

	// Инициализируем Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.CFG.Redis.Addr(),
		Password: config.CFG.Redis.Password,
		DB:       config.CFG.Redis.DB,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Redis connected successfully")

	// API сервер и фоновые задачи — только в режимах full и api
	if config.CFG.AppMode != "bot" {
		// Создаем экземпляр Fiber
		app := fiber.New(fiber.Config{
			AppName: "ITX API",
		})

		// Добавляем middleware
		app.Use(logger.New())
		app.Use(cors.New(cors.Config{
			AllowOrigins: config.CFG.AllowedOrigins,
			AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Telegram-User-Token",
			AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		}))

		app.Use(limiter.New(limiter.Config{
			Max:        100,
			Expiration: 60 * time.Second,
		}))

		// Настраиваем маршруты
		routes.SetupRoutes(app, database.DB, redisClient)

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

		// Запускаем фоновую задачу для выдачи ачивки "Чаттер недели" (каждый понедельник)
		go func() {
			pointsSvc := service.NewPointsService()
			ticker := time.NewTicker(24 * time.Hour)
			defer ticker.Stop()

			// При старте — если сегодня понедельник, запускаем сразу
			if time.Now().UTC().Weekday() == time.Monday {
				pointsSvc.AwardWeeklyChatter()
			}

			for range ticker.C {
				if time.Now().UTC().Weekday() == time.Monday {
					pointsSvc.AwardWeeklyChatter()
				}
			}
		}()

		// Запускаем фоновую задачу для очистки старых сообщений чатов (раз в сутки)
		go func() {
			chatActivitySvc := service.NewChatActivityService()
			ticker := time.NewTicker(24 * time.Hour)
			defer ticker.Stop()

			// Выполняем сразу при старте
			chatActivitySvc.CleanupOldMessages(90)

			for range ticker.C {
				chatActivitySvc.CleanupOldMessages(90)
			}
		}()

		// Запускаем фоновую задачу для розыгрышей (проверка каждые 5 минут)
		go func() {
			raffleSvc := service.NewRaffleService()
			ticker := time.NewTicker(5 * time.Minute)
			defer ticker.Stop()

			raffleSvc.DrawExpiredRaffles()

			for range ticker.C {
				raffleSvc.DrawExpiredRaffles()
			}
		}()

		// Запускаем сервер
		go func() {
			log.Printf("Server starting on port %s", config.CFG.Port)
			if err := app.Listen(":" + config.CFG.Port); err != nil {
				log.Fatalf("Failed to start server: %v", err)
			}
		}()
	}

	// Telegram бот — только в режимах full и bot
	if config.CFG.AppMode != "api" {
		go func() {
			telegramBot, err := bot.NewTelegramBot(redisClient)
			if err != nil {
				log.Printf("Error creating bot: %v", err)
				return
			}

			// Устанавливаем глобальный экземпляр бота
			bot.SetGlobalBot(telegramBot)

			log.Println("Telegram bot started successfully")
			telegramBot.Start()
		}()
	}

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")
}
