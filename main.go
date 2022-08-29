package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/guillermoferraz/data-center-api/controllers"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	port := ":" + os.Getenv("PORT")
	app := fiber.New()
	app.Use(logger.New())
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})
	controllers.UseUsersController(app)
	app.Listen(port)
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}
}
