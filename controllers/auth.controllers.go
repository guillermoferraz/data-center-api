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
type ReturnUser struct {
	Id        string `gorm:"type:uuid;primary_key;"`
	Email     string `gorm:"not null;unique_index"`
	Firstname string `json:"Firstname"`
	Lastname  string `json:"Lastname"`
	Token     string `json:"Token"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Id uuid.UUID `gorm:"type:uuid"`
	jwt.StandardClaims
}

func CheckPasswordHashed(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func UseAuthController(router fiber.Router) {
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

		firstname := len(reqBody.Firstname)
		lastname := len(reqBody.Lastname)
		email := len(reqBody.Email)
		pass := len(reqBody.Password)

		if firstname > 1 && lastname > 1 && email > 3 && pass > 7 {
			// fmt.Printf("%+v\n", reqBody)

			/* save data */
			uuid, err := uuid.NewV4()
			if err != nil {
				return err
			}
			created_at := time.Now()
			user := User{uuid, created_at, reqBody.Firstname, reqBody.Lastname, reqBody.Email, string(hashedPassword)}
			fmt.Println(reqBody.Email)
			model_user := models.User{}
			existEmail := db.DB.Find(&model_user, "email = ?", user.Email)
			if existEmail.RowsAffected == 0 {

				db.DB.Create(&user)
				return c.Status(200).JSON(fiber.Map{
					"message": "Register successfully",
					"status":  200,
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

	/* Login */
	router.Post("/login", func(c *fiber.Ctx) error {
		reqBody := Login{}
		if err := c.BodyParser(&reqBody); err != nil {
			return err
		}
		fmt.Println(reqBody.Email)
		model_user := models.User{}
		existUser := db.DB.Find(&model_user, "email = ?", reqBody.Email)
		fmt.Println("exist user:", existUser.RowsAffected)
		isValidPass := CheckPasswordHashed(reqBody.Password, model_user.Password)
		fmt.Println("is valid pass:", isValidPass)
		if existUser.RowsAffected == 1 && isValidPass {
			mySecret := os.Getenv("JWT_SECRET")
			expirationTime := time.Now().Add(24 * time.Hour)
			claims := &Claims{
				Id: uuid.UUID(model_user.Id),
				StandardClaims: jwt.StandardClaims{
					// In JWT, the expiry time is expressed as unix milliseconds
					ExpiresAt: expirationTime.Unix(),
				},
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString([]byte(mySecret))
			fmt.Println("error:", err)

			dataReturned := ReturnUser{model_user.Id.String(), model_user.Email, model_user.Firstname, model_user.Lastname, tokenString}

			return c.Status(200).JSON(fiber.Map{
				"message": "Login successfully",
				"status":  200,
				"user":    dataReturned,
			})
		}

		return c.Status(401).JSON(fiber.Map{
			"message": "Invalid credentials",
			"status":  401,
		})

	})
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}
}
