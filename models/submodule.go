package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Submodule struct {
	gorm.Model

	Id          uuid.UUID `gorm:"type:uuid;primary_key"`
	UserId      string
	ModuleId    string
	Created_at  time.Time `gorm:"<-:create"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Private     bool      `json:"private"`
}
