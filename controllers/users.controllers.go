package controllers

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/guillermoferraz/data-center-api/db"
	"github.com/guillermoferraz/data-center-api/models"
	"github.com/joho/godotenv"
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
type Claims struct {
	Id uuid.UUID `gorm:"type:uuid"`
	jwt.StandardClaims
}

func UseUsersController(router fiber.Router) {
	loadEnv()
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

				mySecret := os.Getenv("JWT_SECRET")
				expirationTime := time.Now().Add(24 * time.Hour)
				claims := &Claims{
					Id: uuid,
					StandardClaims: jwt.StandardClaims{
						// In JWT, the expiry time is expressed as unix milliseconds
						ExpiresAt: expirationTime.Unix(),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, err := token.SignedString([]byte(mySecret))
				fmt.Println("error:", err)

				db.DB.Create(&user)
				return c.Status(200).JSON(fiber.Map{
					"message": "Register successfully",
					"status":  200,
					"token":   tokenString,
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
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}
}
