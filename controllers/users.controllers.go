package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func UseUsersController(router fiber.Router) {
	router.Post("/register", func(c *fiber.Ctx) error {
		reqBody := User{}
		if err := c.BodyParser(&reqBody); err != nil {
			return err
		}
		// hash password
		password := []byte(reqBody.Password)
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(hashedPassword))

		firstname := len(reqBody.Firstname)
		lastname := len(reqBody.Lastname)
		email := len(reqBody.Email)
		pass := len(reqBody.Password)

		if firstname > 1 && lastname > 1 && email > 3 && pass > 7 {
			fmt.Printf("%+v\n", reqBody)
			return c.Status(200).JSON(fiber.Map{
				"fistname": reqBody.Firstname,
			})
		}
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid received data",
		})
	})
}
