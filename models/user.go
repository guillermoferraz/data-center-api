package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Id         uuid.UUID `gorm:"type:uuid;primary_key;"`
	Firstname  string    `gorm:"not null"`
	Lastname   string    `gorm:"not null"`
	Email      string    `gorm:"not null;unique_index"`
	Password   string    `gorm:"not null"`
	Created_at time.Time `gorm:"<-:create"`
}
