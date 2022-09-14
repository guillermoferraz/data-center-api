package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guillermoferraz/data-center-api/db"
	"github.com/guillermoferraz/data-center-api/middleware"
	"github.com/guillermoferraz/data-center-api/models"
)

type UserData struct {
	Id        string `gorm:"type:uuid;primary_key;"`
	Email     string `gorm:"not null;unique_index"`
	Firstname string `json:"Firstname"`
	Lastname  string `json:"Lastname"`
}
type Request struct {
	Token string `json:"Token"`
}
type ProfileUser struct {
	Id        string `gorm:"type:uuid;primary_key;"`
	Email     string `gorm:"not null;unique_index"`
	Firstname string `json:"Firstname"`
	Lastname  string `json:"Lastname"`
}

func UseUserController(router fiber.Router) {
	loadEnv()
	router.Get("/user", func(c *fiber.Ctx) error {
		reqHeader := c.Request().Header.Peek("Authorization")
		token := string(reqHeader)
		userId := middleware.UseIsAuthorized(token)
		if userId != "Error" {
			var user models.User
			db.DB.Find(&user, "id = ?", userId)
			return c.Status(200).JSON(fiber.Map{
				"Id":        user.Id,
				"Email":     user.Email,
				"Firstname": user.Firstname,
				"Lastname":  user.Lastname,
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"message": "User not found",
		})
	})

}
