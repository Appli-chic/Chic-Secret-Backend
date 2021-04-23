package model

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type LoginToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	Token     int
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	ExpireAt  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (loginToken *LoginToken) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.NewV4()
	return scope.SetColumn("ID", uuid)
}
