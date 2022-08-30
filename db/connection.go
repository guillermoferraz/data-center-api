package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro on loading .env")
	}
}

var DB *gorm.DB

func DBConnection() {
	loadEnv()
	DB_PORT := os.Getenv("DB_PORT")
	DB_NAME := os.Getenv("DB_NAME")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_HOST := os.Getenv("DB_HOST")
	DB_USERNAME := os.Getenv("DB_USERNAME")
	DSN := "host=" + DB_HOST + " " + "user=" + DB_USERNAME + " " + "password=" + DB_PASSWORD + " " + "dbname=" + DB_NAME + " " + "port=" + DB_PORT

	var error error
	DB, error = gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if error != nil {
		log.Fatal(error)
	} else {
		log.Println("\n ")
		log.Println("✅ DB is connected")
	}
}
