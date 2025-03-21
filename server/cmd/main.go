package main

import (
	"fmt"
	"ithozyeva/config"
	"ithozyeva/database"
	"ithozyeva/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config.LoadConfig()

	database.InitDBConnection()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.CFG.CorsUrls,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true, // если нужны куки/авторизация
	}))

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", config.CFG.Serverport)))
}
