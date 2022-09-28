package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Module struct {
	gorm.Model

	Id          uuid.UUID `gorm:"type:uuid;primary_key;"`
	UserId      string
	Name        string    `gorm:"not null"`
	Description string    `gorm:"default:null"`
	Private     bool      `gorm:"default:false"`
	Created_at  time.Time `gorm:"<-:create"`
}
