package controllers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/guillermoferraz/data-center-api/db"
	"github.com/guillermoferraz/data-center-api/models"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id         uuid.UUID `gorm:"type:uuid;primary_key;"`
	Created_at time.Time `gorm:"<-:create"`
	Firstname  string    `json:"firstname"`
	Lastname   string    `json:"lastname"`
	Email      string    `gorm:"not null;unique_index"`
	Password   string    `json:"password"`
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

			/* save data */
			uuid, err := uuid.NewV4()
			if err != nil {
				return err
			}
			created_at := time.Now()
			user := User{uuid, created_at, reqBody.Firstname, reqBody.Lastname, reqBody.Email, string(hashedPassword)}
			fmt.Println(reqBody.Email)

			var model_user models.User

			existEmail := db.DB.Find(&model_user, "email = ?", user.Email)
			if existEmail.RowsAffected == 0 {
				db.DB.Create(&user)
				return c.Status(200).JSON(fiber.Map{
					"message": "Register successfully",
				})
			} else {
				return c.Status(409).JSON((fiber.Map{
					"message": "The email entered alredy exists",
				}))
			}
		}
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid received data",
		})
	})
}
