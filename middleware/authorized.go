package middleware

import (
	"fmt"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/guillermoferraz/data-center-api/db"
	"github.com/guillermoferraz/data-center-api/models"
	"github.com/joho/godotenv"
)

func UseIsAuthorized(token string) string {
	loadEnv()
	var err error
	if err != nil {
		log.Fatal(err)
		return "Error"
	} else {
		mySecret := os.Getenv("JWT_SECRET")

		claims := jwt.MapClaims{}
		tokenVerification, err := jwt.ParseWithClaims(token, claims, func(tokenVerification *jwt.Token) (interface{}, error) {
			return []byte(mySecret), nil
		})
		fmt.Println(tokenVerification)
		fmt.Println(err)
		model_user := models.User{}
		for key, val := range claims {
			if key == "Id" {
				isValidToken := db.DB.Find(&model_user, "id = ?", val)
				if isValidToken.RowsAffected == 1 {
					return model_user.Id.String()
				} else {
					return "Error"
				}
			}
		}
	}
	return "Error"
}
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}
}
