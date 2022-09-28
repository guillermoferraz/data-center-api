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
			submodule := Submodule{uuid, user.Id.String(), reqBody.ModuleId, created_at, reqBody.Name, reqBody.Description, reqBody.Private}
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
}
