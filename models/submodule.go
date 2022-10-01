package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Submodule struct {
	gorm.Model

	Id          uuid.UUID `gorm:"type:uuid;primary_key;"`
	UserId      string
	ModuleId    string
	Name        string `gorm:"not null"`
	Description string `gorm:"default:null"`
	Private     bool   `gorm:"default:false"`
	Content     string
	Created_at  time.Time `gorm:"<-:create"`
}
