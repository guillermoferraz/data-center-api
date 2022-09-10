package middlewares

import (
	"fmt"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/guillermoferraz/data-center-api/db"
	"github.com/guillermoferraz/data-center-api/models"
	"github.com/joho/godotenv"
)

var model_user models.User

func isAuthorized(token string) {
	loadEnv()
	mySecret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{}
	tokenVerification, err := jwt.ParseWithClaims(token, claims, func(tokenVerification *jwt.Token) (interface{}, error) {
		return []byte(mySecret), nil
	})
	fmt.Println(tokenVerification)
	fmt.Println(err)
	for key, val := range claims {
		if key == "id" {
			isValidToken := db.DB.Find(&model_user, "id = ?", val)
			if isValidToken.RowsAffected == 1 {
				return
			}
		}
	}
}
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}
}
