package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/guillermoferraz/data-center-api/db"
	"github.com/guillermoferraz/data-center-api/middleware"
	"github.com/guillermoferraz/data-center-api/models"
)

type Submodule struct {
	Id          uuid.UUID `gorm:"type:uuid;primary_key"`
	UserId      string
	ModuleId    string
	Created_at  time.Time `gorm:"<-:create"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Private     bool      `json:"private"`
	Content     string    `json:"content"`
}

func UseSubmoduleController(router fiber.Router) {
	loadEnv()
	router.Post("/addsubmodule", func(c *fiber.Ctx) error {
		reqBody := Submodule{}
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
			submodule := Submodule{uuid, user.Id.String(), reqBody.ModuleId, created_at, reqBody.Name, reqBody.Description, reqBody.Private, reqBody.Content}
			db.DB.Create(&submodule)
			return c.Status(200).JSON(fiber.Map{
				"message": "New subomdule added successfully",
				"status":  200,
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"message": "Error saving new submodule",
			"status":  500,
		})
	})

	router.Get("/submodulesbymodule", func(c *fiber.Ctx) error {
		moduleId := string(c.Request().URI().QueryString())
		reqHeader := c.Request().Header.Peek("Authorization")
		token := string(reqHeader)
		userId := middleware.UseIsAuthorized(token)
		if userId != "Error" {
			user := models.User{}
			submodules := []models.Submodule{}
			db.DB.Find(&user, "id = ?", userId)
			db.DB.Find(&submodules, "module_id = ?", moduleId)
			return c.Status(200).JSON(submodules)

		}
		return c.Status(500).JSON(fiber.Map{
			"message": "Error get submodule by module id",
			"status":  500,
		})
	})
	router.Delete("/submodule", func(c *fiber.Ctx) error {
		submoduleId := string(c.Request().URI().QueryString())
		reqHeader := c.Request().Header.Peek("Authorization")
		token := string(reqHeader)
		userId := middleware.UseIsAuthorized(token)
		if userId != "Error" {
			user := models.User{}
			submodules := models.Submodule{}
			db.DB.Find(&user, "id = ?", userId)
			db.DB.Unscoped().Delete(&submodules, "Id = ?", submoduleId)
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

	router.Patch("/submodule/edit", func(c *fiber.Ctx) error {
		reqBody := Submodule{}
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
			submodule := Submodule{reqBody.Id, user.Id.String(), reqBody.ModuleId, created_at, reqBody.Name, reqBody.Description, reqBody.Private, reqBody.Content}
			db.DB.Model(&submodule).
				Update("module_id", reqBody.ModuleId).
				Update("name", reqBody.Name).
				Update("description", reqBody.Description).
				Update("private", reqBody.Private).
				Update("content", reqBody.Content)

			return c.Status(200).JSON(fiber.Map{
				"message": "New subomdule added successfully",
				"status":  200,
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"message": "Error saving new submodule",
			"status":  500,
		})
	})
}
