package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/guillermoferraz/data-center-api/db"
	"github.com/guillermoferraz/data-center-api/middleware"
	"github.com/guillermoferraz/data-center-api/models"
)

type Module struct {
	Id          uuid.UUID `gorm:"type:uuid;primary_key"`
	UserId      string
	Created_at  time.Time `gorm:"<-:create"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Private     bool      `json:"private"`
}

func UseModuleController(router fiber.Router) {
	loadEnv()
	router.Post("/addmodule", func(c *fiber.Ctx) error {
		reqBody := Module{}
		if err := c.BodyParser(&reqBody); err != nil {
			return err
		}
		reqHeader := c.Request().Header.Peek("Authorization")
		token := string(reqHeader)
		userId := middleware.UseIsAuthorized(token)
		if userId != "Error" {
			user := models.User{}
			db.DB.Find(&user, "id = ?", userId)
			uuid, err := uuid.NewV4()
			if err != nil {
				return err
			}
			created_at := time.Now()
			module := Module{uuid, user.Id.String(), created_at, reqBody.Name, reqBody.Description, reqBody.Private}
			db.DB.Create(&module)
			return c.Status(200).JSON(fiber.Map{
				"message": "New module saved successfully",
				"status":  200,
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"message": "Error saving new module",
			"status":  500,
		})
	})

	router.Get("/modules", func(c *fiber.Ctx) error {
		reqHeader := c.Request().Header.Peek("Authorization")
		token := string(reqHeader)
		userId := middleware.UseIsAuthorized(token)
		if userId != "Error" {
			modules := []models.Module{}
			db.DB.Find(&modules, "user_id = ?", userId)
			return c.Status(200).JSON(modules)
		}
		return c.Status(500).JSON(fiber.Map{
			"message": "Error on get modules",
			"status":  500,
		})
	})
	router.Delete("/module", func(c *fiber.Ctx) error {
		moduleId := string(c.Request().URI().QueryString())
		reqHeader := c.Request().Header.Peek("Authorization")
		token := string(reqHeader)
		userId := middleware.UseIsAuthorized(token)
		if userId != "Error" {
			user := models.User{}
			modules := models.Module{}
			submodules := models.Submodule{}
			db.DB.Find(&user, "id = ?", userId)
			db.DB.Unscoped().Delete(&modules, "Id = ?", moduleId)
			db.DB.Unscoped().Delete(&submodules, "module_id = ?", moduleId)
			return c.Status(200).JSON(fiber.Map{
				"message": "Submodule deleted successfully",
				"status":  200,
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"message": "Error on delete submodule",
			"status":  500,
		})
	})

	router.Patch("/module/edit", func(c *fiber.Ctx) error {
		reqBody := Module{}
		if err := c.BodyParser(&reqBody); err != nil {
			return err
		}
		reqHeader := c.Request().Header.Peek("Authorization")
		token := string(reqHeader)
		userId := middleware.UseIsAuthorized(token)
		if userId != "Error" {
			user := models.User{}
			db.DB.Find(&user, "id = ?", userId)

			created_at := time.Now()
			module := Module{reqBody.Id, user.Id.String(), created_at, reqBody.Name, reqBody.Description, reqBody.Private}
			db.DB.Model(&module).Where("Id = ?", reqBody.Id).Update("name", reqBody.Name).Update("description", reqBody.Description).Update("private", reqBody.Private)
			return c.Status(200).JSON(fiber.Map{
				"message": "Module updated successfully",
				"status":  200,
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to update module",
			"status":  500,
		})
	})

}
