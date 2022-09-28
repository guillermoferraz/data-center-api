package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/guillermoferraz/data-center-api/controllers"
	"github.com/guillermoferraz/data-center-api/db"
	"github.com/guillermoferraz/data-center-api/models"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	port := ":" + os.Getenv("PORT")

	/* Database connection */
	db.DBConnection()

	/* migrations tables */
	db.DB.AutoMigrate(&models.User{}, &models.Module{}, &models.Submodule{})

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})
	controllers.UseAuthController(app)
	controllers.UseUserController(app)
	controllers.UseModuleController(app)
	controllers.UseSubmoduleController(app)
	app.Listen(port)
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}
}
